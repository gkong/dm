#!/bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# this script must be in the project root directory and must be issued from there.
# set basedir to the full pathname of the directory in which this script resides.
mustdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
curdir=`pwd`
if [ "$curdir" != "$mustdir" ] ; then
	echo current directory must be "$mustdir"
	exit
fi

exe=dm
target=dmprod
destdir=./artifacts

rm $exe

CGO_ENABLED=0 go build

mv $exe $target
tar -c $target static document index.html dm-admin.html config/base.toml config/secret.toml config/prod.toml | gzip > "$destdir"/"$target".tgz
mv $target $exe
