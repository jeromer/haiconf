// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/osutils"
	. "launchpad.net/gocheck"
	"os"
	"os/user"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type CronTestSuite struct {
	c *Cron
}

var (
	_ = Suite(&CronTestSuite{})

	nilSchedule = make([]string, 5)

	dummyRuntimeConfig = haiconf.RuntimeConfig{
		Verbose: false,
		Output:  nil,
	}
)

func (s *CronTestSuite) SetUpTest(c *C) {
	s.c = new(Cron)
	s.c.SetDefault(&dummyRuntimeConfig)
}

func (s *CronTestSuite) TestSetDefault(c *C) {
	expected := &Cron{
		Command:  "",
		Ensure:   haiconf.ENSURE_PRESENT,
		Env:      map[string]string{},
		Schedule: make([]string, 5),
		Owner:    new(user.User),
		rc:       &dummyRuntimeConfig,
	}

	c.Assert(s.c, DeepEquals, expected)
}

func (s *CronTestSuite) TestSetUserConfig_Complete(c *C) {
	o, err := user.Lookup("nobody")
	c.Assert(err, IsNil)

	command := "/foo/bar"

	args := haiconf.CommandArgs{
		"Command": command,
		"Ensure":  haiconf.ENSURE_ABSENT,
		"Env": map[string]interface{}{
			"PATH": "$PATH:/usr/bin/",
		},
		"Schedule": map[string]interface{}{"Predefined": "daily"},
		"Owner":    "nobody",
	}

	err = s.c.SetUserConfig(args)
	c.Assert(err, IsNil)

	c.Assert(s.c.Command, Equals, command)
	c.Assert(s.c.Ensure, Equals, haiconf.ENSURE_ABSENT)
	c.Assert(s.c.Env, DeepEquals, map[string]string{
		"PATH": "$PATH:/usr/bin/",
	})
	c.Assert(s.c.Schedule, DeepEquals, PREDEFINED_SCHEDULES["daily"])
	c.Assert(s.c.Owner, DeepEquals, o)
}

func (s *CronTestSuite) TestSetSchedule_PredefinedInvalid(c *C) {
	args := haiconf.CommandArgs{
		"Schedule": map[string]interface{}{
			"Predefined": "fooo",
		},
	}

	err := s.c.setSchedule(args)
	c.Assert(err, ErrorMatches, "Invalid choice for Predefined. (.*)")
	c.Assert(s.c.Schedule, DeepEquals, nilSchedule)
}

func (s *CronTestSuite) TestSetSchedule_PredefinedValid(c *C) {
	args := haiconf.CommandArgs{
		"Schedule": map[string]interface{}{
			"Predefined": "yearly",
		},
	}

	err := s.c.setSchedule(args)
	c.Assert(err, IsNil)
	c.Assert(s.c.Schedule, DeepEquals, PREDEFINED_SCHEDULES["yearly"])
}

func (s *CronTestSuite) TestSetSchedule_NonPredefinedInvalid(c *C) {
	args := haiconf.CommandArgs{
		"Schedule": map[string]interface{}{
			"Foo": "bar",
		},
	}

	err := s.c.setSchedule(args)
	c.Assert(err, ErrorMatches, "Minute must be provided. (.*)")
	c.Assert(s.c.Schedule, DeepEquals, nilSchedule)
}

func (s *CronTestSuite) TestSetSchedule_NonPredefinedIncomplete(c *C) {
	args := haiconf.CommandArgs{
		"Schedule": map[string]interface{}{
			"Hour":     "*",
			"Minute":   "0",
			"WeekDay":  "*",
			"MonthDay": "0",
		},
	}

	err := s.c.setSchedule(args)
	c.Assert(err, ErrorMatches, "Month must be provided. (.*)")
	c.Assert(s.c.Schedule, DeepEquals, nilSchedule)
}

func (s *CronTestSuite) TestSetSchedule_NonPredefinedValid(c *C) {
	args := haiconf.CommandArgs{
		"Schedule": map[string]interface{}{
			"Hour":     "0",
			"Minute":   "*/2",
			"MonthDay": "1,31",
			"WeekDay":  "0-6",
			"Month":    "*",
		},
	}

	expected := []string{
		"*/2",
		"0",
		"1,31",
		"*",
		"0-6",
	}

	err := s.c.setSchedule(args)
	c.Assert(err, IsNil)
	c.Assert(s.c.Schedule, DeepEquals, expected)
}

func (s *CronTestSuite) TestRun_EnsurePresent(c *C) {
	defer s.cleanCrontab(c)

	u, err := user.Current()
	c.Assert(err, IsNil)

	args := haiconf.CommandArgs{
		"Command": "/foo/bar",
		"Ensure":  haiconf.ENSURE_PRESENT,
		"Env": map[string]interface{}{
			"PATH": "$PATH:/usr/bin/",
		},
		"Schedule": map[string]interface{}{
			"Predefined": "daily",
		},
		"Owner": u.Username,
	}

	err = s.c.SetUserConfig(args)
	c.Assert(err, IsNil)

	err = s.c.Run()
	c.Assert(err, IsNil)

	obtained, err := NewCrontab(u).Read()
	c.Assert(err, IsNil)

	expected := []Cronjob{
		Cronjob{
			Env: map[string]string{
				"PATH": "$PATH:/usr/bin/",
			},
			Command:  "/foo/bar",
			Schedule: PREDEFINED_SCHEDULES["daily"],
		},
	}

	c.Assert(obtained, DeepEquals, expected)
}

func (s *CronTestSuite) TestRun_EnsureAbsent(c *C) {
	defer s.cleanCrontab(c)
	u, err := user.Current()
	c.Assert(err, IsNil)

	args := haiconf.CommandArgs{
		"Command": "/foo/bar",
		"Ensure":  haiconf.ENSURE_ABSENT,
		"Env": map[string]interface{}{
			"PATH": "$PATH:/usr/bin/",
		},
		"Schedule": map[string]interface{}{
			"Predefined": "daily",
		},
		"Owner": u.Username,
	}

	err = s.c.SetUserConfig(args)
	c.Assert(err, IsNil)

	err = s.c.Run()
	c.Assert(err, IsNil)

	obtained, err := NewCrontab(u).Read()
	c.Assert(err, IsNil)

	expected := []Cronjob{}

	c.Assert(obtained, DeepEquals, expected)
}

func (s *CronTestSuite) cleanCrontab(c *C) {
	u, err := user.Current()
	c.Assert(err, IsNil)

	ct := NewCrontab(u)

	sc := osutils.SystemCommand{
		Path:                 ct.Path,
		Args:                 []string{"-u", u.Username, "-r"},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	output := sc.Run()
	c.Assert(output.HasError(), Equals, false)
}
