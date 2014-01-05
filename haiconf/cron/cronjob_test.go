// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"fmt"
	. "launchpad.net/gocheck"
	"strings"
)

type CronjobTestSuite struct {
}

var (
	_ = Suite(&CronjobTestSuite{})
)

func (s *CronjobTestSuite) TestBuildJobLine(c *C) {
	cj := Cronjob{
		Schedule: []string{"0", "0", "1", "1", "*"},
		Command:  "/bin/false",
	}

	obtained := cj.buildJobLine()
	expected := "0 0 1 1 * /bin/false"

	c.Assert(obtained, Equals, expected)
}

func (s *CronjobTestSuite) TestBuildEnvVarsBlock_NoEnvVars(c *C) {
	cj := Cronjob{}
	obtained := cj.buildEnvVarsBlock()
	c.Assert(obtained, Equals, "")
}

func (s *CronjobTestSuite) TestBuildEnvVarsBlock_EnvVars(c *C) {
	cj := Cronjob{
		Env: map[string]string{
			"MAILTO": "foo",
			"SHELL":  "/bin/zsh",
			"PATH":   "/usr/local/bin",
		},
	}

	obtained := cj.buildEnvVarsBlock()
	expected := []string{
		"MAILTO=foo",
		"SHELL=/bin/zsh",
		"PATH=/usr/local/bin",
	}

	c.Assert(obtained, Equals, strings.Join(expected, "\n"))
}

func (s *CronjobTestSuite) TestBuildBlock_NoEnvVars(c *C) {
	cj := Cronjob{
		Schedule: []string{"0", "0", "1", "1", "*"},
		Command:  "/bin/false",
	}

	obtained := cj.BuildBlock()
	expected := "0 0 1 1 * /bin/false\n"
	c.Assert(obtained, Equals, expected)
}

func (s *CronjobTestSuite) TestBuildBlock_EnvVars(c *C) {
	cj := Cronjob{
		Schedule: []string{"0", "0", "1", "1", "*"},
		Command:  "/bin/false",
		Env: map[string]string{
			"MAILTO": "foo",
		},
	}

	obtained := cj.BuildBlock()
	expected := fmt.Sprintf("%s\n%s\n", "MAILTO=foo", "0 0 1 1 * /bin/false")
	c.Assert(obtained, Equals, expected)
}

func (s *CronjobTestSuite) TestHash(c *C) {
	cj := Cronjob{
		Schedule: []string{"0", "0", "1", "1", "*"},
		Command:  "/bin/false",
		Env: map[string]string{
			"MAILTO": "foo",
			"SHELL":  "",
			"PATH":   "/foo",
		},
	}

	h := cj.Hash()
	c.Assert(len(h), Equals, 40)
}
