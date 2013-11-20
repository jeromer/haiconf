package pkg

// Usage in lua configuration file
//
// AptGet({
//         Method = "install",
//
//         -- defined here:
//         Packages = {"vim", "mutt", "cowsay"},
//
//         -- or alternatively:
//         PackagesFromSource = "/path/to/packages.to.install.txt",
//
//         -- automatically added to the apt-get call
//         ExtraOptions = {
//             "--download-only",
//             "--simulate",
//             "--fix-broken",
//         }
//     })

import (
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/osutils"
	"github.com/jeromer/haiconf/haiconf/stringutils"
	"io/ioutil"
	"os"
	"strings"
)

const (
	APT_GET = "/usr/bin/apt-get"
)

var (
	availableMethods = []string{
		"install",
	}

	envVariables = map[string]string{
		// http://www.debian.org/releases/stable/s390/ch05s02.html.en
		"DEBIAN_FRONTEND": "noninteractive",

		// http://www.debianadmin.com/manpages/aptlistbugsmanpage.htm
		"APT_LISTBUGS_FRONTEND": "none",

		// http://www.debianadmin.com/manpages/aptlistchangesmanpage.htm
		"APT_LISTCHANGES_FRONTEND": "none",
	}

	defaultOptions = []string{
		"--yes",
		"--quiet",
	}
)

type AptGet struct {
	method       string
	packages     []string
	extraOptions []string
	shellCmd     string
}

func (ag *AptGet) SetDefault() error {
	return nil
}

func (ag *AptGet) SetUserConfig(args haiconf.CommandArgs) error {
	setters := []func(haiconf.CommandArgs) error{
		ag.setMethod,
		ag.setPackages,
		ag.setExtraOptions,
	}

	var err error
	for _, s := range setters {
		err = s(args)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ag *AptGet) Run() error {
	// XXX : crap
	args := []string{ag.method}
	args = append(args, defaultOptions...)
	args = stringutils.RemoveDuplicates(append(args, ag.extraOptions...))
	args = stringutils.RemoveDuplicates(append(args, ag.packages...))

	sc := osutils.SystemCommand{
		Path:                 APT_GET,
		Args:                 args,
		EnvVars:              envVariables,
		ExecDir:              os.TempDir(),
		EnableShellExpansion: false,
	}

	return sc.Run()
}

func (ag *AptGet) setMethod(args haiconf.CommandArgs) error {
	m, err := haiconf.CheckStringChoice("Method", args, availableMethods)

	if err != nil {
		return err
	}

	ag.method = m

	return nil
}

func (ag *AptGet) setPackages(args haiconf.CommandArgs) error {
	pl, _ := haiconf.CheckStringList("Packages", args)

	if len(pl) > 0 {
		ag.packages = stringutils.RemoveDuplicates(pl)
		return nil
	}

	pfs, _ := haiconf.CheckAbsolutePath("PackagesFromSource", args)
	if pfs != "" {
		buff, err := ioutil.ReadFile(pfs)
		if err != nil {
			return err
		}

		nlSplit := func(r rune) bool {
			if r == '\n' {
				return true
			}
			return false
		}

		pkgs := strings.FieldsFunc(string(buff), nlSplit)
		ag.packages = stringutils.RemoveDuplicates(pkgs)

		return nil
	}

	msg := "You must provide a value for either Packages or PackagesFromSource"
	return haiconf.NewArgError(msg, args)
}

func (ag *AptGet) setExtraOptions(args haiconf.CommandArgs) error {
	xtraOpts, _ := haiconf.CheckStringList("ExtraOptions", args)
	l := len(xtraOpts)

	if l > 0 {
		ag.extraOptions = stringutils.RemoveDuplicates(xtraOpts)
	}

	return nil
}
