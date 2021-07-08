// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package gen

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

const bytesPerInt64 = 8

// convert int64 to bytes and copy into a given byte slice
func itobPut(dest []byte, i int64) {
	binary.LittleEndian.PutUint64(dest, uint64(i))
}

// convert int64 to []byte and return it in a new byte slice
func itob(i int64) []byte {
	b := make([]byte, bytesPerInt64)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// convert []byte to int64
func btoi(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

func randomBytes(size int) []byte {
	b := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("qsldb.randomBytes - cannot read rand.Reader")
	}
	return b
}
