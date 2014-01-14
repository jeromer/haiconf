// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/jeromer/haiconf/haiconf/osutils"
	. "launchpad.net/gocheck"
	"os"
)

type GroupManagerTestSuite struct {
	g *GroupManager
}

var (
	_ = Suite(&GroupManagerTestSuite{})
)

func (s *GroupManagerTestSuite) SetUpTest(c *C) {
	s.g = NewGroupManager()
}

func (s *GroupManagerTestSuite) TestGroupAddCmd(c *C) {
	n := "foo"
	s.g.Name = n

	obtained := s.g.groupAddCmd()
	expected := osutils.SystemCommand{
		Path:                 s.g.groupAddPath,
		Args:                 []string{n},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	c.Assert(obtained, DeepEquals, expected)
}

func (s *GroupManagerTestSuite) TestGroupDelCmd(c *C) {
	n := "foo"
	s.g.Name = n

	obtained := s.g.groupDelCmd()
	expected := osutils.SystemCommand{
		Path:                 s.g.groupDelPath,
		Args:                 []string{n},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	c.Assert(obtained, DeepEquals, expected)
}
