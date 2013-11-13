package fs

import (
	"github.com/jeromer/haiconf/lib"
	. "launchpad.net/gocheck"
)

type CommonTestSuite struct{}

var _ = Suite(&CommonTestSuite{})

func (s *CommonTestSuite) TestCheckPath_PathInvalid(c *C) {
	p, err := CheckPath(haiconf.CommandArgs{})
	c.Assert(err, Equals, ErrNameEmpty)
	c.Assert(p, Equals, "")
}

func (s *CommonTestSuite) TestCheckPath_PathMustBeAbsolute(c *C) {
	path := "./relative/dir"
	p, err := CheckPath(haiconf.CommandArgs{"Path": path})
	c.Assert(err, Equals, ErrPathMustBeAbsolute)
	c.Assert(p, Equals, path)
}

func (s *CommonTestSuite) TestCheckEnsure_NotProvided(c *C) {
	e, err := CheckEnsure(haiconf.CommandArgs{})
	c.Assert(err, Equals, ErrEnsureEmpty)
	c.Assert(e, Equals, "")
}

func (s *CommonTestSuite) TestCheckEnsure_WrongChoice(c *C) {
	_, err := CheckEnsure(haiconf.CommandArgs{"Ensure": "foo"})
	c.Assert(err, Equals, ErrInvalidChoice)
}

func (s *CommonTestSuite) TestCheckOwner_Exists(c *C) {
	username := "nobody"
	o, err := CheckOwner(haiconf.CommandArgs{"Owner": username})

	c.Assert(err, IsNil)
	c.Assert(o.Username, Equals, username)
}

func (s *CommonTestSuite) TestCheckOwner_DoesNotExists(c *C) {
	username := "azertyuiop-1234567890"
	o, err := CheckOwner(haiconf.CommandArgs{"Owner": username})

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "user: unknown user "+username)
	c.Assert(o, IsNil)
}

func (s *CommonTestSuite) TestCheckGroup_Exists(c *C) {
	groupName := "nobody"
	g, err := CheckGroup(haiconf.CommandArgs{"Group": groupName})

	c.Assert(err, IsNil)
	c.Assert(g.Name, Equals, groupName)
}

func (s *CommonTestSuite) TestCheckGroup_DoesNotExists(c *C) {
	groupName := "foo"
	g, err := CheckGroup(haiconf.CommandArgs{"Group": groupName})

	c.Assert(err, NotNil)
	c.Assert(g, IsNil)
}

func (s *CommonTestSuite) TestCheckRecurse_Provided(c *C) {
	r := CheckRecurse(haiconf.CommandArgs{"Recurse": true})
	c.Assert(r, Equals, true)
}

func (s *CommonTestSuite) TestCheckRecurse_NotProvided(c *C) {
	r := CheckRecurse(haiconf.CommandArgs{})
	c.Assert(r, Equals, false)
}

func (s *CommonTestSuite) TestCheckMode_NotProvided(c *C) {
	m, err := CheckMode(haiconf.CommandArgs{})
	c.Assert(err, Equals, ErrModeEmpty)
	c.Assert(m, Equals, int64(0))
}

func (s *CommonTestSuite) TestCheckMode_Provided(c *C) {
	m, err := CheckMode(haiconf.CommandArgs{"Mode": "0750"})
	c.Assert(err, IsNil)
	c.Assert(m, Equals, int64(0750))
}
