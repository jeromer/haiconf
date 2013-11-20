package haiconf

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type CheckersTestSuite struct{}

var _ = Suite(&CheckersTestSuite{})

func (s *CheckersTestSuite) TestCheckAbsolutePath_PathInvalid(c *C) {
	p, err := CheckAbsolutePath("Path", CommandArgs{})
	c.Assert(err, ErrorMatches, "Path must be provided(.*)")
	c.Assert(p, Equals, "")
}

func (s *CheckersTestSuite) TestCheckAbsolutePath_RelativePath(c *C) {
	path := "./relative/dir"
	p, err := CheckAbsolutePath("Path", CommandArgs{"Path": path})
	c.Assert(err, ErrorMatches, "Path must be absolute(.*)")
	c.Assert(p, Equals, path)
}

func (s *CheckersTestSuite) TestCheckEnsure_NotProvided(c *C) {
	e, err := CheckEnsure(CommandArgs{})
	c.Assert(err, ErrorMatches, "Ensure must be provided(.*)")
	c.Assert(e, Equals, "")
}

func (s *CheckersTestSuite) TestCheckEnsure_WrongChoice(c *C) {
	_, err := CheckEnsure(CommandArgs{"Ensure": "foo"})
	c.Assert(err, ErrorMatches, "Invalid choice for Ensure.(.*)")
}

func (s *CheckersTestSuite) TestCheckSystemUser_NotProvided(c *C) {
	o, err := CheckSystemUser("Owner", CommandArgs{})
	c.Assert(err, ErrorMatches, "Owner must be defined(.*)")
	c.Assert(o, IsNil)
}

func (s *CheckersTestSuite) TestCheckSystemUser_Exists(c *C) {
	username := "nobody"
	o, err := CheckSystemUser("Owner", CommandArgs{"Owner": username})

	c.Assert(err, IsNil)
	c.Assert(o.Username, Equals, username)
}

func (s *CheckersTestSuite) TestCheckSystemUser_DoesNotExists(c *C) {
	username := "azertyuiop-1234567890"
	o, err := CheckSystemUser("Owner", CommandArgs{"Owner": username})

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "user: unknown user "+username)
	c.Assert(o, IsNil)
}

func (s *CheckersTestSuite) TestCheckSystemGroup_NotProvided(c *C) {
	g, err := CheckSystemGroup("Group", CommandArgs{})
	c.Assert(err, ErrorMatches, "Group must be defined(.*)")
	c.Assert(g, IsNil)
}

func (s *CheckersTestSuite) TestCheckSystemGroup_Exists(c *C) {
	groupName := "nobody"
	g, err := CheckSystemGroup("Group", CommandArgs{"Group": groupName})

	c.Assert(err, IsNil)
	c.Assert(g.Name, Equals, groupName)
}

func (s *CheckersTestSuite) TestCheckSystemGroup_DoesNotExists(c *C) {
	groupName := "foo"
	g, err := CheckSystemGroup("Group", CommandArgs{"Group": groupName})

	c.Assert(err, NotNil)
	c.Assert(g, IsNil)
}

func (s *CheckersTestSuite) TestCheckBoolean_Provided(c *C) {
	r := CheckBool("Recurse", CommandArgs{"Recurse": true})
	c.Assert(r, Equals, true)
}

func (s *CheckersTestSuite) TestCheckBoolean_NotProvided(c *C) {
	r := CheckBool("Recurse", CommandArgs{})
	c.Assert(r, Equals, false)
}

func (s *CheckersTestSuite) TestCheckInt64_NotProvided(c *C) {
	m, err := CheckInt64("Mode", CommandArgs{})
	c.Assert(err, ErrorMatches, "Mode must be provided(.*)")
	c.Assert(m, Equals, int64(0))
}

func (s *CheckersTestSuite) TestCheckInt64_Provided(c *C) {
	m, err := CheckInt64("Mode", CommandArgs{"Mode": "0750"})
	c.Assert(err, IsNil)
	c.Assert(m, Equals, int64(0750))
}

func (s *CheckersTestSuite) TestCheckString_Empty(c *C) {
	p, err := CheckString("String", CommandArgs{})
	c.Assert(err, ErrorMatches, "String must be provided(.*)")
	c.Assert(p, Equals, "")
}

func (s *CheckersTestSuite) TestCheckString_NonEmpty(c *C) {
	p, err := CheckString("String", CommandArgs{"String": "foo"})
	c.Assert(err, IsNil)
	c.Assert(p, Equals, "foo")
}

func (s *CheckersTestSuite) TestCheckStringList_Empty(c *C) {
	p, err := CheckStringList("StringList", CommandArgs{})
	c.Assert(err, ErrorMatches, "StringList must be provided(.*)")
	var expected []string
	c.Assert(p, DeepEquals, expected)
}

func (s *CheckersTestSuite) TestCheckStringList_NonEmpty_Interface(c *C) {
	// XXX : we sometime receive a list of interface containing strings
	p, err := CheckStringList("StringList", CommandArgs{
		"StringList": []interface{}{"foo"},
	})

	c.Assert(err, IsNil)
	c.Assert(p, DeepEquals, []string{"foo"})
}
