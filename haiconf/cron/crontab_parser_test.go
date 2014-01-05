// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"io/ioutil"
	. "launchpad.net/gocheck"
)

type CrontabParserTestSuite struct {
	cp *CrontabParser
	cl *CrontabLine
}

var (
	_ = Suite(&CrontabParserTestSuite{})
)

func (s *CrontabParserTestSuite) SetUpTest(c *C) {
	s.cp = new(CrontabParser)
}

func (s *CrontabParserTestSuite) TestParse(c *C) {
	s.cp.Buff = s.readFixture("./fixtures/crontab/complete", c)
	obtained, err := s.cp.Parse()
	c.Assert(err, IsNil)

	expected := []Cronjob{
		Cronjob{
			Env:      map[string]string{},
			Schedule: []string{"15", "14", "1", "*", "*"},
			Command:  "/bin/false >> /dev/null 2>&1",
		},

		Cronjob{
			Env:      map[string]string{"MAILTO": "\"\""},
			Schedule: PREDEFINED_SCHEDULES["daily"],
			Command:  "/bin/false",
		},

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
				"SHELL":  "/bin/csh",
				"MAILTO": "bar@example.com",
				"PATH":   "/path/to/foo/bin",
			},
			Schedule: []string{"*", "*", "*", "*", "0"},
			Command:  "/bin/false >> /dev/null 2>&1",
		},

		Cronjob{
			Env: map[string]string{
				"FOO": "/bin/csh",
				"BAR": "bar@example.com",
				"BAZ": "/path/to/foo/bin",
			},
			Schedule: []string{"*", "*", "*", "0", "*"},
			Command:  "/bin/true >> /dev/null 2>&1",
		},
	}

	c.Assert(len(obtained), Equals, len(expected))
	c.Assert(obtained, DeepEquals, expected)
}

func (s *CrontabParserTestSuite) readFixture(fileName string, c *C) []byte {
	buff, err := ioutil.ReadFile(fileName)
	c.Assert(err, IsNil)
	c.Assert(len(buff) > 0, Equals, true)

	return buff
}

func (s *CrontabParserTestSuite) TestCrontabLine_IsComment(c *C) {
	cl := NewCrontabLine([]byte("     # xxxxxx"))
	c.Assert(cl.IsComment(), Equals, true)
}

func (s *CrontabParserTestSuite) TestCrontabLine_IsBlank(c *C) {
	cl := NewCrontabLine([]byte("		"))
	c.Assert(cl.IsBlank(), Equals, true)
}

func (s *CrontabParserTestSuite) TestCrontabLine_IsEnvVar(c *C) {
	fixtures := [][]byte{
		[]byte("Foo=bar"),
		[]byte("FOO=bar"),
	}

	for _, f := range fixtures {
		cl := NewCrontabLine(f)
		c.Assert(cl.IsEnvVar(), Equals, true)
	}
}

func (s *CrontabParserTestSuite) TestCrontabLine_IsCrontabLine(c *C) {
	fixtures := [][]byte{
		[]byte("@daily aaaaa"),
		[]byte("15 14 1 * * yyyy"),
		[]byte("* * * * 0 xxxx"),
	}

	for _, f := range fixtures {
		cl := NewCrontabLine(f)
		c.Assert(cl.IsCrontabLine(), Equals, true)
	}
}

func (s *CrontabParserTestSuite) TestCrontabLine_SplitCrontabLine(c *C) {
	cmd := "/bin/false >> /dev/null 2>&1"

	// ---

	cl := NewCrontabLine([]byte("15 14 1 * * " + cmd))

	schedule, command := cl.SplitCrontabLine()
	c.Assert(schedule, DeepEquals, []string{"15", "14", "1", "*", "*"})
	c.Assert(command, Equals, cmd)

	// ---

	cl = NewCrontabLine([]byte("@daily " + cmd))

	schedule, command = cl.SplitCrontabLine()
	c.Assert(schedule, DeepEquals, PREDEFINED_SCHEDULES["daily"])
	c.Assert(command, Equals, cmd)
}
