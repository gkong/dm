// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package main

import (
	"bufio"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// REST API metrics
// one prometheus "vec" object per metric, with a label per API endpoint.
// individual (per-endpoint) observer objects are cached in maps, for efficiency.
//
// to determine the endpoint for a request, we simply take the first segment
// of the URL path, as determined by helper function pathFirstSegment().

var reqCtrVec *prometheus.CounterVec
var reqCtrs = make(map[string]prometheus.Counter)

var handlerLatSumVec *prometheus.SummaryVec
var handlerLatObservers = make(map[string]prometheus.Observer)

var clientLatSumVec *prometheus.SummaryVec
var clientLatObservers = make(map[string]prometheus.Observer) // also used for validation

// data is reported via an http handler, which prometheus scrapes periodically

var prHandler http.Handler

func metricsSetup() {
	preg := prometheus.NewRegistry()

	prHandler = promhttp.HandlerFor(preg, promhttp.HandlerOpts{
		ErrorLog:      logErrFmt,
		ErrorHandling: promhttp.ContinueOnError,
	})

	reqCtrVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dm_http_requests_total",
			Help: "Count of HTTP requests processed, by API endpoint.",
		},
		[]string{"endpoint"},
	)
	preg.MustRegister(reqCtrVec)

	handlerLatSumVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "dm_http_request_handler_latency_seconds",
			Help:       "HTTP request latency, by API endpoint, as measured by handlers.",
			Objectives: map[float64]float64{0.5: 0.05, 0.99: 0.001, 0.999: 0.0001},
		},
		[]string{"endpoint"},
	)
	preg.MustRegister(handlerLatSumVec)

	clientLatSumVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "dm_http_request_client_latency_seconds",
			Help:       "HTTP request latency, by API endpoint, as reported by clients.",
			Objectives: map[float64]float64{0.5: 0.05, 0.99: 0.001, 0.999: 0.0001},
		},
		[]string{"endpoint"},
	)
	preg.MustRegister(clientLatSumVec)

	preg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{Namespace: "dm"}))
	preg.MustRegister(&linuxCollector{})
}

// register a REST API endpoint
func metricsRegisterEndpoint(endpoint string, handlerLatency, clientLatency bool) {
	if _, ok := reqCtrs[endpoint]; ok {
		// if multiple routes with the same first path segment, only register once
		return
	}

	reqCtrs[endpoint] = reqCtrVec.WithLabelValues(endpoint)

	if handlerLatency {
		handlerLatObservers[endpoint] = handlerLatSumVec.WithLabelValues(endpoint)
	}

	if clientLatency {
		clientLatObservers[endpoint] = clientLatSumVec.WithLabelValues(endpoint)
	}
}

// record a latency report from a handler.
func metricsHTTPRequest(endpoint string, timeUsec int) {
	endpoint = pathFirstSegment(endpoint)

	c, ok := reqCtrs[endpoint]
	if !ok {
		// XXX - log error
		return
	}
	c.Inc()

	obs, ok := handlerLatObservers[endpoint]
	if !ok {
		// we don't track all endpoints, so it's OK if not found
		return
	}
	obs.Observe(float64(timeUsec) / 1000000)
}

// record a latency report from a client
func metricsClientLatency(endpoint string, timeMsec int) {
	if timeMsec < 0 {
		return
	}

	// we don't track all endpoints; check if this is one we're interested in
	obs, ok := clientLatObservers[endpoint]
	if !ok {
		return
	}

	obs.Observe(float64(timeMsec) / 1000)
}

// system metrics

// a prometheus Collector for a minimal set of linux system metrics
type linuxCollector struct{}

const ticksPerSec = 100 // C.sysconf(C._SC_CLK_TCK)
const bytesPerKB = 1024

var cpuDesc = prometheus.NewDesc("dm_cpu_idle_seconds", "aggregate cpu idle time", []string{}, nil)
var memTotalDesc = prometheus.NewDesc("dm_mem_total_bytes", "size of physical memory", []string{}, nil)
var memAvailableDesc = prometheus.NewDesc("dm_mem_available_bytes", "memory available", []string{}, nil)

func (c *linuxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cpuDesc
	ch <- memTotalDesc
	ch <- memAvailableDesc
}

func (c *linuxCollector) Collect(ch chan<- prometheus.Metric) {
	collectCPU(ch)
	collectMem(ch)
}

func collectCPU(ch chan<- prometheus.Metric) {
	ncpu := 0
	idleticks := -1

	procstat, err := os.Open(path.Join(Config.SysProc, "stat"))
	if err != nil {
		// XXX log
		return
	}
	defer procstat.Close()

	scanner := bufio.NewScanner(procstat)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		if len(words) < 1 {
			continue
		}
		if words[0] == "cpu" {
			if len(words) < 5 {
				// err := errors.New("collectCPU - malformed cpu line in /proc/stat")
				// XXX log
				return
			}
			idleticks, err = strconv.Atoi(words[4])
			if err != nil {
				// XXX log
				return
			}
			continue
		}
		if strings.HasPrefix(words[0], "cpu") {
			ncpu++
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		// XXX log
		return
	}

	if ncpu == 0 || idleticks < 0 {
		// XXX log
		return
	}

	ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, float64(idleticks)/float64(ncpu*ticksPerSec))
}

func collectMem(ch chan<- prometheus.Metric) {
	var total, available int
	var seenTotal, seenAvailable bool

	procmeminfo, err := os.Open(path.Join(Config.SysProc, "meminfo"))
	if err != nil {
		// XXX log
		return
	}
	defer procmeminfo.Close()

	scanner := bufio.NewScanner(procmeminfo)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		if len(words) < 1 {
			continue
		}
		if words[0] == "MemTotal:" {
			if len(words) != 3 || words[2] != "kB" {
				// err := errors.New("collectMem - malformed MemTotal line in /proc/meminfo")
				// XXX log
				return
			}
			total, err = strconv.Atoi(words[1])
			if err != nil {
				// XXX log
				return
			}
			seenTotal = true
		}
		if words[0] == "MemAvailable:" {
			if len(words) != 3 || words[2] != "kB" {
				// err := errors.New("collectMem - malformed MemAvailable line in /proc/meminfo")
				// XXX log
				return
			}
			available, err = strconv.Atoi(words[1])
			if err != nil {
				// XXX log
				return
			}
			seenAvailable = true
		}
		if seenTotal && seenAvailable {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		// XXX log
		return
	}

	if (!seenTotal) || (!seenAvailable) {
		// XXX log
		return
	}

	ch <- prometheus.MustNewConstMetric(memTotalDesc, prometheus.GaugeValue, float64(total*bytesPerKB))
	ch <- prometheus.MustNewConstMetric(memAvailableDesc, prometheus.GaugeValue, float64(available*bytesPerKB))
}
