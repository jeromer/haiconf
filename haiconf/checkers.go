package haiconf

import (
	"github.com/jeromer/haiconf/hacks"
	"os/user"
	"path"
	"strconv"
)

const (
	ENSURE_PRESENT = "present"
	ENSURE_ABSENT  = "absent"
)

func CheckAbsolutePath(k string, args CommandArgs) (string, error) {
	p, _ := args[k].(string)

	if p == "" {
		return p, NewArgError(k+" must be provided", args)
	}

	if !path.IsAbs(p) {
		return p, NewArgError(k+" must be absolute", args)
	}

	return p, nil
}

func CheckInt64(k string, args CommandArgs) (int64, error) {
	mStr, _ := args["Mode"].(string)

	if mStr == "" {
		return 0, NewArgError(k+" must be provided", args)
	}

	return strconv.ParseInt(mStr, 8, 0)
}

func CheckSystemUser(k string, args CommandArgs) (*user.User, error) {
	o, _ := args[k].(string)
	if o == "" {
		return nil, NewArgError(k+" must be defined", args)
	}

	return user.Lookup(o)
}

func CheckSystemGroup(k string, args CommandArgs) (*hacks.Group, error) {
	g, _ := args[k].(string)
	if g == "" {
		return nil, NewArgError(k+" must be defined", args)
	}

	return hacks.LookupSystemGroup(g)
}

func CheckBool(k string, args CommandArgs) bool {
	r, _ := args[k].(bool)
	return r
}

func CheckEnsure(args CommandArgs) (string, error) {
	e, _ := args["Ensure"].(string)

	if e != "" {
		if e != ENSURE_PRESENT && e != ENSURE_ABSENT {
			errMsg := "Invalid choice for Ensure. Valid choices are \"" + ENSURE_PRESENT + "\" or \"" + ENSURE_ABSENT + "\""
			return "", NewArgError(errMsg, args)
		}

		return e, nil
	}

	return e, NewArgError("Ensure flag must be provided", args)
}
