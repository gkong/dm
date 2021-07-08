/** @preserve
 *
 * @copyright   Copyright (C) 2017 George S. Kong, All Rights Reserved.
 * Use of this source code is governed by a license that can be found in the LICENSE.txt file.
 */

jQuery = $ = require('jquery');
var Promise = require('es6-promise').Promise;
Tether = require('tether');
require('bootstrap');

window.admincmd = admincmd;

$(function() {
	$('#admin-alert').hide();
	$('#admin-alert-ok').hide();
	$('#dm-admin-cmd').click(function() { admincmd(); })
});

var xhrPostJson = function(url, data) {
	return new Promise(function(resolve, reject) {
		var req = new XMLHttpRequest();
		req.open("POST", url, true);
		req.onreadystatechange = handler;
		req.setRequestHeader("Content-type", "application/json;charset=UTF-8");
		var tok = document.cookie.replace(/(?:(?:^|.*;\s*)dm_csrf\s*\=\s*([^;]*).*$)|^.*$/, "$1");
		if (tok != "")
			req.setRequestHeader("Dm-Csrf", tok);
		req.send(data);
		function handler() {
			if (this.readyState === this.DONE) {
				if (this.status === 200)
					resolve(this);
				else
					reject(this);
			}
		}
	});
};

function formErrorAlert(selector, errorXHR, altMsg) {
	// only display altMsg if response contained no error message
	if (!errorXHR.responseText)
		$(selector).html(altMsg + " - error " + errorXHR.status.toString());
	else
		$(selector).html(errorXHR.responseText);
	$(selector).show();
}

function admincmd() {
	$('#admin-alert').hide();
	$('#admin-alert-ok').hide();
	var params = {
		cmd: $('#adminCmd').val(),
		arg: $('#adminArg').val()
	};
	$('#admin-response-body').html("");

	xhrPostJson("/admin", JSON.stringify(params)).then(function(okXHR) {
		$('#admin-alert-ok').show();
		$('#admin-response-body').html(okXHR.responseText);
	}, function(errorXHR) {
		formErrorAlert('#admin-alert', errorXHR, "password recovery failed");
	});
}
