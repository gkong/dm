// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package gen

import (
	"github.com/gkong/go-qweb/qsess"
)

// verify sessions

func NewVerifySessData() qsess.SessData {
	return &VerifySessData{}
}

func (m *VerifySessData) Marshal() ([]byte, error) {
	return m.MarshalMsg([]byte{})
}

func (m *VerifySessData) Unmarshal(b []byte) error {
	_, err := m.UnmarshalMsg(b)
	return err
}

// email sessions

func NewEmailSessData() qsess.SessData {
	return &EmailSessData{}
}

func (m *EmailSessData) Marshal() ([]byte, error) {
	return m.MarshalMsg([]byte{})
}

func (m *EmailSessData) Unmarshal(b []byte) error {
	_, err := m.UnmarshalMsg(b)
	return err
}
