#! /bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

mkdir -p 0

if [ ! -d 0 ] ; then
	echo ERROR - could not make directory 0
	exit 1
fi


if [ -f 0/schemaBoil.go ] ; then
	echo ERROR - found 0/schemaBoil.go
	exit 1
fi

if [ -f 0/sessdataBoil.go ] ; then
	echo ERROR - found 0/sessdataBoil.go
	exit 1
fi

mv schemaBoil.go sessdataBoil.go 0

rm -f *_generated.go *_ffjson_expose.go

go generate

mv 0/schemaBoil.go 0/sessdataBoil.go .




