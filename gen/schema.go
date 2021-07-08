//go:generate ffjson -w=schema_ffjson_generated.go schema.go

//go:generate zebrapack -o=schema_zebra_generated.go -io=false -no-structnames-onwire -tests=false

//msgp:ignore Plan
// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package gen

const (
	int64size  = 8
	prefixSize = 1 // all prefixes must be exactly one byte long.
)

// key prefixes (correspond to database tables).
// number sequentially, starting from zero, and keep in sync with PrefixNames.
var (
	// managed by qsess
	PfxSessLogin    = []byte{0}
	PfxSessVerify   = []byte{1}
	PfxSessEmail    = []byte{2}
	PfxSessRecovery = []byte{3}

	// structured data types defined below
	PfxUser      = []byte{4}
	PfxActivity  = []byte{5} // one user's execution of a plan
	PfxLastVisit = []byte{6} // contains just Unix secs since 1970-01-01
)

var PrefixNames = [...]string{
	"SessLogin",
	"SessVerify",
	"SessEmail",
	"SessRecovery",
	"User",
	"Activity",
	"LastVisit",
}

//
// User table
//

// use email addr for user key.
// since we're going to do prefix searches, append '/', illegal in email addrs,
// to prevent unwanted hits when one email addr is prefix of another.
func UserKey(email []byte) []byte {
	return bscat(PfxUser, email, []byte{'/'})
}

func EmailFromUserKey(ukey []byte) []byte {
	return ukey[len(PfxUser) : len(ukey)-1]
}

// ffjson: skip
type User struct {
	Password      []byte `zid:"0"` // empty before sending to client!
	FirstName     string `zid:"1"`
	LastName      string `zid:"2"`
	BibleVersion  string `zid:"3"`
	BibleProvider string `zid:"4"`

	LastNotificationSeen int64 `zid:"5"` // in Unix secs since 1970-01-01

	CreatedTime int64  `zid:"6"` // in Unix secs since 1970-01-01
	OrigEmail   string `zid:"7"` // email addr when acct was created
}

//
// Activity table
//

// activity key = userkey + PlanName
func ActivityKey(userkey []byte, planname []byte) []byte {
	return bscat(PfxActivity, userkey, planname)
}

// XXX - delete
type Activity00 struct {
	PlanName                string `zid:"0" json:"plan"`   // redundantly stored here and in key for convenience
	CurrentDay              int    `zid:"1" json:"curday"` // day (relative to plan) currently working on (starting with 1)
	Done                    []bool `zid:"2" json:"done"`   // done state for each stream (as of current day)
	BibleVersion            string `zid:"3" json:"version"`
	BibleProvider           string `zid:"4" json:"provider"`
	AccountabilityStartDate int    `zid:"5" json:"accstartdate"` // client's idea of days since 1970-01-01
	AccountabilityVisible   bool   `zid:"6" json:"accvisible"`
}

type Activity struct {
	PlanName string `zid:"0" json:"plan"` // redundantly stored here and in key for convenience
	// Day is an array which stores, for each stream in the plan,
	// the day (relative to plan) currently on display (starting with 1).
	// Its value may be one more than the number of days in the plan,
	// indicating that the given stream has been completed.
	Day                     []int  `zid:"1" json:"day"`
	BibleVersion            string `zid:"2" json:"version"`
	BibleProvider           string `zid:"3" json:"provider"`
	AccountabilityStartDate int    `zid:"4" json:"accstartdate"` // client's idea of days since 1970-01-01
	AccountabilityVisible   bool   `zid:"5" json:"accvisible"`
}

//
// LastVisit table
//

// last-visit key is same as user key, except with a different prefix.
// these functions convert one to the other, in-place.

func ConvertUKeyToLVKey(key []byte) {
	copy(key, PfxLastVisit)
}

func ConvertLVKeyToUKey(key []byte) {
	copy(key, PfxUser)
}

// last-visit table content is just an int64 = Unix secs since 1970-01-01.

//
// Plan files
//

// plans are JSON files which live in the document store.
// we only need to deserialize them in order to pre-compute some docqueries.

// ffjson: noencoder
type Plan struct {
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	TotalDays int      `json:"days"`
	Streams   []string `json:"streams"`
}

// helpers

// concatenate byte slices, returning a new one containing the contents of all
func bscat(bb ...[]byte) []byte {
	size := 0
	for _, b := range bb {
		size += len(b)
	}

	ret := make([]byte, size)
	pos := 0
	for _, b := range bb {
		copy(ret[pos:], b)
		pos += len(b)
	}

	return ret
}
