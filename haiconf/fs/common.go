// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"os"
	"os/user"
	"strconv"
)

const (
	DEFAULT_MODE_DIRECTORY = os.FileMode(0755)
	DEFAULT_MODE_FILE      = os.FileMode(0644)
)

func MkDir(path string, recurse bool, mode os.FileMode) error {
	// XXX : symlink support ?
	_, err := os.Stat(path)

	// directory already exists
	if err == nil {
		return nil
	}

	if recurse {
		return os.MkdirAll(path, mode)
	}

	return os.Mkdir(path, mode)
}

func RmDir(path string, recurse bool) error {
	// XXX : symlink support ?

	_, err := os.Stat(path)

	// directory does not exists
	if os.IsNotExist(err) {
		return nil
	}

	if recurse {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}

func Chmod(path string, mode os.FileMode) error {
	err := os.Chmod(path, mode)
	if err != nil {
		return err
	}

	return nil
}

func Chown(path string, usr *user.User, grp *hacks.Group) error {
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		return err
	}

	return os.Chown(path, uid, gid)
}
