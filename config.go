// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package main

import (
	"os"
	"path/filepath"

	. "github.com/gkong/dm/gen"

	"github.com/gkong/go-qweb/qsess"
	"github.com/gkong/go-qweb/qsess/qsldb"

	"github.com/BurntSushi/toml"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

// toml file schema

var Config tomlConfig

type tomlConfig struct {

	// secrets

	AdminUser      string
	AdminEmailAddr string

	MetricsUser     string
	MetricsPassword string

	EmailSender   string
	EmailUserName string
	EmailPassword string

	QsessKeys []string // a separate key for each qsess store would allow independent key rotation intervals...

	// differ between debug / test / production

	DebugPanicStack bool
	DebugEmail      bool
	LogAllHTTP      bool

	ClientVersion string

	TCPAddr      string
	ReverseProxy bool
	SSL          bool
	SSLTCPAddr   string
	SSLCertDir   string
	SSLHosts     []string
	BaseURL      string

	CSRFCookieSecure bool

	StaticDir string
	DocDir    string

	DbDir        string
	BckDir       string
	BckLatestDir string
	LogDest      string

	// more or less fixed

	ClientVersionReqHdr string
	ClientUpdateReqHdr  string

	EmailServer string
	EmailPort   string

	DailyBackupTime int

	AlertEmailIntervalSecs int
	AlertEmailMsgError     string
	AlertEmailMsgPanic     string

	VerifLimitPerMin    int
	VerifLimitBurst     int
	LoginLimitPerMin    int
	LoginLimitBurst     int
	ContactLimitPerDay  int
	ContactLimitBurst   int
	NotFoundLimitPerMin int
	NotFoundLimitBurst  int
	NonsslLimitPerMin   int
	NonsslLimitBurst    int

	BcryptCost int

	DocTTL        int
	DocQueriesTTL int

	DefaultBibleVersion  string
	DefaultBibleProvider string

	SysProc string

	LumberjackMaxSize    int
	LumberjackMaxBackups int
	LumberjackMaxAge     int

	CSRFCookieName     string
	CSRFHeader         string
	CSRFRandLen        int
	CSRFMaxAgeSecs     int
	CSRFMinRefreshSecs int

	QsLogin    ConfigQsess
	QsVerify   ConfigQsess
	QsEmail    ConfigQsess
	QsRecovery ConfigQsess
}

type ConfigQsess struct {
	AuthType           string
	MaxAgeSecs         int
	MinRefreshSecs     int
	LongMaxAgeSecs     int
	LongMinRefreshSecs int
	CookieName         string
	CookieSecure       bool
	CookieHTTPOnly     bool
	CookieSameSite     string
}

func doConfig(args []string) {
	if len(args) == 0 {
		panic("doConfig - must give at least one command-line arg")
	}

	var dir string
	var arg int

	fi, err := os.Stat(args[0])
	if err != nil {
		panic("doConfig - cannot Stat first command-line arg - " + err.Error())
	}

	if fi.Mode().IsDir() {
		dir = args[0]
		arg = 1
		if len(args) < 2 {
			panic("doConfig - need at least one config file")
		}
	}

	for ; arg < len(args); arg++ {
		fname := args[arg]
		if (dir != "") && (args[arg][0] != '/') {
			fname = filepath.Join(dir, args[arg])
		}
		// multiple sparse files decoded into one struct
		if _, err := toml.DecodeFile(fname, &Config); err != nil {
			panic("doConfig - toml.Decode failed - " + err.Error())
		}
	}
}

// rate limit PER IP+port for things with emailed verifications - signup and pwrecover
var verifLimiter *throttled.GCRARateLimiterCtx

// rate limit PER IP+port for login attempts
var loginLimiter *throttled.GCRARateLimiterCtx

// rate limit PER USER for contact form submission
var contactLimiter *throttled.GCRARateLimiterCtx

// rate limit PER IP+PORT for URL not found (to discourage excessive crawling)
var notfoundLimiter *throttled.GCRARateLimiterCtx

func throttledSetup() {
	store, err := memstore.New(0)
	if err != nil {
		logPanic("throttledSetup - memstore.New", err, nil)
	}

	quota := throttled.RateQuota{throttled.PerMin(Config.VerifLimitPerMin), Config.VerifLimitBurst}
	verifLimiter, err = throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		logPanic("throttledSetup - verif NewGCRARateLimiter", err, nil)
	}

	quota = throttled.RateQuota{throttled.PerMin(Config.LoginLimitPerMin), Config.LoginLimitBurst}
	loginLimiter, err = throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		logPanic("throttledSetup - login NewGCRARateLimiter", err, nil)
	}

	quota = throttled.RateQuota{throttled.PerDay(Config.ContactLimitPerDay), Config.ContactLimitBurst}
	contactLimiter, err = throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		logPanic("throttledSetup - contact NewGCRARateLimiter", err, nil)
	}

	quota = throttled.RateQuota{throttled.PerMin(Config.NotFoundLimitPerMin), Config.NotFoundLimitBurst}
	notfoundLimiter, err = throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		logPanic("throttledSetup - notfound NewGCRARateLimiter", err, nil)
	}
}

// database setup

var gldb *leveldb.DB

func gldbSetup() {
	var err error

	if gldb != nil {
		return
	}

	if gldb, err = leveldb.OpenFile(Config.DbDir, nil); err != nil {
		logPanic("gldbSetup - leveldb.OpenFile failed", err, nil)
	}
}

// session stores

var qsLogin *qsess.Store    // login sessions
var qsVerify *qsess.Store   // sign-up email verifications
var qsEmail *qsess.Store    // email change verifications
var qsRecovery *qsess.Store // password recovery verifications

// fill in a qsess.Store struct with data we've obtained from a TOML file
func doQsess(store *qsess.Store, config ConfigQsess) {
	var err error

	store.AuthType, err = qsess.ParseAuthType(config.AuthType)
	if err != nil {
		logPanic("doQsess - unrecognized value for qsess.Store.AuthType in config file", err, nil)
	}

	store.MaxAgeSecs = config.MaxAgeSecs
	store.MinRefreshSecs = config.MinRefreshSecs

	if config.AuthType == "cookie" {
		store.CookieName = config.CookieName
		store.CookieSecure = config.CookieSecure
		store.CookieHTTPOnly = config.CookieHTTPOnly
		store.CookieSameSite = strToSameSite(config.CookieSameSite)
	}
}

// setup session stores for login, sign-up email verification, email change, and password recovery
func sessSetup() {
	var err error

	// store for login sessions

	qsLogin, err = qsldb.NewGldbStore(gldb, PfxSessLogin, logErrWriter, stringsToByteSlices(Config.QsessKeys)...)
	if err != nil {
		logPanic("sessSetup - failed to create login store", err, nil)
	}

	if Config.QsLogin.AuthType != "cookie" {
		logPanic("sessSetup - qsLogin AuthType must be cookie", nil, nil)
	}

	doQsess(qsLogin, Config.QsLogin)
	qsLogin.NewSessData = nil
	qsLogin.SessionSaved = dbRecordVisit

	// store for sign-up email verifications

	qsVerify, err = qsldb.NewGldbStore(gldb, PfxSessVerify, logErrWriter, stringsToByteSlices(Config.QsessKeys)...)
	if err != nil {
		logPanic("sessSetup - failed to create verify store", err, nil)
	}

	if Config.QsVerify.AuthType != "token" {
		logPanic("sessSetup - qsVerify AuthType must be token", nil, nil)
	}

	doQsess(qsVerify, Config.QsVerify)
	qsVerify.NewSessData = NewVerifySessData

	// store for email change verifications

	qsEmail, err = qsldb.NewGldbStore(gldb, PfxSessEmail, logErrWriter, stringsToByteSlices(Config.QsessKeys)...)
	if err != nil {
		logPanic("sessSetup - failed to create email store", err, nil)
	}

	if Config.QsEmail.AuthType != "token" {
		logPanic("sessSetup - qsEmail AuthType must be token", nil, nil)
	}

	doQsess(qsEmail, Config.QsEmail)
	qsEmail.NewSessData = NewEmailSessData

	// store for password recovery verifications

	qsRecovery, err = qsldb.NewGldbStore(gldb, PfxSessRecovery, logErrWriter, stringsToByteSlices(Config.QsessKeys)...)
	if err != nil {
		logPanic("sessSetup - failed to create recovery store", err, nil)
	}

	if Config.QsRecovery.AuthType != "token" {
		logPanic("sessSetup - qsRecovery AuthType must be token", nil, nil)
	}

	doQsess(qsRecovery, Config.QsRecovery)
	qsRecovery.NewSessData = nil
}
