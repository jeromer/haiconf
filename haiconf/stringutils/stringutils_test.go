// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stringutils

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type StringUtilsTestSuite struct {
}

var _ = Suite(&StringUtilsTestSuite{})

func (s *StringUtilsTestSuite) TestRemoveDuplicates_NoDuplicates(c *C) {
	obtained := RemoveDuplicates([]string{"a", "b"})
	expected := []string{"a", "b"}
	c.Assert(obtained, DeepEquals, expected)
}

func (s *StringUtilsTestSuite) TestRemoveDuplicates_Duplicates(c *C) {
	obtained := RemoveDuplicates([]string{"a", "b", "b", "a"})
	expected := []string{"a", "b"}
	c.Assert(obtained, DeepEquals, expected)
}
