package fs

import (
	"fmt"
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/lib"
	"os"
	"os/user"
	"path"
	"strconv"
)

const (
	DEFAULT_MODE_DIRECTORY = os.FileMode(0755)
	ENSURE_PRESENT         = "present"
	ENSURE_ABSENT          = "absent"
)

type FsError struct {
	Msg  string
	Args haiconf.CommandArgs
}

func NewFsError(m string, args haiconf.CommandArgs) *FsError {
	return &FsError{Msg: m, Args: args}
}

func (err *FsError) Error() string {
	return fmt.Sprintf("%s. Received args : %+v", err.Msg, err.Args)
}

func CheckPath(args haiconf.CommandArgs) (string, error) {
	k := "Path"
	p, _ := args[k].(string)

	if p == "" {
		return p, NewFsError("Path must be provided", args)
	}

	if !path.IsAbs(p) {
		return p, NewFsError("Path must be absolute", args)
	}

	return p, nil
}

func CheckEnsure(args haiconf.CommandArgs) (string, error) {
	e, _ := args["Ensure"].(string)

	if e != "" {
		if e != ENSURE_PRESENT && e != ENSURE_ABSENT {
			errMsg := "Invalid choice for Ensure. Valid choices are \"" + ENSURE_PRESENT + "\" or \"" + ENSURE_ABSENT + "\""
			return "", NewFsError(errMsg, args)
		}

		return e, nil
	}

	return e, NewFsError("Ensure flag must be provided", args)
}

func CheckMode(args haiconf.CommandArgs) (int64, error) {
	mStr, _ := args["Mode"].(string)

	if mStr == "" {
		return 0, NewFsError("Mode must be provided", args)
	}

	return strconv.ParseInt(mStr, 8, 0)
}

func CheckOwner(args haiconf.CommandArgs) (*user.User, error) {
	o, _ := args["Owner"].(string)
	if o == "" {
		return nil, NewFsError("Owner must be defined", args)
	}

	return user.Lookup(o)
}

func CheckGroup(args haiconf.CommandArgs) (*hacks.Group, error) {
	g, _ := args["Group"].(string)
	if g == "" {
		return nil, NewFsError("Group must be defined", args)
	}

	return hacks.LookupSystemGroup(g)
}

func CheckRecurse(args haiconf.CommandArgs) bool {
	r, _ := args["Recurse"].(bool)
	return r
}

func MkDir(path string, recurse bool, mode os.FileMode) error {
	// XXX : symlink support ?
	_, err := os.Stat(path)

	// directory already exists
	if err == nil {
		return nil
	}

	if recurse {
		return os.MkdirAll(path, mode)
	}

	return os.Mkdir(path, mode)
}

func RmDir(path string, recurse bool) error {
	// XXX : symlink support ?

	_, err := os.Stat(path)

	// directory does not exists
	if os.IsNotExist(err) {
		return nil
	}

	if recurse {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}

func Chmod(path string, mode os.FileMode) error {
	err := os.Chmod(path, mode)
	if err != nil {
		return err
	}

	return nil
}

func Chown(path string, usr *user.User, grp *hacks.Group) error {
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		return err
	}

	return os.Chown(path, uid, gid)
}
