package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/jeromer/haiconf/haiconf"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type TarGzTestSuite struct {
	t *TarGz
}

var (
	_ = Suite(&TarGzTestSuite{})

	dummyRuntimeConfig = haiconf.RuntimeConfig{
		Verbose: false,
		Output:  nil,
	}
)

func (s *TarGzTestSuite) SetUpTest(c *C) {
	s.t = new(TarGz)
	err := s.t.SetDefault(&dummyRuntimeConfig)
	c.Assert(err, IsNil)
}

func (s *TarGzTestSuite) TearDownTest(c *C) {
	s.t.source = ""
	s.t.dest = ""
}

func (s *TarGzTestSuite) TestSetSource_Empty(c *C) {
	args := haiconf.CommandArgs{
		"Source": "",
	}

	err := s.t.setSource(args)
	c.Assert(err, ErrorMatches, "Source must be provided. (.*)")
	c.Assert(s.t.source, Equals, "")
}

func (s *TarGzTestSuite) TestSetSource_FileDoesNotExist(c *C) {
	args := haiconf.CommandArgs{
		"Source": "/foo/bar/tarball.tar.gz",
	}

	err := s.t.setSource(args)
	c.Assert(err, ErrorMatches, "Source does not exist. (.*)")
	c.Assert(s.t.source, Equals, "")
}

func (s *TarGzTestSuite) TestSetSource_FileExists(c *C) {
	tmpDir := c.MkDir()
	args := haiconf.CommandArgs{
		"Source": tmpDir,
	}

	err := s.t.setSource(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.source, Equals, tmpDir)
}

func (s *TarGzTestSuite) TestSetDest_Empty(c *C) {
	args := haiconf.CommandArgs{
		"Dest": "",
	}

	err := s.t.setDest(args)
	c.Assert(err, ErrorMatches, "Dest must be provided. (.*)")
	c.Assert(s.t.dest, Equals, "")
}

func (s *TarGzTestSuite) TestSetDest_DirDoesNotExist(c *C) {
	tmpDir := c.MkDir() + "/foo"
	args := haiconf.CommandArgs{
		"Dest": tmpDir + "/tarball.tar.gz",
	}

	err := s.t.setDest(args)
	c.Assert(err, ErrorMatches, "stat "+tmpDir+": no such file or directory")
	c.Assert(s.t.dest, Equals, "")
}

func (s *TarGzTestSuite) TestSetDest_Exists(c *C) {
	tmpDir := c.MkDir()
	tarball := "tarball.tar.gz"
	fullPath := tmpDir + "/" + tarball
	args := haiconf.CommandArgs{
		"Dest": fullPath,
	}

	err := s.t.setDest(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.dest, Equals, fullPath)
}

func (s *TarGzTestSuite) TestSetDest_NoTarballName(c *C) {
	args := haiconf.CommandArgs{
		"Dest": c.MkDir(),
	}

	err := s.t.setDest(args)
	c.Assert(err, ErrorMatches, "No tarball name provided. (.*)")
	c.Assert(s.t.dest, Equals, "")
}

func (s *TarGzTestSuite) TestSetUserConfig_Complete(c *C) {
	tmpDir := c.MkDir()

	source := tmpDir
	dest := tmpDir + "/foo.tar.gz"

	args := haiconf.CommandArgs{
		"Source": source,
		"Dest":   dest,
	}

	err := s.t.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.source, Equals, source)
	c.Assert(s.t.dest, Equals, dest)
}

func (s *TarGzTestSuite) TestCreateTar(c *C) {
	buff, err := createTar("./fixtures/")
	c.Assert(err, IsNil)
	c.Assert(len(buff) > 0, Equals, true)

	rdr := bytes.NewReader(buff)
	tr := tar.NewReader(rdr)

	var obtained []*tar.Header
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		obtained = append(obtained, hdr)
	}

	expected := []*tar.Header{
		&tar.Header{
			Name:     "./fixtures/",
			Typeflag: tar.TypeDir,
		},
		&tar.Header{
			Name:     "fixtures/dir",
			Typeflag: tar.TypeDir,
		},
		&tar.Header{
			Name:     "fixtures/dir/file.txt",
			Typeflag: tar.TypeReg,
		},
		&tar.Header{
			Name:     "fixtures/dir/hardlink",
			Typeflag: tar.TypeReg,
		},
		&tar.Header{
			Name:     "fixtures/dir/symlink",
			Typeflag: tar.TypeSymlink,
		},
	}

	c.Assert(len(obtained), Equals, len(expected))
	for i, hdr := range obtained {
		c.Assert(hdr.Name, Equals, expected[i].Name)
		c.Assert(hdr.Typeflag, Equals, expected[i].Typeflag)
	}
}

func (s *TarGzTestSuite) TestGz(c *C) {
	expected := []byte("foo")

	gzipped, err := gz(expected)
	c.Assert(err, IsNil)
	c.Assert(len(gzipped) > 0, Equals, true)

	reader, err := gzip.NewReader(bytes.NewBuffer(gzipped))
	c.Assert(err, IsNil)
	defer reader.Close()

	buff, err := ioutil.ReadAll(reader)
	c.Assert(err, IsNil)
	c.Assert(buff, DeepEquals, expected)
}

func (s *TarGzTestSuite) TestRun_TarGz(c *C) {
	source := "./fixtures"
	dest := c.MkDir() + "/fixtures.tar.gz"

	args := haiconf.CommandArgs{
		"Source": source,
		"Dest":   dest,
	}

	err := s.t.SetUserConfig(args)
	c.Assert(err, IsNil)
	c.Assert(s.t.source, Equals, source)
	c.Assert(s.t.dest, Equals, dest)

	err = s.t.Run()
	c.Assert(err, IsNil)

	f, err := os.Open(dest)
	c.Assert(err, IsNil)
	defer f.Close()
}
