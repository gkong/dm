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

if [ "$DM_TEST_SECRETS_TOML" = "" -o ! -f "$DM_TEST_SECRETS_TOML" -o ! -r "$DM_TEST_SECRETS_TOML" ] ; then
	echo need environment variable DM_TEST_SECRETS_TOML - path to a secrets file similar to config/example-secrets.toml
	exit
else
	echo using secrets file "$DM_TEST_SECRETS_TOML"
fi

targs="-P --transform=s|$DM_TEST_SECRETS_TOML|config/secret.toml|"

exe=dm
target=dmtest
destdir=./artifacts

rm $exe
CGO_ENABLED=0 go build
mv $exe $target
tar $targs -c $target static document index.html dm-admin.html config/base.toml config/test.toml "$DM_TEST_SECRETS_TOML" | gzip > "$destdir"/"$target".tgz
mv $target $exe
