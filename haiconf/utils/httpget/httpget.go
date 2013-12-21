// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpget

import (
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/utils"
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
	From string
	To   string

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
	haiconf.Output(h.rc, "Downloading %s to %s", h.From, h.To)

	resp, err := http.Get(h.From)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(h.To)
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
		h.From = f
		return nil
	}

	return haiconf.NewArgError("From must be http or https", args)
}

func (h *HttpGet) setTo(args haiconf.CommandArgs) error {
	t, _ := haiconf.CheckString("To", args)

	if len(t) == 0 {
		return haiconf.NewArgError("To must be provided", args)
	}

	if !utils.HasFileName(t) {
		return haiconf.NewArgError(t+" must be a file name", args)
	}

	dir := filepath.Dir(t)
	id, err := utils.IsDir(dir)
	if err != nil {
		return haiconf.NewArgError(err.Error(), args)
	}

	if !id {
		return haiconf.NewArgError(dir+" is not a directory", args)
	}

	h.To = t

	return nil
}
