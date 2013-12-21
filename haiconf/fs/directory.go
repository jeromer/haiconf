// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Usage in lua configuration file
//
// Directory({
//         Path    = "/tmp/haiconf/testdirectory",
//         Mode    = "0755",
//         Owner   = "jerome",
//         Group   = "wheel",
//         Recurse = true,
//         Ensure  = "present",
//     })
//
package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/haiconf"
	"os"
	"os/user"
)

type Directory struct {
	Path    string
	Mode    os.FileMode
	Owner   *user.User
	Recurse bool
	Ensure  string

	// no user.Group in golang yet
	// (https://code.google.com/p/go/issues/detail?id=2617)
	// Let's use a temporary one
	Group *hacks.Group

	rc *haiconf.RuntimeConfig
}

func (d *Directory) SetDefault(rc *haiconf.RuntimeConfig) error {
	*d = Directory{
		Path:    "",
		Mode:    DEFAULT_MODE_DIRECTORY,
		Owner:   new(user.User),
		Group:   new(hacks.Group),
		Recurse: false,
		Ensure:  haiconf.ENSURE_PRESENT,
		rc:      rc,
	}

	return nil
}

func (d *Directory) SetUserConfig(args haiconf.CommandArgs) error {
	err := d.setPath(args)
	if err != nil {
		return err
	}

	err = d.setEnsure(args)
	if err != nil {
		return err
	}

	err = d.setRecurse(args)
	if err != nil {
		return err
	}

	if d.Ensure == haiconf.ENSURE_ABSENT {
		return nil
	}

	err = d.setMode(args)
	if err != nil {
		return err
	}

	err = d.setOwner(args)
	if err != nil {
		return err
	}

	err = d.setGroup(args)
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) Run() error {
	// XXX : acquire/release lock
	if d.Ensure == haiconf.ENSURE_ABSENT {
		haiconf.Output(d.rc, "Removing directory %s", d.Path)
		return RmDir(d.Path, d.Recurse)
	}

	haiconf.Output(d.rc, "Creating directory %s", d.Path)
	err := MkDir(d.Path, d.Recurse, d.Mode)
	if err != nil {
		return err
	}

	haiconf.Output(d.rc, "Chmod %s on %s", d.Mode, d.Path)
	err = Chmod(d.Path, d.Mode)
	if err != nil {
		return err
	}

	haiconf.Output(d.rc, "Chown %s:%s on %s", d.Owner.Username, d.Group.Name, d.Path)
	err = Chown(d.Path, d.Owner, d.Group)
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) setPath(args haiconf.CommandArgs) error {
	p, err := haiconf.CheckAbsolutePath("Path", args)
	if err != nil {
		return err
	}

	d.Path = p

	return nil
}

func (d *Directory) setEnsure(args haiconf.CommandArgs) error {
	e, err := haiconf.CheckEnsure(args)
	if err != nil {
		return err
	}

	d.Ensure = e

	return nil
}

func (d *Directory) setMode(args haiconf.CommandArgs) error {
	m, err := haiconf.CheckInt64("Mode", args)
	if err != nil {
		return err
	}

	d.Mode = os.FileMode(m)

	return nil
}

func (d *Directory) setOwner(args haiconf.CommandArgs) error {
	u, err := haiconf.CheckSystemUser("Owner", args)
	if err != nil {
		return err
	}

	d.Owner = u
	return nil
}

func (d *Directory) setGroup(args haiconf.CommandArgs) error {
	grp, err := haiconf.CheckSystemGroup("Group", args)
	if err != nil {
		return err
	}

	d.Group = grp
	return nil
}

func (d *Directory) setRecurse(args haiconf.CommandArgs) error {
	d.Recurse = haiconf.CheckBool("Recurse", args)
	return nil
}
