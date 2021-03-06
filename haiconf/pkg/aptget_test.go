// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

var (
	_ = Suite(&AptGetTestSuite{})

	dummyRuntimeConfig = haiconf.RuntimeConfig{
		Verbose: false,
		Output:  nil,
	}
)

func (s *AptGetTestSuite) SetUpTest(c *C) {
	s.ag = new(AptGet)
	err := s.ag.SetDefault(&dummyRuntimeConfig)
	c.Assert(err, IsNil)
}

func (s *AptGetTestSuite) TestSetDefault(c *C) {
	// XXX : s.d.SetDefault() called in Setuptest
	c.Assert(s.ag.rc, Equals, &dummyRuntimeConfig)
}

func (s *AptGetTestSuite) TestSetPackages_NotFromList(c *C) {
	args := haiconf.CommandArgs{
		"Packages": []interface{}{"a", "b"},
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
	c.Assert(s.ag.Packages, DeepEquals, []string{"vim", "cowsay"})
}

func (s *AptGetTestSuite) TestSetPackages_DuplicateRemoved(c *C) {
	args := haiconf.CommandArgs{
		"Packages": []interface{}{"a", "b", "a"},
	}

	err := s.ag.setPackages(args)
	c.Assert(err, IsNil)
	c.Assert(s.ag.Packages, DeepEquals, []string{"a", "b"})
}

func (s *AptGetTestSuite) TestSetPackages_PackagesHasPrecedence(c *C) {
	cwd, err := os.Getwd()
	c.Assert(err, IsNil)

	args := haiconf.CommandArgs{
		"Packages":           []interface{}{"foo", "bar"},
		"PackagesFromSource": cwd + "/testdata/packages.txt",
	}

	err = s.ag.setPackages(args)
	c.Assert(err, IsNil)
	// Packages has precedence over Packagesfromsource
	c.Assert(s.ag.Packages, DeepEquals, []string{"foo", "bar"})
}

func (s *AptGetTestSuite) TestSetExtraOptions_Empty(c *C) {
	err := s.ag.setExtraOptions(haiconf.CommandArgs{})
	c.Assert(err, IsNil)
	c.Assert(s.ag.ExtraOptions, DeepEquals, []string(nil))
}

func (s *AptGetTestSuite) TestSetExtraOptions_Provided(c *C) {
	err := s.ag.setExtraOptions(haiconf.CommandArgs{
		"ExtraOptions": []interface{}{"a", "b"},
	})
	c.Assert(err, IsNil)
	c.Assert(s.ag.ExtraOptions, DeepEquals, []string{"a", "b"})
}

func (s *AptGetTestSuite) TestSetExtraOptions_DuplicateRemoved(c *C) {
	err := s.ag.setExtraOptions(haiconf.CommandArgs{
		"ExtraOptions": []interface{}{"a", "b", "a", "a"},
	})
	c.Assert(err, IsNil)
	c.Assert(s.ag.ExtraOptions, DeepEquals, []string{"a", "b"})
}

func (s *AptGetTestSuite) TestSetUserConfig_Install(c *C) {
	args := haiconf.CommandArgs{
		"Method":       METHOD_INSTALL,
		"Packages":     []interface{}{"a", "b"},
		"ExtraOptions": []interface{}{"foo", "bar"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.ag.Method, Equals, args["Method"])
	c.Assert(s.ag.Packages, DeepEquals, []string{"a", "b"})
	c.Assert(s.ag.ExtraOptions, DeepEquals, []string{"foo", "bar"})
}

func (s *AptGetTestSuite) TestSetUserConfig_Remove(c *C) {
	args := haiconf.CommandArgs{
		"Method":       METHOD_REMOVE,
		"Packages":     []interface{}{"a", "b"},
		"ExtraOptions": []string{"foo", "bar"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.ag.Method, Equals, args["Method"])
	c.Assert(s.ag.Packages, DeepEquals, []string{"a", "b"})
	c.Assert(s.ag.ExtraOptions, DeepEquals, []string(nil))
}

func (s *AptGetTestSuite) TestSetUserConfig_Update(c *C) {
	args := haiconf.CommandArgs{
		"Method":       METHOD_UPDATE,
		"Packages":     []interface{}{"a", "b"},
		"ExtraOptions": []string{"foo", "bar"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.ag.Method, Equals, args["Method"])
	c.Assert(s.ag.Packages, DeepEquals, []string(nil))
	c.Assert(s.ag.ExtraOptions, DeepEquals, []string(nil))
}

func (s *AptGetTestSuite) TestRun_Failed(c *C) {
	args := haiconf.CommandArgs{
		"Method":       "install",
		"Packages":     []interface{}{"azertyuiop"},
		"ExtraOptions": []string{"-o", "DPkg::Options::=--force-confnew"},
	}

	err := s.ag.SetUserConfig(args)
	c.Assert(err, IsNil)

	err = s.ag.Run()
	c.Assert(err, NotNil)
}
