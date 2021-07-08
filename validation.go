// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package main

import (
	"errors"
	"html"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

func sanitizeName(s string) string {
	return html.EscapeString(strings.TrimSpace(s))
}

// enable user code to invoke multiple validators with only one error check.
// the price we pay for this user code simplification is non-lazy evaluation.
func firstError(errors ...error) error {
	for i := range errors {
		if errors[i] != nil {
			return errors[i]
		}
	}
	return nil
}

func validEmailAddr(s string) error {
	if !govalidator.IsEmail(s) {
		return errors.New("not a valid email address")
	}
	return nil
}

func validPassword(s string) error {
	if len(s) < 8 || len(s) > 100 {
		return errors.New("password must be 8 to 100 chars long")
	}
	return nil
}

func validFirstName(s string) error {
	if len(s) == 0 {
		return errors.New("first name may not be empty")
	}
	if len(s) > 50 {
		return errors.New("first name may not be longer than 50 characters")
	}
	return nil
}

func validLastName(s string) error {
	if len(s) == 0 {
		return errors.New("last name may not be empty")
	}
	if len(s) > 50 {
		return errors.New("last name may not be longer than 50 characters")
	}
	return nil
}

func validEmailMsg(s string) error {
	if len(s) == 0 {
		return errors.New("message may not be empty")
	}
	if len(s) > 1000 {
		return errors.New("message may not be longer than 1000 characters")
	}
	return nil
}

var pathSegmentRegexp = regexp.MustCompile("^[a-zA-Z0-9~:=@_\\+\\-\\.\\$\\&]+$")

func validBibleVersion(s string) error {
	if len(s) > 50 {
		return errors.New("Bible version code is too long")
	}
	if !pathSegmentRegexp.MatchString(s) {
		return errors.New("illegal characters in Bible version code")
	}
	return nil
}
