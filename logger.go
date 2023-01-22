// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// guidelines for choosing severity levels
//	info - metrics, in early production - every HTTP request
//	warning - a problem which should be investigated (e.g. unexpected
//			client behavior, which might indicate a client bug)
//	error - a definite server problem, which must be fixed
//  panic - can't continue, generally because of a problem during initialization

package main

import (
	"os"
	"strings"
	"sync"
	"time"

	. "github.com/gkong/dm/gen"

	"github.com/gkong/go-qweb/qctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// code outside this file typically only calls the following 4 functions:

func logInfo(msg string, c *qctx.Ctx) {
	zlogger.Info(msg, getFields(nil, c)...)
}

func logWarn(msg string, err error, c *qctx.Ctx) {
	zlogger.Warn(msg, getFields(err, c)...)
}

func logErr(msg string, err error, c *qctx.Ctx) {
	zlogger.Error(msg, getFields(err, c)...)
	emailAlertThrottled(alertLevelError)
}

func logPanic(msg string, err error, c *qctx.Ctx) {
	// email first, because nothing after zlog.Panic gets executed
	emailAlertThrottled(alertLevelPanic)
	zlogger.Panic(msg, getFields(err, c)...)
}

var zlogger *zap.Logger
var zsugar *zap.SugaredLogger

// give these to libraries that want an io.Writer for logging.
// these often result in wasteful conversions of string->[]byte->string,
// so they should not be used in any high-volume situations.
var logInfoWriter zInfoWriter
var logErrWriter zErrWriter
var logFilterWriter zFilterWriter

type zInfoWriter struct {
	logger *zap.Logger
}

type zErrWriter struct {
	logger *zap.Logger
}

type zFilterWriter struct {
	logger *zap.Logger
}

func (z zInfoWriter) Write(p []byte) (n int, err error) {
	z.logger.Info(string(p))
	return len(p), nil
}

func (z zErrWriter) Write(p []byte) (n int, err error) {
	z.logger.Error(string(p))
	emailAlertThrottled(alertLevelError)
	return len(p), nil
}

// std library doesn't distinguish levels.
// classify each invocation as error or info.
func (z zFilterWriter) Write(p []byte) (n int, err error) {
	s := string(p)

	// only the default case results in an error
	switch {
	case strings.Contains(s, "TLS handshake error"):
	case strings.Contains(s, "http2: server: error reading preface from client"):
	case strings.Contains(s, "URL query contains semicolon"):
	default:
		z.logger.Error("std log - " + s)
		emailAlertThrottled(alertLevelError)
		return len(p), nil
	}

	z.logger.Info("std log - " + s)
	return len(p), nil
}

// give this to libraries that want a fmt-style interface for logging
var logErrFmt zErrSugared

type zErrSugared struct {
	sugar *zap.SugaredLogger
}

func (z zErrSugared) Println(v ...interface{}) {
	z.sugar.Error(v...)
	emailAlertThrottled(alertLevelError)
}

// zap's "Production" paradigm: level, msg, ts (floating-point secs since 1970)
func logSetup() {
	if Config.LogDest == "stdout" {
		c := zap.NewProductionConfig()
		c.Sampling = nil
		c.DisableCaller = true
		c.DisableStacktrace = true
		c.OutputPaths = []string{"stdout"}
		c.ErrorOutputPaths = []string{"stdout"}
		zlogger, _ = c.Build()
	} else {
		// verify we can write to the specified log file
		f, err := os.OpenFile(Config.LogDest, os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic("logSetup - cannot open log file - " + err.Error())
		}
		f.Close()

		ljlog := &lumberjack.Logger{
			Filename:   Config.LogDest,
			MaxSize:    Config.LumberjackMaxSize,
			MaxBackups: Config.LumberjackMaxBackups,
			MaxAge:     Config.LumberjackMaxAge,
		}

		ljWS := zapcore.Lock(zapcore.AddSync(ljlog))
		enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(enc, ljWS, zap.NewAtomicLevel())
		zlogger = zap.New(core)
		zsugar = zlogger.Sugar()
	}

	logErrWriter = zErrWriter{zlogger}
	logInfoWriter = zInfoWriter{zlogger}
	logFilterWriter = zFilterWriter{zlogger}
	logErrFmt = zErrSugared{zsugar}
}

// simple facility to email alerts

const (
	alertLevelError = iota
	alertLevelPanic
	alertNumLevels
)

type alertLev struct {
	lastTime int64
	mu       sync.Mutex
}

var alertLevs [alertNumLevels]alertLev

// for security, emailed messages are generic, with no error-specific content.
// send same message as both subject and body.
func emailAlertThrottled(level int) {
	doit := false
	alertLevs[level].mu.Lock()
	{
		now := time.Now().Unix()
		if alertLevs[level].lastTime == 0 || now-alertLevs[level].lastTime > int64(Config.AlertEmailIntervalSecs) {
			alertLevs[level].lastTime = now
			doit = true
		}
	}
	alertLevs[level].mu.Unlock()
	if doit {
		var msg string
		switch level {
		case alertLevelError:
			msg = Config.AlertEmailMsgError
		case alertLevelPanic:
			msg = Config.AlertEmailMsgPanic
		default:
			logErr("emailAlertThrottled - bad level", nil, nil)
			return
		}
		if Config.DebugEmail {
			logInfo("emailAlertThrottled - DEBUG - "+msg, nil)
		} else {
			go sendTextEmail(Config.AdminEmailAddr, msg, msg)
		}
	}
}

// after checking logs, call this to reset alert intervals,
// allowing alerts to be sent immediately again.
func alertReset() {
	for level := range alertLevs {
		alertLevs[level].mu.Lock()
		alertLevs[level].lastTime = 0
		alertLevs[level].mu.Unlock()
	}
}

// return a slice of logger fields, if any, from the given error and context,
// both of which can be nil.
func getFields(err error, c *qctx.Ctx) []zapcore.Field {
	fields := make([]zapcore.Field, 0, 5)

	if err != nil {
		fields = append(fields, zErr(err))
	}

	if c != nil {
		fields = append(fields, zIP(c.R.RemoteAddr), zMethod(c.R.Method), zPath(c.R.URL.Path))
		if c.Sess != nil {
			userID := c.Sess.UserID()
			if len(userID) > 0 {
				fields = append(fields, zUser(string(EmailFromUserKey(userID))))
			}
		}
	}

	return fields
}

// application-specific wrappers

func zErr(val error) zapcore.Field {
	return zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: val}
}

func zIP(val string) zapcore.Field {
	return zapcore.Field{Key: "ip", Type: zapcore.StringType, String: val}
}

func zMethod(val string) zapcore.Field {
	return zapcore.Field{Key: "method", Type: zapcore.StringType, String: val}
}

func zReferer(val string) zapcore.Field {
	return zapcore.Field{Key: "referer", Type: zapcore.StringType, String: val}
}

func zPath(val string) zapcore.Field {
	return zapcore.Field{Key: "path", Type: zapcore.StringType, String: val}
}

func zStatus(val int) zapcore.Field {
	return zapcore.Field{Key: "status", Type: zapcore.Int64Type, Integer: int64(val)}
}

func zUser(val string) zapcore.Field {
	return zapcore.Field{Key: "user", Type: zapcore.StringType, String: val}
}

func zLatencyUs(val int64) zapcore.Field {
	return zapcore.Field{Key: "latencyus", Type: zapcore.Int64Type, Integer: val}
}
