// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

// very simple application-specific database abstraction layer

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	. "github.com/gkong/dm/gen"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func dbHas(key []byte) (bool, error) {
	has, err := gldb.Has(key, nil)
	if err != nil {
		return false, wrapErr{"dbHas - Has", err}
	}
	return has, nil
}

func dbDelete(key []byte) error {
	if err := gldb.Delete(key, nil); err != nil {
		return wrapErr{"dbDelete - Delete", err}
	}
	return nil
}

func dbUserPut(key []byte, u *User) error {
	ubytes, err := u.Marshal()
	if err != nil {
		return wrapErr{"dbUserPut - Marshal", err}
	}
	if err := gldb.Put(key, ubytes, nil); err != nil {
		return wrapErr{"dbUserPut - Put", err}
	}
	return nil
}

func dbUserGet(key []byte) (*User, error) {
	data, err := gldb.Get(key, nil)
	if err != nil {
		return nil, wrapErr{"dbUserGet - Get", err}
	}
	u := &User{}
	if err := u.Unmarshal(data); err != nil {
		return nil, wrapErr{"dbUserGet - Unmarshal", err}
	}
	return u, nil
}

func dbActivityPut(key []byte, a *Activity) error {
	abytes, err := a.Marshal()
	if err != nil {
		return wrapErr{"dbActivityPut - Marshal", err}
	}
	if err := gldb.Put(key, abytes, nil); err != nil {
		return wrapErr{"dbActivityPut - Put", err}
	}
	return nil
}

func dbActivityGet(key []byte) (*Activity, error) {
	data, err := gldb.Get(key, nil)
	if err != nil {
		return nil, wrapErr{"dbActivityGet - Get", err}
	}
	a := &Activity{}
	if err := a.Unmarshal(data); err != nil {
		return nil, wrapErr{"dbActivityGet - Unmarshal", err}
	}
	return a, nil
}

func dbUserActivities(userkey []byte) ([]*Activity, error) {
	key := ActivityKey(userkey, []byte{})
	activities := []*Activity{}
	iter := gldb.NewIterator(nil, nil)
	for ok := iter.Seek(key); ok; ok = iter.Next() {
		if !bytes.HasPrefix(iter.Key(), key) {
			break
		}
		a := &Activity{}
		if err := a.Unmarshal(iter.Value()); err != nil {
			return nil, wrapErr{"dbUserActivities - Unmarshal", err}
		}
		activities = append(activities, a)
	}
	return activities, nil
}

// called on every Session.Save to update LastVisit table
func dbRecordVisit(key []byte, t time.Time) error {
	// the in-place key conversions are NOT goroutine-safe!
	ConvertUKeyToLVKey(key)
	err := gldb.Put(key, u64tob(uint64(t.Unix())), nil)
	ConvertLVKeyToUKey(key)
	if err != nil {
		return wrapErr{"dbRecordVisit - Put", err}
	}
	return nil
}

// backup and restore

const secsPerDay = 86400

// make a backup at a specified time every day.
// maintains a directory which contains only the latest daily backup.
// dailyTime is seconds after midnight UTC.
func dailyBackupSetup(bckDir string, latestDir string, dailyTime int, ilog, elog io.Writer) {
	go func() {
		for {
			// sleep until daily backup time

			now := time.Now().UTC()
			nowSecs := (((now.Hour() * 60) + now.Minute()) * 60) + now.Second()
			waitSecs := dailyTime - nowSecs
			if waitSecs <= 0 {
				waitSecs += secsPerDay
			}
			time.Sleep(time.Duration(waitSecs) * time.Second)

			// make the backup

			fname := time.Now().Format("2006-01-02")
			fullname := filepath.Join(bckDir, fname)
			n, err := dbBackup(fullname)
			if err != nil {
				io.WriteString(elog, "daily backup failed - "+err.Error())
				continue
			}
			io.WriteString(ilog, fmt.Sprintf("daily backup - saved %d records to %s", n, fullname))

			// clean out latest backup dir

			ldir, err := os.Open(latestDir)
			if err != nil {
				io.WriteString(elog, "daily backup - Open(latestDir) - "+err.Error())
				continue
			}
			filenames, err := ldir.Readdirnames(-1)
			if err != nil {
				io.WriteString(elog, "daily backup - Readdirnames - "+err.Error())
				continue
			}
			for _, f := range filenames {
				if err = os.Remove(filepath.Join(latestDir, f)); err != nil {
					io.WriteString(elog, "daily backup - Remove - "+err.Error())
					continue
				}
			}

			// link the backup we just made into latest backup dir

			if err = os.Link(fullname, filepath.Join(latestDir, fname)); err != nil {
				io.WriteString(elog, "daily backup - Link - "+err.Error())
				continue
			}
		}
	}()
}

// back up long-term, high-value data (User and Activity tables).
// do NOT back up any of the session stores.
func dbBackup(filename string) (int, error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return 0, wrapErr{"dbBackup - OpenFile", err}
	}
	defer f.Close()

	snap, err := gldb.GetSnapshot()
	if err != nil {
		return 0, wrapErr{"dbBackup - GetSnapshot", err}
	}
	defer snap.Release()

	nuser, err := dbBackupTable(snap, PfxUser, f)
	if err != nil {
		return 0, wrapErr{"dbBackup - User table", err}
	}

	nact, err := dbBackupTable(snap, PfxActivity, f)
	if err != nil {
		return 0, wrapErr{"dbBackup - Activity table", err}
	}

	nlv, err := dbBackupTable(snap, PfxLastVisit, f)
	if err != nil {
		return 0, wrapErr{"dbBackup - LastVisit table", err}
	}

	return nuser + nact + nlv, nil
}

// emit binary data of the form: keylen | key | valuelen | value
func dbBackupTable(snap *leveldb.Snapshot, prefix []byte, w io.Writer) (int, error) {
	iter := snap.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	count := 0
	for iter.Next() {
		n, err := w.Write(u32tob(uint32(len(iter.Key()))))
		if n != bytesPerUint32 || err != nil {
			return 0, wrapErr{"dbBackupTable - Write keylen", err}
		}

		n, err = w.Write(iter.Key())
		if n != len(iter.Key()) || err != nil {
			return 0, wrapErr{"dbBackupTable - Write key", err}
		}

		n, err = w.Write(u32tob(uint32(len(iter.Value()))))
		if n != bytesPerUint32 || err != nil {
			return 0, wrapErr{"dbBackupTable - Write keylen", err}
		}

		n, err = w.Write(iter.Value())
		if n != len(iter.Value()) || err != nil {
			return 0, wrapErr{"dbBackupTable - Write key", err}
		}

		count++
	}

	return count, nil
}

// delete all contents of db, then restore from the given file.
// admin user must be logged in, so can't start from a completely empty db.
//
// restore doesn't need to know which tables are in the backup,
// because all data looks the same to it, tables only differing by key prefix.
func dbRestore(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, wrapErr{"dbRestore - OpenFile", err}
	}
	defer f.Close()

	dbNuke()

	keylenbuf := make([]byte, bytesPerUint32)
	keybuf := make([]byte, 10000)
	vallenbuf := make([]byte, bytesPerUint32)
	valbuf := make([]byte, 10000)
	count := 0
	for {
		n, err := f.Read(keylenbuf)
		if n == 0 && err == io.EOF {
			break
		}
		if n != bytesPerUint32 || err != nil {
			return 0, wrapErr{"dbRestore - short read - keylen", err}
		}
		keylen := btou32(keylenbuf)

		n, err = f.Read(keybuf[:keylen])
		if n != int(keylen) || err != nil {
			return 0, wrapErr{"dbRestore - short read - key", err}
		}

		n, err = f.Read(vallenbuf)
		if n != bytesPerUint32 || err != nil {
			return 0, wrapErr{"dbRestore - short read - vallen", err}
		}
		vallen := btou32(vallenbuf)

		n, err = f.Read(valbuf[:vallen])
		if n != int(vallen) || err != nil {
			return 0, wrapErr{"dbRestore - short read - val", err}
		}

		if err := gldb.Put(keybuf[:keylen], valbuf[:vallen], nil); err != nil {
			return 0, wrapErr{"dbRestore - Put", err}
		}

		count++
	}

	return count, nil
}

// admin and debug functions

// delete all the records for a given prefix
func dbDropTable(prefix []byte) {
	iter := gldb.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	for iter.Next() {
		gldb.Delete(iter.Key(), nil)
	}
}

// delete EVERYTHING!
func dbNuke() {
	iter := gldb.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		gldb.Delete(iter.Key(), nil)
	}
}

// change all records with a given prefix to another prefix
func dbChangePrefix(oldprefix []byte, newprefix []byte) {
	iter := gldb.NewIterator(util.BytesPrefix(oldprefix), nil)
	defer iter.Release()
	for iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		key[0] = newprefix[0]
		gldb.Put(key, iter.Value(), nil)
		gldb.Delete(iter.Key(), nil)
	}
}

var (
	ansiDefault = "\033[0m"
	ansiHex     = "\033[36m"    // cyan on black
	ansiText    = "\033[37;40m" // white on dark gray
	ansiDash    = "\033[31m"    // red on black

	ansiDefaultByte = []byte(ansiDefault)
	ansiHexByte     = []byte(ansiHex)
	ansiTextByte    = []byte(ansiText)
)

// display entire contents of the database, by table.
func dbDisplay(w io.Writer) {
	for i, name := range PrefixNames {
		prefix := []byte{byte(i)}
		first := true
		iter := gldb.NewIterator(util.BytesPrefix(prefix), nil)
		for iter.Next() {
			if first {
				first = false
				fmt.Fprintf(w, "\n\n%s\n", name)
			}
			fmt.Fprint(w, dbViewable(iter.Key())+ansiDash+" -- "+ansiDefault+dbViewable(iter.Value())+ansiDefault+"\n\n")
		}
		iter.Release()
	}

	// display any records after the last one with a legitimate prefix
	iter := gldb.NewIterator(nil, nil)
	defer iter.Release()
	first := true
	for ok := iter.Seek([]byte{byte(len(PrefixNames))}); ok; ok = iter.Next() {
		if first {
			first = false
			fmt.Fprintf(w, "\n\n========== ORPHANS ==========\n\n")
		}
		fmt.Fprint(w, dbViewable(iter.Key())+ansiDash+" -- "+ansiDefault+dbViewable(iter.Value())+ansiDefault+"\n\n")
	}
}

// display entire contents of the database, all in one list.
func dbDisplayAll(w io.Writer) {
	iter := gldb.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		fmt.Fprint(w, dbViewable(iter.Key())+ansiDash+" -- "+ansiDefault+dbViewable(iter.Value())+ansiDefault+"\n\n")
	}
}

// make a byte slice of mixed ascii/binary data into a string to be printed.

func dbViewable(in []byte) string {
	space := []byte{' '}

	out := &bytes.Buffer{}

	first := true
	lastPrintable := false
	for i, b := range in {
		if (b >= 33) && (b <= 126) {
			// printable (don't consider the space character to be printable)
			if !lastPrintable {
				out.Write(ansiDefaultByte)
				if !first {
					out.Write(space)
				}
				out.Write(ansiTextByte)
			}
			lastPrintable = true
			out.Write(in[i : i+1])
		} else {
			// not printable; print one byte in hex
			out.Write(ansiDefaultByte)
			if !first {
				out.Write(space)
			}
			out.Write(ansiHexByte)
			lastPrintable = false
			fmt.Fprintf(out, "%2.2x", in[i])
		}
		first = false
	}

	return out.String()
}

////////////////////////////////////////////////////////
// XXX - delete

// migrate all Activity records from Activity00 to Activity
func dbMigrate() {
	iter := gldb.NewIterator(util.BytesPrefix(PfxActivity), nil)
	defer iter.Release()
	for iter.Next() {

		key := iter.Key()

		fmt.Fprint(os.Stderr, "\n")

		a00, err := dbActivity00Get(key)
		if err == nil {
			fmt.Fprintln(os.Stderr, dbViewable(key))
		} else {
			fmt.Fprintln(os.Stderr, "ERROR -- "+dbViewable(key))
		}

		fmt.Fprint(os.Stderr, "OLD -- ")
		fmt.Fprintln(os.Stderr, a00)

		a := Activity{
			PlanName:                a00.PlanName,
			Day:                     make([]int, len(a00.Done)),
			BibleVersion:            a00.BibleVersion,
			BibleProvider:           a00.BibleProvider,
			AccountabilityStartDate: a00.AccountabilityStartDate,
			AccountabilityVisible:   a00.AccountabilityVisible,
		}

		for d := 0; d < len(a00.Done); d++ {
			a.Day[d] = a00.CurrentDay
		}

		if err := dbActivityPut(key, &a); err != nil {
			fmt.Fprintln(os.Stderr, "==ERROR== "+err.Error())
		}

		fmt.Fprint(os.Stderr, "NEW -- ")
		fmt.Fprintln(os.Stderr, a)
		fmt.Fprint(os.Stderr, "\n")

	}
}

func dbActivity00Get(key []byte) (*Activity00, error) {
	data, err := gldb.Get(key, nil)
	if err != nil {
		return nil, wrapErr{"dbActivity00Get - Get", err}
	}
	a := &Activity00{}
	if err := a.Unmarshal(data); err != nil {
		return nil, wrapErr{"dbActivity00Get - Unmarshal", err}
	}
	return a, nil
}
