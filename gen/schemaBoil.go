// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// hand-written schema boilerplate functions.
// should live in schema.go, but zebrapack chokes on them.

package gen

func (a *Activity) Marshal() ([]byte, error) {
	return a.MarshalMsg([]byte{})
}

func (a *Activity) Unmarshal(b []byte) error {
	_, err := a.UnmarshalMsg(b)
	return err
}

func (u *User) Marshal() ([]byte, error) {
	return u.MarshalMsg([]byte{})
}

func (u *User) Unmarshal(b []byte) error {
	_, err := u.UnmarshalMsg(b)
	return err
}

////////////////////////////////////////////////////////
// XXX - delete

func (a *Activity00) Marshal() ([]byte, error) {
	return a.MarshalMsg([]byte{})
}

func (a *Activity00) Unmarshal(b []byte) error {
	_, err := a.UnmarshalMsg(b)
	return err
}
