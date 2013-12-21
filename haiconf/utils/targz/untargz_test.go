// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package targz

import (
	//"bytes"
	//"compress/gzip"
	"github.com/dotcloud/tar"
	"github.com/jeromer/haiconf/haiconf"
	//"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"runtime"
)

type UnTarGzTestSuite struct {
	t *UnTarGz
}

var (
	_ = Suite(&UnTarGzTestSuite{})
)

func (s *UnTarGzTestSuite) SetUpTest(c *C) {
	s.t = new(UnTarGz)
	err := s.t.SetDefault(&dummyRuntimeConfig)
	c.Assert(err, IsNil)
}

func (s *UnTarGzTestSuite) TearDownTest(c *C) {
	s.t.Source = ""
	s.t.Dest = ""
}

func (s *UnTarGzTestSuite) TestSetSource_Empty(c *C) {
	args := haiconf.CommandArgs{
		"Source": "",
	}

	err := s.t.setSource(args)
	c.Assert(err, ErrorMatches, "Source must be provided. (.*)")
	c.Assert(s.t.Source, Equals, "")
}

func (s *UnTarGzTestSuite) TestSetSource_FileDoesNotExist(c *C) {
	args := haiconf.CommandArgs{
		"Source": "/foo/bar/tarball.tar.gz",
	}

	err := s.t.setSource(args)
	c.Assert(err, ErrorMatches, "stat /foo/bar/tarball.tar.gz: no such file or directory")
	c.Assert(s.t.Source, Equals, "")
}

func (s *UnTarGzTestSuite) TestSetSource_FileExists(c *C) {
	tmpFile := c.MkDir() + "/foo.tar.gz"

	f, err := os.Create(tmpFile)
	c.Assert(err, IsNil)
	f.Close()

	args := haiconf.CommandArgs{
		"Source": tmpFile,
	}

	err = s.t.setSource(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.Source, Equals, tmpFile)
}

func (s *UnTarGzTestSuite) TestSetDest_Empty(c *C) {
	args := haiconf.CommandArgs{
		"Dest": "",
	}

	err := s.t.setDest(args)
	c.Assert(err, ErrorMatches, "Dest must be provided. (.*)")
	c.Assert(s.t.Dest, Equals, "")
}

func (s *UnTarGzTestSuite) TestSetDest_DirDoesNotExist(c *C) {
	tmpDir := c.MkDir() + "/foo"
	args := haiconf.CommandArgs{
		"Dest": tmpDir,
	}

	err := s.t.setDest(args)
	c.Assert(err, ErrorMatches, "stat "+tmpDir+": no such file or directory")
	c.Assert(s.t.Dest, Equals, "")
}

func (s *UnTarGzTestSuite) TestSetDest_Exists(c *C) {
	tmpDir := c.MkDir()
	args := haiconf.CommandArgs{
		"Dest": tmpDir,
	}

	err := s.t.setDest(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.Dest, Equals, tmpDir)
}

func (s *UnTarGzTestSuite) TestSetUserConfig_Complete(c *C) {
	source := "./fixtures.tar.gz"
	dest := c.MkDir()

	args := haiconf.CommandArgs{
		"Source": source,
		"Dest":   dest,
	}

	err := s.t.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.Source, Equals, source)
	c.Assert(s.t.Dest, Equals, dest)
}

func (s *UnTarGzTestSuite) TestGunzip(c *C) {
	gunzipped, err := gunzip("./fixtures.tar.gz")
	c.Assert(err, IsNil)

	fi, err := os.Stat("./fixtures.tar")
	c.Assert(err, IsNil)
	c.Assert(int64(len(gunzipped)), Equals, fi.Size())
}

func (s *UnTarGzTestSuite) TestUntar(c *C) {
	buff, err := ioutil.ReadFile("./fixtures.tar")
	c.Assert(err, IsNil)
	c.Assert(len(buff) > 0, Equals, true)

	obtained, err := untar(buff)
	c.Assert(err, IsNil)
	c.Assert(len(obtained) > 0, Equals, true)

	fileContents, err := ioutil.ReadFile("./fixtures/dir/file.txt")
	c.Assert(err, IsNil)

	expected := []tarItem{
		tarItem{
			header: &tar.Header{
				Name:     "fixtures/",
				Typeflag: tar.TypeDir,
			},
			body: []byte(""),
		},
		tarItem{
			header: &tar.Header{
				Name:     "fixtures/dir/",
				Typeflag: tar.TypeDir,
			},
			body: []byte(""),
		},
		tarItem{
			header: &tar.Header{
				Name:     "fixtures/dir/file.txt",
				Typeflag: tar.TypeReg,
			},
			body: fileContents,
		},
		tarItem{
			header: &tar.Header{
				Name:     "fixtures/dir/hardlink",
				Typeflag: tar.TypeReg,
			},
			body: fileContents,
		},
		tarItem{
			header: &tar.Header{
				Name:     "fixtures/dir/symlink",
				Typeflag: tar.TypeSymlink,
			},
			body: []byte(""),
		},
	}

	c.Assert(len(obtained), Equals, len(expected))
	for i, ti := range obtained {
		c.Assert(ti.header.Name, Equals, expected[i].header.Name)
		c.Assert(ti.header.Typeflag, Equals, expected[i].header.Typeflag)
		c.Assert(ti.body, DeepEquals, expected[i].body)
	}
}

func (s *UnTarGzTestSuite) TestWriteFiles(c *C) {
	dest := c.MkDir()
	/*
		dest = "/tmp/testwritefiles"
		e := os.RemoveAll(dest)
		c.Assert(e, IsNil)
		e = os.Mkdir(dest, 0755)
		c.Assert(e, IsNil)
	*/

	dir := "somedir"
	file := "somefile"
	link := "symlink"

	fileContents, err := ioutil.ReadFile("./fixtures/dir/file.txt")
	c.Assert(err, IsNil)

	linkMode := int64(0750)
	if runtime.GOOS == "linux" {
		linkMode = 0777
	}

	items := []tarItem{
		tarItem{
			header: &tar.Header{
				Name:     dir,
				Typeflag: tar.TypeDir,
				Mode:     0750,
			},
		},
		tarItem{
			header: &tar.Header{
				Name:     file,
				Typeflag: tar.TypeReg,
				Mode:     0640,
			},
			body: fileContents,
		},
		tarItem{
			header: &tar.Header{
				Name:     link,
				Typeflag: tar.TypeSymlink,
				Mode:     linkMode,
				Linkname: "somefile",
			},
		},
	}

	err = writeFiles(items, dest)
	c.Assert(err, IsNil)

	for _, it := range items {
		stat, err := os.Lstat(dest + "/" + it.header.Name)
		c.Assert(err, IsNil)

		fi := it.header.FileInfo()
		c.Assert(stat.Name(), Equals, fi.Name())
		c.Assert(stat.Mode(), Equals, fi.Mode())
	}

	obtained, err := ioutil.ReadFile(dest + "/" + file)
	c.Assert(err, IsNil)
	c.Assert(obtained, DeepEquals, fileContents)
}

func (s *UnTarGzTestSuite) TestRun_UnTarGz(c *C) {
	source := "./fixtures.tar.gz"

	dest := c.MkDir()
	/*
		dest = "/tmp/testrun_untargz"
		e := os.RemoveAll(dest)
		c.Assert(e, IsNil)
		e = os.Mkdir(dest, 0755)
		c.Assert(e, IsNil)
	*/

	args := haiconf.CommandArgs{
		"Source": source,
		"Dest":   dest,
	}

	err := s.t.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.Source, Equals, source)
	c.Assert(s.t.Dest, Equals, dest)

	err = s.t.Run()
	c.Assert(err, IsNil)

	f, err := os.Open(dest)
	c.Assert(err, IsNil)
	defer f.Close()

	n, err := f.Readdirnames(-1)
	c.Assert(err, IsNil)
	c.Assert(len(n) > 0, Equals, true)
}
