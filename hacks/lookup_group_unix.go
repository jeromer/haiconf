// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd linux netbsd openbsd
// +build cgo

package hacks

import (
	"fmt"
	"runtime"
	"strconv"
	"syscall"
	"unsafe"
)

/*
#include <grp.h>
#include <unistd.h>
#include <stdlib.h>
*/
import "C"

type unknownGroupError string

type Group struct {
	Gid  string
	Name string
}

// This is a local implementation of system group lookups
// This function MUST be removed when the following bug is fixed
// -> https://code.google.com/p/go/issues/detail?id=2617
func LookupSystemGroup(groupName string) (*Group, error) {
	return lookupGroup(groupName)
}

func lookupGroup(groupName string) (*Group, error) {
	var grp C.struct_group
	var result *C.struct_group

	var bufSize C.long
	if runtime.GOOS == "freebsd" {
		// FreeBSD doesn't have _SC_GETPW_R_SIZE_MAX
		// and just returns -1.  So just use the same
		// size that Linux returns
		bufSize = 1024
	} else {
		bufSize = C.sysconf(C._SC_GETPW_R_SIZE_MAX)
		if bufSize <= 0 || bufSize > 1<<20 {
			return nil, fmt.Errorf("user: unreasonable _SC_GETPW_R_SIZE_MAX of %d", bufSize)
		}
	}

	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)
	var rv C.int

	nameC := C.CString(groupName)
	defer C.free(unsafe.Pointer(nameC))

	rv = C.getgrnam_r(
		nameC,
		&grp,
		(*C.char)(buf),
		C.size_t(bufSize),
		&result)

	if rv != 0 {
		return nil, fmt.Errorf("user: lookup group name %s: %s", groupName, syscall.Errno(rv))
	}

	if result == nil {
		return nil, unknownGroupError(groupName)
	}

	g := &Group{
		Gid:  strconv.Itoa(int(grp.gr_gid)),
		Name: C.GoString(grp.gr_name),
	}

	return g, nil
}

func (e unknownGroupError) Error() string {
	return "group: unknown group" + string(e)
}
