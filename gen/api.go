// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// declarations of request and response structs for rest api
//
// due to an ffjson constraint, this can't be in package main.

//go:generate ffjson -w=api_generated.go api.go

package gen

// ffjson: noencoder
type AdminRequest struct {
	Cmd string `json:"cmd"`
	Arg string `json:"arg"`
}

// ffjson: noencoder
type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// ffjson: noencoder
type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	KeepLoggedIn bool   `json:"keep"`
}

// ffjson: noencoder
type PWRecoverRequest struct {
	Email string `json:"email"`
}

// ffjson: noencoder
type DoRecoverRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ffjson: noencoder
type ContactRequest struct {
	Msg string `json:"msg"`
}

// StateResponse is returned by multiple REST API endpoints

// ffjson: nodecoder
type StateResponse struct {
	// this response is returned by multiple endpoints
	FirstName  string      `json:"firstname"`
	Activities []*Activity `json:"activities"`
}

// both a request and a response
type UserReqResp struct {
	Email         string `zid:"0" json:"email"`
	OldPassword   string `zid:"1" json:"oldpass"` // empty before sending to client!
	NewPassword   string `zid:"2" json:"newpass"`
	FirstName     string `zid:"3" json:"firstname"`
	LastName      string `zid:"4" json:"lastname"`
	BibleVersion  string `zid:"5" json:"version"`
	BibleProvider string `zid:"6" json:"provider"`
}

// ffjson: nodecoder
type UserSetResponse struct {
	Logout          bool   `json:"logout"`
	EmailSent       bool   `json:"emailsent"`
	PasswordChanged bool   `json:"pwchanged"`
	FirstName       string `json:"firstname"`
}

// ffjson: noencoder
type ActivityAddRequest struct {
	Plan  string `json:"plan"`
	Today int    `json:"today"` // client's idea of days since 1970-01-01
}

// ffjson: noencoder
type ActivityDeleteRequest struct {
	Plan string `json:"plan"`
}

// change the day value for one stream of an activity

// ffjson: noencoder
type DayChangeRequest struct {
	Plan        string `json:"plan"`
	StreamIndex int    `json:"streamindex"`
	PrevDay     int    `json:"prevday"` // client's idea of value before this change is applied
	Delta       int    `json:"delta"`   // currently, only 1 or -1 are allowed
	Today       int    `json:"today"`   // in case need to reset AccountabilityStartDate
}

// ffjson: noencoder
type ActivityJumpRequest struct {
	Plan  string `json:"plan"`
	Day   int    `json:"day"`
	Today int    `json:"today"` // to reset AccountabilityStartDate
}

// ffjson: noencoder
type ActivityVersionRequest struct {
	Plan          string `json:"plan"`
	BibleProvider string `json:"provider"`
	BibleVersion  string `json:"version"`
}

// DayResponse is returned by multiple REST API endpoints,
// to update a given activity with a day change.

// ffjson: nodecoder
type DayResponse struct {
	Day                     []int `json:"day"`
	AccountabilityStartDate int   `json:"accstartdate"`
	AccountabilityVisible   bool  `json:"accvisible"`
}

// ffjson: noencoder
type AccountabilityResetRequest struct {
	Plan  string `json:"plan"`
	Today int    `json:"today"` // to reset AccountabilityStartDate
}

// ffjson: noencoder
type AccountabilityEnabledRequest struct {
	Plan    string `json:"plan"`
	Enabled bool   `json:"enabled"`
}

// ffjson: noencoder
type LatencyRequest struct {
	Endpoint string `json:"endpoint"`
	TimeMsec int    `json:"timemsec"`
}

// ffjson: nodecoder
type DocWrapper struct {
	TTLSecs int    `json:"ttl"`
	Data    string `json:"data"`
}

// documents in the doc store, returned by /getdoc API endpoint

// ffjson: nodecoder
type PlanDesc struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// ffjson: nodecoder
type PlandescsDocument struct {
	PlanDescs []PlanDesc `json:"plandescs"`
}

// ffjson: nodecoder
type Provider struct {
	Provider    string `json:"provider"`
	URLTemplate string `json:"url"`
}

// ffjson: nodecoder
type ProvidersDocument struct {
	Providers []Provider `json:"providers"`
}
