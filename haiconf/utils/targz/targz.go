package targz

import (
	"bytes"
	"compress/gzip"
	"github.com/dotcloud/tar"
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Usage in a haiconf file
//
// TarGz({
//    Source = "/path/to/dir/",
//    Dest = "/path/to/tarball.tar.gz",
// })

type TarGz struct {
	source string
	dest   string

	rc *haiconf.RuntimeConfig
}

func (t *TarGz) SetDefault(rc *haiconf.RuntimeConfig) error {
	t.rc = rc
	return nil
}

func (t *TarGz) SetUserConfig(args haiconf.CommandArgs) error {
	err := t.setSource(args)
	if err != nil {
		return err
	}

	err = t.setDest(args)
	if err != nil {
		return err
	}

	return nil
}

func (t *TarGz) Run() error {
	haiconf.Output(t.rc, "Archiving %s to %s", t.source, t.dest)

	return tarGz(t.source, t.dest)
}

func (t *TarGz) setSource(args haiconf.CommandArgs) error {
	s, _ := haiconf.CheckString("Source", args)
	if len(s) == 0 {
		return haiconf.NewArgError("Source must be provided", args)
	}

	f, err := os.Open(s)
	if err != nil {
		if os.IsNotExist(err) {
			return haiconf.NewArgError("Source does not exist", args)
		}
	}
	defer f.Close()

	t.source = s
	return nil
}

func (t *TarGz) setDest(args haiconf.CommandArgs) error {
	d, _ := haiconf.CheckString("Dest", args)
	if len(d) == 0 {
		return haiconf.NewArgError("Dest must be provided", args)
	}

	if !strings.HasSuffix(d, ".tar.gz") {
		return haiconf.NewArgError("No tarball name provided", args)
	}

	id, err := utils.IsDir(filepath.Dir(d))
	if err != nil {
		return err
	}

	if !id {
		return haiconf.NewArgError(d+" is not a directory", args)
	}

	t.dest = d
	return nil
}

func tarGz(source string, dest string) error {
	archive, err := createTar(source)
	if err != nil {
		return err
	}

	gzBuff, err := gz(archive)
	if err != nil {
		return err
	}

	return writeFile(gzBuff, dest)
}

func createTar(source string) ([]byte, error) {
	buff := new(bytes.Buffer)
	nilBuff := []byte(nil)
	var err error

	tarWriter := tar.NewWriter(buff)

	// texas ranger :D
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fi, err := os.Lstat(path)
		if err != nil {
			return err
		}

		mode := fi.Mode()

		link := ""
		if mode&os.ModeSymlink != 0 {
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		hdr.Name = path

		err = tarWriter.WriteHeader(hdr)
		if err != nil {
			return err
		}

		if !mode.IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		fbuff, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		f.Close()

		_, err = tarWriter.Write(fbuff)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(source, walker)
	if err != nil {
		return nilBuff, err
	}

	err = tarWriter.Close()
	if err != nil {
		return nilBuff, err
	}

	return buff.Bytes(), nil
}

func gz(buff []byte) ([]byte, error) {
	gzBuff := new(bytes.Buffer)
	nilBuff := []byte(nil)

	writer := gzip.NewWriter(gzBuff)
	_, err := writer.Write(buff)
	if err != nil {
		return nilBuff, err
	}

	err = writer.Flush()
	if err != nil {
		return nilBuff, err
	}

	err = writer.Close()
	if err != nil {
		return nilBuff, err
	}

	return gzBuff.Bytes(), nil
}

func writeFile(buff []byte, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buff)
	return err
}
