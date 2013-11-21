package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/haiconf"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"os/user"
	"path"
)

type FileTestSuite struct {
	f *File
}

var _ = Suite(&FileTestSuite{})

func (s *FileTestSuite) SetUpTest(c *C) {
	s.f = new(File)
	s.f.SetDefault(&dummyRuntimeConfig)
}

func (s *FileTestSuite) TestSetDefault(c *C) {
	expected := &File{
		path:              "",
		mode:              DEFAULT_MODE_FILE,
		owner:             new(user.User),
		group:             new(hacks.Group),
		ensure:            haiconf.ENSURE_PRESENT,
		templateVariables: nil,
		rc:                &dummyRuntimeConfig,
	}

	c.Assert(s.f, DeepEquals, expected)
}

func (s *FileTestSuite) TestSetUserConfig_Present(c *C) {
	args := haiconf.CommandArgs{
		"Path":   "/foo.txt",
		"Ensure": haiconf.ENSURE_PRESENT,
		"Mode":   "0777",
		"Owner":  "nobody",
		"Group":  "nogroup",
		"Source": "/foo",
	}

	err := s.f.SetUserConfig(args)
	c.Assert(err, IsNil)

	u, err := user.Lookup("nobody")
	c.Assert(err, IsNil)

	g, err := hacks.LookupSystemGroup("nogroup")
	c.Assert(err, IsNil)

	c.Assert(s.f.path, Equals, args["Path"])
	c.Assert(s.f.mode, Equals, os.FileMode(0777))
	c.Assert(s.f.owner, DeepEquals, u)
	c.Assert(s.f.group, DeepEquals, g)
	c.Assert(s.f.source, DeepEquals, args["Source"])
	c.Assert(s.f.ensure, Equals, args["Ensure"])
}

func (s *FileTestSuite) TestSetUserConfig_Absent(c *C) {
	args := haiconf.CommandArgs{
		"Path":   "/foo",
		"Ensure": haiconf.ENSURE_ABSENT,
		"Mode":   0777,
		"Owner":  "nobody",
		"Group":  "nogroup",
	}

	err := s.f.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.f.path, Equals, args["Path"])
	c.Assert(s.f.ensure, Equals, args["Ensure"])
	c.Assert(s.f.mode, Equals, DEFAULT_MODE_FILE)

	// Since we want to remove the directory we do not care about
	// the value of attributes below
	c.Assert(s.f.owner, DeepEquals, new(user.User))
	c.Assert(s.f.group, DeepEquals, new(hacks.Group))
	c.Assert(s.f.templateVariables, IsNil)
	c.Assert(s.f.source, Equals, "")
}

func (s *FileTestSuite) TestRun_CreateRecursive(c *C) {
	cwd, err := os.Getwd()
	c.Assert(err, IsNil)

	tmpFile := c.MkDir() + "/foo/bar/baz/foo.txt"
	sourceFile := cwd + "/testdata/nontemplate.txt"

	err = s.f.SetUserConfig(haiconf.CommandArgs{
		"Path":   tmpFile,
		"Ensure": haiconf.ENSURE_PRESENT,
		"Mode":   "0600",
		"Owner":  currentUser.Username,
		"Group":  dummyGroup,
		"Source": sourceFile,
	})
	c.Assert(err, IsNil)

	err = s.f.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpFile)
	c.Assert(err, IsNil)
	c.Assert(f.Name(), Equals, path.Base(tmpFile))
	c.Assert(f.Mode().Perm(), Equals, os.FileMode(0600).Perm())
	c.Assert(f.IsDir(), Equals, false)

	obtained, err := ioutil.ReadFile(tmpFile)
	c.Assert(err, IsNil)

	expected, err := ioutil.ReadFile(sourceFile)
	c.Assert(err, IsNil)

	c.Assert(obtained, DeepEquals, expected)
}

func (s *FileTestSuite) TestRun_CreateNonRecursive(c *C) {
	cwd, err := os.Getwd()
	c.Assert(err, IsNil)

	tmpFile := c.MkDir() + "/foo.txt"
	sourceFile := cwd + "/testdata/nontemplate.txt"

	err = s.f.SetUserConfig(haiconf.CommandArgs{
		"Path":   tmpFile,
		"Ensure": haiconf.ENSURE_PRESENT,
		"Mode":   "0644",
		"Owner":  currentUser.Username,
		"Group":  dummyGroup,
		"Source": sourceFile,
	})
	c.Assert(err, IsNil)

	err = s.f.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpFile)
	c.Assert(err, IsNil)
	c.Assert(f.Name(), Equals, path.Base(tmpFile))
	c.Assert(f.Mode().Perm(), Equals, os.FileMode(0644).Perm())
	c.Assert(f.IsDir(), Equals, false)

	obtained, err := ioutil.ReadFile(tmpFile)
	c.Assert(err, IsNil)

	expected, err := ioutil.ReadFile(sourceFile)
	c.Assert(err, IsNil)

	c.Assert(obtained, DeepEquals, expected)
}

func (s *FileTestSuite) TestRun_CreateApplyTemplates(c *C) {
	cwd, err := os.Getwd()
	c.Assert(err, IsNil)

	tmpFile := c.MkDir() + "/foo.txt"
	sourceFilePath := cwd + "/testdata/template.txt"

	err = s.f.SetUserConfig(haiconf.CommandArgs{
		"Path":   tmpFile,
		"Ensure": haiconf.ENSURE_PRESENT,
		"Mode":   "0644",
		"Owner":  currentUser.Username,
		"Group":  dummyGroup,
		"Source": sourceFilePath,
		"TemplateVariables": map[string]interface{}{
			"Int":    123,
			"String": "foo",
			"Float":  3.14,
			"Slice":  []int{1, 2, 3},
		},
	})
	c.Assert(err, IsNil)

	err = s.f.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpFile)
	c.Assert(err, IsNil)
	c.Assert(f.Name(), Equals, path.Base(tmpFile))
	c.Assert(f.Mode().Perm(), Equals, os.FileMode(0644).Perm())
	c.Assert(f.IsDir(), Equals, false)

	obtained, err := ioutil.ReadFile(tmpFile)
	c.Assert(err, IsNil)

	expected := "foo | 123 | 3.14 | [1 2 3]\n"
	c.Assert(string(obtained), DeepEquals, expected)
}

func (s *FileTestSuite) TestRun_Remove(c *C) {
	fileName := "foo.txt"
	tmpFile := c.MkDir() + "/" + fileName

	err := ioutil.WriteFile(tmpFile, []byte{}, 0644)
	c.Assert(err, IsNil)

	err = s.f.SetUserConfig(haiconf.CommandArgs{
		"Path":   tmpFile,
		"Ensure": haiconf.ENSURE_ABSENT,
	})
	c.Assert(err, IsNil)

	err = s.f.Run()
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpFile)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
}

func (s *FileTestSuite) TestSetTemplateVariables_NoVariables(c *C) {
	err := s.f.setTemplateVariables(haiconf.CommandArgs{})
	c.Assert(err, IsNil)
	c.Assert(s.f.templateVariables, IsNil)
}

func (s *FileTestSuite) TestSetTemplateVariables_VariablesConverted(c *C) {
	err := s.f.setTemplateVariables(haiconf.CommandArgs{
		"TemplateVariables": map[string]interface{}{
			"BoolTrue":       "true",
			"BoolFalse":      "false",
			"StandardString": "foo",
		},
	})
	c.Assert(err, IsNil)

	expected := map[string]interface{}{
		"BoolTrue":       true,
		"BoolFalse":      false,
		"StandardString": "foo",
	}
	c.Assert(s.f.templateVariables, DeepEquals, expected)
}
