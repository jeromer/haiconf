// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	. "launchpad.net/gocheck"
	"os"
)

type CommonTestSuite struct{}

var _ = Suite(&CommonTestSuite{})

func (s *CommonTestSuite) TestMkDir_NonRecursive(c *C) {
	tmpDir := c.MkDir() + "/foo"

	err := MkDir(tmpDir, false, 0755)
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)
}

func (s *CommonTestSuite) TestMkDir_Recursive(c *C) {
	tmpDir := c.MkDir() + "/foo/bar/baz"

	err := MkDir(tmpDir, true, 0755)
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)
}

func (s *CommonTestSuite) TestRmDir_NonRecursive(c *C) {
	tmpDir := c.MkDir() + "/foo"

	err := MkDir(tmpDir, false, 0755)
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)

	err = RmDir(tmpDir, false)
	c.Assert(err, IsNil)

	f, err = os.Stat(tmpDir)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
}

func (s *CommonTestSuite) TestRmDir_Recursive(c *C) {
	suffix := "/foo/bar/baz"
	tmpDir := c.MkDir()

	err := MkDir(tmpDir+suffix, true, 0755)
	c.Assert(err, IsNil)

	f, err := os.Stat(tmpDir)
	c.Assert(err, IsNil)
	c.Assert(f.IsDir(), Equals, true)

	err = RmDir(tmpDir, true)
	c.Assert(err, IsNil)

	f, err = os.Stat(tmpDir + suffix)
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(f, IsNil)
}
