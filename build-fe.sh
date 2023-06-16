#! /bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# i'd prefer make, but i'm not confident i'd get the js dependencies right,
# so just barrel ahead and build everything every time.

# for the sake of caching, it's important not to change file modified times.
# to accomplish this, we write final minified versions to a temp directory,
# and only copy to dest dir if contents differ.

# change version numbers here, AND IN index.html, to force clients to reload
dm=dm-v25
dmdeps=dmdeps-v7

# this script must be in the parent to dirs "js" and "static".
# set basedir to the full pathname of the directory in which this script resides.
basedir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

jsdir=$basedir/js
dest=$basedir/static
tmpdir=`mktemp -d`

( cd $jsdir; browserify $jsdir/dmadmin.js -o $dest/dmadmin.js )

# dummy allows us to factor out all the dependencies into a separate bundle
( cd $jsdir; browserify $jsdir/dm.js $jsdir/dummy.js -p [ factor-bundle -o $dest/dm.js -o /tmp/dummy.js ] -o $dest/dmdeps.js )

terser $dest/dm.js --lint --compress warnings=false --mangle --comments --output $tmpdir/"$dm".min.js
terser $dest/dmdeps.js --compress warnings=false --mangle --comments --output $tmpdir/"$dmdeps".min.js

if cmp -s $tmpdir/"$dm".min.js $dest/"$dm".min.js ; then
	# they're the same; do nothing
	echo > /dev/null
else
	mv $tmpdir/"$dm".min.js $dest/"$dm".min.js
fi

if cmp -s $tmpdir/"$dmdeps".min.js $dest/"$dmdeps".min.js ; then
	# they're the same; do nothing
	echo > /dev/null
else
	mv $tmpdir/"$dmdeps".min.js $dest/"$dmdeps".min.js
fi

rm -rf $tmpdir

echo base JavaScript bundles built
echo invoking make to build template JavaScript bundle and CSS bundle

# now run make, to make the stuff it's responsible for

make
