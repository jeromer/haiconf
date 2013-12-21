// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Usage in lua configuration file
//
// File({
//     Path     = "/etc/ssh/ssh_config",
//     Mode     = "0644",
//     Owner    = "root",
//     Group    = "root",
//     Ensure   = "present",
//     Source = "/absolute/path/to/templates/etc/ssh_config",
//     TemplateVariables = {
//         "VarString" = "some string",
//         "VarBoolean" = false,
//         "VarInt" = 1234,
//         "VarFloat" = 3.14,
//         "VarTable" = {"one", "two", "three"},
//         "VarMap" = {"a":"1", "b": "2"},
//     }
// })

package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/haiconf"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"text/template"
)

type File struct {
	Path   string
	Mode   os.FileMode
	Ensure string
	Owner  *user.User
	Group  *hacks.Group
	Source string

	TemplateVariables map[string]interface{}

	rc *haiconf.RuntimeConfig
}

func (f *File) SetDefault(rc *haiconf.RuntimeConfig) error {
	*f = File{
		Path:              "",
		Mode:              DEFAULT_MODE_FILE,
		Owner:             new(user.User),
		Group:             new(hacks.Group),
		Ensure:            haiconf.ENSURE_PRESENT,
		Source:            "",
		TemplateVariables: nil,
		rc:                rc,
	}

	return nil
}

func (f *File) SetUserConfig(args haiconf.CommandArgs) error {
	err := f.setPath(args)
	if err != nil {
		return err
	}

	err = f.setEnsure(args)
	if err != nil {
		return err
	}

	if f.Ensure == haiconf.ENSURE_ABSENT {
		return nil
	}

	err = f.setMode(args)
	if err != nil {
		return err
	}

	err = f.setOwner(args)
	if err != nil {
		return err
	}

	err = f.setGroup(args)
	if err != nil {
		return err
	}

	err = f.setSource(args)
	if err != nil {
		return err
	}

	err = f.setTemplateVariables(args)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Run() error {
	// XXX : acquire/release lock
	if f.Ensure == haiconf.ENSURE_ABSENT {
		haiconf.Output(f.rc, "Removing file %s", f.Path)
		return os.Remove(f.Path)
	}

	err := MkDir(path.Dir(f.Path), true, 0755)
	if err != nil {
		return err
	}

	haiconf.Output(f.rc, "Creating file %s", f.Path)
	err = f.storeFile()
	if err != nil {
		return err
	}

	haiconf.Output(f.rc, "Chmod %s on %s", f.Mode, f.Path)
	err = Chmod(f.Path, f.Mode)
	if err != nil {
		return err
	}

	haiconf.Output(f.rc, "Chown %s:%s on %s", f.Owner.Username, f.Group.Name, f.Path)
	err = Chown(f.Path, f.Owner, f.Group)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) setPath(args haiconf.CommandArgs) error {
	p, err := haiconf.CheckAbsolutePath("Path", args)
	if err != nil {
		return err
	}

	f.Path = p

	return nil
}

func (f *File) setEnsure(args haiconf.CommandArgs) error {
	e, err := haiconf.CheckEnsure(args)
	if err != nil {
		return err
	}

	f.Ensure = e

	return nil
}

func (f *File) setMode(args haiconf.CommandArgs) error {
	m, err := haiconf.CheckInt64("Mode", args)
	if err != nil {
		return err
	}

	f.Mode = os.FileMode(m)

	return nil
}

func (f *File) setOwner(args haiconf.CommandArgs) error {
	u, err := haiconf.CheckSystemUser("Owner", args)
	if err != nil {
		return err
	}

	f.Owner = u
	return nil
}

func (f *File) setGroup(args haiconf.CommandArgs) error {
	grp, err := haiconf.CheckSystemGroup("Group", args)
	if err != nil {
		return err
	}

	f.Group = grp
	return nil
}

func (f *File) setTemplateVariables(args haiconf.CommandArgs) error {
	tv, _ := args["TemplateVariables"].(map[string]interface{})
	l := len(tv)

	if l <= 0 {
		return nil
	}

	tmp := make(map[string]interface{}, l)

	for k, v := range tv {
		switch v.(type) {
		case string:
			if v == "true" || v == "false" {
				b, err := strconv.ParseBool(v.(string))
				if err != nil {
					return err
				}

				tmp[k] = b
				continue
			}

			tmp[k] = v
		default:
			tmp[k] = v
		}
	}

	f.TemplateVariables = tmp
	return nil
}

func (f *File) setSource(args haiconf.CommandArgs) error {
	src, err := haiconf.CheckAbsolutePath("Source", args)
	if err != nil {
		return err
	}

	f.Source = src
	return nil
}

func (f *File) storeFile() error {
	from, err := os.Open(f.Source)
	if err != nil {
		return err
	}

	to, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode)
	if err != nil {
		return err
	}

	if f.TemplateVariables == nil {
		// XXX : check bytes written ?
		_, err = io.Copy(to, from)
		return err
	}

	buff, err := ioutil.ReadAll(from)
	if err != nil {
		return err
	}

	tpl := template.New(path.Base(f.Path) + "-template")
	t := template.Must(tpl.Parse(string(buff)))
	return t.Execute(to, f.TemplateVariables)
}
