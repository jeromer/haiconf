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
	groupName := "nobody"
	g, err := LookupSystemGroup(groupName)

	c.Assert(err, IsNil)

	c.Assert(g.Name, Equals, groupName)
	c.Assert(len(g.Gid) > 0, Equals, true)
}
