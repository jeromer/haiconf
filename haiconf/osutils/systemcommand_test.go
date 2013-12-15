// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package osutils

import (
	. "launchpad.net/gocheck"
	"os"
)

type SystemCommandTestSuite struct{}

var _ = Suite(&SystemCommandTestSuite{})

func (s *SystemCommandTestSuite) TestBuildEnvVars(c *C) {
	sc := &SystemCommand{
		EnvVars: map[string]string{
			"FOO": "bar",
			"BAR": "baz",
		},
	}

	obtained := sc.buildEnvVars()
	expected := []string{
		"FOO=bar",
		"BAR=baz",
		"PATH=" + os.Getenv("PATH"),
	}
	c.Assert(obtained, DeepEquals, expected)
}

func (s *SystemCommandTestSuite) TestBuildCmd_ShellExpansionDisabled(c *C) {
	path := "/foo/bar"
	args := []string{"a", "b"}

	sc := &SystemCommand{
		Path:                 path,
		Args:                 args,
		EnableShellExpansion: false,
	}

	cmd := sc.buildCmd()
	c.Assert(cmd.Path, Equals, path)
	c.Assert(cmd.Args, DeepEquals, args)
}

func (s *SystemCommandTestSuite) TestBuildCmd_ShellExpansionEnabled(c *C) {
	path := "/foo/bar"
	args := []string{"a", "b"}

	sc := &SystemCommand{
		Path:                 path,
		Args:                 args,
		EnableShellExpansion: true,
	}

	cmd := sc.buildCmd()
	c.Assert(cmd.Path, Equals, "/bin/sh")
	c.Assert(cmd.Args, DeepEquals, []string{"sh", "-c", "/foo/bar a b"})
}

func (s *SystemCommandTestSuite) TestRun_CommandFailedWrongPath(c *C) {
	sc := &SystemCommand{
		Path:                 "/path/to/inexistant/command",
		Args:                 []string{"a", "b"},
		EnableShellExpansion: false,
	}

	err := sc.Run()
	c.Assert(err, NotNil)
	expected := `Error with command "/path/to/inexistant/command a b". Error message was "fork/exec /path/to/inexistant/command: no such file or directory".`
	c.Assert(err.Error(), Equals, expected)
}

func (s *SystemCommandTestSuite) TestRun_CommandFailed(c *C) {
	sc := &SystemCommand{
		Path:                 "/usr/bin/tr",
		Args:                 []string{"--xxx"},
		EnableShellExpansion: false,
	}

	err := sc.Run()
	// TODO : improve error msg check
	c.Assert(err, NotNil)
}
