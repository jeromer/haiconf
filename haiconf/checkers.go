package haiconf

import (
	"github.com/jeromer/haiconf/hacks"
	"os/user"
	"path"
	"strconv"
	"strings"
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
	return checkStringChoice("Ensure", args, []string{ENSURE_PRESENT, ENSURE_ABSENT})
}

func checkStringChoice(k string, args CommandArgs, choices []string) (string, error) {
	s, _ := args[k].(string)

	if s == "" {
		return s, NewArgError(k+" must be provided", args)
	}

	for _, c := range choices {
		if s == c {
			return s, nil
		}
	}

	errMsg := "Invalid choice for " + k + ". Valid choices are " + strings.Join(choices, ", ")
	return "", NewArgError(errMsg, args)
}
