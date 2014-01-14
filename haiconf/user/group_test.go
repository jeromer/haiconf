// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/jeromer/haiconf/haiconf"
	. "launchpad.net/gocheck"
	"math/rand"
	"strconv"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type GroupTestSuite struct {
	g *Group
}

var (
	_ = Suite(&GroupTestSuite{})

	dummyRuntimeConfig = haiconf.RuntimeConfig{
		Verbose: false,
		Output:  nil,
	}
)

func (s *GroupTestSuite) SetUpTest(c *C) {
	s.g = new(Group)
	s.g.SetDefault(&dummyRuntimeConfig)
}

func (s *GroupTestSuite) TestSetDefault(c *C) {
	expected := &Group{
		Name:   "",
		Ensure: haiconf.ENSURE_PRESENT,
		action: ACTION_NOOP,
		rc:     &dummyRuntimeConfig,
	}

	c.Assert(s.g, DeepEquals, expected)
}

func (s *GroupTestSuite) TestSetUserConfig_PresentGroupAlreadyExists(c *C) {
	args := haiconf.CommandArgs{
		"Name":   "nogroup",
		"Ensure": haiconf.ENSURE_PRESENT,
	}

	err := s.g.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.g.action, Equals, ACTION_NOOP)
}

func (s *GroupTestSuite) TestSetUserConfig_PresentGroupMissing(c *C) {
	args := haiconf.CommandArgs{
		"Name":   strconv.Itoa(rand.Int()),
		"Ensure": haiconf.ENSURE_PRESENT,
	}

	err := s.g.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.g.action, Equals, ACTION_CREATE)
}

func (s *GroupTestSuite) TestSetUserConfig_AbsentGroupAlreadyExists(c *C) {
	args := haiconf.CommandArgs{
		"Name":   "nogroup",
		"Ensure": haiconf.ENSURE_ABSENT,
	}

	err := s.g.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.g.action, Equals, ACTION_REMOVE)
}

func (s *GroupTestSuite) TestSetUserConfig_AbsenttGroupMissing(c *C) {
	args := haiconf.CommandArgs{
		"Name":   strconv.Itoa(rand.Int()),
		"Ensure": haiconf.ENSURE_ABSENT,
	}

	err := s.g.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.g.action, Equals, ACTION_NOOP)
}

func (s *GroupTestSuite) TestSetUserConfig_Complete(c *C) {
	n := "foo"
	e := haiconf.ENSURE_ABSENT

	args := haiconf.CommandArgs{
		"Name":   n,
		"Ensure": e,
	}

	err := s.g.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.g.Name, Equals, n)
	c.Assert(s.g.Ensure, Equals, e)
}
