package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/haiconf"
	"os"
	"os/user"
	"path"
	"strconv"
)

const (
	DEFAULT_MODE_DIRECTORY = os.FileMode(0755)
	DEFAULT_MODE_FILE      = os.FileMode(0644)
	ENSURE_PRESENT         = "present"
	ENSURE_ABSENT          = "absent"
)

func CheckPath(args haiconf.CommandArgs) (string, error) {
	return checkAbsolutePath(args, "Path")
}

func CheckSource(args haiconf.CommandArgs) (string, error) {
	return checkAbsolutePath(args, "Source")
}

func checkAbsolutePath(args haiconf.CommandArgs, k string) (string, error) {
	p, _ := args[k].(string)

	if p == "" {
		return p, haiconf.NewArgError(k+" must be provided", args)
	}

	if !path.IsAbs(p) {
		return p, haiconf.NewArgError(k+" must be absolute", args)
	}

	return p, nil
}

func CheckEnsure(args haiconf.CommandArgs) (string, error) {
	e, _ := args["Ensure"].(string)

	if e != "" {
		if e != ENSURE_PRESENT && e != ENSURE_ABSENT {
			errMsg := "Invalid choice for Ensure. Valid choices are \"" + ENSURE_PRESENT + "\" or \"" + ENSURE_ABSENT + "\""
			return "", haiconf.NewArgError(errMsg, args)
		}

		return e, nil
	}

	return e, haiconf.NewArgError("Ensure flag must be provided", args)
}

func CheckMode(args haiconf.CommandArgs) (int64, error) {
	mStr, _ := args["Mode"].(string)

	if mStr == "" {
		return 0, haiconf.NewArgError("Mode must be provided", args)
	}

	return strconv.ParseInt(mStr, 8, 0)
}

func CheckOwner(args haiconf.CommandArgs) (*user.User, error) {
	o, _ := args["Owner"].(string)
	if o == "" {
		return nil, haiconf.NewArgError("Owner must be defined", args)
	}

	return user.Lookup(o)
}

func CheckGroup(args haiconf.CommandArgs) (*hacks.Group, error) {
	g, _ := args["Group"].(string)
	if g == "" {
		return nil, haiconf.NewArgError("Group must be defined", args)
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