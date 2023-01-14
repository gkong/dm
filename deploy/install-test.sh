#! /bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# run this on TEST server, to install/upgrade the dm TEST instance.

# assumes dmtest.service has been installed on TEST server.


prog=dmtest

dest=/dm/test/app


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
chown -R ubuntu:ubuntu $dest

systemctl start $prog

sleep 2

systemctl status $prog
