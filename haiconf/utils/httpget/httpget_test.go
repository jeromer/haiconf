package httpget

import (
	"github.com/jeromer/haiconf/haiconf"
	. "launchpad.net/gocheck"
	"os"
	"strings"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type HttpGetTestSuite struct {
	h *HttpGet
}

var (
	_ = Suite(&HttpGetTestSuite{})

	dummyRuntimeConfig = haiconf.RuntimeConfig{
		Verbose: false,
		Output:  nil,
	}
)

func (s *HttpGetTestSuite) SetUpTest(c *C) {
	s.h = new(HttpGet)
	err := s.h.SetDefault(&dummyRuntimeConfig)
	c.Assert(err, IsNil)
}

func (s *HttpGetTestSuite) TearDownTest(c *C) {
	s.h.from = ""
	s.h.to = ""
}

func (s *HttpGetTestSuite) TestSetDefault(c *C) {
	// XXX : s.d.SetDefault() called in Setuptest
	c.Assert(s.h.rc, Equals, &dummyRuntimeConfig)
}

func (s *HttpGetTestSuite) TestSetFrom_Empty(c *C) {
	args := haiconf.CommandArgs{
		"From": "",
	}

	err := s.h.setFrom(args)
	c.Assert(err, ErrorMatches, "From must be provided. (.*)")
	c.Assert(s.h.from, Equals, "")
}

func (s *HttpGetTestSuite) TestSetFrom_NotHttp(c *C) {
	args := haiconf.CommandArgs{
		"From": "foo://example.com",
	}

	err := s.h.setFrom(args)
	c.Assert(err, ErrorMatches, "From must be http or https. (.*)")
	c.Assert(s.h.from, Equals, "")
}

func (s *HttpGetTestSuite) TestSetFrom_IsHttp(c *C) {
	url := "HTTP://example.com"
	args := haiconf.CommandArgs{
		"From": url,
	}

	err := s.h.setFrom(args)
	c.Assert(err, IsNil)
	c.Assert(s.h.from, Equals, strings.ToLower(url))
}

func (s *HttpGetTestSuite) TestSetFrom_IsHttpS(c *C) {
	url := "HTTPS://example.com"

	args := haiconf.CommandArgs{
		"From": url,
	}

	err := s.h.setFrom(args)
	c.Assert(err, IsNil)
	c.Assert(s.h.from, Equals, strings.ToLower(url))
}

func (s *HttpGetTestSuite) TestSetTo_Empty(c *C) {
	args := haiconf.CommandArgs{
		"To": "",
	}

	err := s.h.setTo(args)
	c.Assert(err, ErrorMatches, "To must be provided. (.*)")
	c.Assert(s.h.to, Equals, "")
}

func (s *HttpGetTestSuite) TestSetTo_TargetDirDoesNotExists(c *C) {
	args := haiconf.CommandArgs{
		"To": "/foo/bar/baz.html",
	}

	err := s.h.setTo(args)
	c.Assert(err, ErrorMatches, "stat /foo/bar: no such file or directory. (.*)")
	c.Assert(s.h.to, Equals, "")
}

func (s *HttpGetTestSuite) TestSetTo_NotAFileName(c *C) {
	to := c.MkDir()

	args := haiconf.CommandArgs{
		"To": to,
	}

	err := s.h.setTo(args)
	c.Assert(err, ErrorMatches, to+" must be a file name. (.*)")
	c.Assert(s.h.to, Equals, "")
}

func (s *HttpGetTestSuite) TestSetUserConfig_Complete(c *C) {
	tmpDir := c.MkDir()
	from := "http://example.com/"
	to := tmpDir + "/example.html"

	args := haiconf.CommandArgs{
		"From": from,
		"To":   to,
	}

	err := s.h.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.h.from, Equals, from)
	c.Assert(s.h.to, Equals, to)
}

func (s *HttpGetTestSuite) TestRun(c *C) {
	from := "http://example.com/"
	to := c.MkDir() + "/example.html"

	args := haiconf.CommandArgs{
		"From": from,
		"To":   to,
	}

	err := s.h.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.h.from, Equals, from)
	c.Assert(s.h.to, Equals, to)

	err = s.h.Run()
	c.Assert(err, IsNil)

	err = s.h.Run()
	c.Assert(err, IsNil)

	f, err := os.Open(to)
	c.Assert(err, IsNil)
	defer f.Close()

	buff := make([]byte, 10)
	bytesRead, err := f.Read(buff)
	c.Assert(err, IsNil)
	c.Assert(bytesRead > 0, Equals, true)
}
