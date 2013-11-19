package osutils

import (
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

var (
	dummyEnVars = map[string]string{
		"FOO": "none",
		"BAR": "quux",
	}
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type OsUtilsTestSuite struct{}

var _ = Suite(&OsUtilsTestSuite{})

func (s *OsUtilsTestSuite) SetUpTest(c *C) {
	var err error

	for name, value := range dummyEnVars {
		err = os.Setenv(name, value)
		if err != nil {
			panic(err)
		}
	}
}

func (s *OsUtilsTestSuite) TearDownTest(c *C) {
	var err error

	for name, _ := range dummyEnVars {
		err = os.Setenv(name, "")
		if err != nil {
			panic(err)
		}
	}
}

func (s *OsUtilsTestSuite) Test_GetEnvList(c *C) {
	bucket := GetEnvList([]string{"FOO", "BAR"})

	c.Assert(bucket, DeepEquals, dummyEnVars)
}

func (s *OsUtilsTestSuite) Test_SetEnvList(c *C) {
	envList := map[string]string{
		"FOO": "a",
		"BAR": "b",
	}

	err := SetEnvList(envList)
	c.Assert(err, IsNil)

	bucket := GetEnvList([]string{"FOO", "BAR"})
	c.Assert(bucket, DeepEquals, envList)
}
