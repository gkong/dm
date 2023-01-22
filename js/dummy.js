/** @preserve
 *
 * @copyright   Copyright (C) 2017 George S. Kong, All Rights Reserved.
 * Use of this source code is governed by a license that can be found in the LICENSE.txt file.
 */

// dummy so that we can factor out dm.js dependencies into dmdeps.js.
// presumably, the latter changes rarely, can be cached by browsers,
// put onto a CDN, etc.
//
// this should be maintained as a clone of the corresponding lines from dm.js

jQuery = $ = require('jquery');
var Rlite = require('rlite-router');
Handlebars = require('handlebars/runtime');  // put into global scope, so dmt.js can see it
require('bootstrap');
var spa = require('spa-components');
