package targz

import (
	"bytes"
	"compress/gzip"
	"github.com/dotcloud/tar"
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/utils"
	"io"
	"io/ioutil"
	"os"
)

// Usage in a haiconf file
//
// UnTarGz({
//    Source = "/path/to/tarball.tar.gz",
//    Dest = "/path/to/dir",
// })

type UnTarGz struct {
	source string
	dest   string

	rc *haiconf.RuntimeConfig
}

type tarItem struct {
	header *tar.Header
	body   []byte
}

func (t *UnTarGz) SetDefault(rc *haiconf.RuntimeConfig) error {
	t.rc = rc
	return nil
}

func (t *UnTarGz) SetUserConfig(args haiconf.CommandArgs) error {
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

func (t *UnTarGz) Run() error {
	haiconf.Output(t.rc, "Extracting %s to %s", t.source, t.dest)

	return unTarGz(t.source, t.dest)
}

func (t *UnTarGz) setSource(args haiconf.CommandArgs) error {
	s, _ := haiconf.CheckString("Source", args)
	if len(s) == 0 {
		return haiconf.NewArgError("Source must be provided", args)
	}

	fi, err := os.Stat(s)
	if err != nil {
		return err
	}

	if !fi.Mode().IsRegular() {
		return haiconf.NewArgError(s+" is not a file", args)
	}

	t.source = s
	return nil
}

func (t *UnTarGz) setDest(args haiconf.CommandArgs) error {
	d, _ := haiconf.CheckString("Dest", args)
	if len(d) == 0 {
		return haiconf.NewArgError("Dest must be provided", args)
	}

	id, err := utils.IsDir(d)
	if err != nil {
		return err
	}

	if !id {
		return haiconf.NewArgError(d+" is not a directory", args)
	}

	t.dest = d
	return nil
}

func unTarGz(source string, dest string) error {
	gunzipped, err := gunzip(source)
	if err != nil {
		return err
	}

	archive, err := untar(gunzipped)
	if err != nil {
		return err
	}

	return writeFiles(archive, dest)
}

func untar(buff []byte) ([]tarItem, error) {
	tr := tar.NewReader(bytes.NewReader(buff))

	var items []tarItem

	buff = make([]byte, 32*1024)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		item := tarItem{
			header: hdr,
			body:   []byte{},
		}

		if err != nil {
			return []tarItem{}, err
		}

		nr, err := tr.Read(buff)
		if err != nil && err != io.EOF {
			return []tarItem{}, err
		}

		if nr > 0 {
			item.body = buff[0:nr]
		}

		items = append(items, item)
	}

	return items, nil
}

func gunzip(source string) ([]byte, error) {
	nilBuff := []byte(nil)

	f, err := os.Open(source)
	if err != nil {
		return nilBuff, err
	}
	defer f.Close()

	reader, err := gzip.NewReader(f)
	if err != nil {
		return nilBuff, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func writeFiles(items []tarItem, dest string) error {
	for _, it := range items {
		typeFlag := it.header.Typeflag
		mode := os.FileMode(it.header.Mode)
		name := dest + "/" + it.header.Name

		if typeFlag == tar.TypeDir {
			err := os.Mkdir(name, mode.Perm())
			if err != nil {
				return err
			}

			continue
		}

		if typeFlag == tar.TypeReg || typeFlag == tar.TypeRegA {
			f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}

			_, err = f.Write(it.body)
			if err != nil {
				f.Close()
				return err
			}
			f.Close()

			continue
		}

		if typeFlag == tar.TypeSymlink {
			err := os.Symlink(it.header.Linkname, name)
			if err != nil {
				return err
			}

			continue
		}
	}

	return nil
}
