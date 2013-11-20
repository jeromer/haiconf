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
	"github.com/jeromer/haiconf/haiconf"
	"os"
	"os/user"
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

	rc *haiconf.RuntimeConfig
}

func (d *Directory) SetDefault(rc *haiconf.RuntimeConfig) error {
	*d = Directory{
		path:    "",
		mode:    DEFAULT_MODE_DIRECTORY,
		owner:   new(user.User),
		group:   new(hacks.Group),
		recurse: false,
		ensure:  haiconf.ENSURE_PRESENT,
		rc:      rc,
	}

	return nil
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

	if d.ensure == haiconf.ENSURE_ABSENT {
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
	if d.ensure == haiconf.ENSURE_ABSENT {
		haiconf.Output(d.rc, "Removing directory %s", d.path)
		return RmDir(d.path, d.recurse)
	}

	haiconf.Output(d.rc, "Creating directory %s", d.path)
	err := MkDir(d.path, d.recurse, d.mode)
	if err != nil {
		return err
	}

	haiconf.Output(d.rc, "Chmod %s on %s", d.mode, d.path)
	err = Chmod(d.path, d.mode)
	if err != nil {
		return err
	}

	haiconf.Output(d.rc, "Chown %s:%s on %s", d.owner.Username, d.group.Name, d.path)
	err = Chown(d.path, d.owner, d.group)
	if err != nil {
		return err
	}

	return nil
}

func (d *Directory) setPath(args haiconf.CommandArgs) error {
	p, err := haiconf.CheckAbsolutePath("Path", args)
	if err != nil {
		return err
	}

	d.path = p

	return nil
}

func (d *Directory) setEnsure(args haiconf.CommandArgs) error {
	e, err := haiconf.CheckEnsure(args)
	if err != nil {
		return err
	}

	d.ensure = e

	return nil
}

func (d *Directory) setMode(args haiconf.CommandArgs) error {
	m, err := haiconf.CheckInt64("Mode", args)
	if err != nil {
		return err
	}

	d.mode = os.FileMode(m)

	return nil
}

func (d *Directory) setOwner(args haiconf.CommandArgs) error {
	u, err := haiconf.CheckSystemUser("Owner", args)
	if err != nil {
		return err
	}

	d.owner = u
	return nil
}

func (d *Directory) setGroup(args haiconf.CommandArgs) error {
	grp, err := haiconf.CheckSystemGroup("Group", args)
	if err != nil {
		return err
	}

	d.group = grp
	return nil
}

func (d *Directory) setRecurse(args haiconf.CommandArgs) error {
	d.recurse = haiconf.CheckBool("Recurse", args)
	return nil
}
