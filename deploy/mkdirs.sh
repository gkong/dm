#! /bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# run this on production server, to make directories

destRoot=/dm

if [ ! -d "$destRoot"  -o  ! -w "$destRoot" ] ; then
	echo "$destRoot" must exist and be writable
	exit 1
fi

cd $destRoot

mkdir -m 755 active bin incoming prod shared shared/cert test tmp

( cd prod ; mkdir -m 755 db log bck bcklatest )

( cd test ; mkdir -m 755 db log bck bcklatest )

# app goes into /dm/prod/app and /dm/test/app, which are immutable
# and are deleted and reconstituted on each update.

