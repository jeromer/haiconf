package utils

import (
	"github.com/jeromer/haiconf/haiconf"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Usage in a haiconf file
//
// HttpGet({
//    From = "http://some.url/path.to.file.ext",
//    To = "/tmp/my.file.ext",
// })

type HttpGet struct {
	from string
	to   string

	rc *haiconf.RuntimeConfig
}

func (h *HttpGet) SetDefault(rc *haiconf.RuntimeConfig) error {
	h.rc = rc
	return nil
}

func (h *HttpGet) SetUserConfig(args haiconf.CommandArgs) error {
	err := h.setFrom(args)
	if err != nil {
		return err
	}

	err = h.setTo(args)
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpGet) Run() error {
	haiconf.Output(h.rc, "Downloading %s to %s", h.from, h.to)

	resp, err := http.Get(h.from)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(h.to)
	if err != nil {
		return err
	}

	f.Write(buff)
	return f.Close()
}

func (h *HttpGet) setFrom(args haiconf.CommandArgs) error {
	f, _ := haiconf.CheckString("From", args)

	if len(f) == 0 {
		return haiconf.NewArgError("From must be provided", args)
	}

	f = strings.ToLower(f)

	if f[:7] == "http://" || f[:8] == "https://" {
		h.from = f
		return nil
	}

	return haiconf.NewArgError("From must be http or https", args)
}

func (h *HttpGet) setTo(args haiconf.CommandArgs) error {
	t, _ := haiconf.CheckString("To", args)

	if len(t) == 0 {
		return haiconf.NewArgError("To must be provided", args)
	}

	if !hasFileName(t) {
		return haiconf.NewArgError(t+" must be a file name", args)
	}

	dir := filepath.Dir(t)
	if !isDir(dir) {
		return haiconf.NewArgError(dir+" does not exists", args)
	}

	h.to = t

	return nil
}

func hasFileName(f string) bool {
	ext := filepath.Ext(f)
	return len(ext) > 0
}

func isDir(d string) bool {
	fd, err := os.Open(d)
	if err == nil {
		fd.Close()
		return true
	}

	return os.IsExist(err)
}
