package directory

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/lib"
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
	s.d.SetDefault()
}

func (s *DirectoryTestSuite) TestSetDefault(c *C) {
	// XXX : s.d.SetDefault() called in Setuptest

	expected := &Directory{
		path:    "",
		mode:    0755,
		owner:   new(user.User),
		group:   new(hacks.Group),
		recurse: false,
		ensure:  ENSURE_PRESENT,
	}

	c.Assert(s.d, DeepEquals, expected)
}

func (s *DirectoryTestSuite) TestSetPath_PathInvalid(c *C) {
	err := s.d.setPath(haiconf.CommandArgs{})
	c.Assert(err, Equals, ErrNameEmpty)
	c.Assert(s.d.path, Equals, "")
}

func (s *DirectoryTestSuite) TestSetPath_PathMustBeAbsolute(c *C) {
	err := s.d.setPath(haiconf.CommandArgs{"Path": "./relative/dir"})
	c.Assert(err, Equals, ErrPathMustBeAbsolute)
	c.Assert(s.d.path, Equals, "")
}

func (s *DirectoryTestSuite) TestSetEnsure_FallbackToDefault(c *C) {
	err := s.d.setEnsure(haiconf.CommandArgs{})
	c.Assert(err, IsNil)
	c.Assert(s.d.ensure, Equals, ENSURE_PRESENT)
}

func (s *DirectoryTestSuite) TestSetEnsure_WrongChoice(c *C) {
	err := s.d.setEnsure(haiconf.CommandArgs{"Ensure": "foo"})
	c.Assert(err, Equals, ErrInvalidChoice)
}

func (s *DirectoryTestSuite) TestSetOwner_Exists(c *C) {
	o := "nobody"
	err := s.d.setOwner(haiconf.CommandArgs{"Owner": o})

	c.Assert(err, IsNil)
	c.Assert(s.d.owner.Username, Equals, o)
}

func (s *DirectoryTestSuite) TestSetOwner_DoesNotExists(c *C) {
	o := "azertyuiop-1234567890"
	err := s.d.setOwner(haiconf.CommandArgs{"Owner": o})

	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "user: unknown user "+o)
	c.Assert(s.d.owner, DeepEquals, new(user.User))
}

func (s *DirectoryTestSuite) TestSetGroup_Exists(c *C) {
	g := "nobody"
	err := s.d.setGroup(haiconf.CommandArgs{"Group": g})

	c.Assert(err, IsNil)
	c.Assert(s.d.group.Name, Equals, g)
}

func (s *DirectoryTestSuite) TestSetGroup_DoesNotExists(c *C) {
	g := "foo"
	err := s.d.setGroup(haiconf.CommandArgs{"Group": g})

	c.Assert(err, NotNil)
	c.Assert(s.d.group, DeepEquals, new(hacks.Group))
}

func (s *DirectoryTestSuite) TestSetRecurse_Provided(c *C) {
	err := s.d.setRecurse(haiconf.CommandArgs{"Recurse": true})
	c.Assert(err, IsNil)
	c.Assert(s.d.recurse, Equals, true)
}

func (s *DirectoryTestSuite) TestSetRecurse_NotProvided(c *C) {
	err := s.d.setRecurse(haiconf.CommandArgs{})
	c.Assert(err, IsNil)
	c.Assert(s.d.recurse, Equals, false)
}

func (s *DirectoryTestSuite) TestSetMode_FallbackToDefault(c *C) {
	err := s.d.setMode(haiconf.CommandArgs{})
	c.Assert(err, IsNil)
	c.Assert(s.d.mode, Equals, DEFAULT_MODE)
}

func (s *DirectoryTestSuite) TestSetMode_Provided(c *C) {
	err := s.d.setMode(haiconf.CommandArgs{"Mode": "0750"})
	c.Assert(err, IsNil)
	c.Assert(s.d.mode, Equals, os.FileMode(0750))
}

func (s *DirectoryTestSuite) TestSetUserConfig_Complete(c *C) {
	args := haiconf.CommandArgs{
		"Path":    "/foo",
		"Ensure":  ENSURE_PRESENT,
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
		"Ensure":  ENSURE_ABSENT,
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
	c.Assert(s.d.mode, Equals, DEFAULT_MODE)

	// Since we want to remove the directory we do not care about
	// the value of attributes below
	c.Assert(s.d.owner, DeepEquals, new(user.User))
	c.Assert(s.d.group, DeepEquals, new(hacks.Group))
}

func (s *DirectoryTestSuite) TestMkDir_NonRecursive(c *C) {
	tmpDir := c.MkDir() + "/foo"
	s.d.path = tmpDir

	err := s.d.mkDir()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)
}

func (s *DirectoryTestSuite) TestMkDir_Recursive(c *C) {
	tmpDir := c.MkDir() + "/foo/bar/baz"
	s.d.path = tmpDir
	s.d.recurse = true

	err := s.d.mkDir()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)
}

func (s *DirectoryTestSuite) TestRmDir_NonRecursive(c *C) {
	tmpDir := c.MkDir() + "/foo"
	s.d.path = tmpDir

	err := s.d.mkDir()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)

	err = s.d.rmDir()
	c.Assert(err, IsNil)

	f, err = os.Stat(tmpDir)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
}

func (s *DirectoryTestSuite) TestRmDir_Recursive(c *C) {
	suffix := "/foo/bar/baz"
	tmpDir := c.MkDir()
	s.d.path = tmpDir + suffix
	s.d.recurse = true

	err := s.d.mkDir()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)

	err = s.d.rmDir()
	c.Assert(err, IsNil)

	f, err = os.Stat(tmpDir + suffix)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
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
		"Ensure":  ENSURE_PRESENT,
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
		"Ensure":  ENSURE_ABSENT,
	})
	c.Assert(err, IsNil)

	err = s.d.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
}
