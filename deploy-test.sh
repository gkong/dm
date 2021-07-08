#!/bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# assumes dmtest.service has been installed on production server.

# this script must be in the project root directory and must be issued from there.
# set basedir to the full pathname of the directory in which this script resides.
mustdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
curdir=`pwd`
if [ "$curdir" != "$mustdir" ] ; then
	echo current directory must be "$mustdir"
	exit
fi

rcp deploy/install-test.sh artifacts/dmtest.tgz dm:/dm/incoming/
echo running install script...
time ssh dm sudo /dm/incoming/install-test.sh
