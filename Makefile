
# figure out a way to purge obsolete versions but leave behind desired versions:
# rm $(static)/dm*.min.js $(static)/dm*.min.css

# leave non-versioned, non-minified versions in static, for the benefit of dm-admin.html

# change version numbers here, AND IN index.html, to force clients to reload
dmt = dmt-v18
css = dm-v11

# this Makefile must reside in the parent of "template", "js", "css", and "static".
# set basedir to the full pathname of the directory in which this Makefile resides.
basedir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

tdir = $(basedir)/template
jsdir = $(basedir)/js
cssdir = $(basedir)/css
static = $(basedir)/static

all: \
	$(static)/dm.css \
	$(static)/$(css).min.css \
	$(static)/$(dmt).min.js

#	$(static)/$(dm).min.js
#	$(static)/$(dmdeps).min.js
#	$(static)/dmadmin.js

$(static)/dm.css: $(cssdir)/dm.css
	cp $< $@

$(static)/$(css).min.css: $(cssdir)/dm.css
	npx postcss $< > $@

$(static)/dmt.js: $(tdir)/*.handlebars
	handlebars $(tdir) -f $@ -k each -k if -k unless

$(static)/$(dmt).min.js: $(static)/dmt.js
	terser $< --lint --compress warnings=false --mangle --comments --output $@
