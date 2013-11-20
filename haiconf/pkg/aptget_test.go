package pkg

import (
	"github.com/jeromer/haiconf/haiconf"
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type AptGetTestSuite struct {
	ag *AptGet
}

var _ = Suite(&AptGetTestSuite{})

func (s *AptGetTestSuite) SetUpTest(c *C) {
	s.ag = new(AptGet)
	err := s.ag.SetDefault()
	c.Assert(err, IsNil)
}

func (s *AptGetTestSuite) TestSetPackages_NotFromList(c *C) {
	args := haiconf.CommandArgs{
		"Packages": []string{"a", "b"},
	}

	err := s.ag.setPackages(args)
	c.Assert(err, IsNil)
}

func (s *AptGetTestSuite) TestSetPackages_FromList(c *C) {
	cwd, err := os.Getwd()
	c.Assert(err, IsNil)

	args := haiconf.CommandArgs{
		"PackagesFromSource": cwd + "/testdata/packages.txt",
	}

	err = s.ag.setPackages(args)
	c.Assert(err, IsNil)
	c.Assert(s.ag.packages, DeepEquals, []string{"vim", "cowsay"})
}

func (s *AptGetTestSuite) TestSetPackages_DuplicateRemoved(c *C) {
	args := haiconf.CommandArgs{
		"Packages": []string{"a", "b", "a"},
	}

	err := s.ag.setPackages(args)
	c.Assert(err, IsNil)
	c.Assert(s.ag.packages, DeepEquals, []string{"a", "b"})
}

func (s *AptGetTestSuite) TestSetPackages_PackagesHasPrecedence(c *C) {
	cwd, err := os.Getwd()
	c.Assert(err, IsNil)

	args := haiconf.CommandArgs{
		"Packages":           []string{"foo", "bar"},
		"PackagesFromSource": cwd + "/testdata/packages.txt",
	}

	err = s.ag.setPackages(args)
	c.Assert(err, IsNil)
	// Packages has precedence over Packagesfromsource
	c.Assert(s.ag.packages, DeepEquals, []string{"foo", "bar"})
}

func (s *AptGetTestSuite) TestSetExtraOptions_Empty(c *C) {
	err := s.ag.setExtraOptions(haiconf.CommandArgs{})
	c.Assert(err, IsNil)
	c.Assert(s.ag.extraOptions, DeepEquals, []string(nil))
}

func (s *AptGetTestSuite) TestSetExtraOptions_Provided(c *C) {
	xtraOpts := []string{"a", "b"}
	err := s.ag.setExtraOptions(haiconf.CommandArgs{"ExtraOptions": xtraOpts})
	c.Assert(err, IsNil)
	c.Assert(s.ag.extraOptions, DeepEquals, xtraOpts)
}

func (s *AptGetTestSuite) TestSetExtraOptions_DuplicateRemoved(c *C) {
	xtraOpts := []string{"a", "b", "a", "a"}
	err := s.ag.setExtraOptions(haiconf.CommandArgs{"ExtraOptions": xtraOpts})
	c.Assert(err, IsNil)
	c.Assert(s.ag.extraOptions, DeepEquals, []string{"a", "b"})
}

func (s *AptGetTestSuite) TestSetUserConfig_Complete(c *C) {
	args := haiconf.CommandArgs{
		"Method":       "install",
		"Packages":     []string{"a", "b"},
		"ExtraOptions": []string{"foo", "bar"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.ag.method, Equals, args["Method"])
	c.Assert(s.ag.packages, DeepEquals, args["Packages"])
	c.Assert(s.ag.extraOptions, DeepEquals, args["ExtraOptions"])
}

func (s *AptGetTestSuite) TestSetUserConfig_UpdateMethod(c *C) {
	args := haiconf.CommandArgs{
		"Method":       "update",
		"Packages":     []string{"a", "b"},
		"ExtraOptions": []string{"foo", "bar"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.ag.method, Equals, args["Method"])
	c.Assert(s.ag.packages, DeepEquals, []string(nil))
	c.Assert(s.ag.extraOptions, DeepEquals, []string(nil))
}

func (s *AptGetTestSuite) TestRun(c *C) {
	args := haiconf.CommandArgs{
		"Method":       "install",
		"Packages":     []string{"azertyuiop"},
		"ExtraOptions": []string{"-o", "DPkg::Options::=--force-confnew"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	err = s.ag.Run()
	c.Assert(err, NotNil)
}