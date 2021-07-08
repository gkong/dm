// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// queries derived from documents from the doc store.
// some are saved here in global maps; some in the doc store "queries" bucket.

// XXX - currently only runs during program initialization. need to support on-going updates.

package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	. "github.com/gkong/dm/gen"
)

// maps for use by server code

var streamDays = make(map[string]int)  // number of days each stream requires
var planStreams = make(map[string]int) // number of streams each plan has
var planDays = make(map[string]int)    // number of days each plan requires
var providers = make(map[string]int)   // just use for key existence

func docQueriesSetup(ttlsecs int) {
	var err error

	// docstore - queries/providers, providers

	provs := ProvidersDocument{Providers: []Provider{}}
	provNames := docNames("provider")
	for _, provName := range provNames {
		providers[provName] = 0
		pdoc, _, err := docGet("provider", provName)
		dqerrcheck(err, "queries/providers docGet")

		p := Provider{
			Provider:    provName,
			URLTemplate: string(pdoc),
		}
		provs.Providers = append(provs.Providers, p)
	}
	jdata, err := provs.MarshalJSON()
	dqerrcheck(err, "queries/providers MarshalJSON")
	err = docPut("queries", "providers", ttlsecs, jdata)
	dqerrcheck(err, "queries/providers docPut")

	// docstore - queries/plandescs

	r := PlandescsDocument{PlanDescs: []PlanDesc{}}
	names := docNames("plan")
	for _, planName := range names {
		plandoc, _, err := docGet("plan", planName)
		dqerrcheck(err, "planDescsJSON - docGet")

		p := &Plan{}
		err = p.UnmarshalJSON(plandoc)
		dqerrcheck(err, "planDescsJSON UnmarshalJSON")

		pd := PlanDesc{}
		pd.Name = planName
		pd.Title = p.Title
		pd.Desc = p.Desc

		r.PlanDescs = append(r.PlanDescs, pd)
	}
	jdata, err = r.MarshalJSON()
	dqerrcheck(err, "planDescs MarshalJSON")
	err = docPut("queries", "plandescs", ttlsecs, jdata)
	dqerrcheck(err, "queries/plandescs docPut")

	// streamDays

	streamNames := docNames("stream")
	for _, streamName := range streamNames {
		s, _, err := docGet("stream", streamName)
		dqerrcheck(err, "streamDays - docGet")

		// count lines in stream document
		streamDays[streamName] = bytes.Count(s, []byte{'\n'})

		// compare to number at the end of the stream name
		nameNumber, _ := strconv.Atoi(streamName[strings.LastIndexByte(streamName, '-')+1:])
		if nameNumber != streamDays[streamName] {
			logPanic(fmt.Sprintf("docQueriesSetup - stream %s line count (%d) does not match name", streamName, streamDays[streamName]), nil, nil)
		}
	}

	// planStreams, planDays - streamDays must be calculated BEFORE planDays

	planNames := docNames("plan")
	for _, planName := range planNames {
		plandoc, _, err := docGet("plan", planName)
		dqerrcheck(err, "planStreams - docGet")

		p := &Plan{}
		err = p.UnmarshalJSON(plandoc)
		dqerrcheck(err, "Plan UnmarshalJSON")
		planStreams[planName] = len(p.Streams)

		// go thru all streams in the plan and take the max day count
		days := 0
		for _, stream := range p.Streams {
			if streamDays[stream] > days {
				days = streamDays[stream]
			}
		}
		planDays[planName] = days
	}
}

func dqerrcheck(err error, msg string) {
	if err != nil {
		logPanic("docQueriesSetup - "+msg, err, nil)
	}
}
