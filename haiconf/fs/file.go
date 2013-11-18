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
	path   string
	mode   os.FileMode
	ensure string
	owner  *user.User
	group  *hacks.Group
	source string

	templateVariables map[string]interface{}
}

func (f *File) SetDefault() error {
	*f = File{
		path:              "",
		mode:              DEFAULT_MODE_FILE,
		owner:             new(user.User),
		group:             new(hacks.Group),
		ensure:            ENSURE_PRESENT,
		source:            "",
		templateVariables: nil,
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

	if f.ensure == ENSURE_ABSENT {
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
	if f.ensure == ENSURE_ABSENT {
		return os.Remove(f.path)
	}

	err := MkDir(path.Dir(f.path), true, 0755)
	if err != nil {
		return err
	}

	err = f.storeFile()
	if err != nil {
		return err
	}

	err = Chmod(f.path, f.mode)
	if err != nil {
		return err
	}

	err = Chown(f.path, f.owner, f.group)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) setPath(args haiconf.CommandArgs) error {
	p, err := CheckPath(args)
	if err != nil {
		return err
	}

	f.path = p

	return nil
}

func (f *File) setEnsure(args haiconf.CommandArgs) error {
	e, err := CheckEnsure(args)
	if err != nil {
		return err
	}

	f.ensure = e

	return nil
}

func (f *File) setMode(args haiconf.CommandArgs) error {
	m, err := CheckMode(args)
	if err != nil {
		return err
	}

	f.mode = os.FileMode(m)

	return nil
}

func (f *File) setOwner(args haiconf.CommandArgs) error {
	u, err := CheckOwner(args)
	if err != nil {
		return err
	}

	f.owner = u
	return nil
}

func (f *File) setGroup(args haiconf.CommandArgs) error {
	grp, err := CheckGroup(args)
	if err != nil {
		return err
	}

	f.group = grp
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

	f.templateVariables = tmp
	return nil
}

func (f *File) setSource(args haiconf.CommandArgs) error {
	src, err := CheckSource(args)
	if err != nil {
		return err
	}

	f.source = src
	return nil
}

func (f *File) storeFile() error {
	from, err := os.Open(f.source)
	if err != nil {
		return err
	}

	to, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.mode)
	if err != nil {
		return err
	}

	if f.templateVariables == nil {
		// XXX : check bytes written ?
		_, err = io.Copy(to, from)
		return err
	}

	buff, err := ioutil.ReadAll(from)
	if err != nil {
		return err
	}

	tpl := template.New(path.Base(f.path) + "-template")
	t := template.Must(tpl.Parse(string(buff)))
	return t.Execute(to, f.templateVariables)
}
