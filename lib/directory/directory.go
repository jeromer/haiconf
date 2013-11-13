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
package directory

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/lib"
	"os"
	"os/user"
	"path"
	"strconv"
)

const (
	DEFAULT_MODE   = os.FileMode(0755)
	ENSURE_PRESENT = "present"
	ENSURE_ABSENT  = "absent"
)

var (
	ErrNameEmpty          = &DirectoryError{"Path must have a value"}
	ErrInvalidChoice      = &DirectoryError{"Invalid choice for Ensure. Valid choices are \"" + ENSURE_PRESENT + "\" or \"" + ENSURE_ABSENT + "\""}
	ErrOwnerEmpty         = &DirectoryError{"Owner must be defined"}
	ErrGroupEmpty         = &DirectoryError{"Group must be defined"}
	ErrPathMustBeAbsolute = &DirectoryError{"Path must be absolute"}
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

func ApplyCommand(args haiconf.CommandArgs) error {
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
		mode:    DEFAULT_MODE,
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
		return d.rmDir()
	}

	err := d.mkDir()
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
	p, _ := args["Path"].(string)

	if p == "" {
		return ErrNameEmpty
	}

	if !path.IsAbs(p) {
		return ErrPathMustBeAbsolute
	}

	d.path = p

	return nil
}

func (d *Directory) setEnsure(args haiconf.CommandArgs) error {
	e, _ := args["Ensure"].(string)
	if e != "" {
		if e != ENSURE_PRESENT && e != ENSURE_ABSENT {
			return ErrInvalidChoice
		}

		d.ensure = e
	}

	return nil
}

func (d *Directory) setMode(args haiconf.CommandArgs) error {
	mStr, _ := args["Mode"].(string)

	if mStr == "" {
		return nil
	}

	mOct, err := strconv.ParseInt(mStr, 8, 0)
	if err != nil {
		return err
	}

	// The rest of FileMode checking is left to os.Mkdir and others
	if mOct > 0 {
		d.mode = os.FileMode(mOct)
	}

	return nil
}

func (d *Directory) setOwner(args haiconf.CommandArgs) error {
	o, _ := args["Owner"].(string)
	if o == "" {
		return ErrOwnerEmpty
	}

	u, err := user.Lookup(o)
	if err != nil {
		return err
	}

	d.owner = u
	return nil
}

func (d *Directory) setGroup(args haiconf.CommandArgs) error {
	g, _ := args["Group"].(string)
	if g == "" {
		return ErrGroupEmpty
	}

	grp, err := hacks.LookupSystemGroup(g)
	if err != nil {
		return err
	}

	d.group = grp
	return nil
}

func (d *Directory) setRecurse(args haiconf.CommandArgs) error {
	r, _ := args["Recurse"].(bool)
	d.recurse = r
	return nil
}

func (d *Directory) mkDir() error {
	// XXX : symlink support ?

	_, err := os.Stat(d.path)

	// directory already exists
	if err == nil {
		return nil
	}

	if d.recurse {
		return os.MkdirAll(d.path, os.FileMode(d.mode))
	}

	return os.Mkdir(d.path, os.FileMode(d.mode))
}

func (d *Directory) rmDir() error {
	// XXX : symlink support ?

	_, err := os.Stat(d.path)

	// directory does not exists
	if os.IsNotExist(err) {
		return nil
	}

	if d.recurse {
		return os.RemoveAll(d.path)
	}

	return os.Remove(d.path)
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
