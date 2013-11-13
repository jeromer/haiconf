package fs

import (
	"fmt"
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/lib"
	"os/user"
	"path"
	"strconv"
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
