// Usage in lua configuration file
//
// Directory({
//         Path    = "/tmp/haiconf/testdirectory",
//         Mode    = "0755",
//         Owner   = "jerome",
//         Group   = "wheel",
//         Recurse = true,
//         Ensure  = "present",
//     })
//
package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/lib"
	"os"
	"os/user"
	"strconv"
)

type Directory struct {
	path    string
	mode    os.FileMode
	owner   *user.User
	recurse bool
	ensure  string

	// no user.Group in golang yet
	// (https://code.google.com/p/go/issues/detail?id=2617)
	// Let's use a temporary one
	group *hacks.Group
}

type DirectoryError struct {
	errorString string
}

func ApplyDirectory(args haiconf.CommandArgs) error {
	var err error

	d := new(Directory)
	d.SetDefault()
	err = d.SetUserConfig(args)
	if err != nil {
		return err
	}

	return d.Run()
}

func (d *Directory) SetDefault() {
	*d = Directory{
		path:    "",
		mode:    DEFAULT_MODE_DIRECTORY,
		owner:   new(user.User),
		group:   new(hacks.Group),
		recurse: false,
		ensure:  ENSURE_PRESENT,
	}
}

func (d *Directory) SetUserConfig(args haiconf.CommandArgs) error {
	err := d.setPath(args)
	if err != nil {
		return err
	}

	err = d.setEnsure(args)
	if err != nil {
		return err
	}

	err = d.setRecurse(args)
	if err != nil {
		return err
	}

	if d.ensure == ENSURE_ABSENT {
		return nil
	}

	err = d.setMode(args)
	if err != nil {
		return err
	}

	err = d.setOwner(args)
	if err != nil {
		return err
	}

	err = d.setGroup(args)
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) Run() error {
	// XXX : acquire/release lock
	if d.ensure == ENSURE_ABSENT {
		return RmDir(d.path, d.recurse)
	}

	err := MkDir(d.path, d.recurse, d.mode)
	if err != nil {
		return err
	}

	err = d.chmod()
	if err != nil {
		return err
	}

	err = d.chown()
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) setPath(args haiconf.CommandArgs) error {
	p, err := CheckPath(args)
	if err != nil {
		return err
	}

	d.path = p

	return nil
}

func (d *Directory) setEnsure(args haiconf.CommandArgs) error {
	e, err := CheckEnsure(args)
	if err != nil {
		return err
	}

	d.ensure = e

	return nil
}

func (d *Directory) setMode(args haiconf.CommandArgs) error {
	m, err := CheckMode(args)
	if err != nil {
		return err
	}

	d.mode = os.FileMode(m)

	return nil
}

func (d *Directory) setOwner(args haiconf.CommandArgs) error {
	u, err := CheckOwner(args)
	if err != nil {
		return err
	}

	d.owner = u
	return nil
}

func (d *Directory) setGroup(args haiconf.CommandArgs) error {
	grp, err := CheckGroup(args)
	if err != nil {
		return err
	}

	d.group = grp
	return nil
}

func (d *Directory) setRecurse(args haiconf.CommandArgs) error {
	d.recurse = CheckRecurse(args)
	return nil
}

func (d *Directory) chmod() error {
	err := os.Chmod(d.path, d.mode)
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) chown() error {
	uid, err := strconv.Atoi(d.owner.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(d.group.Gid)
	if err != nil {
		return err
	}

	err = os.Chown(d.path, uid, gid)
	return err
}

func (err *DirectoryError) Error() string {
	return err.errorString
}
