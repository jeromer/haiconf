// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"github.com/jeromer/haiconf/haiconf/osutils"
	. "launchpad.net/gocheck"
	"os"
	"os/user"
)

type CrontabTestSuite struct {
	c *Crontab
}

var (
	_ = Suite(&CrontabTestSuite{})
)

func (s *CrontabTestSuite) SetUpTest(c *C) {
	u, err := user.Current()
	c.Assert(err, IsNil)

	s.c = NewCrontab(u)
}

func (s *CrontabTestSuite) TestRead(c *C) {
	buff, err := s.c.Read()
	c.Assert(err, IsNil)
	c.Assert(len(buff), Equals, 0)
}

func (s *CrontabTestSuite) TestSave(c *C) {
	defer s.cleanCrontab(c)

	cronjobs := []Cronjob{
		Cronjob{
			Env: map[string]string{
				"SHELL":  "/bin/zsh",
				"MAILTO": "foo@example.com",
			},
			Schedule: PREDEFINED_SCHEDULES["hourly"],
			Command:  "/bin/true",
		},

		Cronjob{
			Env: map[string]string{
				"MAILTO": "bar@example.com",
				"SHELL":  "/bin/csh",
			},
			Schedule: PREDEFINED_SCHEDULES["daily"],
			Command:  "/bin/false",
		},
	}

	err := s.c.Save(cronjobs)
	c.Assert(err, IsNil)

	obtained, err := s.c.Read()
	c.Assert(err, IsNil)
	c.Assert(obtained, DeepEquals, cronjobs)
}

func (s *CrontabTestSuite) TestRemove(c *C) {
	defer s.cleanCrontab(c)

	cj1 := Cronjob{
		Schedule: PREDEFINED_SCHEDULES["hourly"],
		Command:  "/bin/false",
		Env: map[string]string{
			"MAILTO": "foo",
			"SHELL":  "/bin/ksh",
			"PATH":   "/foo",
		},
	}

	err := s.c.Add(cj1)
	c.Assert(err, IsNil)

	cj2 := Cronjob{
		Schedule: PREDEFINED_SCHEDULES["daily"],
		Command:  "/bin/false",
		Env: map[string]string{
			"MAILTO": "foo",
			"SHELL":  "/bin/ksh",
			"PATH":   "/foo",
		},
	}

	err = s.c.Add(cj2)
	c.Assert(err, IsNil)

	cronjobs, err := s.c.Read()
	c.Assert(cronjobs, DeepEquals, []Cronjob{cj1, cj2})

	err = s.c.Remove(cj1)
	c.Assert(err, IsNil)

	cronjobs, err = s.c.Read()
	c.Assert(cronjobs, DeepEquals, []Cronjob{cj2})
}

func (s *CrontabTestSuite) TestAdd(c *C) {
	defer s.cleanCrontab(c)

	cj := Cronjob{
		Schedule: []string{"0", "0", "1", "1", "*"},
		Command:  "/bin/false",
		Env: map[string]string{
			"MAILTO": "foo",
			"SHELL":  "/bin/ksh",
			"PATH":   "/foo",
		},
	}

	err := s.c.Add(cj)
	c.Assert(err, IsNil)

	cronjobs, err := s.c.Read()
	c.Assert(cronjobs, DeepEquals, []Cronjob{cj})
}

func (s *CrontabTestSuite) TestRemoveDuplicates(c *C) {
	fixtures := []Cronjob{
		Cronjob{
			Env: map[string]string{
				"SHELL":  "/bin/zsh",
				"MAILTO": "foo@example.com",
			},
			Schedule: []string{"15", "14", "1", "*", "*"},
			Command:  "$HOME/bin/false",
		},

		Cronjob{
			Env: map[string]string{
				"MAILTO": "foo@example.com",
				"SHELL":  "/bin/zsh",
			},
			Schedule: []string{"15", "14", "1", "*", "*"},
			Command:  "$HOME/bin/false",
		},
	}

	expected := []Cronjob{
		Cronjob{
			Env: map[string]string{
				"SHELL":  "/bin/zsh",
				"MAILTO": "foo@example.com",
			},
			Schedule: []string{"15", "14", "1", "*", "*"},
			Command:  "$HOME/bin/false",
		},
	}

	obtained := s.c.RemoveDuplicates(fixtures)

	c.Assert(obtained, DeepEquals, expected)
}

func (s *CrontabTestSuite) cleanCrontab(c *C) {
	sc := osutils.SystemCommand{
		Path:                 s.c.Path,
		Args:                 []string{"-u", s.c.User.Username, "-r"},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	output := sc.Run()
	c.Assert(output.HasError(), Equals, false)
}
