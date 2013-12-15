// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hacks

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type SystemGroupTestSuite struct{}

var _ = Suite(&SystemGroupTestSuite{})

func (s *SystemGroupTestSuite) TestLookupSystemGroup_UnknownGroup(c *C) {
	g, err := LookupSystemGroup("foobar-1234")
	c.Assert(err, NotNil)
	c.Assert(g, IsNil)
}

func (s *SystemGroupTestSuite) TestLookupSystemGroup_KnownGroup(c *C) {
	groupName := "nogroup"
	g, err := LookupSystemGroup(groupName)

	c.Assert(err, IsNil)

	c.Assert(g.Name, Equals, groupName)
	c.Assert(len(g.Gid) > 0, Equals, true)
}
