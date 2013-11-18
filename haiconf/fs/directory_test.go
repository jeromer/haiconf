package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/haiconf"
	. "launchpad.net/gocheck"
	"os"
	"os/user"
	"path"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type DirectoryTestSuite struct {
	d *Directory
}

var _ = Suite(&DirectoryTestSuite{})

func (s *DirectoryTestSuite) SetUpTest(c *C) {
	s.d = new(Directory)
	err := s.d.SetDefault()
	c.Assert(err, IsNil)
}

func (s *DirectoryTestSuite) TestSetDefault(c *C) {
	// XXX : s.d.SetDefault() called in Setuptest

	expected := &Directory{
		path:    "",
		mode:    0755,
		owner:   new(user.User),
		group:   new(hacks.Group),
		recurse: false,
		ensure:  haiconf.ENSURE_PRESENT,
	}

	c.Assert(s.d, DeepEquals, expected)
}

func (s *DirectoryTestSuite) TestSetUserConfig_Complete(c *C) {
	args := haiconf.CommandArgs{
		"Path":    "/foo",
		"Ensure":  haiconf.ENSURE_PRESENT,
		"Recurse": true,
		"Mode":    "0777",
		"Owner":   "nobody",
		"Group":   "nogroup",
	}

	err := s.d.SetUserConfig(args)
	c.Assert(err, IsNil)

	u, err := user.Lookup("nobody")
	c.Assert(err, IsNil)

	g, err := hacks.LookupSystemGroup("nogroup")
	c.Assert(err, IsNil)

	c.Assert(s.d.path, Equals, args["Path"])
	c.Assert(s.d.mode, Equals, os.FileMode(0777))
	c.Assert(s.d.owner, DeepEquals, u)
	c.Assert(s.d.group, DeepEquals, g)
	c.Assert(s.d.recurse, Equals, args["Recurse"])
	c.Assert(s.d.ensure, Equals, args["Ensure"])
}

func (s *DirectoryTestSuite) TestSetUserConfig_Absent(c *C) {
	args := haiconf.CommandArgs{
		"Path":    "/foo",
		"Ensure":  haiconf.ENSURE_ABSENT,
		"Recurse": true,
		"Mode":    0777,
		"Owner":   "nobody",
		"Group":   "nogroup",
	}

	err := s.d.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.d.path, Equals, args["Path"])
	c.Assert(s.d.recurse, Equals, args["Recurse"])
	c.Assert(s.d.ensure, Equals, args["Ensure"])
	c.Assert(s.d.mode, Equals, DEFAULT_MODE_DIRECTORY)

	// Since we want to remove the directory we do not care about
	// the value of attributes below
	c.Assert(s.d.owner, DeepEquals, new(user.User))
	c.Assert(s.d.group, DeepEquals, new(hacks.Group))
}

func (s *DirectoryTestSuite) TestRun_Create(c *C) {
	tmpDir := c.MkDir() + "/foo/bar/baz"

	cu, err := user.Current()
	c.Assert(err, IsNil)

	err = s.d.SetUserConfig(haiconf.CommandArgs{
		"Path":    tmpDir,
		"Owner":   cu.Username,
		"Group":   "nogroup",
		"Recurse": true,
		"Mode":    "0777",
		"Ensure":  haiconf.ENSURE_PRESENT,
	})
	c.Assert(err, IsNil)

	err = s.d.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.Name(), Equals, path.Base(tmpDir))
	c.Assert(f.Mode().Perm(), Equals, os.FileMode(0777).Perm())
	c.Assert(f.IsDir(), Equals, true)

	/*
		sys := f.Sys().(*syscall.Stat_t)
		fmt.Printf("%+v\n", sys)
		c.Assert(sys.Uid, Equals, XXX)
		c.Assert(sys.Gid, Equals, XXX)
	*/
}

func (s *DirectoryTestSuite) TestRun_Remove(c *C) {
	tmpDir := c.MkDir() + "/foo/bar/baz"

	err := os.MkdirAll(tmpDir, os.FileMode(0755))
	c.Assert(err, IsNil)
	_, err = os.Stat(tmpDir)
	c.Assert(os.IsNotExist(err), Equals, false)

	err = s.d.SetUserConfig(haiconf.CommandArgs{
		"Path":    tmpDir,
		"Recurse": true,
		"Ensure":  haiconf.ENSURE_ABSENT,
	})
	c.Assert(err, IsNil)

	err = s.d.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
}
