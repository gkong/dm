//go:generate zebrapack -o=sessdata_generated.go -no-structnames-onwire -tests=false

// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package gen

// NOTE: for Login and Recovery sessions, the only data we need to store is
// the user id, which can be persisted simply by supplying it as argument
// to NewSession and retrieved by calling UserID, so we don't need any
// implementation of SessData

// track emailed tokens for new account signup
type VerifySessData struct {
	UserKey   []byte `zid:"0"`
	FirstName string `zid:"1"`
	LastName  string `zid:"2"`
	Password  []byte `zid:"3"` // hash, not plaintext, of course
}

// track emailed tokens for email address changes
type EmailSessData struct {
	OldUserKey []byte `zid:"0"`
	NewUserKey []byte `zid:"1"`
}
