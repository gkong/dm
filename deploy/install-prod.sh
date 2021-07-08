#! /bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# run this on production server, to install/upgrade the dm production instance.

# assumes dmprod.service has been installed on production server.


prog=dmprod

dest=/dm/prod/app


euid="`id -u`"

if [ "$euid" != "0" ]
	then echo "\nmust be root\n"
	exit 1
fi

systemctl stop $prog

rm -rf $dest
mkdir -p $dest
cd $dest
tar xf /dm/incoming/$prog.tgz
chown -R core:core $dest

systemctl start $prog

sleep 2

systemctl status $prog
