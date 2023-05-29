// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package main

import (
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/gkong/go-qweb/qctx"
)

// CSRF token = base64-encode(secs-since-1970) + base64-encode(random-data)

var (
	csrfDateLen          int
	csrfRandLen          int
	csrfLen              int
	csrfDateMaxDecodeLen int
)

func csrfSetup() {
	csrfDateLen = base64.URLEncoding.EncodedLen(bytesPerUint64)
	csrfRandLen = base64.URLEncoding.EncodedLen(Config.CSRFRandLen)
	csrfLen = csrfDateLen + csrfRandLen
	csrfDateMaxDecodeLen = base64.URLEncoding.DecodedLen(csrfDateLen)
}

func csrfSetCookie(w http.ResponseWriter) {
	nowBytes := make([]byte, bytesPerUint64)
	binary.LittleEndian.PutUint64(nowBytes, uint64(time.Now().Unix()))

	rnd := make([]byte, Config.CSRFRandLen)
	randomBytes(rnd)

	encoded := make([]byte, csrfLen)
	base64.URLEncoding.Encode(encoded, nowBytes)
	base64.URLEncoding.Encode(encoded[csrfDateLen:], rnd)

	http.SetCookie(w, &http.Cookie{
		Name:     Config.CSRFCookieName,
		Value:    string(encoded),
		Path:     "/",
		Domain:   "",
		MaxAge:   Config.CSRFMaxAgeSecs,
		Secure:   Config.CSRFCookieSecure,
		HttpOnly: false, // client JS must be able to read this for "cookie to header" CSRF protection
		SameSite: http.SameSiteLaxMode,
	})
}

// middleware to set a CSRF cookie
func mwSetCSRF(next qctx.CtxHandler) qctx.CtxHandler {
	return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
		csrfSetCookie(c.W)
		next.CtxServeHTTP(c)
	})
}

func csrfErr(c *qctx.Ctx, msg string, err error) {
	logWarn(msg, err, c)
	c.Error("CSRF", http.StatusForbidden)
	csrfSetCookie(c.W)
}

// middleware to guard against CSRF attacks.
// check that CSRF cookie and token match, and that they are not expired.
// also perform CSRF cookie refresh, if required.
// if Method == GET, allow the request.
//
// to avoid false positives, never let CSRF cookie expire before session:
// set csrf max age >= session max age and csrf min refresh <= session min refresh.
func mwCheckCSRF(next qctx.CtxHandler) qctx.CtxHandler {
	return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
		if c.R.Method == "GET" {
			next.CtxServeHTTP(c)
			return
		}

		cookie, err := c.R.Cookie(Config.CSRFCookieName)
		if err != nil {
			csrfErr(c, "no CSRF cookie", err)
			return
		}

		if len(cookie.Value) != csrfLen {
			csrfErr(c, "CSRF cookie wrong size", err)
			return
		}

		dateBytes := make([]byte, csrfDateMaxDecodeLen)
		_, err = base64.URLEncoding.Decode(dateBytes, ([]byte(cookie.Value))[:csrfDateLen])
		if err != nil {
			csrfErr(c, "bad CSRF cookie", err)
			return
		}
		date := int64(btou64(dateBytes))

		now := time.Now().Unix()
		if now-date > int64(Config.CSRFMaxAgeSecs) {
			// this check might be unnecessary; client will expire the cookie
			c.Error("CSRF", http.StatusForbidden)
			csrfSetCookie(c.W)
			return
		}
		if now-date > int64(Config.CSRFMinRefreshSecs) {
			// refresh
			csrfSetCookie(c.W)
		}

		hdr := c.R.Header.Get(Config.CSRFHeader)
		if hdr == "" {
			logWarn("no CSRF header - POSSIBLE ATTACK", nil, c)
			c.Error("CSRF", http.StatusForbidden)
			return
		}
		if hdr != cookie.Value {
			logWarn("CSRF header cookie mismatch", nil, c)
			c.Error("CSRF", http.StatusForbidden)
			return
		}

		next.CtxServeHTTP(c)
	})
}
