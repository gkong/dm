// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// Delight/Meditate is a daily Bible reading web app.
//
// an example of the kind of scheme this program is designed for is the
// carson/m'cheyne system, which mandates 4 readings daily, two from
// the OT and two from the NT + psalms. it gets through the OT once
// and NT + psalms twice each year.
//
// concepts:
// a "stream" is a sequence of daily readings.
// a "plan" is a set of one or more streams, to be done concurrently,
// like the carson/m'cheyne system.
// an "activity" is a user's execution of a plan.
//
// persistence:
// streams and plans are immutable objects stored in the doc store.
// users and activities, which are mutable, reside in a goleveldb database.
// qsess sessions are also stored in the goleveldb database.
//
// URL schema:
// metrics are bucketed based in first segment of URL path.
// in a few cases, noted here, we tweak URL path, for the benefit of metrics.
//     / or /index.html
//         serve index.html
//         prepend /index to path for logging & metrics
//     /c/...
//         routes handled by the client (which is a single-page-app).
//         if the user saves a deep link or does a page refresh, the request
//         will come here, and we simply serve index.html in response.
//         index.html (using page.js) then displays the appropriate "page".
//         prepend /index to path for logging & metrics
//     /static/...
//         static files, served from ./static
//     everything else
//         treat as a REST API request.
//         if we recognize it, serve it, otherwise return 404.
//
// scalability and high-availability options:
//   (goleveldb can NOT scale out!)
//   - observe logs and metrics and optimize
//   - scale up
//   - move static files and doc store to a CDN
//   - convert to cassandra/scylla (just rewrite dbops.go, which is 300 lines)

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/gkong/dm/gen"

	"github.com/gkong/go-qweb/qctx"

	"github.com/julienschmidt/httprouter"
	"github.com/throttled/throttled/v2"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
)

// middleware which implements HTTP basic auth, with a static username/password,
// only used by a metrics data scraper
func mwBasicAuth(username, password string) qctx.MwMaker {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return func(next qctx.CtxHandler) qctx.CtxHandler {
		return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
			authHdr := c.R.Header.Get("Authorization")
			if len(authHdr) < 7 || strings.ToLower(authHdr[0:6]) != "basic " || authHdr[6:] != encoded {
				http.Error(c.W, "not authorized", http.StatusUnauthorized)
				return
			}
			next.CtxServeHTTP(c)
		})
	}
}

// middleware which requires that you be logged in as a given user,
// which we use to authorize one "admin" user.
func mwRequireUser(user []byte) qctx.MwMaker {
	return func(next qctx.CtxHandler) qctx.CtxHandler {
		return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
			if bytes.Compare(EmailFromUserKey(c.Sess.UserID()), user) != 0 {
				http.Error(c.W, "not authorized", http.StatusUnauthorized)
				return
			}
			next.CtxServeHTTP(c)
		})
	}
}

// middleware which applies a given throttled GCRA rate limiter, keyed by
// Request.RemoteAddr, which is IP:port.
//
// NOTE: if run behind a reverse proxy, need to check X-Forwarded-For,
// but not if NOT behind a reverse proxy, because can be spoofed!
func mwThrottledIPPort(rl *throttled.GCRARateLimiterCtx) qctx.MwMaker {
	return func(next qctx.CtxHandler) qctx.CtxHandler {
		return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
			ip := c.R.RemoteAddr
			if Config.ReverseProxy && c.R.Header.Get("X-Forwarded-For") != "" {
				ip = c.R.Header.Get("X-Forwarded-For")
			}
			limited, _, err := rl.RateLimit(ip, 1)
			if err != nil {
				internalErr("mwThrottledIPPort", "RateLimit", err, c)
				return
			}
			if limited {
				c.Error("too many requests", http.StatusTooManyRequests)
				c.R.URL.Path = "/throttled" + c.R.URL.Path
				return
			}
			next.CtxServeHTTP(c)
		})
	}
}

/*
// middleware which applies a given throttled GCRA rate limiter, keyed by IP.
func mwThrottledIP(rl *throttled.GCRARateLimiterCtx) qctx.MwMaker {
	return func(next qctx.CtxHandler) qctx.CtxHandler {
		return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
			var s string
			ndx := strings.IndexByte(c.R.RemoteAddr, ':')
			if ndx == -1 {
				s = c.R.RemoteAddr
			} else {
				s = c.R.RemoteAddr[:ndx]
			}
			limited, _, err := rl.RateLimit(s, 1)
			if err != nil {
				internalErr("mwThrottledIP", "RateLimit", err, c)
				return
			}
			if limited {
				c.Error("too many requests", http.StatusTooManyRequests)
				c.R.URL.Path = "/throttled" + c.R.URL.Path
				return
			}
			next.CtxServeHTTP(c)
		})
	}
}
*/

// middleware which applies a given throttled GCRA rate limiter, keyed by user ID.
// assumes it is downstream from middleware that requires an active session.
func mwThrottledUserID(rl *throttled.GCRARateLimiterCtx) qctx.MwMaker {
	return func(next qctx.CtxHandler) qctx.CtxHandler {
		return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
			limited, _, err := rl.RateLimit(string(EmailFromUserKey(c.Sess.UserID())), 1)
			if err != nil {
				internalErr("mwThrottledUserID", "RateLimit", err, c)
				return
			}
			if limited {
				c.Error("too many requests, please try again later", http.StatusTooManyRequests)
				return
			}
			next.CtxServeHTTP(c)
		})
	}
}

// mwRoot is the outermost middleware. it is applied to ALL incoming requests.
//   - log panics and do the best we can to inform the client
//   - log HTTP requests
//   - collect metrics
//   - check client version and reply with an error, if necessary,
//     which causes the client to reload itself.
//
// logging and metrics are done after calling and returning from downstream
// handlers. metrics are bucketed based on the first segment of the URL path.
func mwRoot(log bool, printStack bool, printAll bool) qctx.MwMaker {
	return func(next qctx.CtxHandler) qctx.CtxHandler {
		return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
			var start time.Time

			if log {
				start = time.Now()
			}

			// detect obsolete clients and tell them to reload themselves.
			// returning a 418 status code forces an immediate, disruptive update.
			// returning a "client update request" header causes the client to
			// update itself the next time it visits the home page or login page.
			version := c.R.Header.Get(Config.ClientVersionReqHdr)
			if version != "" {
				if version != Config.ClientVersion {
					// logInfo("mwRoot - client version <"+version+"> server expects "+Config.ClientVersion, c) // XXX - DEBUG
					c.W.Header().Set(Config.ClientUpdateReqHdr, "true")
				}
			}

			if Config.SSL {
				c.W.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			c.W.Header().Set("X-XSS-Protection", "1; mode=block")
			c.W.Header().Set("X-Frame-Options", "DENY")
			c.W.Header().Set("X-Content-Type-Options", "nosniff")

			// wrap W with WrappedRW (if not already wrapped)
			wrapped, ok := c.W.(*qctx.WrappedRW)
			if !ok {
				wrapped = &qctx.WrappedRW{c.W, false, false, 0}
				c.W = wrapped
			}

			defer func() {
				if x := recover(); x != nil {
					e := errors.New(fmt.Sprintf("%v", x))
					zlogger.Error("mwRoot - PANIC", getFields(e, c)...)
					emailAlertThrottled(alertLevelPanic)
					if printStack {
						bsize := 10000
						if printAll {
							bsize = 1 << 20
						}
						b := make([]byte, bsize)
						n := runtime.Stack(b, printAll)
						zlogger.Error(string(b[0:n]))
					}

					if !wrapped.WroteHeader && !wrapped.WroteBody {
						c.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
						c.Error("!PANIC", http.StatusInternalServerError)
						return
					}
					// if can't send a response code and error message to the
					// client, panic again, which results in the client seeing
					// a closed connection. Panic with ErrAbortHandler, to
					// suppress net/http error logging, since we log here.
					panic(http.ErrAbortHandler)
				}
			}()

			next.CtxServeHTTP(c)

			if log {
				// log this HTTP request
				code := http.StatusOK
				if wrapped.WroteHeader {
					code = wrapped.StatusCode
				}
				logRoot(c, code, time.Now().Sub(start).Nanoseconds()/1000)
			}

			metricsHTTPRequest(c.R.URL.Path, int(time.Now().Sub(start).Nanoseconds()/1000))
		})
	}
}

// logRoot is called by mwRoot, when logging of every HTTP request is enabled
func logRoot(c *qctx.Ctx, code int, lat int64) {
	if c.R.URL.Path == "/latency" || c.R.URL.Path == "/metrics" || c.R.URL.Path == "/malbot" {
		return
	}

	ip := c.R.RemoteAddr
	if Config.ReverseProxy && c.R.Header.Get("X-Forwarded-For") != "" {
		ip = c.R.Header.Get("X-Forwarded-For")
	}

	var user string
	haveUser := false
	if c.Sess != nil {
		userID := c.Sess.UserID()
		if len(userID) > 0 {
			haveUser = true
			user = string(EmailFromUserKey(userID))
		}
	}

	if haveUser {
		zlogger.Info("HTTP-req", zIP(ip), zMethod(c.R.Method),
			zPath(c.R.URL.Path), zStatus(code), zLatencyUs(lat), zUser(user))
	} else {
		ref := c.R.Header.Get("Referer")
		// log Referer if non-logged-in user visiting "/",
		// which indexHandler re-writes to "/index/".
		if c.R.URL.Path == "/index/" && ref != "" {
			zlogger.Info("HTTP-req", zIP(ip), zMethod(c.R.Method),
				zPath(c.R.URL.Path), zStatus(code), zLatencyUs(lat), zReferer(ref))
		} else {
			zlogger.Info("HTTP-req", zIP(ip), zMethod(c.R.Method),
				zPath(c.R.URL.Path), zStatus(code), zLatencyUs(lat))
		}
	}
}

var srv, srvtls *http.Server

var hr *httprouter.Router

// helpers to make route/handler declarations and metrics registration DRY.
// the latency flags indicate whether or not the respective metrics should be tracked.

func get(path string, mw qctx.MakerStack, handler qctx.CtxHandlerFunc, handlerLatency, clientLatency bool) {
	hr.GET(path, mw.HRHandle(handler))
	metricsRegisterEndpoint(pathFirstSegment(path), handlerLatency, clientLatency)
}

func post(path string, mw qctx.MakerStack, handler qctx.CtxHandlerFunc, handlerLatency, clientLatency bool) {
	hr.POST(path, mw.HRHandle(handler))
	metricsRegisterEndpoint(pathFirstSegment(path), handlerLatency, clientLatency)
}

func main() {
	flag.Parse()
	doConfig(flag.Args())

	logSetup()
	log.SetOutput(logFilterWriter) // capture anything logged by packages to the std logger
	metricsSetup()
	gldbSetup()

	/*
		////////////////////////////////////////////////////////////////////////
		dbMigrate()  // XXX
		time.Sleep(2 * time.Second) // XXX
		////////////////////////////////////////////////////////////////////////
	*/

	sessSetup()
	throttledSetup()
	csrfSetup()
	if err := docReadDir(Config.DocDir, Config.DocTTL); err != nil {
		logPanic("main - docReadDir", err, nil)
	}
	docQueriesSetup(Config.DocQueriesTTL)
	dailyBackupSetup(Config.BckDir, Config.BckLatestDir, Config.DailyBackupTime, logInfoWriter, logErrWriter)

	// middleware stacks

	root := qctx.MwStack(mwRoot(Config.LogAllHTTP, Config.DebugPanicStack, false))

	notfound := root.Append(mwThrottledIPPort(notfoundLimiter))
	metrics := root.Append(mwBasicAuth(Config.MetricsUser, Config.MetricsPassword)) // external metrics scraper
	static := root.Append(qctx.MwHeader("Cache-Control", "public, max-age=31536000"))
	index := root.Append(mwSetCSRF, qctx.MwHeader(
		"Cache-Control", "max-age=86400",
		// would prefer just "default-src 'self'" - had to add "img-src data:" for stupid bootstrap SVGs
		"Content-Security-Policy", "default-src 'self'; img-src 'self' data:",
	))
	adminindex := index.Append(qctx.MwRequireSess(qsLogin), mwRequireUser([]byte(strings.ToLower(Config.AdminUser))))

	csrfCheck := root.Append(mwCheckCSRF)

	// all stacks beyond this point check for CSRFs

	plain := csrfCheck.Append(qctx.MwHeader("Content-Type", "text/plain; charset=utf-8"))

	verif := plain.Append(mwThrottledIPPort(verifLimiter))
	login := plain.Append(mwThrottledIPPort(loginLimiter))
	sess := plain.Append(qctx.MwRequireSess(qsLogin))

	contact := sess.Append(mwThrottledUserID(contactLimiter))

	admin := sess.Append(mwRequireUser([]byte(strings.ToLower(Config.AdminUser))))

	hr = httprouter.New()

	metricsRegisterEndpoint("throttled", false, false)
	metricsRegisterEndpoint("malbot", false, false) // known malicious bots

	// routes

	// NOTE: javascript redundantly maintains a table of endpoints for which the
	// client should report latency. must be kept in sync with this code.

	hr.NotFound = notfound.Handle(notfoundHandler)
	metricsRegisterEndpoint("notfound", false, false)

	// routes starting with /c/ are handled by the client. we only see them
	// when a user does a page refresh or re-visits a saved link.
	// in all cases, we serve index.html, and let its JS decide what to do.
	hr.GET("/c/*anything", index.HRHandle(indexHandler))
	hr.GET("/", index.HRHandle(indexHandler))
	metricsRegisterEndpoint("index", true, false)

	get("/static/*filepath", static, fileServer(Config.StaticDir), true, false)

	// routes visited by users in response to verif emails.
	// Handlers redirect to paths which cause our client to be loaded.
	get("/verify/:token", root, verifyHandler, false, false)
	get("/recover/:token", root, recoverHandler, false, false)
	get("/email/:token", root, echangeHandler, false, false)

	// REST endpoints visited by users who do NOT have a login session
	post("/dorecover", plain, dorecoverHandler, false, false)
	post("/signup", verif, signupHandler, false, false)     // rate limited
	post("/pwforgot", verif, pwforgotHandler, false, false) // rate limited
	post("/login", login, loginHandler, true, true)         // rate limited

	get("/metrics", metrics, metricsHandler, false, false)

	// REST endpoint visted by both logged-in and not-logged-in users.
	// this is a GET, because we want browsers and proxies to cache.
	get("/getdoc/:dir/:name", plain, getdocHandler, true, true)

	// REST endpoints only accessible to logged-in users
	post("/getstate", sess, getstateHandler, true, true)
	post("/actadd", sess, activityaddHandler, true, true)
	post("/actdel", sess, activitydelHandler, true, true)
	post("/actjump", sess, activityjumpHandler, true, true)
	post("/actver", sess, activityversionHandler, true, true)
	post("/daychange", sess, daychangeHandler, true, true)
	post("/accreset", sess, accresetHandler, true, true)
	post("/accenab", sess, accenabHandler, true, true)
	post("/userget", sess, usergetHandler, true, true)
	post("/userset", sess, usersetHandler, true, true)
	post("/contact", contact, contactHandler, false, false) // rate limited, slow - don't track latency

	post("/logout", sess, logoutHandler, true, false) // no client latency - logged out!
	post("/latency", sess, latencyHandler, true, false)

	// routes only accessible to the admin user.

	get("/admin", adminindex, func(c *qctx.Ctx) { http.ServeFile(c.W, c.R, "./dm-admin.html") }, false, false)
	post("/admin", admin, adminHandler, false, false)

	// now get down to the business of listening and serving

	logInfo("main - starting up", nil)

	if Config.SSL {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(Config.SSLHosts...),
			Cache:      autocert.DirCache(Config.SSLCertDir),
		}

		srvtls = &http.Server{
			Addr:      Config.SSLTCPAddr,
			Handler:   hr,
			TLSConfig: certManager.TLSConfig(),
		}

		// start the SSL server
		go func() {
			err := srvtls.ListenAndServeTLS("", "")
			if err != nil && err != http.ErrServerClosed {
				logErr("main - ListenAndServeTLS goroutine", err, nil)
				byebye(2)
			}
		}()

		// when SSL is ENabled, the non-SSL server just serves redirects
		srv = &http.Server{
			Addr:    Config.TCPAddr,
			Handler: root.Handle(nonSSLRedirectHandler),
		}
		metricsRegisterEndpoint("nonssl", false, false)
	} else {
		// when SSL is DISabled, the non-SSL server serves everything
		srv = &http.Server{
			Addr:    Config.TCPAddr,
			Handler: hr,
		}
	}

	// start the non-SSL server
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logErr("main - ListenAndServe goroutine", err, nil)
			byebye(1)
		}
	}()

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)
	<-done
	logInfo("main - shutdown signal received", nil)
	byebye(0)
}

var byemu sync.Mutex

func byebye(code int) {
	byemu.Lock()
	defer byemu.Unlock() // not necessary, but i can't bring myself to delete it. :-)

	logInfo("byebye - waiting for outstanding requests to complete", nil)
	if Config.SSL {
		srvtls.Shutdown(context.Background())
	}
	srv.Shutdown(context.Background())
	logInfo("byebye - exiting", nil)
	os.Exit(code)
}

// modify the URL path, for logging and metrics.
// if request is from a well-known malicious bot, change URL to /malbot,
// for logging and metrics.
func notfoundHandler(c *qctx.Ctx) {
	c.R.URL.Path = "/notfound" + c.R.URL.Path
	if pathLastSegment(c.R.URL.Path) == "xmlrpc.php" {
		c.R.URL.Path = "/malbot"
	}
	c.Error("not found", http.StatusNotFound)
}

// serve index.html.
// prepend "/index" to the URL path, for logging and metrics.
func indexHandler(c *qctx.Ctx) {
	if sess, _, err := qsLogin.GetSession(c.W, c.R); err == nil {
		c.Sess = sess // just for logging
	}
	http.ServeFile(c.W, c.R, "./index.html")
	c.R.URL.Path = "/index" + c.R.URL.Path
}

// called for all non-SSL requests.
// if path is one we recognize, send a redirect to the SSL version,
// otherwise, call notfoundHandler, which sends a 404.
// modify the URL path, for logging and metrics.
func nonSSLRedirectHandler(c *qctx.Ctx) {
	if h, _, _ := hr.Lookup(c.R.Method, c.R.URL.Path); h == nil {
		c.R.URL.Path = "/nonssl" + c.R.URL.Path
		notfoundHandler(c)
		return
	}

	c.W.Header().Set("Connection", "close")
	url := Config.BaseURL + c.R.URL.String()
	http.Redirect(c.W, c.R, url, http.StatusMovedPermanently)
	c.R.URL.Path = "/nonssl" + c.R.URL.Path
}

// serve static files from a given directory.
func fileServer(dir string) qctx.CtxHandlerFunc {
	return qctx.CtxHandlerFunc(func(c *qctx.Ctx) {
		if sess, _, err := qsLogin.GetSession(c.W, c.R); err == nil {
			c.Sess = sess // just for logging
		}
		// save path and restore it when done, so log and metrics code see the "/static" prefix
		savedPath := c.R.URL.Path
		c.R.URL.Path = c.Params.ByName("filepath")
		http.FileServer(http.Dir(dir)).ServeHTTP(c.W, c.R)
		c.R.URL.Path = savedPath
	})
}

// user submitted a signup form
func signupHandler(c *qctx.Ctx) {
	if oldsess, _, err := qsLogin.GetSession(c.W, c.R); err == nil {
		logWarn("signupHandler - invoked by logged-in user", nil, c)
		if err := oldsess.Delete(c.W); err != nil {
			logErr("signupHandler - delete old", err, c)
		}
	}

	var req SignupRequest
	if !c.ReadJSON(&req) {
		logWarn("signupHandler - bad request", nil, c)
		return
	}

	first := sanitizeName(req.FirstName)
	last := sanitizeName(req.LastName)

	if err := firstError(
		validEmailAddr(req.Email),
		validPassword(req.Password),
		validFirstName(first),
		validLastName(last),
	); err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	ukey := UserKey([]byte(strings.ToLower(req.Email)))
	exists, err := dbHas(ukey)
	if err != nil {
		internalErr("signupHandler", "existence check error", err, c)
		return
	}
	if exists {
		c.Error("this email address is already registered", http.StatusBadRequest)
		return
	}

	// make a session for verification and email its token to the user

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), Config.BcryptCost)
	if err != nil {
		internalErr("signupHandler", "encryption error", err, c)
		return
	}

	sess := qsVerify.NewSession(nil)
	sd := sess.Data.(*VerifySessData)
	sd.UserKey = ukey
	sd.FirstName = first
	sd.LastName = last
	sd.Password = hashedPass

	if err := sess.Save(c.W); err != nil {
		internalErr("signupHandler", "error saving verify session", err, c)
		return
	}

	token, _, err := sess.Token()
	if err != nil {
		internalErr("signupHandler", "token error", err, c)
		return
	}

	if Config.DebugEmail {
		logInfo("signupHandler - DEBUG - token "+token, c)
	} else {
		if err := sendActivationEmail(req.Email, Config.BaseURL+"/verify/"+token); err != nil {
			internalErr("signupHandler", "email error", err, c)
			return
		}
	}
}

// user clicked the link in a signup verification email
func verifyHandler(c *qctx.Ctx) {
	token := c.Params.ByName("token")

	sess, _, err := qsVerify.GetTokenSession(token)
	if err != nil {
		logWarn("verifyHandler - timeout or bad token", err, c)
		c.Error("timeout or bad token, please try again", http.StatusBadRequest)
		return
	}
	sd := sess.Data.(*VerifySessData)

	exists, err := dbHas(sd.UserKey)
	if err != nil {
		internalErr("verifyHandler", "existence check error", err, c)
		return
	}
	if exists {
		c.Error("this email address is already registered", http.StatusBadRequest)
		return
	}

	u := User{
		Password:      sd.Password,
		FirstName:     sd.FirstName,
		LastName:      sd.LastName,
		BibleVersion:  Config.DefaultBibleVersion,
		BibleProvider: Config.DefaultBibleProvider,
		CreatedTime:   time.Now().Unix(),
		OrigEmail:     string(EmailFromUserKey(sd.UserKey)),
	}
	if err := dbUserPut(sd.UserKey, &u); err != nil {
		internalErr("verifyHandler", "cannot add new user", err, c)
		return
	}

	// it's OK if we return due to error and never explicitly delete the
	// session. it will be deleted automatically, when it expires.
	sess.Delete(c.W)

	// redirect to a path which will load our client, which will then render
	http.Redirect(c.W, c.R, "/c/verifyOK", http.StatusSeeOther)
}

// user submitted a login form
func loginHandler(c *qctx.Ctx) {
	if oldsess, _, err := qsLogin.GetSession(c.W, c.R); err == nil {
		logWarn("loginHandler - invoked by logged-in user", nil, c)
		if err := oldsess.Delete(c.W); err != nil {
			logErr("loginHandler - delete old", err, c)
		}
	}

	var req LoginRequest
	if !c.ReadJSON(&req) {
		logWarn("loginHandler - bad request", nil, c)
		return
	}

	ukey := UserKey([]byte(strings.ToLower(req.Email)))
	u, err := dbUserGet(ukey)
	if err != nil {
		c.Error("incorrect email address or password", http.StatusBadRequest)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		c.Error("incorrect email address or password", http.StatusBadRequest)
		return
	}

	c.Sess = qsLogin.NewSession(ukey)
	if req.KeepLoggedIn {
		c.Sess.MaxAgeSecs = Config.QsLogin.LongMaxAgeSecs
		c.Sess.MinRefreshSecs = Config.QsLogin.LongMinRefreshSecs
	}
	if err := c.Sess.Save(c.W); err != nil {
		internalErr("loginHandler", "session save error", err, c)
		return
	}

	csrfSetCookie(c.W)

	// send StateResponse

	var resp = StateResponse{FirstName: u.FirstName}
	resp.Activities, err = dbUserActivities(ukey)
	if err != nil {
		internalErr("loginHandler", "problem getting activities", err, c)
		return
	}
	c.WriteJSON(&resp)
}

func getstateHandler(c *qctx.Ctx) {
	ukey := c.Sess.UserID()
	u, err := dbUserGet(ukey)
	if err != nil {
		internalErr("getstateHandler", "user not found", err, c)
		return
	}

	var resp = StateResponse{FirstName: u.FirstName}
	resp.Activities, err = dbUserActivities(ukey)
	if err != nil {
		internalErr("getstateHandler", "problem getting activities", err, c)
		return
	}
	c.WriteJSON(&resp)
}

func logoutHandler(c *qctx.Ctx) {
	if err := c.Sess.Delete(c.W); err != nil {
		internalErr("logoutHandler", "error deleting session", err, c)
		return
	}
}

// user submitted a "forgot my password" form (contains only email addr).
func pwforgotHandler(c *qctx.Ctx) {
	if oldsess, _, err := qsLogin.GetSession(c.W, c.R); err == nil {
		logWarn("pwforgotHandler - invoked by logged-in user", nil, c)
		if err := oldsess.Delete(c.W); err != nil {
			logErr("pwforgotHandler - delete old", err, c)
		}
	}

	var req PWRecoverRequest
	if !c.ReadJSON(&req) {
		logWarn("pwforgotHandler - bad request", nil, c)
		return
	}

	ukey := UserKey([]byte(strings.ToLower(req.Email)))
	exists, err := dbHas(ukey)
	if err != nil {
		internalErr("pwforgotHandler", "problem finding email address", err, c)
		return
	}
	if !exists {
		logWarn("pwforgotHandler - nonexistent email addr", nil, c)
		c.Error("no account for this email address", http.StatusBadRequest)
		return
	}

	// make a Recovery session and email its token to the user

	sess := qsRecovery.NewSession(ukey)

	if err := sess.Save(c.W); err != nil {
		internalErr("pwforgotHandler", "error saving recovery session", err, c)
		return
	}

	token, _, err := sess.Token()
	if err != nil {
		internalErr("pwforgotHandler", "token error", err, c)
		return
	}

	if Config.DebugEmail {
		logInfo("pwforgotHandler - DEBUG - token "+token, c)
	} else {
		if err := sendRecoveryEmail(req.Email, Config.BaseURL+"/recover/"+token); err != nil {
			internalErr("pwforgotHandler", "email error", err, c)
			return
		}
	}
}

// user clicked the link in a "forgot my password" verification email.
// redirect them to a form into which to enter their new password.
func recoverHandler(c *qctx.Ctx) {
	token := c.Params.ByName("token")

	sess, _, err := qsRecovery.GetTokenSession(token)
	if err != nil {
		logWarn("recoverHandler - timeout or bad token", err, c)
		c.Error("timeout or bad token, please try again", http.StatusBadRequest)
		return
	}
	ukey := sess.UserID()

	exists, err := dbHas(ukey)
	if err != nil {
		internalErr("recoverHandler", "problem finding user", err, c)
		return
	}
	if !exists {
		logWarn("recoverHandler - no account", nil, c)
		c.Error("no account for this email address", http.StatusBadRequest)
		return
	}

	// redirect to a path which will load our client, which will then render
	http.Redirect(c.W, c.R, "/c/pwreset/"+token, http.StatusSeeOther)
}

// user submitted password form (final step of "forgot my password" sequence).
func dorecoverHandler(c *qctx.Ctx) {
	var req DoRecoverRequest
	if !c.ReadJSON(&req) {
		logWarn("dorecoverHandler - bad request", nil, c)
		return
	}

	sess, _, err := qsRecovery.GetTokenSession(req.Token)
	if err != nil {
		logWarn("dorecoverHandler - timeout or bad token", err, c)
		c.Error("timeout or bad token, please try again", http.StatusBadRequest)
		return
	}
	ukey := sess.UserID()

	u, err := dbUserGet(ukey)
	if err != nil {
		internalErr("dorecoverHandler", "user not found", err, c)
		return
	}

	if err := validPassword(req.Password); err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), Config.BcryptCost)
	if err != nil {
		internalErr("dorecoverHandler", "encryption error", err, c)
		return
	}
	u.Password = hashedPass

	if err := dbUserPut(ukey, u); err != nil {
		internalErr("dorecoverHandler", "cannot update user profile", err, c)
		return
	}

	// create (but don't Save) a login session, just to invoke DeleteByUserID
	lsess := qsLogin.NewSession(ukey)
	err = lsess.DeleteByUserID(c.W)
	if err != nil {
		logErr("dorecoverHandler - DeleteByUserID", err, c)
	}

	sess.Delete(c.W)
}

// user submitted a "contact us" message.
func contactHandler(c *qctx.Ctx) {
	var req ContactRequest
	if !c.ReadJSON(&req) {
		logWarn("contactHandler - bad request", nil, c)
		return
	}

	if err := validEmailMsg(req.Msg); err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	u, err := dbUserGet(c.Sess.UserID())
	if err != nil {
		internalErr("contactHandler", "user not found", err, c)
		return
	}

	if Config.DebugEmail {
		logInfo("contactHandler - DEBUG - << "+strings.Replace(req.Msg, "\n", " ", -1)+" >>", c)
	} else {
		if err := sendContactEmail(string(EmailFromUserKey(c.Sess.UserID())), u.FirstName+" "+u.LastName, req.Msg); err != nil {
			internalErr("contactHandler", "email error", err, c)
			return
		}
	}
}

func usergetHandler(c *qctx.Ctx) {
	ukey := c.Sess.UserID()
	u, err := dbUserGet(ukey)
	if err != nil {
		internalErr("usergetHandler", "user not found", err, c)
		return
	}

	c.WriteJSON(&UserReqResp{
		Email:         string(EmailFromUserKey(ukey)),
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		BibleVersion:  u.BibleVersion,
		BibleProvider: u.BibleProvider,
	})
}

// user submitted a profile update.
// check password (which uses a deliberately slow hash).
// if changing password, must hash the new one, which is equally slow.
// if changing email, send an email verification, before completing the change.
func usersetHandler(c *qctx.Ctx) {
	var newukey []byte
	var emailChange, pwChange bool
	var pwChangeErr error

	var req UserReqResp
	if !c.ReadJSON(&req) {
		logWarn("usersetHandler - bad request", nil, c)
		return
	}

	first := sanitizeName(req.FirstName)
	last := sanitizeName(req.LastName)

	if req.NewPassword != "" {
		pwChange = true
		pwChangeErr = validPassword(req.NewPassword)
	}

	if err := firstError(
		pwChangeErr,
		validFirstName(first),
		validLastName(last),
		validBibleVersion(req.BibleVersion),
	); err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := providers[req.BibleProvider]; !ok {
		logWarn("usersetHandler - provider does not exist", nil, c)
		c.Error("provider does not exist", http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	u, err := dbUserGet(ukey)
	if err != nil {
		internalErr("usersetHandler", "user not found", err, c)
		return
	}

	newEmail := strings.ToLower(req.Email)
	newEmailBytes := []byte(newEmail)
	if bytes.Compare(EmailFromUserKey(ukey), newEmailBytes) != 0 {
		emailChange = true
		if err := validEmailAddr(newEmail); err != nil {
			c.Error(err.Error(), http.StatusBadRequest)
			return
		}
		newukey = UserKey(newEmailBytes)
		exists, err := dbHas(newukey)
		if err != nil {
			internalErr("usersetHandler", "error looking up new email address", err, c)
			return
		}
		if exists {
			c.Error("there is already a user with the new email address", http.StatusBadRequest)
			return
		}
	}

	u.FirstName = first
	u.LastName = last
	u.BibleVersion = req.BibleVersion
	u.BibleProvider = req.BibleProvider

	// do expensive operations (password hashing) as late as possible,
	// so we don't have to do them if abort due to earlier errors.

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)) != nil {
		c.Error("incorrect password", http.StatusBadRequest)
		return
	}

	if pwChange {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), Config.BcryptCost)
		if err != nil {
			internalErr("usersetHandler", "encryption error", err, c)
			return
		}
		u.Password = hashedPass
	}

	// update the user record
	if err := dbUserPut(ukey, u); err != nil {
		internalErr("usersetHandler", "cannot update user profile", err, c)
		return
	}
	if pwChange {
		// do this as soon as we know user record has been updated successfully
		defer c.Sess.DeleteByUserID(c.W)
	}

	// if email is being updated, make an email session and email token to client.
	// email is not actually changed until client presents the token.
	if emailChange {
		esess := qsEmail.NewSession(nil)
		esd := esess.Data.(*EmailSessData)
		esd.OldUserKey = ukey
		esd.NewUserKey = newukey

		if err := esess.Save(c.W); err != nil {
			internalErr("usersetHandler", "error saving email session", err, c)
			return
		}

		token, _, err := esess.Token()
		if err != nil {
			internalErr("usersetHandler", "token error", err, c)
			return
		}

		if Config.DebugEmail {
			logInfo("usersetHandler - DEBUG - token "+token, c)
		} else {
			if err := sendEChangeEmail(req.Email, Config.BaseURL+"/email/"+token); err != nil {
				internalErr("usersetHandler", "email error", err, c)
				return
			}
		}
	}

	c.WriteJSON(&UserSetResponse{
		Logout:          pwChange,
		EmailSent:       emailChange,
		PasswordChanged: pwChange,
		FirstName:       first,
	})
}

// user clicked the link in an "email change verification" email.
// since email addr is user id:
// (1) must re-key all activity records for the user.
// (2) must delete all sessions for the user, since qsess maintains an index
// by userid, which would be invalidated. this is also good security practice.
func echangeHandler(c *qctx.Ctx) {
	token := c.Params.ByName("token")

	sess, _, err := qsEmail.GetTokenSession(token)
	if err != nil {
		logWarn("echangeHandler - timeout or bad token", err, c)
		c.Error("timeout or bad token, please try again", http.StatusBadRequest)
		return
	}
	sd := sess.Data.(*EmailSessData)

	u, err := dbUserGet(sd.OldUserKey)
	if err != nil {
		internalErr("echangeHandler", "user not found", err, c)
		return
	}

	exists, err := dbHas(sd.NewUserKey)
	if err != nil {
		internalErr("echangeHandler", "error looking up new email address", err, c)
		return
	}
	if exists {
		c.Error("there is already a user with the new email address", http.StatusBadRequest)
		return
	}

	// switch User record to new email addr (i.e. new key).
	// if we encounter errors after this point, we will log them
	// but keep going, trying to complete as much as possible.

	if err := dbUserPut(sd.NewUserKey, u); err != nil {
		internalErr("echangeHandler", "cannot add new user", err, c)
		return
	}
	if err := dbDelete(sd.OldUserKey); err != nil {
		logErr("echangeHandler - delete old user "+string(EmailFromUserKey(sd.OldUserKey)), err, c)
	}

	// revoke all active sessions

	// create (but don't Save) a login session, just to invoke DeleteByUserID
	lsess := qsLogin.NewSession(sd.OldUserKey)
	err = lsess.DeleteByUserID(c.W)
	if err != nil {
		logErr("echangeHandler - DeleteByUserID", err, c)
	}

	// activities are keyed by user key (= email addr); convert them to new key

	activities, err := dbUserActivities(sd.OldUserKey)
	if err != nil {
		logErr("echangeHandler - DeleteByUserID", err, c)
	} else {
		for i := range activities {
			oldakey := ActivityKey(sd.OldUserKey, []byte(activities[i].PlanName))
			if err := dbDelete(oldakey); err != nil {
				logErr("echangeHandler - delete old activity", err, c)
			}
			newakey := ActivityKey(sd.NewUserKey, []byte(activities[i].PlanName))
			if err := dbActivityPut(newakey, activities[i]); err != nil {
				logErr("echangeHandler - create new activity", err, c)
			}
		}
	}

	sess.Delete(c.W)

	// redirect to a path which will load our client, which will then render
	http.Redirect(c.W, c.R, "/c/ecdone", http.StatusSeeOther)
}

func activityaddHandler(c *qctx.Ctx) {
	var req ActivityAddRequest
	if !c.ReadJSON(&req) {
		logWarn("activityaddHandler - bad request", nil, c)
		return
	}

	if _, ok := planDays[req.Plan]; !ok {
		logWarn("activityaddHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	u, err := dbUserGet(ukey)
	if err != nil {
		internalErr("activityaddHandler", "user not found", err, c)
		return
	}

	a := &Activity{
		PlanName:                req.Plan,
		Day:                     make([]int, planStreams[req.Plan]),
		BibleVersion:            u.BibleVersion,
		BibleProvider:           u.BibleProvider,
		AccountabilityStartDate: req.Today,
		AccountabilityVisible:   true,
	}
	for i := range a.Day {
		a.Day[i] = 1
	}

	akey := ActivityKey(ukey, []byte(req.Plan))
	if err := dbActivityPut(akey, a); err != nil {
		internalErr("activityaddHandler", "failed to create activity", err, c)
		return
	}

	// send StateResponse
	getstateHandler(c)
}

func activitydelHandler(c *qctx.Ctx) {
	var req ActivityDeleteRequest
	if !c.ReadJSON(&req) {
		logWarn("activitydelHandler - bad request", nil, c)
		return
	}

	if _, ok := planDays[req.Plan]; !ok {
		logWarn("activitydelHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	akey := ActivityKey(ukey, []byte(req.Plan))

	has, err := dbHas(akey)
	if err != nil {
		internalErr("activitydelHandler", "problem locating plan", err, c)
		return
	}
	if !has {
		logWarn("activitydelHandler - activity does not exist", nil, c)
		c.Error("plan not installed", http.StatusBadRequest)
		return
	}
	if err := dbDelete(akey); err != nil {
		internalErr("activitydelHandler", "plan not deleted", err, c)
		return
	}
}

func activityjumpHandler(c *qctx.Ctx) {
	var req ActivityJumpRequest
	if !c.ReadJSON(&req) {
		logWarn("activityjumpHandler - bad request", nil, c)
		return
	}

	totaldays, ok := planDays[req.Plan]
	if !ok {
		logWarn("activityjumpHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	if req.Day < 1 || req.Day > totaldays {
		c.Error("day must be between 1 and "+strconv.Itoa(totaldays), http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	akey := ActivityKey(ukey, []byte(req.Plan))
	a, err := dbActivityGet(akey)
	if err != nil {
		internalErr("activityjumpHandler", "problem retrieving plan", err, c)
		return
	}

	for i := range a.Day {
		a.Day[i] = req.Day
	}
	a.AccountabilityStartDate = req.Today - req.Day + 1

	if err := dbActivityPut(akey, a); err != nil {
		internalErr("activityjumpHandler", "problem saving new status", err, c)
		return
	}

	c.WriteJSON(&DayResponse{
		Day:                     a.Day,
		AccountabilityStartDate: a.AccountabilityStartDate,
		AccountabilityVisible:   a.AccountabilityVisible,
	})
}

func activityversionHandler(c *qctx.Ctx) {
	var req ActivityVersionRequest
	if !c.ReadJSON(&req) {
		logWarn("activityversionHandler - bad request", nil, c)
		return
	}

	if _, ok := planDays[req.Plan]; !ok {
		logWarn("activityversionHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	if _, ok := providers[req.BibleProvider]; !ok {
		logWarn("activityversionHandler - provider does not exist", nil, c)
		c.Error("provider does not exist", http.StatusBadRequest)
		return
	}

	if err := validBibleVersion(req.BibleVersion); err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	akey := ActivityKey(ukey, []byte(req.Plan))
	a, err := dbActivityGet(akey)
	if err != nil {
		logWarn("activityversionHandler - activity does not exist", nil, c)
		c.Error("plan is not installed", http.StatusBadRequest)
		return
	}

	a.BibleProvider = req.BibleProvider
	a.BibleVersion = req.BibleVersion
	if err := dbActivityPut(akey, a); err != nil {
		internalErr("activityversionHandler", "problem saving new status", err, c)
		return
	}
}

// Set accountability start date so that user has at least one and at most
// one day's worth of readings to do.
// The currently-displayed items remain the same;
// the only thing modified by this function is Activity.AccountabilityStartDate.
func accresetHandler(c *qctx.Ctx) {
	var req AccountabilityResetRequest
	if !c.ReadJSON(&req) {
		logWarn("accresetHandler - bad request", nil, c)
		return
	}

	if _, ok := planDays[req.Plan]; !ok {
		logWarn("accresetHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	akey := ActivityKey(ukey, []byte(req.Plan))
	a, err := dbActivityGet(akey)
	if err != nil {
		logWarn("accresetHandler - activity does not exist", nil, c)
		c.Error("plan is not installed", http.StatusBadRequest)
		return
	}

	readingsDone := 0
	for _, day := range a.Day {
		readingsDone += day
	}

	a.AccountabilityStartDate = req.Today - (readingsDone / len(a.Day)) + 1

	if err := dbActivityPut(akey, a); err != nil {
		internalErr("accresetHandler", "problem saving new status", err, c)
		return
	}

	c.WriteJSON(&DayResponse{
		Day:                     a.Day,
		AccountabilityStartDate: a.AccountabilityStartDate,
		AccountabilityVisible:   a.AccountabilityVisible,
	})
}

func accenabHandler(c *qctx.Ctx) {
	var req AccountabilityEnabledRequest
	if !c.ReadJSON(&req) {
		logWarn("accenabHandler - bad request", nil, c)
		return
	}

	if _, ok := planDays[req.Plan]; !ok {
		logWarn("accenabHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	akey := ActivityKey(ukey, []byte(req.Plan))
	a, err := dbActivityGet(akey)
	if err != nil {
		logWarn("accenabHandler - activity does not exist", nil, c)
		c.Error("plan is not installed", http.StatusBadRequest)
		return
	}

	a.AccountabilityVisible = req.Enabled

	if err := dbActivityPut(akey, a); err != nil {
		internalErr("accenabHandler", "problem saving new status", err, c)
		return
	}

	c.WriteJSON(&DayResponse{
		Day:                     a.Day,
		AccountabilityStartDate: a.AccountabilityStartDate,
		AccountabilityVisible:   a.AccountabilityVisible,
	})
}

func daychangeHandler(c *qctx.Ctx) {
	var req DayChangeRequest
	if !c.ReadJSON(&req) {
		logWarn("daychangeHandler - bad request", nil, c)
		return
	}

	pdays, ok := planDays[req.Plan]
	if !ok {
		logWarn("daychangeHandler - not a valid plan", nil, c)
		c.Error("not a valid plan", http.StatusBadRequest)
		return
	}

	ukey := c.Sess.UserID()
	akey := ActivityKey(ukey, []byte(req.Plan))
	a, err := dbActivityGet(akey)
	if err != nil {
		logWarn("daychangeHandler - activity does not exist", nil, c)
		c.Error("plan is not installed", http.StatusBadRequest)
		return
	}

	if req.StreamIndex < 0 || req.StreamIndex >= len(a.Day) {
		logWarn("daychangeHandler - bad streamindex", nil, c)
		c.Error("bad streamindex", http.StatusBadRequest)
		return
	}

	if req.Delta != 1 && req.Delta != -1 {
		logWarn("daychangeHandler - bad delta", nil, c)
		c.Error("bad delta", http.StatusBadRequest)
		return
	}

	// if client's idea of day value does not match that in the database,
	// then client is out-of-date (i.e. logged in on multiple devices and
	// has advanced on some other device), and his state must be updated.

	if req.PrevDay == a.Day[req.StreamIndex] {
		if req.Delta == -1 && req.PrevDay > 1 {
			a.Day[req.StreamIndex]--
		}
		if req.Delta == 1 && req.PrevDay <= pdays {
			a.Day[req.StreamIndex]++

			// if all streams completed, start over with day 1 !!!
			alldone := true
			for _, d := range a.Day {
				if d != pdays+1 {
					alldone = false
					break
				}
			}
			if alldone {
				for i := range a.Day {
					a.Day[i] = 1
					a.AccountabilityStartDate = req.Today
				}
			}
		}
		// in some cases, this is unnecessary, but not worth checking
		if err := dbActivityPut(akey, a); err != nil {
			internalErr("daychangeHandler", "problem saving status", err, c)
			return
		}
	}

	c.WriteJSON(&DayResponse{
		Day:                     a.Day,
		AccountabilityStartDate: a.AccountabilityStartDate,
		AccountabilityVisible:   a.AccountabilityVisible,
	})
}

func getdocHandler(c *qctx.Ctx) {
	// we do NOT require a session, since non-logged-in users can browse plans,
	// however, if there is a session, put it into the context, for logging.
	if sess, _, err := qsLogin.GetSession(c.W, c.R); err == nil {
		c.Sess = sess
	}

	dir := c.Params.ByName("dir")
	name := c.Params.ByName("name")

	doc, ttl, err := docGet(dir, name)
	if err != nil {
		logWarn("getdocHandler - not found - "+dir+"/"+name, err, c)
		c.Error("document not found", http.StatusNotFound)
		return
	}
	c.W.Header().Set("Cache-control", "max-age="+strconv.Itoa(ttl))
	c.WriteJSON(&DocWrapper{TTLSecs: ttl, Data: string(doc)})
}

// a client is reporting latency for one XHR.
func latencyHandler(c *qctx.Ctx) {
	var req LatencyRequest
	if !c.ReadJSON(&req) {
		logWarn("latencyHandler - bad request", nil, c)
		return
	}

	metricsClientLatency(req.Endpoint, req.TimeMsec) // does its own validation
}

// for prometheus
func metricsHandler(c *qctx.Ctx) {
	prHandler.ServeHTTP(c.W, c.R)
}

func adminHandler(c *qctx.Ctx) {
	var req AdminRequest
	if !c.ReadJSON(&req) {
		logWarn("adminHandler - bad request", nil, c)
		return
	}

	switch req.Cmd {
	case "alertreset":
		alertReset()

	case "metrics":
		prHandler.ServeHTTP(c.W, c.R)

	case "backup":
		fname := filepath.Join(Config.BckDir, time.Now().Format("dmbackup-20060102-150405"))
		logInfo("database backup - saving to "+fname, nil)
		n, err := dbBackup(fname)
		if err != nil {
			internalErr("backupHandler", "dbBackup failed", err, c)
			return
		}
		// could define structured logging fields for filename and record count, but not worth the trouble
		logInfo(fmt.Sprintf("database backup - saved %d records to %s", n, fname), nil)
		c.W.Write([]byte("created " + fname))

	case "restore":
		fname := filepath.Join(Config.BckDir, req.Arg)
		logInfo("database restore - restoring from "+fname, nil)
		n, err := dbRestore(fname)
		if err != nil {
			internalErr("restoreHandler", "dbRestore failed", err, c)
			return
		}
		// could define structured logging fields for filename and record count, but not worth the trouble
		logInfo(fmt.Sprintf("database restore - restored %d records from %s", n, fname), nil)
		c.W.Write([]byte("OK"))

	case "profcpu":
		pprof.Profile(c.W, c.R)

	case "profmem":
		pprof.Handler("heap").ServeHTTP(c.W, c.R)

	case "error":
		logErr("test error", nil, nil)

	case "panic":
		logPanic("test panic", nil, nil)

	case "stdlog":
		log.Printf("this is a test msg to the std logger")

	case "sess":
		// find session cookie in your browser, paste contents, see its TTL
		if len(req.Arg) == 0 {
			c.Error("token required", http.StatusBadRequest)
			return
		}
		s, ttl, err := qsLogin.GetTokenSession(req.Arg)
		if err != nil {
			c.Error("GetTokenSession error - "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(c.W, "ttl %d, maxage %d, minrefresh %d\n", ttl, s.MaxAgeSecs, s.MinRefreshSecs)

	case "debug":
		// dbDropTable(PfxSessLogin)
		// dbDropTable(PfxSessVerify)
		// dbDropTable(PfxSessEmail)
		// dbChangePrefix([]byte{2}, []byte{4})
		// dbDelete(UserKey([]byte("george@alreadynotyet.com")))
		// gldb.Put([]byte{ 0, 2, 0, 0 }, []byte{4, 5, 6}, nil)  // malformed session exp key
		// dbNuke()
		dbDisplay(os.Stderr)

	default:
		c.Error("bad cmd", http.StatusBadRequest)
	}
}

// helper that logs an error and sends http.StatusInternalServerError.
// usrMsg is logged and may be seen by the user - doesn't need to be
// comprehensible, but should not reveal anything sensitive.
func internalErr(context string, usrMsg string, err error, c *qctx.Ctx) {
	logErr(context+" - "+usrMsg, err, c)
	c.Error(usrMsg, http.StatusInternalServerError)
}
