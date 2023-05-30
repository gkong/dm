require=(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
/** @preserve
 *
 * @copyright   Copyright (C) 2017 George S. Kong, All Rights Reserved.
 * Use of this source code is governed by a license that can be found in the LICENSE.txt file.
 */

jQuery = $ = require('jquery');
var Rlite = require('rlite-router');
Handlebars = require('handlebars/runtime');  // put into global scope, so dmt.js can see it
require('bootstrap');
var spa = require('spa-components');


var state;                  // state information obtained from server (with some client-side additions)
var route;                  // rlite router
var clientVersion;          // client version string read from DOM (defined in index.html)
var clientUpdateRequested;  // reload client as soon as it can be done gracefully
var dialogPage;             // short alias for Handlebars.templates.dialog function

const lrkey = "dmLastReload";      // localStorage key for last reload time, in secs since 1970

$(function() {
	// jquery idiom for delaying until DOM is built and all scripts parsed

	state = {};
	state.loggedIn = false;
	state.firstName = "";
	state.activities = [];  // from server, homePage() enhances

	clientUpdateRequested = false;

	$('#dm-client-version').addClass("dm-gone");
	clientVersion = $('#dm-client-version').text();

	dialogPage = Handlebars.templates.dialog;
	Handlebars.registerPartial('signupPartial', Handlebars.templates.partialsignup);
	Handlebars.registerPartial('versionPartial', Handlebars.templates.partialversion);

	$('#dm-navbar').html(Handlebars.templates.navbar({}));
	$('#dm-menu-logout').click(function() { logout(); });

	// collapse bootstrap mobile menus after they've been clicked.
	$(".navbar-nav li a.dropdown-item").click(function() {
		$(".navbar-collapse").collapse('hide');
	});
	$(".navbar-nav li a.collapser").click(function() {
		$(".navbar-collapse").collapse('hide');
	});

	// collapse bootstrap navbar and/or dropdown when you click outside it.
	$('html').on('click, touchend', function (e) {

		// if not in the navbar, collapse the navbar.
		if ($(e.target).closest('.navbar').length == 0) {
			$('.navbar-collapse').collapse('hide');
		}

		// if not in a dropdown, collapse dropdowns.
		// XXX - probably fails if one dropdown is open and you click on a DIFFERENT dropdown's toggle!
		if ($(e.target).closest('.dropdown').length == 0) {
			$('.dropdown-menu').filter(':visible').first().siblings('.dropdown-toggle').dropdown('toggle');
		}

	});

	const route = Rlite(notfoundPage, {
		'/':                    rootPage,
		'/c/login':             loginPage,
		'/c/signup':            signupPage,
		'/c/pending':           signupPendingPage,
		'/c/verifyOK':          verifyOkPage,
		'/c/forgot':            forgotPage,
		'/c/forgotpending':     forgotpendingPage,
		'/c/pwreset/:token':    pwresetPage,

		'/c/home':              homePage,
		'/c/about':             aboutPage,
		'/c/testdrive':         testdrivePage,
		'/c/help':              helpPage,
		'/c/contact':           contactPage,
		'/c/plans':             plansPage,
		'/c/plandetails/:plan': plandetailsPage,
		'/c/planday/:plan':     plandayPage,
		'/c/version/:plan':     versionPage,
		'/c/profile':           profilePage,
		'/c/ecpending':         ecpendingPage,
		'/c/ecdone':            ecdonePage,
		'/c/badprovider':       badProviderPage
	});

	var path = window.location.pathname;

	// must call spa.init, BEFORE xhrPostJson, and tell it NOT to render,
	// because xhrPostJson's error handling can call spa.visit,
	// to redirect to a different page.

	var prevState = spa.init({
		router:   route,
		logging:  false,
	});

	xhrPostJson("/getstate", "").then(function(okXHR) {
		state.loggedIn = true;
		receiveStateResponse(JSON.parse(okXHR.responseText));
		return(okXHR)
	}, function(errorXHR) {
		return(errorXHR);
	}).then(function(XHR) {
		if (window.location.pathname === path) {
			// we haven't been redirected to a different path by xhrPostJson,
			// so render the page specified by the URL by which we were invoked.
			// history has already been set up by spaInit; just need to render.
			route(path);
			if (prevState != null)
				spa.scrollTo(prevState.scrollx, prevState.scrolly);
		}
		navSetup();
		cron();
	});
});

function mwReq(req, method, url) {
	req.setRequestHeader("Dm-Client-Version", clientVersion);

	var tok = document.cookie.replace(/(?:(?:^|.*;\s*)dm_csrf\s*\=\s*([^;]*).*$)|^.*$/, "$1");
	if (tok != "")
		req.setRequestHeader("Dm-Csrf", tok);
		var endtime = performance.now();

	req.myXhrStartTime = performance.now();  // stick this onto the xhr for latency calculation
}

function mwReqJson(req, method, url) {
	mwReq(req);
	req.setRequestHeader("Content-type", "application/json;charset=UTF-8");
}

function mwBefore(resp, method, url) {
	// graceful reload request
	if (resp.getResponseHeader('Dm-Client-Update') === "true")
		clientUpdateRequested = true;
}

function mwFailure(resp, method, url) {
	if (resp.status === 418) {
		// cancel graceful reload and perform an immediate, disruptive reload
		clientUpdateRequested = false;
		modalReload("Due to a server software upgrade, this page will now be reloaded.");
		resp.dmhandled = true;  // on some browsers, callers will finish executing before the reload happens
	} else if (resp.status === 403) {
		// can't do anything, must reload to get a new CSRF token
		modalReload("CSRF error. This page will now be reloaded.");
		resp.dmhandled = true;  // on some browsers, callers will finish executing before the reload happens
	} else if (resp.status === 401) {
		stateLoggedOut();
		if (!ok401(window.location.pathname))
			spa.visit('/c/login');
		resp.dmhandled = true;
	}
}

function mwAfter(resp, method, url) {
	// tell the server how long this request took
	var endtime = performance.now();
	if (state.loggedIn)
		postLatency(url, Math.floor(endtime-resp.myXhrStartTime));
}

var xhrGet = spa.httpReqFunc("GET", mwReq, mwBefore, undefined, mwFailure, mwAfter);
var xhrPostJson = spa.httpReqFunc("POST", mwReqJson, mwBefore, undefined, mwFailure, mwAfter);

// return true if it's ok to stay on the current page when we receive a
// 401 (unauthorized) status code, rather than redirecting to the login page.
//
// this can happen during an emailed verification sequence, when the server
// redirects us to a client route. the resulting fresh start of the client 
// should NOT redirect to the login page.
//
// it can also happen during a test drive on chrome, when the user hits the
// back button from an external provider page, resulting in a fresh start
// of the client. when that happens, we want to stay on the test drive page.

function ok401(path) {
	if (path === '/c/verifyOK' || path === '/c/ecdone' || path.indexOf('/c/pwreset/') === 0 || path === '/c/testdrive')
		return true;
	return false;
}

// endpoints for which we should report latency
// NOTE: this should be kept in sync with main.go
var latencyEndpoints = {
	"login": true,
	"getdoc": true,
	"getstate": true,
	"actadd": true,
	"actdel": true,
	"actjump": true,
	"actver": true,
	"daychange": true,
	"accreset": true,
	"accenab": true,
	"userget": true,
	"userset": true,
}

// tell the server how long a request took
function postLatency(url, timemsec) {
	var endpoint = pathFirstSegment(url);

	if (!latencyEndpoints[endpoint])
		return;

	var lParams = {
		endpoint: pathFirstSegment(url),
		timemsec: timemsec
	};

	// don't bother with Dm-Client-Version, since latency always immediately
	// follows another request, which will have sent it

	var req = new XMLHttpRequest();
	req.open("POST", "/latency", true);
	req.setRequestHeader("Content-type", "application/json;charset=UTF-8");
	var tok = document.cookie.replace(/(?:(?:^|.*;\s*)dm_csrf\s*\=\s*([^;]*).*$)|^.*$/, "$1");
	if (tok != "")
		req.setRequestHeader("Dm-Csrf", tok);
	req.send(JSON.stringify(lParams));
}

// display an error message from the server in a form's alert area
function formErrorAlert(selector, errorXHR, altMsg) {
	if (errorXHR.dmhandled)
		return;

	// only display altMsg if response contained no error message
	if (errorXHR.responseText)
		$(selector).html(errorXHR.responseText);
	else
		$(selector).html(altMsg + " - error " + errorXHR.status.toString());
	$(selector).show();
}

// handle a StateResponse from the server
function receiveStateResponse(responseObj) {
	if (responseObj.hasOwnProperty("firstname"))
		state.firstName = responseObj.firstname;
	if (responseObj.hasOwnProperty("activities"))
		state.activities = responseObj.activities;
	navSetup();
}

function stateLoggedOut() {
	state.loggedIn = false;
	navSetup();
}

function navSetup() {
	if (state.loggedIn) {
		$("#user-dropdown").html(state.firstName);
		$("#nav-user-dropdown").removeClass("dm-gone");
		$("#nav-login").addClass("dm-gone");
		$("#nav-contact").removeClass("dm-gone");
	} else {
		$("#nav-user-dropdown").addClass("dm-gone");
		$("#nav-login").removeClass("dm-gone");
		$("#nav-contact").addClass("dm-gone");
	}
}

// every midnight: if showing home page, re-draw, to update temporal contexts.
// theoretically, this should generate zero server traffic, but we add a
// random delta, just in case.
//
// XXX - could re-fetch state from server, as partial solution to the problem
// of multiple logins on different devices getting out of date.
function cron() {
	if ($("#dm-home-visible").length > 0) {
		spa.visit('/c/home');
	}
	var d = new Date();
	var msTillMidnight = ((24*60*60) - ((((d.getHours() * 60) + d.getMinutes()) * 60) + d.getSeconds())) * 1000;
	var randDeltaMs = 1 * 60000 + Math.floor(Math.random() * 19 * 60000);  // add from 1 to 20 mins
	setTimeout(cron, msTillMidnight + randDeltaMs);
}

function showPage(html) {
	$('#main').html(html);
}

// return a promise, which displays an alert and resolves when dismissed.
function modalAlert(msg) {
	return new Promise(function(resolve) {
		$('#dm-modal').html(Handlebars.templates.modalalert({ msg: msg }));
		$('.modal').modal();
		$('.modal').on('hidden.bs.modal', function() {
			$('#dm-modal').html("");
			resolve();
		});
	});
}

// return a promise, which asks a yes/no question and resolves if yes, rejects if no.
function modalConfirm(msg) {
	return new Promise(function(resolve, reject) {
		var confirmed = false;
		$('#dm-modal').html(Handlebars.templates.modalconfirm({ msg: msg }));
		$('#dm-modal-ok').click(function() {
			confirmed = true;
		});
		$('.modal').modal();
		$('.modal').on('hidden.bs.modal', function() {
			$('#dm-modal').html("");
			if (confirmed)
				resolve();
			else
				reject();
		});
	});
}

// page handlers

function rootPage() {
	if (state.loggedIn)
		spa.visit('/c/home');
	else
		spa.visit('/c/login');
}

function notfoundPage() {
	document.title = "Delight/Meditate - page not found";
	showPage('<br/><h1>PAGE NOT FOUND</h1>');
}

function signupPage() {
	document.title = "Delight/Meditate - signup";
	showPage(Handlebars.templates.signup({}));
	$('#signup-alert').hide();
	$('#dm-signup-button').click(function() { signup(); });
	$('#dm-login-button').click(function() { spa.visit('/c/login'); });
	$('#dm-signup-panel').on("keydown", function(e) {
		if (e.keyCode === 13) {
			e.preventDefault();
			signup();
		}
	});
}

function signupPendingPage() {
	document.title = "Delight/Meditate - email sent";
	showPage(dialogPage({		
		msg: "You will receive an email shortly. Please follow its instructions, to activate your account."
	}));
	$('#dm-dialog-button').addClass("dm-gone");
}

function verifyOkPage() {
	document.title = "Delight/Meditate - signup complete";
	showPage(dialogPage({		
		msg: "Your account has been activated and is ready for you to log in!",
		buttonlabel: "Go to login page"
	}));
	$('#dm-dialog-button button').click(function() { spa.visit('/c/login'); });
	$('#dm-dialog-button').removeClass("dm-gone");
}

function loginPage() {
	if (clientUpdateRequested) {
		clientUpdateRequested = false;
		quietReload("Due to a server software upgrade, this page will now be reloaded.");
	}

	document.title = "Delight/Meditate - login";
	showPage(Handlebars.templates.login({}));
	$('#login-alert').hide();
	$('#dm-about-button').click(function() { spa.visit('/c/about'); });
	$('#dm-test-drive-button').click(function() {
		spa.visit('/c/testdrive');
	});
	$('#dm-login-button').click(function() { login(); });
	$('#dm-signup-button').click(function() { spa.visit('/c/signup'); });
	$('#dm-login-panel').on("keydown", function(e) {
		if (e.keyCode === 13) {
			e.preventDefault();
			login();
		}
	});
}

function forgotPage() {
	document.title = "Delight/Meditate - forgot password";
	showPage(Handlebars.templates.forgot({}));
	$('#forgot-alert').hide();
	$('#dm-forgot-button').click(function() { forgot(); });
}

function forgotpendingPage() {
	document.title = "Delight/Meditate - password reset email sent";
	showPage(dialogPage({		
		msg: "You will receive an email shortly. Please follow its instructions, to reset your password. If you don't see the email, check your spam folder."
	}));
	$('#dm-dialog-button').addClass("dm-gone");
}

function pwresetPage(params) {
	document.title = "Delight/Meditate - password reset";
	var token = params.token;
	showPage(Handlebars.templates.pwreset(token));
	$('#pwreset-alert').hide();
	$('#dm-button-pwreset').click(function() { pwreset(); });
}

function aboutPage() {
	document.title = "Delight/Meditate - about";
	showPage(Handlebars.templates.about({}));
	if (state.loggedIn) {
		$('.dm-visitor-only').addClass("dm-gone");
	} else {
		$('#dm-test-drive-button').click(function() {
			spa.visit('/c/testdrive');
		});
	}
	$('#dm-help-button').click(function() { spa.visit('/c/help'); });
}

function helpPage() {
	document.title = "Delight/Meditate - help";
	showPage(Handlebars.templates.help({}));
}

function contactPage() {
	document.title = "Delight/Meditate - contact";
	showPage(Handlebars.templates.contact({}));
	$('#contact-alert').hide();
	$('#dm-contact-button').click(function() { contact(); })
	$('#dm-contact-cancel').click(function() { window.history.back(); });
}

// home page
//
// in general, when changes are needed, like moving to a new reading,
// we re-render the page, rather than tweaking the DOM in-place.
// this simplifies things and does not hurt performance, since rendering is client-side.
//
// if args (aindex, sindex, prevday, dir) are present, an animated citation
// transition is required. in this case, state.activities contains
// the latest state; the args specify the citation and its previous state.
//
// the transition is accomplished by putting two elements side-by-side
// with a clip-path only large enough for one of them, and animating
// a horizontal move, which slides one out of the clip area and the other into it.
//
// for activities that do NOT currently have an animation playing,
// the left-hand element is visible, and the right-hand element is clipped out
// and has css "pointer-events" set to "none", to work around a webkit bug,
// which causes clipped-out elements to eat touch events.

function homePage(ignore1, ignore2, ignore3, testdrivearg, aindex, sindex, prevday, dir) {

	if (clientUpdateRequested) {
		clientUpdateRequested = false;
		quietReload("Due to a server software upgrade, this page will now be reloaded.");
	}

	var testdrive = false;
	if (testdrivearg !== undefined && testdrivearg)
		testdrive = true;

	if (state.activities.length == 0) {
		document.title = "Delight/Meditate - home";
		showPage(dialogPage({		
			msg: "You don't currently have any reading plans.",
			buttonlabel: "Add a plan"
		}));
		$('#dm-dialog-button button').click(function() { spa.visit('/c/plans'); });
		$('#dm-dialog-button').removeClass("dm-gone");
		return
	}

	var animate = false;
	if (dir !== undefined)
		animate = true;

	// array of URLs that take us to bible readings
	var citurls = [];
	var cits = 0;

	// fetch any required documents, as asynchronously as possible.
	// each plan must be received before its streams can be fetched, so we have
	// a two-level hierarchy of promises, plus one more for queries/providers.
	//
	// rather than use the array of return values from Promise.all,
	// we rely on Promise.all to have cached everything, and read docs from
	// cache as needed. this is simpler and is not noticeably inefficient,
	// because the docs are very small.
	var promises = state.activities.map(function (activity) {
			return doc("plan", activity.plan).then(function(s) {
				var plan = JSON.parse(s);
				return Promise.all(plan.streams.map(streamdoc))
			});
		});
	promises.push(doc("queries", "providers"));
	Promise.all(promises).then(function() {

		// now have all docs, can access docs synchronously (from cache) from here on

		var provobj = JSON.parse(docfromcache("queries", "providers"));

		// add stuff to state.activities, for use by handlebars in rendering the home page.

		for (var a = 0;  a < state.activities.length;  a++) {
			var act = state.activities[a];
			var plan = JSON.parse(docfromcache("plan", act.plan));

			act.plantitle = plan.title;
			act.actindex = a;
			act.status = "";
			act.statusgone = "";
			act.ontrack = "";

			if (act.accvisible) {
				act.showduedate = "Turn off Accountability";
				act.armenuhide = "";

				var totaldue = (daysSince1970() - act.accstartdate + 1) * act.day.length;
				var duenow = totaldue - readingsDone(act);
	
				if (duenow > 0)
					act.status = "Readings due now: "+duenow;
				else if (duenow === 0)
					act.status = "You are up to date.";
				else
					act.status = "You are ahead by: "+(-duenow);
	
				if (duenow <= 0)
					act.ontrack = "dm-panel-ontrack";
			} else {
				act.showduedate = "Turn on Accountability";
				act.armenuhide = "dm-gone";
				act.statusgone = "dm-gone";
			}

			act.daysleft = Math.ceil(((plan.days * act.day.length) - readingsDone(act)) / act.day.length);

			act.readings = [];
			for (var r = 0;  r < act.day.length;  r++) {
				var day = act.day[r];
				var lines = docfromcache("stream", plan.streams[r]).split('\n');
				var cit = (day == plan.days + 1  ?  "[done]"  :  lines[day-1]);

				var reading = {};
				reading.plan = act.plan;
				reading.plandays = plan.days;
				reading.actindex = a;
				reading.streamindex = r;

				reading.citleft = reading.citright = citpretty(cit);

				reading.tctxl = reading.tctxr = "Day " + day;

				if (animate  &&  a === aindex  &&  r === sindex) {
					var prevtctx = "Day " + prevday + ':';
					if (dir === "L") {
						reading.citleft = citpretty(lines[prevday-1]);
						reading.tctxl = prevtctx;
					} else {
						reading.citright = citpretty(lines[prevday-1]);
						reading.tctxr = prevtctx;
					}
				}

				reading.ltriangle = reading.rtriangle = "dm-triangle";
				if (day == 1)
					reading.ltriangle = "dm-triangle-greyed";
				if (day == plan.days + 1)
					reading.rtriangle = "dm-triangle-greyed";

				var gotprov = false;
				for (var p = 0; p < provobj.providers.length; p++) {
					if (provobj.providers[p].provider == act.provider) {
						reading.citurlindex = cits;
						citurls[cits++] = citurl(provobj.providers[p].url, act.version, cit);
						gotprov = true;
						break;
					}
				}
				if (!gotprov)
					citurls[cits++] = '/c/badprovider';

				act.readings[r] = reading;
			}

			state.activities[a] = act;
		}

		// render the home page.
		// now can insert click handlers and do other DOM operations.
		// DON'T make the mistake of trying to access the DOM for the home page before this!  :-)

		showPage(Handlebars.templates.home(state));

		if (testdrive) {
			document.title = "Delight/Meditate - test drive";
			$('.dm-home-plan-dropdown').remove();
			testdriveMsgDisplay();
			$('#dm-testdrive').removeClass("dm-gone");
			$('#dm-testdrive-ok').click(function() { window.history.back(); });
		} else {
			document.title = "Delight/Meditate - home";
			$('.dm-home-testdrive-dropdown').remove();
		}

		$('.dm-cit-left').click(function() {
			// PARENT: id="dm-cit-{{citurlindex}}"
			visitExternal(citurls[Number($(this).closest('.dm-citation').prop("id").substr(7))]);
		});

		if (animate) {
			var wipe = "";
			if (dir === "L")
				wipe = "dm-wipeleft";   // user clicked the RIGHT arrow
			else
				wipe = "dm-wiperight";  // user clicked the LEFT arrow

			var lold = $('#dm-cit-L-'+aindex+'-'+sindex)
			var rold = $('#dm-cit-R-'+aindex+'-'+sindex)

			var lnew = lold.clone(true).addClass(wipe);
			var rnew = rold.clone(true).addClass(wipe);

			// work-around for webkit bug (clipped-out elements eat touch events)
			if (dir === "L") {
				lnew.addClass("dm-noclick");
				rnew.removeClass("dm-noclick");
				rnew.click(function() {
					// PARENT: id="dm-cit-{{citurlindex}}"
					visitExternal(citurls[Number($(this).closest('.dm-citation').prop("id").substr(7))]);
				});
			}

			// stuff to do after animation ends
			lnew.one('webkitAnimationEnd oanimationend msAnimationEnd animationend', function(e) {
				if (testdrive) {
					delay(500).then(function() {
						switch (readingsDone(state.activities[0])) {
							case 0:
								state.testdrivemsg = 1;
								break;
							case 1:
								state.testdrivemsg = 2;
								break;
							default:
								state.testdrivemsg = 3;
								break;
						}
						testdriveMsgDisplay();
					});
				}
			});

			// DOM surgery that initiates the animation
			rold.before(rnew);
			rold.remove();
			lold.before(lnew);
			lold.remove();

			// animate the temporal context. (ignore webkit bug, since nothing clickable).

			var tctxlold = $('#dm-tctx-L-'+aindex+'-'+sindex)
			var tctxrold = $('#dm-tctx-R-'+aindex+'-'+sindex)

			var tctxlnew = tctxlold.clone(true).addClass(wipe);
			var tctxrnew = tctxrold.clone(true).addClass(wipe);

			tctxlold.before(tctxlnew);
			tctxlold.remove();
			tctxrold.before(tctxrnew);
			tctxrold.remove();
		}

		$('.dm-acctoggle').click(function() {
			// id="dm-acctoggle-{{plan}}"
			acctoggle($(this).prop("id").substr(13));
		});

		$('.dm-accreset').click(function() {
			// id="dm-accreset-{{plan}}"
			accreset($(this).prop("id").substr(12));
		});

		$('.dm-delplan').click(function() {
			// id="dm-delplan-{{plan}}"
			delplan($(this).prop("id").substr(11));
		});

		// user clicked right or left arrow
		$('.dm-move').click(function() {
			// id="dm-L-{{actindex}}-{{streamindex}}-{{plandays}}"
			// id="dm-R-{{actindex}}-{{streamindex}}-{{plandays}}"
			var idfields = $(this).prop("id").split("-");

			var delta = (idfields[1] === "L" ? -1 : 1);
			var actindex = Number(idfields[2]);
			var streamindex = Number(idfields[3]);
			var plandays = Number(idfields[4]);

			if (!testdrive) {
				daychange(actindex, streamindex, delta);
			} else {
				var day = state.activities[actindex].day[streamindex];
				if (delta == -1 && day > 1) {
					state.activities[actindex].day[streamindex]--;
					homePage({}, {}, {}, true, actindex, streamindex, day, "R");
				}
				if (delta == 1 && day <= plandays) {
					state.activities[actindex].day[streamindex]++;
					homePage({}, {}, {}, true, actindex, streamindex, day, "L");
				}
			}
		});

	});
}

// how many readings have been completed, since the beginning?
function readingsDone(act) {
	var readingsdone = 0;
	for (var d = 0; d < act.day.length; d++)
		readingsdone += act.day[d] - 1;
	return readingsdone;
}

// this is a home page mock-up, using the real home page template.
function testdrivePage() {
	var day = 1;

	state.activities = [
		{
			plan:         "mcheyne-2yr-2perday",
			day:          [ day, day ],
			version:      "ESV",
			provider:     "BibleGateway.com-PRINT",
			accstartdate: daysSince1970() - day + 1,
			accvisible:   true,
		},
	];

	state.testdrivemsg = 1;

	homePage({}, {}, {}, true);
}

// hide all testdrive messages except the currently-selected one.
function testdriveMsgDisplay() {
	var msgs = 3;

	// div ids are sequential, starting with 1
	for (var i = 1;  i <= msgs;  i++) {
		if (state.testdrivemsg === i)
			$('#dm-testdrive-msg-'+i).removeClass('dm-gone');
		else
			$('#dm-testdrive-msg-'+i).addClass('dm-gone');
	}
}

function plansPage() {
	doc("queries", "plandescs").then(function(s) {
		var descs = JSON.parse(s);

		document.title = "Delight/Meditate - plans";
		showPage(Handlebars.templates.plans(descs));

		$('.dm-plan-details-button').click(function() {
			// id="dm-pdb-{{name}}"
			spa.visit('/c/plandetails/' + $(this).prop("id").substr(7))
		});
		if (state.loggedIn) {
			$('.dm-plan-add-button button').click(function() {
				// id="dm-addplan-{{name}}"
				addplan($(this).prop("id").substr(11));
			});

			// hide "add" buttons for plans that are currently active
			for (var a = 0;  a < state.activities.length; a++) {
				$('#pb-'+state.activities[a].plan).addClass("dm-gone");
			}
		} else {
			$('.dm-plan-add-button').addClass('dm-gone');
			$('.dm-plan-login-button').removeClass('dm-gone');
			$('.dm-plan-login-button button').click(function() { spa.visit('/c/login'); });
		}
	});
}

function plandetailsPage(params) {
	var planName = params.plan;
	var planDays;
	var planTitle;
	doc("plan", planName).then(function(s) {
		var plan = JSON.parse(s);
		planDays = plan.days;
		planTitle = plan.title;
		return Promise.all(plan.streams.map(streamdoc))
	}).then(function(streamdocs) {
		var streamlines = [];
		for (var i = 0; i < streamdocs.length; i++)
			streamlines[i] = streamdocs[i].split('\n');

		var details = {};
		details.title = planTitle;
		details.days = planDays;
		details.rows = [];

		for (var day = 0; day < planDays; day++) {
			var row = {};
			row.day = day + 1;
			row.readings = [];
			for (var col = 0; col < streamdocs.length; col++) {
				row.readings[col] = streamlines[col][day];
			}
			details.rows[day] = row;
		}
		document.title = "Delight/Meditate - plan details";
		showPage(Handlebars.templates.plandetails(details));

		$('#dm-plandetails-ok').click(function() { window.history.back(); });
	});
}

function plandayPage(params) {
	var planName = params.plan;
	var aindex = actindex(planName);
	if (aindex === -1) {
		return;
	}
	doc("plan", planName).then(function(s) {
		var plan = JSON.parse(s);
		var pd = {};
		pd.plan = planName;
		pd.totaldays = plan.days;
		pd.plantitle = plan.title;

		// "current day" is defined as earliest day of all streams
		pd.day = state.activities[aindex].day[0];
		for (var i = 1; i < state.activities[aindex].day.length; i++) {
			var day = state.activities[aindex].day[i];
			if (day < pd.day)
				pd.day = day;
		}

		document.title = "Delight/Meditate - jump to a specific day";
		showPage(Handlebars.templates.planday(pd));

		$('.dm-planday').click(function() {
			// id="dm-planday-{{plan}}"
			newplanday($(this).prop("id").substr(11));
		});
		$('#dm-planday-panel').on("keydown", function(e) {
			if (e.keyCode === 13) {
				e.preventDefault();
				newplanday($(this).find('.dm-planday').prop("id").substr(11));
			}
		});

		$('.dm-plandet').click(function() {
			// id="dm-plandet-{{plan}}"
			spa.visit('/c/plandetails/' + $(this).prop("id").substr(11));
		});

		$('#dm-planday-cancel').click(function() { window.history.back(); });

		$('#planday-alert').hide();
	});
}

function versionPage(params) {
	var planName = params.plan;
	var aindex = actindex(planName);
	if (aindex === -1) {
		return;
	}
	Promise.all([doc("plan", planName), doc("queries", "providers")]).then(function(dd) {
		var plan = JSON.parse(dd[0]);
		var bv = {};
		bv.plan = planName;
		bv.plantitle = plan.title;
		bv.version = state.activities[aindex].version;
		bv.choices = provChoices(dd[1], state.activities[aindex].provider);
		document.title = "Delight/Meditate - change Bible version";
		showPage(Handlebars.templates.version(bv));

		$('.dm-newver').click(function() {
			// id="dm-newver-{{plan}}"
			newversion($(this).prop("id").substr(10));
		});
		$('#dm-version-panel').on("keydown", function(e) {
			if (e.keyCode === 13) {
				e.preventDefault();
				newversion($(this).find('.dm-newver').prop("id").substr(10));
			}
		});

		$('#dm-newver-cancel').click(function() { window.history.back(); });

		$('#version-alert').hide();
	});
}

function profilePage() {
	Promise.all([xhrPostJson("/userget", ""), doc("queries", "providers")]).then(function(pp) {
		var u = JSON.parse(pp[0].responseText);
		var p = {};
		p.firstname = u.firstname;
		p.lastname = u.lastname;
		p.email = u.email;
		p.version = u.version;
		p.choices = provChoices(pp[1], u.provider);
		document.title = "Delight/Meditate - profile";
		showPage(Handlebars.templates.profile(p));
		$('#profile-alert').hide();
		$('#dm-user-profile').click(function() { userprofile(); });
		$('#dm-profile-panel').on("keydown", function(e) {
			if (e.keyCode === 13) {
				e.preventDefault();
				userprofile();
			}
		});
		$('#dm-profile-cancel').click(function() { window.history.back(); });
	}, function(mixed) {
		if (mixed[0].dmhandled)
			return;
		if (mixed[0].hasOwnProperty('responseText'))
			modalAlert("profile fetch failed - " + mixed.responseText);
		else if (mixed[0].hasOwnProperty('status'))
			modalAlert("profile fetch failed - " + mixed.status.toString());
		else
			modalAlert("profile fetch failed");
	});
}

function ecpendingPage() {
	document.title = "Delight/Meditate - email change sent";
	showPage(dialogPage({		
		msg: "You will receive an email shortly. Please follow its instructions, to verify your new email address.",
		buttonlabel: "OK"
	}));
	$('#dm-dialog-button button').click(function() { spa.visit('/c/home'); });
	$('#dm-dialog-button').removeClass("dm-gone");
}

function ecdonePage() {
	document.title = "Delight/Meditate - email change completed";
	showPage(dialogPage({		
		msg: "Your email address change has been confirmed! Please log in with your new email address.",
		buttonlabel: "Go to login page"
	}));
	$('#dm-dialog-button button').click(function() { spa.visit('/c/login'); });
	$('#dm-dialog-button').removeClass("dm-gone");
}

function badProviderPage() {
	document.title = "Delight/Meditate - problem with Bible provider";
	showPage(dialogPage({		
		msg: "There is a problem with this Bible provider. Please click the 3-dot menu for this plan, choose Change Bible Version, and select a new Bible provider.",
		buttonlabel: "OK"
	}));
	$('#dm-dialog-button button').click(function() { window.close(); });
	$('#dm-dialog-button').removeClass("dm-gone");
}


// actions invoked by UI events, which interact with the server REST API

function signup() {
	$('#signup-alert').hide();
	var signupParams = {
		firstname: $('#sfFirst').val(),
		lastname: $('#sfLast').val(),
		email: $('#sfEmail').val(),
		password: $('#sfPass').val(),
	};

	xhrPostJson("/signup", JSON.stringify(signupParams)).then(function(/* okXHR */) {
		spa.replace('/c/pending');
	}, function(errorXHR) {
		formErrorAlert('#signup-alert', errorXHR, "signup failed");
	});
}

function login() {
	$('#login-alert').hide();
	var loginParams = {
		email: $('#lfEmail').val(),
		password: $('#lfPass').val(),
		keep: $('#lfKeep').is(':checked'),
	};
	xhrPostJson("/login", JSON.stringify(loginParams)).then(function(okXHR) {
		state.loggedIn = true;
		receiveStateResponse(JSON.parse(okXHR.responseText));
		spa.visit('/c/home');
	}, function(errorXHR) {
		formErrorAlert('#login-alert', errorXHR, "login failed");
	});
}

function logout() {
	xhrPostJson("/logout","").then(function(/* okXHR */) {
		stateLoggedOut();
		spa.replace('/c/login');
	}, function(errorXHR) {
		if (errorXHR.dmhandled)
			return;
		modalAlert("logout error "+errorXHR.status.toString()).then(function() {
			stateLoggedOut();
			spa.replace('/c/login');
		});
	});
}

function forgot() {
	$('#forgot-alert').hide();
	var pwrecoverParams = {
		email: $('#prEmail').val()
	};

	xhrPostJson("/pwforgot", JSON.stringify(pwrecoverParams)).then(function(/* okXHR */) {
		spa.visit('/c/forgotpending');
	}, function(errorXHR) {
		formErrorAlert('#forgot-alert', errorXHR, "password recovery failed");
	});
}

function pwreset() {
	$('#pwreset-alert').hide();
	var pwresetParams = {
		token: $('#prToken').val(),
		password: $('#prPass').val()
	};

	xhrPostJson("/dorecover", JSON.stringify(pwresetParams)).then(function(/* okXHR */) {
		modalAlert("Your password has been changed. Please log in with your new password.").then(function() {
			spa.visit('/c/login');
		});
	}, function(errorXHR) {
		formErrorAlert('#pwreset-alert', errorXHR, "password reset failed");
	});
}

function contact() {
	$('#contact-alert').hide();
	var contactParams = {
		msg: $('#contactText').val(),
	};
	xhrPostJson("/contact", JSON.stringify(contactParams)).then(function(/* okXHR */) {
		modalAlert("Your message has been sent.").then(function() {
			window.history.back();
		});
	}, function(errorXHR) {
		formErrorAlert('#contact-alert', errorXHR, "problem sending message");
	});	
}

function userprofile() {
	var firstName = $('#sfFirst').val()
	var uParams = {
		email: $('#sfEmail').val(),
		oldpass: $('#oldPass').val(),
		newpass: $('#newPass').val(),
		firstname: firstName,
		lastname: $('#sfLast').val(),
		version: $('#bvVersion').val(),
		provider: $('#bvProvider').val(),
	};
	xhrPostJson("/userset", JSON.stringify(uParams)).then(function(okXHR) {
		var resp = JSON.parse(okXHR.responseText);
		state.firstName = resp.firstname;
		navSetup();
		if (resp.logout)
			stateLoggedOut();
		if (resp.emailsent)
			spa.visit('/c/ecpending');
		else if (resp.pwchanged)
			modalAlert("Your password has been changed. Please log in with your new password.").then(function() {
				spa.visit('/c/login');
			});
		else
			spa.visit('/c/home');
	}, function(errorXHR) {
		formErrorAlert('#profile-alert', errorXHR, "problem updating profile");
	});
}

function addplan(name) {
	var aaParams = {
		plan: name,
		today: daysSince1970(),
	};
	xhrPostJson("/actadd", JSON.stringify(aaParams)).then(function(okXHR) {
		receiveStateResponse(JSON.parse(okXHR.responseText));
		spa.visit('/c/home');
	}, function(errorXHR) {
		if (!errorXHR.dmhandled)
			modalAlert("couldn't add plan - "+errorXHR.status.toString());
	});
}

function delplan(planName) {
	modalConfirm("Delete this plan? (If you change your mind, you can add it again at any time.)").then(function() {
		var adParams = {
			plan: planName,
		};
		xhrPostJson("/actdel", JSON.stringify(adParams)).then(function(/* okXHR */) {
			var aindex = actindex(planName);
			if (aindex === -1) {
				return;
			}
			state.activities.splice(aindex, 1);
			spa.replace('/c/home');
		}, function(errorXHR) {
			if (!errorXHR.dmhandled)
				modalAlert("couldn't delete plan - "+errorXHR.status.toString());
		});
	});
}

function daychange(aindex, sindex, delta) {
	var prevday = state.activities[aindex].day[sindex];
	var planName = state.activities[aindex].plan;
	var dcParams = {
		plan: planName,
		streamindex: sindex,
		prevday: prevday,
		delta: delta,
		today: daysSince1970(),
	};
	xhrPostJson("/daychange", JSON.stringify(dcParams)).then(function(okXHR) {
		receiveDayResponse(planName, okXHR.responseText);
		var newday = state.activities[aindex].day[sindex];
		if (newday === prevday + 1)
			homePage({}, {}, {}, false, aindex, sindex, prevday, "L");  // render with an animated transition
		else if(newday === prevday - 1)
			homePage({}, {}, {}, false, aindex, sindex, prevday, "R");  // render with an animated transition
		else
			homePage(); // no animated transition
	}, function(errorXHR) {
		if (!errorXHR.dmhandled)
			modalAlert("move failed - "+errorXHR.status.toString()+" - "+errorXHR.responseText);
	});
}

function newplanday(planName) {
	var ajParams = {
		plan: planName,
		day: Number($('#pdfDay').val()),
		today: daysSince1970(),
	};
	xhrPostJson("/actjump", JSON.stringify(ajParams)).then(function(okXHR) {
		receiveDayResponse(planName, okXHR.responseText);
		window.history.back();
	}, function(errorXHR) {
		formErrorAlert('#planday-alert', errorXHR, "could not change day");
	});
}

function accreset(planName) {
	modalConfirm("Reset the due dates for this plan, to make your current reading(s) due today?").then(function() {
		var params = {
			plan: planName,
			today: daysSince1970(),
		};
		xhrPostJson("/accreset", JSON.stringify(params)).then(function(okXHR) {
			receiveDayResponse(planName, okXHR.responseText);
			homePage();  // re-render without touching history or scroll
		}, function(errorXHR) {
			if (!errorXHR.dmhandled)
				modalAlert("couldn't reset accountability - "+errorXHR.status.toString());
		});
	});
}

function acctoggle(planName) {
	var aindex = actindex(planName);
	if (aindex === -1) {
		return;
	}

	var enable = !(state.activities[aindex].accvisible);

	var params = {
		plan: planName,
		enabled: enable,
	};
	xhrPostJson("/accenab", JSON.stringify(params)).then(function(okXHR) {
		receiveDayResponse(planName, okXHR.responseText);
		homePage();  // re-render without touching history or scroll
	}, function(errorXHR) {
		if (!errorXHR.dmhandled)
			modalAlert("couldn't toggle accountability - "+errorXHR.status.toString());
	});
}

function newversion(planName) {
	var params = {
		plan: planName,
		provider: $('#bvProvider').val(),
		version: $('#bvVersion').val(),
	};
	xhrPostJson("/actver", JSON.stringify(params)).then(function(/* okXHR */) {
		var aindex = actindex(planName);
		if (aindex === -1) {
			return;
		}
		state.activities[aindex].provider = params.provider;
		state.activities[aindex].version = params.version;
		window.history.back();
	}, function(errorXHR) {
		formErrorAlert('#version-alert', errorXHR, "could not change version");
	});
}

// helpers

// return the index of the activity corresponding to the given plan name.
// it should always be found. if not, something is seriously wrong.
function actindex(planName) {
	for (var i = 0; i < state.activities.length; i++) {
		if (state.activities[i].plan === planName)
			return i;
	}
	console.log("delightmeditate - actindex(" + planName + ") - not found");
	return -1;
}

// several API endpoints return a DayResponse.
// this helper updates the relevant activity.
function receiveDayResponse(planName, responseText) {
	var resp = JSON.parse(responseText);
	var aindex = actindex(planName);
	if (aindex === -1) {
		return 0;
	}
	state.activities[aindex].day = resp.day;
	state.activities[aindex].accstartdate = resp.accstartdate;
	state.activities[aindex].accvisible = resp.accvisible;
}

// given the queries/providers document in unparsed json form,
// return an array of { provider, selected } objects, for a pulldown template
function provChoices(provdoc, curprov) {
	var provobj = JSON.parse(provdoc);
	var choices = [];
	for (var i = 0; i < provobj.providers.length; i++) {
		choices[i] = {};
		choices[i].provider = provobj.providers[i].provider;
		choices[i].selected = "";
		if (choices[i].provider === curprov)
			choices[i].selected = "selected";
	}

	return choices;
}

// visit an external web page in a new tab
function visitExternal(s) {
	window.open(s, '_blank');
}

function pathFirstSegment(s) {
	if (s.charAt(0) === '/')
		s = s.substr(1);
	var end = s.indexOf("/");
	if (end != -1)
		s = s.substr(0,end);
	return s;
}

// helper to assemble a citation URL
function citurl(url, version, citation) {
	var s = url.replace('{{version}}', version);
	return s.replace('{{citation}}', citation);
}

function delay(msec) {
	return new Promise(function(resolve) {
		setTimeout(function() { resolve(); }, msec);
	});
}

function secsSince1970() {
	var tzoffset = new Date().getTimezoneOffset();  // in mins
	return Math.floor((Date.now() / 1000) - (tzoffset * 60));
}

function daysSince1970() {
	return Math.floor(secsSince1970() / (60*60*24));
}


/*
function temporalcontext(day, desired) {
	if (day == desired)
		return "today";
	if (desired - day === 1)
		return "tomorrow";
	if (desired - day > 1)
		return (desired - day).toString() + " days from now";
	if (desired - day === -1)
		return "yesterday";
	if (desired - day < -1)
		return (day - desired).toString() + " days ago";
}
*/


// tell the browser to reload, which re-fetches index.html
// and re-initializes javascript.
//
// the usual reason for reloading the client is that the server has
// indicated that a newer version of the client is available.
//
// reload() is a low-level function and should not be called by user code.
// it saves the time of the most recent call reload in localstorage.
// this is just for the purpose of rate-limiting calls to reload;
// it does not pay attention to user-initiated page loads or reloads.
function reload() {
	localStorage.setItem(lrkey, secsSince1970().toString());
	history.replaceState({}, "", "/");
	window.location.reload(true);
}

// call modalReload when a reload is liable to abort a user action,
// to alert the user before performing the reload.
function modalReload(msg) {
	modalAlert(msg).then(function() {
		reload();
	});
}

const minQuietReloadIntervalSecs = 60*60;  // one hour

// call quietReload if the reload can be performed without disrupting the
// user beyond a momentary flash. it attempts to perform the reload
// without alerting the user.
//
// a non-alerting reload carries the potential danger that application-level
// bugs can result in an infinite reload loop, turning each user into the
// unwilling source of a DDOS attack. to protect against this possibility,
// we check to see how long it's been since the last reload, and, if below
// a threshold, alert the user first. thus, the (rare) case of multiple
// legitimate client updates in a short period is still handled correctly,
// with only the inconvenience of the user having to dismiss an alert.
//
// this possibility of alerting the user requires that a message be supplied
// as an argument, even though it will rarely or never be used.
function quietReload(msg) {
	var s = localStorage.getItem(lrkey);
	if (s == null) {
		reload();
		return;
	}

	var last = parseInt(s, 10);

	if (secsSince1970() - last < minQuietReloadIntervalSecs) {
		modalAlert(msg).then(function() {
			reload();
		});
	} else
		reload();
}


// in-memory document cache

var docs = {};

// return the contents of a document that is known to be in cache
function docfromcache(dir, name) {
	var fullname = dir + '/' + name;
	return docs[fullname];
}

// return a promise that resolves to the contents of a document.
// use and maintain caches in memory and in localStorage.
function doc(dir, name) {
	var fullname = dir + '/' + name;
	var expname = 'exp/' + fullname;

	if (docs.hasOwnProperty(expname)) {
		if (docs[expname] > secsSince1970()  &&  docs.hasOwnProperty(fullname)) {
			return Promise.resolve(docs[fullname]);
		}
	}

	var exp = localStorage.getItem(expname);
	if (exp != null  &&  Number(exp) > secsSince1970()) {
		var data = localStorage.getItem(fullname);
		if (data != null) {
			docs[fullname] = data;
			docs[expname] = Number(exp);
			return Promise.resolve(data);
		}
	}

	return xhrGet("/getdoc/"+dir+"/"+name).then(function(okXHR) {
		var resp = JSON.parse(okXHR.responseText);
		var exp = resp.ttl + secsSince1970();
		localStorage.setItem(fullname, resp.data);
		localStorage.setItem(expname, exp.toString());
		docs[fullname] = resp.data;
		docs[expname] = exp;
		return resp.data;
	}, function(errorXHR) {
		if (!errorXHR.dmhandled)
			return modalAlert("doc not found - " + fullname + " - " + errorXHR.status.toString()).then(function() { return ""; });
		return "";
	});
}

// directory-specific document function, to facilitate mapping
function streamdoc(name) {
	return doc("stream", name);
}

var booknames = [
	// put Philippians before Philemon, so it will match "Phil".
	//
	// use "Song", "1Thess", and "2Thess", to save screen real estate
	// so ABBREVIATIONS OF THESE 3 BOOKS MAY NOT BE LONGER THAN THESE STRINGS.

	"1Chronicles", "1Corinthians", "1John", "1Kings", "1Peter", "1Samuel",
	"1Thess", "1Timothy", "2Chronicles", "2Corinthians", "2John",
	"2Kings", "2Peter", "2Samuel", "2Thess", "2Timothy", "3John",
	"Acts", "Amos", "Colossians", "Daniel", "Deuteronomy", "Ecclesiastes",
	"Ephesians", "Esther", "Exodus", "Ezekiel", "Ezra", "Galatians", "Genesis",
	"Habakkuk", "Haggai", "Hebrews", "Hosea", "Isaiah", "James", "Jeremiah",
	"Job", "Joel", "John", "Jonah", "Joshua", "Jude", "Judges", "Lamentations",
	"Leviticus", "Luke", "Malachi", "Mark", "Matthew", "Micah", "Nahum",
	"Nehemiah", "Numbers", "Obadiah", "Philippians", "Philemon", "Proverbs",
	"Psalm", "Revelation", "Romans", "Ruth", "Song", "Titus",
	"Zechariah", "Zephaniah"
];

// given an abbreviation of a bible book name, return the full name.
// this is implemented simply by returning the first element in booknames
// for which the given abbreviation is a prefix.
function fullname(name) {
	for (var i = 0; i < booknames.length; i++) {
		if (booknames[i].lastIndexOf(name, 0) === 0)
			return booknames[i];
	}
	return name;
}

// change a citation's book name to its full name and change space to <br/>
function citpretty(cit) {
	var spaceNdx = cit.indexOf(" ")
	if (spaceNdx === -1)
		return fullname(cit);
	else
		return fullname(cit.slice(0, spaceNdx))+"<br/>"+cit.slice(spaceNdx+1);
}

},{"bootstrap":3,"handlebars/runtime":4,"jquery":25,"rlite-router":27,"spa-components":28}]},{},[1]);
