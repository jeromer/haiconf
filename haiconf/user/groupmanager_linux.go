// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/jeromer/haiconf/haiconf/osutils"
	"os"
)

type GroupManager struct {
	groupAddPath string
	groupDelPath string
	Name         string
}

func NewGroupManager() *GroupManager {
	return &GroupManager{
		groupAddPath: "/usr/sbin/groupadd",
		groupDelPath: "/usr/sbin/groupdel",
	}
}

func (mgr *GroupManager) Add() error {
	return mgr.run(mgr.groupAddCmd())
}

func (mgr *GroupManager) Remove() error {
	return mgr.run(mgr.groupDelCmd())
}

func (mgr *GroupManager) groupAddCmd() osutils.SystemCommand {
	return osutils.SystemCommand{
		Path:                 mgr.groupAddPath,
		Args:                 []string{mgr.Name},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}
}

func (mgr *GroupManager) groupDelCmd() osutils.SystemCommand {
	return osutils.SystemCommand{
		Path:                 mgr.groupDelPath,
		Args:                 []string{mgr.Name},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}
}

func (mgr *GroupManager) run(c osutils.SystemCommand) error {
	output := c.Run()
	if output.HasError() {
		return output
	}

	return nil

}
