// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Usage in lua configuration file
//
//  Group({
//      Name="testgroup",
//      Ensure="present",
//  })

package user

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/haiconf"
)

const (
	ACTION_NOOP   = 0
	ACTION_CREATE = 1
	ACTION_REMOVE = 2
)

type Group struct {
	Name   string
	Ensure string
	action int
	rc     *haiconf.RuntimeConfig
}

func (g *Group) SetDefault(rc *haiconf.RuntimeConfig) error {
	*g = Group{
		Name:   "",
		Ensure: haiconf.ENSURE_PRESENT,
		action: ACTION_NOOP,
		rc:     rc,
	}

	return nil
}

func (g *Group) SetUserConfig(args haiconf.CommandArgs) error {
	// XXX : must be first as setName() uses g.Ensure
	err := g.setEnsure(args)
	if err != nil {
		return err
	}

	err = g.setName(args)
	if err != nil {
		return err
	}

	return nil
}

func (g *Group) Run() error {
	if g.action == ACTION_NOOP {
		return nil
	}

	mgr := NewGroupManager()
	mgr.Name = g.Name

	if g.Ensure == haiconf.ENSURE_PRESENT {
		haiconf.Output(g.rc, "Adding group %s", g.Name)
		return mgr.Add()
	}

	haiconf.Output(g.rc, "Removing group %s", g.Name)
	return mgr.Remove()
}

func (g *Group) setName(args haiconf.CommandArgs) error {
	n, err := haiconf.CheckString("Name", args)
	if err != nil {
		return err
	}

	grp, err := hacks.LookupSystemGroup(n)
	exists := err == nil && grp.Gid != ""

	g.Name = n

	if g.Ensure == haiconf.ENSURE_PRESENT {
		if !exists {
			g.action = ACTION_CREATE
		}
		// ACTION_NOOP in setDefault()
	}

	if g.Ensure == haiconf.ENSURE_ABSENT {
		if exists {
			g.action = ACTION_REMOVE
		}
		// ACTION_NOOP in setDefault()
	}

	return nil
}

func (g *Group) setEnsure(args haiconf.CommandArgs) error {
	e, err := haiconf.CheckEnsure(args)
	if err != nil {
		return err
	}

	g.Ensure = e

	return nil
}
