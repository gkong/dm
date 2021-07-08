// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// manage documents that are few, small and (mostly) immutable.
//
// there is a single set of subdirectories or buckets.
// documents are initially read from files during program initialization.
//
// a "time to live" is maintained for each document.
// this TTL is NOT how long it lives here, but how long users may cache it.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type doc struct {
	contents []byte
	ttlsecs  int
}

var docs = make(map[string]map[string]doc) // keys: bucketname / filename

var docMutex sync.RWMutex

// return a list of filenames in the given bucket.
//
// XXX - makes an array from docs map and sorts it every time. optimize.
func docNames(bucket string) []string {
	docMutex.RLock()
	defer docMutex.RUnlock()

	var names []string
	for name := range docs[bucket] {
		names = append(names, name)
	}

	sort.Sort(stringSlice(names))
	return names
}

func docGet(bucket string, filename string) (contents []byte, ttlsecs int, err error) {
	docMutex.RLock()
	defer docMutex.RUnlock()

	d, ok := docs[bucket][filename]
	if !ok {
		return []byte{}, 0, wrapErr{"not found", err}
	}
	return d.contents, d.ttlsecs, nil
}

// if document already exists, it is replaced.
func docPut(bucket string, filename string, ttlsecs int, data []byte) error {
	docMutex.Lock()
	defer docMutex.Unlock()

	if _, ok := docs[bucket]; !ok {
		docs[bucket] = make(map[string]doc)
	}

	docs[bucket][filename] = doc{data, ttlsecs}
	return nil
}

// all files get the same value for ttlsecs.
func docReadDir(docdir string, ttlsecs int) error {
	parentDir, err := os.Open(docdir)
	if err != nil {
		return wrapErr{"docReadDir - parentDir Open", err}
	}

	subDirNames, err := parentDir.Readdirnames(-1)
	if err != nil {
		return wrapErr{"docReadDir - parentDir.Readdirnames", err}
	}

	for _, dirName := range subDirNames {
		dirPath := filepath.Join(docdir, dirName)
		dir, err := os.Open(dirPath)
		if err != nil {
			return wrapErr{"docReadDir - subdir Open", err}
		}

		fileDirNames, err := dir.Readdirnames(-1)
		if err != nil {
			return wrapErr{"docReadDir - dir.Readdirnames", err}
		}

		for _, fileName := range fileDirNames {
			filePath := filepath.Join(dirPath, fileName)
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return wrapErr{"docReadDir - file read", err}
			}

			if err := docPut(dirName, fileName, ttlsecs, data); err != nil {
				return wrapErr{"docReadDir - docPut", err}
			}
		}
	}

	return nil
}

// for case-independent sorting
type stringSlice []string

func (s stringSlice) Len() int      { return len(s) }
func (s stringSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s stringSlice) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}
