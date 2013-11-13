package fs

import (
	"github.com/jeromer/haiconf/hacks"
	"github.com/jeromer/haiconf/lib"
	"os/user"
	"path"
	"strconv"
)

func CheckPath(args haiconf.CommandArgs) (string, error) {
	p, _ := args["Path"].(string)

	if p == "" {
		return p, ErrNameEmpty
	}

	if !path.IsAbs(p) {
		return p, ErrPathMustBeAbsolute
	}

	return p, nil
}

func CheckEnsure(args haiconf.CommandArgs) (string, error) {
	e, _ := args["Ensure"].(string)

	if e != "" {
		if e != ENSURE_PRESENT && e != ENSURE_ABSENT {
			return "", ErrInvalidChoice
		}

		return e, nil
	}

	return e, ErrEnsureEmpty
}

func CheckMode(args haiconf.CommandArgs) (int64, error) {
	mStr, _ := args["Mode"].(string)

	if mStr == "" {
		return 0, ErrModeEmpty
	}

	return strconv.ParseInt(mStr, 8, 0)
}

func CheckOwner(args haiconf.CommandArgs) (*user.User, error) {
	o, _ := args["Owner"].(string)
	if o == "" {
		return new(user.User), ErrOwnerEmpty
	}

	return user.Lookup(o)
}

func CheckGroup(args haiconf.CommandArgs) (*hacks.Group, error) {
	g, _ := args["Group"].(string)
	if g == "" {
		return new(hacks.Group), ErrGroupEmpty
	}

	return hacks.LookupSystemGroup(g)
}

func CheckRecurse(args haiconf.CommandArgs) bool {
	r, _ := args["Recurse"].(bool)
	return r
}
