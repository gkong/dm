(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
templates['about'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div class=\"dm-panel\">\n	<h1>About</h1>\n	<div class=\"dm-quote-container\">\n		<hr/>\n			<div class=\"dm-quote\">\n				... his <b>delight</b> is in the law of the Lord,\n				and on his law he <b>meditates</b> day and night.\n				<div class=\"dm-quote-citation\">Psalm 1</div>\n			</div>\n		<hr/>\n	</div>\n	<p>\n		<b>Delight/Meditate</b> is a free daily Bible reading web app\n		for mobile devices and desktop web browsers.\n		It takes you through a daily reading plan\n		and is designed to be simple, quick, and frictionless.\n		For each reading, you can read your physical Bible\n		or, with a single click, view the reading\n		at one of the popular online Bible web sites,\n		in your chosen translation.\n	</p>\n	<p>\n		<b>Delight/Meditate</b> is a web app - to use it, just visit DelightMeditate.com.\n		You don't need to install anything on your device.\n	</p>\n	<p class=\"dm-visitor-only\">\n		For a quick introduction, click <b>Test drive</b> below.\n		To learn more, sign up for a free account.\n		You can experiment as much as you like.\n		It's easy to add or remove plans and to jump around within plans.\n	</p>\n	<div class=\"dm-vspacer dm-visitor-only\"></div>\n	<button type=\"button\" class=\"btn btn-primary dm-visitor-only\" id=\"dm-test-drive-button\">Test drive</button>\n	<div class=\"dm-vspacer\"></div>\n	<button type=\"button\" class=\"btn btn-primary\" id=\"dm-help-button\">How it works</button>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['contact'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div class=\"dm-panel\">\n	<div class=\"dm-title\">Contact Us</div>\n	<p>\n		Please enter a brief message and click <b>Send</b>\n	</p>\n	<form class=\"add-border dm-number-form\">\n		<div class=\"form-group\">\n			<textarea id=\"contactText\" class=\"form-control\" name=\"contact123\" rows=\"8\" maxlength=\"1000\"></textarea>\n		</div>\n		<div id=\"contact-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-contact-button\">Send</button>\n		<div class=\"dm-vspacer\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-contact-cancel\">Cancel</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['dialog'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel-fixed\">\n	<div class=\"dm-vspacer\"></div>\n	<p>\n		"
    + alias4(((helper = (helper = lookupProperty(helpers,"msg") || (depth0 != null ? lookupProperty(depth0,"msg") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"msg","hash":{},"data":data,"loc":{"start":{"line":4,"column":2},"end":{"line":4,"column":9}}}) : helper)))
    + "\n	</p>\n	<div id=\"dm-dialog-button\">\n		<div class=\"dm-vspacer\"></div>\n		<button type=\"button\" class=\"btn btn-primary\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"buttonlabel") || (depth0 != null ? lookupProperty(depth0,"buttonlabel") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"buttonlabel","hash":{},"data":data,"loc":{"start":{"line":8,"column":48},"end":{"line":8,"column":63}}}) : helper)))
    + "</button>\n	</div>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['forgot'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div class=\"dm-panel-fixed\">\n	<h1>Forgot Password</h1>\n	<p>\n		Please enter the email address for your account.\n	</p>\n	<form class=\"add-border\">\n		<div class=\"form-group\">\n			<label for=\"prEmail\">email address</label>\n			<input id=\"prEmail\" class=\"form-control\" \"type=\"text\" name=\"forgot123\" />\n		</div>\n		<div class=\"dm-vspacer\"></div>\n		<div id=\"forgot-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-forgot-button\">Continue</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['help'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div class=\"dm-panel\">\n	<h1>Help</h1>\n\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-small-title\">Daily Use</div>\n	<p>\n		The home screen lists one or more readings, typically a chapter of the Bible, like \"Genesis 1.\"\n	</p>\n	<p>\n		To read a reading, tap it.\n		You'll be taken to an external web page which displays the text.\n		You can read it there or, if the provider supports it, click a button\n		and hear the text read to you. When you're done, close the tab.\n	</p>\n	<p>\n		Back at the home page, click the right-facing arrow, which takes you to the next reading.\n		Or you can go back to a previous day by clicking the left-facing arrow.\n		As you progress through readings by clicking these arrows, the app remembers your position.\n	</p>\n\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-small-title\">Plans</div>\n	<p>\n		Click <b>Plans</b>, to browse available reading plans and add one or more plans.\n	</p>\n	<p>\n		A plan can have multiple readings per day.\n		For example, click <b>Plans</b> in the menu,\n		find <b>M'Cheyne 1-Year (3/day)</b>, and click <b>See plan details</b>.\n		This plan has 3 streams: two from the Old Testament\n		and one from the New Testament and Psalms.\n		There is one reading from each stream for every day,\n		and the plan takes one year to complete.\n	</p>\n\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-small-title\">Tracking and Accountability</div>\n	<p>\n		If you are using the plan named <b>M'Cheyne 1-Year (3/day)</b>,\n		three readings are due each day.\n		Normally, you would read one reading from each stream every day,\n		but you are not required to follow this order -\n		you can read ahead in one stream and then catch up in the other streams later.\n		The app considers you up to date if you read the expected number\n		of readings for your plan each day, regardless of which stream(s) they come from.\n	</p>\n	<p>\n		If you get behind in your reading, you can do extra reading to catch up.\n		If you get hopelessly behind, click the three-dot menu in the\n		upper-right corner of the plan, and choose <b>Reset Due Dates</b>.\n		The app will reset its accountability start date,\n		so that only one day's worth of readings are currently due.\n	</p>\n	<p>\n		To turn off accountability completely, click a plan's\n		three-dot menu and choose <b>Turn off Accountability</b>.\n		Then the app won't say anything about what is currently due;\n		it will simply track your progress, as you proceed at whatever pace you prefer.\n	</p>\n\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-small-title\">Bible Versions</div>\n	<p>\n		Bible texts are displayed by free providers, like <b>BibleGateway.com</b>.\n		Each provider supports a set of Bible versions, specified by short codes,\n		which are usually standard abbreviations, like <b>ESV</b>, <b>KJV</b>, etc.\n		You can find the codes for less commonly-used versions\n		by exploring providers' websites, selecting versions and examining the\n		web addresses in your browser's location bar.\n		</p>\n	<p>\n		You can choose the provider and version for a plan you've added,\n		by clicking the plan's three-dot menu and choosing <b>Change Bible Version</b>.\n		You can change your default settings for new plans,\n		by clicking your name in the menu and choosing <b>Profile</b>.\n	</p>\n\n	<div class=\"dm-vspacer\"></div>\n</div>\n<div class=\"dm-vspacer\"></div>\n<div class=\"dm-vspacer\"></div>\n";
},"useData":true});
templates['home'] = template({"1":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return ((stack1 = lookupProperty(helpers,"if").call(alias1,(data && lookupProperty(data,"index")),{"name":"if","hash":{},"fn":container.program(2, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":3,"column":1},"end":{"line":5,"column":8}}})) != null ? stack1 : "")
    + "	<div class=\"dm-plan dm-panel-fixed "
    + alias4(((helper = (helper = lookupProperty(helpers,"ontrack") || (depth0 != null ? lookupProperty(depth0,"ontrack") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"ontrack","hash":{},"data":data,"loc":{"start":{"line":6,"column":36},"end":{"line":6,"column":47}}}) : helper)))
    + "\">\n		<div class=\"dm-header\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"plantitle") || (depth0 != null ? lookupProperty(depth0,"plantitle") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plantitle","hash":{},"data":data,"loc":{"start":{"line":7,"column":25},"end":{"line":7,"column":38}}}) : helper)))
    + "</div>\n		<div class=\"dropdown dm-hamburger-container\">\n			<a class=\"dropdown-toggle dm-dropdown-no-triangle\" href=\"\" id=\"dropdownMenuLink\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">\n				<div class=\"dm-hamburger\"></div>	\n			</a>\n			<div class=\"dropdown-menu dropdown-menu-right dm-home-testdrive-dropdown\" aria-labelledby=\"dropdownMenuLink\">\n				<a class=\"dropdown-item\">[ menu disabled during test drive ]</a>\n			</div>\n			<div class=\"dropdown-menu dropdown-menu-right dm-home-plan-dropdown\" aria-labelledby=\"dropdownMenuLink\">\n				<a class=\"dropdown-item\" href=\"/c/plandetails/"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":16,"column":50},"end":{"line":16,"column":58}}}) : helper)))
    + "\">Plan Details</a>\n				<a class=\"dropdown-item\" href=\"/c/planday/"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":17,"column":46},"end":{"line":17,"column":54}}}) : helper)))
    + "\">Jump to a Specific Day</a>\n				<a class=\"dropdown-item dm-acctoggle\" id=\"dm-acctoggle-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":18,"column":59},"end":{"line":18,"column":67}}}) : helper)))
    + "\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"showduedate") || (depth0 != null ? lookupProperty(depth0,"showduedate") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"showduedate","hash":{},"data":data,"loc":{"start":{"line":18,"column":69},"end":{"line":18,"column":84}}}) : helper)))
    + "</a>\n				<a class=\"dropdown-item dm-accreset "
    + alias4(((helper = (helper = lookupProperty(helpers,"armenuhide") || (depth0 != null ? lookupProperty(depth0,"armenuhide") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"armenuhide","hash":{},"data":data,"loc":{"start":{"line":19,"column":40},"end":{"line":19,"column":54}}}) : helper)))
    + "\" id=\"dm-accreset-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":19,"column":72},"end":{"line":19,"column":80}}}) : helper)))
    + "\">Reset Due Dates</a>\n				<a class=\"dropdown-item\" href=\"/c/version/"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":20,"column":46},"end":{"line":20,"column":54}}}) : helper)))
    + "\">Change Bible Version</a>\n				<a class=\"dropdown-item dm-delplan\" id=\"dm-delplan-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":21,"column":55},"end":{"line":21,"column":63}}}) : helper)))
    + "\">Delete This Plan</a>\n			</div>\n		</div>\n		<div class=\"dm-status "
    + alias4(((helper = (helper = lookupProperty(helpers,"statusgone") || (depth0 != null ? lookupProperty(depth0,"statusgone") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"statusgone","hash":{},"data":data,"loc":{"start":{"line":24,"column":24},"end":{"line":24,"column":38}}}) : helper)))
    + "\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"status") || (depth0 != null ? lookupProperty(depth0,"status") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"status","hash":{},"data":data,"loc":{"start":{"line":24,"column":40},"end":{"line":24,"column":50}}}) : helper)))
    + "</div>\n"
    + ((stack1 = lookupProperty(helpers,"each").call(alias1,(depth0 != null ? lookupProperty(depth0,"readings") : depth0),{"name":"each","hash":{},"fn":container.program(4, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":25,"column":2},"end":{"line":61,"column":11}}})) != null ? stack1 : "")
    + "		<div id=\"dm-testdrive\" class=\"dm-gone\">\n			<div class=\"dm-vspacer\"></div>\n			<div id=\"dm-testdrive-msg-1\" class=\"alert alert-warning\">\n				<b>TEST DRIVE</b>\n				<br/>\n				1. Click one of the Bible references above to view a reading in a new browser tab.\n				<br/>\n				2. After finishing the reading, close its tab.\n				<br/>\n				3. Click the right-pointing arrow next to the reading, to advance to the next day's reading.\n			</div>\n			<div id=\"dm-testdrive-msg-2\" class=\"alert alert-warning\">\n				<b>TEST DRIVE</b>\n				<br/><br/>\n				You have completed a reading!\n				<br/><br/>\n				When you are logged in, the app remembers where you are,\n				based on your arrow clicks.\n				<br/><br/>\n				Now read the other reading for today and click its right-pointing arrow.\n 			</div>\n			<div id=\"dm-testdrive-msg-3\" class=\"alert alert-warning\">\n				<b>TEST DRIVE</b>\n				<br/><br/>\n				You have completed today's readings.\n				The green background indicates there is nothing more due today.\n				<br/><br/>\n				This completes the test drive.\n				You can try clicking the arrows some more,\n				to become more familiar with the app.\n				<div class=\"dm-vspacer\"></div>\n				<button type=\"button\" class=\"btn btn-primary\" id=\"dm-testdrive-ok\">OK</button>\n 			</div>\n 		</div>\n		<div class=\"dm-footer\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"daysleft") || (depth0 != null ? lookupProperty(depth0,"daysleft") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"daysleft","hash":{},"data":data,"loc":{"start":{"line":96,"column":25},"end":{"line":96,"column":37}}}) : helper)))
    + " days left</div>\n	</div>\n";
},"2":function(container,depth0,helpers,partials,data) {
    return "		<div class=\"dm-vspacer\"></div>\n";
},"4":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "			<div class=\"dm-vspacer\"></div>\n\n			<div class=\"dm-temporal-context\">\n				<div class=\"dm-tctx-left\" id=\"dm-tctx-L-"
    + alias4(((helper = (helper = lookupProperty(helpers,"actindex") || (depth0 != null ? lookupProperty(depth0,"actindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"actindex","hash":{},"data":data,"loc":{"start":{"line":29,"column":44},"end":{"line":29,"column":56}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"streamindex") || (depth0 != null ? lookupProperty(depth0,"streamindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"streamindex","hash":{},"data":data,"loc":{"start":{"line":29,"column":57},"end":{"line":29,"column":72}}}) : helper)))
    + "\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"tctxl") || (depth0 != null ? lookupProperty(depth0,"tctxl") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"tctxl","hash":{},"data":data,"loc":{"start":{"line":29,"column":74},"end":{"line":29,"column":83}}}) : helper)))
    + "</div><!--\n				--><div class=\"dm-tctx-right\" id=\"dm-tctx-R-"
    + alias4(((helper = (helper = lookupProperty(helpers,"actindex") || (depth0 != null ? lookupProperty(depth0,"actindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"actindex","hash":{},"data":data,"loc":{"start":{"line":30,"column":48},"end":{"line":30,"column":60}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"streamindex") || (depth0 != null ? lookupProperty(depth0,"streamindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"streamindex","hash":{},"data":data,"loc":{"start":{"line":30,"column":61},"end":{"line":30,"column":76}}}) : helper)))
    + "\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"tctxr") || (depth0 != null ? lookupProperty(depth0,"tctxr") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"tctxr","hash":{},"data":data,"loc":{"start":{"line":30,"column":78},"end":{"line":30,"column":87}}}) : helper)))
    + "</div>\n			</div>\n\n			<div class=\"dm-row\">\n				<div class=\"dm-col dm-move-left dm-move\" id=\"dm-L-"
    + alias4(((helper = (helper = lookupProperty(helpers,"actindex") || (depth0 != null ? lookupProperty(depth0,"actindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"actindex","hash":{},"data":data,"loc":{"start":{"line":34,"column":54},"end":{"line":34,"column":66}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"streamindex") || (depth0 != null ? lookupProperty(depth0,"streamindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"streamindex","hash":{},"data":data,"loc":{"start":{"line":34,"column":67},"end":{"line":34,"column":82}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plandays") || (depth0 != null ? lookupProperty(depth0,"plandays") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plandays","hash":{},"data":data,"loc":{"start":{"line":34,"column":83},"end":{"line":34,"column":95}}}) : helper)))
    + "\">\n					<div class=\"vcenter-container\">\n						<div class=\"vcenter dm-line-height-hack\">\n							<svg>\n								<polygon class=\""
    + alias4(((helper = (helper = lookupProperty(helpers,"ltriangle") || (depth0 != null ? lookupProperty(depth0,"ltriangle") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"ltriangle","hash":{},"data":data,"loc":{"start":{"line":38,"column":24},"end":{"line":38,"column":37}}}) : helper)))
    + "\" points=\"16 0 0 25 16 50\">\n							</svg>\n						</div>\n					</div>\n				</div><!--\n				--><div class=\"dm-col dm-citation\" id=\"dm-cit-"
    + alias4(((helper = (helper = lookupProperty(helpers,"citurlindex") || (depth0 != null ? lookupProperty(depth0,"citurlindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"citurlindex","hash":{},"data":data,"loc":{"start":{"line":43,"column":50},"end":{"line":43,"column":65}}}) : helper)))
    + "\">\n					<div class=\"vcenter-container dm-cit-left\" id=\"dm-cit-L-"
    + alias4(((helper = (helper = lookupProperty(helpers,"actindex") || (depth0 != null ? lookupProperty(depth0,"actindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"actindex","hash":{},"data":data,"loc":{"start":{"line":44,"column":61},"end":{"line":44,"column":73}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"streamindex") || (depth0 != null ? lookupProperty(depth0,"streamindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"streamindex","hash":{},"data":data,"loc":{"start":{"line":44,"column":74},"end":{"line":44,"column":89}}}) : helper)))
    + "\">\n						<div class=\"vcenter\">"
    + ((stack1 = ((helper = (helper = lookupProperty(helpers,"citleft") || (depth0 != null ? lookupProperty(depth0,"citleft") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"citleft","hash":{},"data":data,"loc":{"start":{"line":45,"column":27},"end":{"line":45,"column":40}}}) : helper))) != null ? stack1 : "")
    + "</div>\n					</div><!--\n					--><div class=\"vcenter-container dm-cit-right dm-noclick\" id=\"dm-cit-R-"
    + alias4(((helper = (helper = lookupProperty(helpers,"actindex") || (depth0 != null ? lookupProperty(depth0,"actindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"actindex","hash":{},"data":data,"loc":{"start":{"line":47,"column":76},"end":{"line":47,"column":88}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"streamindex") || (depth0 != null ? lookupProperty(depth0,"streamindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"streamindex","hash":{},"data":data,"loc":{"start":{"line":47,"column":89},"end":{"line":47,"column":104}}}) : helper)))
    + "\">\n						<div class=\"vcenter\">"
    + ((stack1 = ((helper = (helper = lookupProperty(helpers,"citright") || (depth0 != null ? lookupProperty(depth0,"citright") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"citright","hash":{},"data":data,"loc":{"start":{"line":48,"column":27},"end":{"line":48,"column":41}}}) : helper))) != null ? stack1 : "")
    + "</div>\n					</div>\n				</div><!--\n				--><div class=\"dm-col dm-move-right dm-move\" id=\"dm-R-"
    + alias4(((helper = (helper = lookupProperty(helpers,"actindex") || (depth0 != null ? lookupProperty(depth0,"actindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"actindex","hash":{},"data":data,"loc":{"start":{"line":51,"column":58},"end":{"line":51,"column":70}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"streamindex") || (depth0 != null ? lookupProperty(depth0,"streamindex") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"streamindex","hash":{},"data":data,"loc":{"start":{"line":51,"column":71},"end":{"line":51,"column":86}}}) : helper)))
    + "-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plandays") || (depth0 != null ? lookupProperty(depth0,"plandays") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plandays","hash":{},"data":data,"loc":{"start":{"line":51,"column":87},"end":{"line":51,"column":99}}}) : helper)))
    + "\">\n					<div class=\"vcenter-container\">\n						<div class=\"vcenter dm-line-height-hack\">\n							<svg>\n								<polygon class=\""
    + alias4(((helper = (helper = lookupProperty(helpers,"rtriangle") || (depth0 != null ? lookupProperty(depth0,"rtriangle") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"rtriangle","hash":{},"data":data,"loc":{"start":{"line":55,"column":24},"end":{"line":55,"column":37}}}) : helper)))
    + "\" points=\"0 0 16 25 0 50\">\n							</svg>\n						</div>\n					</div>\n				</div>\n			</div>\n";
},"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div id=\"dm-home-visible\"></div>\n"
    + ((stack1 = lookupProperty(helpers,"each").call(depth0 != null ? depth0 : (container.nullContext || {}),(depth0 != null ? lookupProperty(depth0,"activities") : depth0),{"name":"each","hash":{},"fn":container.program(1, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":2,"column":0},"end":{"line":98,"column":9}}})) != null ? stack1 : "");
},"useData":true});
templates['login'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div class=\"dm-panel\">\n	<p class=\"dm-tag-line\">\n		Delight-Meditate is a free, simple, quick, frictionless way to follow\n		a daily Bible reading plan.\n	</p>\n	<button type=\"button\" class=\"btn btn-primary\" id=\"dm-about-button\">About</button>\n	<button type=\"button\" class=\"btn btn-primary dm-button-inline\" id=\"dm-test-drive-button\">Test drive</button>\n	<div class=\"dm-vspacer\"></div>\n</div>\n<div class=\"dm-vspacer\"></div>\n<div class=\"dm-panel-fixed\" id=\"dm-login-panel\">\n	<h1>Login</h1>\n	<div class=\"dm-vspacer\"></div>\n	<form class=\"add-border\">\n		<div class=\"form-group\">\n			<label for=\"lfEmail\">email address</label>\n			<input id=\"lfEmail\" class=\"form-control\" \"type=\"text\" name=\"email987\" autocomplete=\"on\" />\n		</div>\n		<div class=\"form-group\">\n			<label for=\"lfPass\">password</label>\n			<input id=\"lfPass\" class=\"form-control\" type=\"password\" name=\"pass987\" />\n		</div>\n		<div class=\"form-check\">\n			<label class=\"form-check-label\">\n				<input id=\"lfKeep\" class=\"form-check-input\" type=\"checkbox\" value=\"\" checked>\n				Keep me logged in\n			</label>\n		</div>\n		<div class=\"dm-vspacer\"></div>\n		<div id=\"login-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-login-button\">Log in</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n	<button type=\"button\" class=\"btn btn-primary\" id=\"dm-signup-button\">Sign up</button>\n	<div class=\"dm-vspacer\"></div>\n	<a href=\"/c/forgot\">Forgot your password?</a>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['modalalert'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var helper, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"modal fade\">\n	<div class=\"modal-dialog\" role=\"document\">\n		<div class=\"modal-content\">\n			<div class=\"modal-body\">\n				<p>"
    + container.escapeExpression(((helper = (helper = lookupProperty(helpers,"msg") || (depth0 != null ? lookupProperty(depth0,"msg") : depth0)) != null ? helper : container.hooks.helperMissing),(typeof helper === "function" ? helper.call(depth0 != null ? depth0 : (container.nullContext || {}),{"name":"msg","hash":{},"data":data,"loc":{"start":{"line":5,"column":7},"end":{"line":5,"column":14}}}) : helper)))
    + "</p>\n			</div>\n			<div class=\"modal-footer\">\n				<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\">OK</button>\n			</div>\n		</div>\n	</div>\n</div>\n";
},"useData":true});
templates['modalconfirm'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var helper, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"modal fade\">\n	<div class=\"modal-dialog\" role=\"document\">\n		<div class=\"modal-content\">\n			<div class=\"modal-body\">\n				<p>"
    + container.escapeExpression(((helper = (helper = lookupProperty(helpers,"msg") || (depth0 != null ? lookupProperty(depth0,"msg") : depth0)) != null ? helper : container.hooks.helperMissing),(typeof helper === "function" ? helper.call(depth0 != null ? depth0 : (container.nullContext || {}),{"name":"msg","hash":{},"data":data,"loc":{"start":{"line":5,"column":7},"end":{"line":5,"column":14}}}) : helper)))
    + "</p>\n			</div>\n			<div class=\"modal-footer\">\n				<button type=\"button\" class=\"btn btn-secondary\" data-dismiss=\"modal\" id=\"dm-modal-cancel\">cancel</button>\n				<button type=\"button\" class=\"btn btn-primary\" data-dismiss=\"modal\" id=\"dm-modal-ok\">OK</button>\n			</div>\n		</div>\n	</div>\n</div>\n";
},"useData":true});
templates['navbar'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<nav class=\"navbar navbar-expand-md navbar-light\">\n	<a class=\"navbar-brand\" href=\"/\">Delight/Meditate</a>\n	<button class=\"navbar-toggler\" type=\"button\" data-toggle=\"collapse\" data-target=\"#dmNavbar\" aria-controls=\"dmNavbar\" aria-expanded=\"false\" aria-label=\"Toggle navigation\">\n		<span class=\"navbar-toggler-icon\"></span>\n	</button>\n	<div class=\"collapse navbar-collapse\" id=\"dmNavbar\">\n		<ul class=\"navbar-nav mr-auto mt-2 mt-lg-0\">\n			<li class=\"nav-item\">\n				<a class=\"nav-link collapser\" href=\"/c/about\">About</a>\n			</li>\n			<li class=\"nav-item\">\n				<a class=\"nav-link collapser\" href=\"/c/plans\">Plans</a>\n			</li>\n			<li class=\"nav-item\">\n				<a class=\"nav-link collapser\" href=\"/c/help\">Help</a>\n			</li>\n			<li id=\"nav-contact\" class=\"nav-item\">\n				<a class=\"nav-link collapser\" href=\"/c/contact\">Contact</a>\n			</li>\n			<li id=\"nav-login\" class=\"nav-item\">\n				<a class=\"nav-link collapser\" href=\"/c/login\">Login/Signup</a>\n			</li>\n			<li id=\"nav-user-dropdown\" class=\"nav-item dropdown dm-gone\">\n				<a class=\"nav-link dropdown-toggle\" href=\"#\" role=\"button\" id=\"user-dropdown\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">User</a>\n				<div class=\"dropdown-menu\" aria-labelledby=\"user-dropdown\">\n					<a class=\"dropdown-item\" href=\"/c/profile\">Profile</a>\n					<a class=\"dropdown-item\" id=\"dm-menu-logout\">Logout</a>\n				</div>\n			</li>\n		</ul>\n	</div>\n</nav>\n<div class=\"dm-vspacer\"></div>\n";
},"useData":true});
templates['partialsignup'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"form-group\">\n	<label for=\"sfFirst\">first name</label>\n	<input id=\"sfFirst\" class=\"form-control\" \"type=\"text\" value=\""
    + alias4(((helper = (helper = lookupProperty(helpers,"firstname") || (depth0 != null ? lookupProperty(depth0,"firstname") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"firstname","hash":{},"data":data,"loc":{"start":{"line":3,"column":62},"end":{"line":3,"column":75}}}) : helper)))
    + "\" name=\"fname123\" />\n</div>\n<div class=\"form-group\">\n	<label for=\"sfLast\">last name</label>\n	<input id=\"sfLast\" class=\"form-control\" \"type=\"text\" value=\""
    + alias4(((helper = (helper = lookupProperty(helpers,"lastname") || (depth0 != null ? lookupProperty(depth0,"lastname") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"lastname","hash":{},"data":data,"loc":{"start":{"line":7,"column":61},"end":{"line":7,"column":73}}}) : helper)))
    + "\" name=\"lname123\" />\n</div>\n<div class=\"form-group\">\n	<label for=\"sfEmail\">email address</label>\n	<input id=\"sfEmail\" class=\"form-control\" \"type=\"text\" value=\""
    + alias4(((helper = (helper = lookupProperty(helpers,"email") || (depth0 != null ? lookupProperty(depth0,"email") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"email","hash":{},"data":data,"loc":{"start":{"line":11,"column":62},"end":{"line":11,"column":71}}}) : helper)))
    + "\" name=\"email123\" />\n</div>\n";
},"useData":true});
templates['partialversion'] = template({"1":function(container,depth0,helpers,partials,data) {
    var helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "			<option "
    + alias4(((helper = (helper = lookupProperty(helpers,"selected") || (depth0 != null ? lookupProperty(depth0,"selected") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"selected","hash":{},"data":data,"loc":{"start":{"line":5,"column":11},"end":{"line":5,"column":23}}}) : helper)))
    + ">"
    + alias4(((helper = (helper = lookupProperty(helpers,"provider") || (depth0 != null ? lookupProperty(depth0,"provider") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"provider","hash":{},"data":data,"loc":{"start":{"line":5,"column":24},"end":{"line":5,"column":36}}}) : helper)))
    + "</option>\n";
},"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"form-group\">\n    <label for=\"bvProvider\">Bible provider</label>\n    <select class=\"form-control\" id=\"bvProvider\">\n"
    + ((stack1 = lookupProperty(helpers,"each").call(alias1,(depth0 != null ? lookupProperty(depth0,"choices") : depth0),{"name":"each","hash":{},"fn":container.program(1, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":4,"column":2},"end":{"line":6,"column":11}}})) != null ? stack1 : "")
    + "    </select>\n</div>\n<div class=\"form-group\">\n	<label for=\"bvVersion\">Bible version code</label>\n	<input id=\"bvVersion\" class=\"form-control\" type=\"text\" value=\""
    + container.escapeExpression(((helper = (helper = lookupProperty(helpers,"version") || (depth0 != null ? lookupProperty(depth0,"version") : depth0)) != null ? helper : container.hooks.helperMissing),(typeof helper === "function" ? helper.call(alias1,{"name":"version","hash":{},"data":data,"loc":{"start":{"line":11,"column":63},"end":{"line":11,"column":74}}}) : helper)))
    + "\" />\n</div>\n";
},"useData":true});
templates['planday'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel\" id=\"dm-planday-panel\">\n	<div class=\"dm-title\">Jump to a Specific Day</div>\n	<p>\n		You are currently on day "
    + alias4(((helper = (helper = lookupProperty(helpers,"day") || (depth0 != null ? lookupProperty(depth0,"day") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"day","hash":{},"data":data,"loc":{"start":{"line":4,"column":27},"end":{"line":4,"column":34}}}) : helper)))
    + " of "
    + alias4(((helper = (helper = lookupProperty(helpers,"totaldays") || (depth0 != null ? lookupProperty(depth0,"totaldays") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"totaldays","hash":{},"data":data,"loc":{"start":{"line":4,"column":38},"end":{"line":4,"column":51}}}) : helper)))
    + " for plan <b>"
    + alias4(((helper = (helper = lookupProperty(helpers,"plantitle") || (depth0 != null ? lookupProperty(depth0,"plantitle") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plantitle","hash":{},"data":data,"loc":{"start":{"line":4,"column":64},"end":{"line":4,"column":77}}}) : helper)))
    + ".</b>\n	</p>\n	<form class=\"add-border dm-number-form\">\n		<div class=\"form-group\">\n			<label for=\"pdfDay\">new day number</label>\n			<input id=\"pdfDay\" class=\"form-control\" \"type=\"number\" />\n		</div>\n		<div id=\"planday-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary dm-planday\" id=\"dm-planday-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":12,"column":74},"end":{"line":12,"column":82}}}) : helper)))
    + "\">Jump to this day</button>\n		<div class=\"dm-vspacer\"></div>\n		<button type=\"button\" class=\"btn btn-primary dm-plandet\" id=\"dm-plandet-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":14,"column":74},"end":{"line":14,"column":82}}}) : helper)))
    + "\">See plan details</button>\n		<div class=\"dm-vspacer\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-planday-cancel\">Cancel</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['plandetails'] = template({"1":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "		<hr/>\n		<span class=\"dm-details-day\">"
    + container.escapeExpression(((helper = (helper = lookupProperty(helpers,"day") || (depth0 != null ? lookupProperty(depth0,"day") : depth0)) != null ? helper : container.hooks.helperMissing),(typeof helper === "function" ? helper.call(alias1,{"name":"day","hash":{},"data":data,"loc":{"start":{"line":6,"column":31},"end":{"line":6,"column":38}}}) : helper)))
    + "</span>\n"
    + ((stack1 = lookupProperty(helpers,"each").call(alias1,(depth0 != null ? lookupProperty(depth0,"readings") : depth0),{"name":"each","hash":{},"fn":container.program(2, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":7,"column":2},"end":{"line":9,"column":11}}})) != null ? stack1 : "");
},"2":function(container,depth0,helpers,partials,data) {
    return "			- "
    + container.escapeExpression(container.lambda(depth0, depth0))
    + "\n";
},"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel\">\n	<div class=\"dm-title\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"title") || (depth0 != null ? lookupProperty(depth0,"title") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"title","hash":{},"data":data,"loc":{"start":{"line":2,"column":23},"end":{"line":2,"column":32}}}) : helper)))
    + "</div>\n	"
    + alias4(((helper = (helper = lookupProperty(helpers,"days") || (depth0 != null ? lookupProperty(depth0,"days") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"days","hash":{},"data":data,"loc":{"start":{"line":3,"column":1},"end":{"line":3,"column":9}}}) : helper)))
    + " days\n"
    + ((stack1 = lookupProperty(helpers,"each").call(alias1,(depth0 != null ? lookupProperty(depth0,"rows") : depth0),{"name":"each","hash":{},"fn":container.program(1, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":4,"column":1},"end":{"line":10,"column":10}}})) != null ? stack1 : "")
    + "	<div class=\"dm-vspacer\"></div>\n	<button type=\"button\" class=\"btn btn-primary\" id=\"dm-plandetails-ok\">OK</button>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['plans'] = template({"1":function(container,depth0,helpers,partials,data) {
    var helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "		<div class=\"dm-vspacer\"></div>\n		<div class=\"dm-small-title\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"title") || (depth0 != null ? lookupProperty(depth0,"title") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"title","hash":{},"data":data,"loc":{"start":{"line":5,"column":30},"end":{"line":5,"column":39}}}) : helper)))
    + "</div>\n		"
    + alias4(((helper = (helper = lookupProperty(helpers,"desc") || (depth0 != null ? lookupProperty(depth0,"desc") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"desc","hash":{},"data":data,"loc":{"start":{"line":6,"column":2},"end":{"line":6,"column":10}}}) : helper)))
    + "\n		<div class=\"dm-vspacer\"></div>\n		<div class=\"dm-plans-button-container\">\n			<button type=\"button\" class=\"btn btn-primary dm-plan-details-button\" id=\"dm-pdb-"
    + alias4(((helper = (helper = lookupProperty(helpers,"name") || (depth0 != null ? lookupProperty(depth0,"name") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"name","hash":{},"data":data,"loc":{"start":{"line":9,"column":83},"end":{"line":9,"column":91}}}) : helper)))
    + "\">See plan details</button>\n		</div>\n		<div id=\"pb-"
    + alias4(((helper = (helper = lookupProperty(helpers,"name") || (depth0 != null ? lookupProperty(depth0,"name") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"name","hash":{},"data":data,"loc":{"start":{"line":11,"column":14},"end":{"line":11,"column":22}}}) : helper)))
    + "\" class=\"dm-plans-button-container dm-plan-add-button\">\n			<button type=\"button\" class=\"btn btn-primary\" id=\"dm-addplan-"
    + alias4(((helper = (helper = lookupProperty(helpers,"name") || (depth0 != null ? lookupProperty(depth0,"name") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"name","hash":{},"data":data,"loc":{"start":{"line":12,"column":64},"end":{"line":12,"column":72}}}) : helper)))
    + "\">Add this plan</button>\n		</div>\n		<div class=\"dm-plans-button-container dm-plan-login-button dm-gone\">\n			<button type=\"button\" class=\"btn btn-primary\">Log in to add this plan</button>\n		</div>\n		<div class=\"dm-vspacer\"></div>\n		<hr/>\n";
},"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel\">\n	<div class=\"dm-title\">Available Plans</div>\n"
    + ((stack1 = lookupProperty(helpers,"each").call(depth0 != null ? depth0 : (container.nullContext || {}),(depth0 != null ? lookupProperty(depth0,"plandescs") : depth0),{"name":"each","hash":{},"fn":container.program(1, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":3,"column":1},"end":{"line":19,"column":10}}})) != null ? stack1 : "")
    + "	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['profile'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel\" id=\"dm-profile-panel\">\n	<div class=\"dm-title\">User Profile</div>\n	<p>\n		The Bible provider and version code will apply to any plans you add in the future.\n		For help choosing them, choose \"Help\" from the menu above.\n	</p>\n	<form class=\"add-border dm-number-form\">\n		<div class=\"dm-vspacer\"></div>\n		<div class=\"form-group\">\n			<label for=\"oldPass\">current password (required)</label>\n			<input id=\"oldPass\" class=\"form-control\" type=\"password\" name=\"pass789\" />\n		</div>\n		<div class=\"dm-vspacer\"></div>\n		<div class=\"form-group\">\n			<label for=\"newPass\">new password (optional)</label>\n			<input id=\"newPass\" class=\"form-control\" type=\"password\" name=\"pass456\" />\n		</div>\n		<div class=\"dm-vspacer\"></div>\n"
    + ((stack1 = container.invokePartial(lookupProperty(partials,"signupPartial"),depth0,{"name":"signupPartial","data":data,"indent":"\t\t","helpers":helpers,"partials":partials,"decorators":container.decorators})) != null ? stack1 : "")
    + ((stack1 = container.invokePartial(lookupProperty(partials,"versionPartial"),depth0,{"name":"versionPartial","data":data,"indent":"\t\t","helpers":helpers,"partials":partials,"decorators":container.decorators})) != null ? stack1 : "")
    + "		<div id=\"profile-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-user-profile\">Update</button>\n		<div class=\"dm-vspacer\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-profile-cancel\">Cancel</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"usePartial":true,"useData":true});
templates['pwreset'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div class=\"dm-panel-fixed\">\n	<h1>New Password</h1>\n	<p>\n		Please enter a new password for your account.\n	</p>\n	<form class=\"add-border\">\n		<div class=\"form-group\">\n			<label for=\"prPass\">password</label>\n			<input id=\"prPass\" class=\"form-control\" type=\"password\" name=\"pwreset123\" />\n		</div>\n		<input id=\"prToken\" type=\"hidden\" value=\""
    + container.escapeExpression(container.lambda(depth0, depth0))
    + "\" name=\"prhid123\"/>\n		<div class=\"dm-vspacer\"></div>\n		<div id=\"pwreset-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-button-pwreset\">Submit</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"useData":true});
templates['signup'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel-fixed\"  id=\"dm-signup-panel\">\n	<h1>Sign Up</h1>\n	<div class=\"dm-vspacer\"></div>\n	<form class=\"add-border\">\n"
    + ((stack1 = container.invokePartial(lookupProperty(partials,"signupPartial"),depth0,{"name":"signupPartial","data":data,"indent":"\t\t","helpers":helpers,"partials":partials,"decorators":container.decorators})) != null ? stack1 : "")
    + "		<div class=\"form-group\">\n			<label for=\"sfPass\">password</label>\n			<input id=\"sfPass\" class=\"form-control\" type=\"password\" name=\"pass123\" />\n		</div>\n		<div id=\"signup-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-signup-button\">Sign up</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n	<button type=\"button\" class=\"btn btn-primary\" id=\"dm-login-button\">Log in</button>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"usePartial":true,"useData":true});
templates['version'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div class=\"dm-panel\" id=\"dm-version-panel\">\n	<div class=\"dm-title\">Bible Version</div>\n	<p>\n		Choose a provider and version code for plan <b>"
    + alias4(((helper = (helper = lookupProperty(helpers,"plantitle") || (depth0 != null ? lookupProperty(depth0,"plantitle") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plantitle","hash":{},"data":data,"loc":{"start":{"line":4,"column":49},"end":{"line":4,"column":62}}}) : helper)))
    + ".</b>\n		The default choices are provider <b>BibleGateway.com-PRINT</b> and version code <b>ESV</b>.\n		For more information, choose \"Help\" from the menu above.\n	</p>\n	<form class=\"add-border dm-number-form\">\n"
    + ((stack1 = container.invokePartial(lookupProperty(partials,"versionPartial"),depth0,{"name":"versionPartial","data":data,"indent":"\t\t","helpers":helpers,"partials":partials,"decorators":container.decorators})) != null ? stack1 : "")
    + "		<div id=\"version-alert\" class=\"alert alert-danger\"></div>\n		<button type=\"button\" class=\"btn btn-primary dm-newver\" id=\"dm-newver-"
    + alias4(((helper = (helper = lookupProperty(helpers,"plan") || (depth0 != null ? lookupProperty(depth0,"plan") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"plan","hash":{},"data":data,"loc":{"start":{"line":11,"column":72},"end":{"line":11,"column":80}}}) : helper)))
    + "\">Update</button>\n		<div class=\"dm-vspacer\"></div>\n		<button type=\"button\" class=\"btn btn-primary\" id=\"dm-newver-cancel\">Cancel</button>\n	</form>\n	<div class=\"dm-vspacer\"></div>\n	<div class=\"dm-vspacer\"></div>\n</div>\n";
},"usePartial":true,"useData":true});
})();