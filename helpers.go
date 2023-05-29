// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package main

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"net/http"
	"strings"
)

// fill a given byte slice with random data
func randomBytes(b []byte) error {
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return wrapErr{"randomBytes - io.ReadFull", err}
	}
	return nil
}

// return the first segment of a URL path
func pathFirstSegment(s string) string {
	start := 0
	if s[0] == '/' {
		start = 1
	}

	end := strings.IndexByte(s[start:], '/')
	if end == -1 {
		end = len(s)
	} else {
		end += start
	}

	return s[start:end]
}

// return the last segment of a URL path
func pathLastSegment(s string) string {
	start := strings.LastIndexByte(s, '/')
	if start == -1 {
		return s
	} else {
		return s[start+1:]
	}
}

// tell me if a slice of bools are all true
func alltrue(b []bool) bool {
	for _, val := range b {
		if !val {
			return false
		}
	}
	return true
}

func stringsToByteSlices(s []string) [][]byte {
	b := make([][]byte, len(s))
	for i := range s {
		b[i] = []byte(s[i])
	}
	return b
}

// serialize integer types

const bytesPerUint32 = 4

// convert uint32 to []byte
func u32tob(u uint32) []byte {
	dest := make([]byte, bytesPerUint32)
	binary.LittleEndian.PutUint32(dest, u)
	return dest
}

// convert []byte to uint32
func btou32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

const bytesPerUint64 = 8

// convert uint64 to []byte
func u64tob(u uint64) []byte {
	dest := make([]byte, bytesPerUint64)
	binary.LittleEndian.PutUint64(dest, u)
	return dest
}

// convert []byte to uint64
func btou64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

// error message wrapping with lazy evaluation of error message string

type wrapErr struct {
	msg string
	err error
}

func (e wrapErr) Error() string {
	if e.err != nil {
		return e.msg + " - " + e.err.Error()
	}
	return e.msg
}

// convert string to http.SameSite

func strToSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}
