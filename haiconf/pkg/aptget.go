// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

	METHOD_INSTALL = "install"
	METHOD_UPDATE  = "update"
	METHOD_REMOVE  = "remove"
)

var (
	availableMethods = []string{
		METHOD_INSTALL,
		METHOD_UPDATE,
		METHOD_REMOVE,
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
	Method       string
	Packages     []string
	ExtraOptions []string
	shellCmd     string

	rc *haiconf.RuntimeConfig
}

func (ag *AptGet) SetDefault(rc *haiconf.RuntimeConfig) error {
	ag.rc = rc
	return nil
}

func (ag *AptGet) SetUserConfig(args haiconf.CommandArgs) error {
	err := ag.setMethod(args)
	if err != nil {
		return err
	}

	if ag.Method == METHOD_UPDATE {
		return nil
	}

	err = ag.setPackages(args)
	if err != nil {
		return err
	}

	err = ag.setExtraOptions(args)
	if err != nil {
		return err
	}

	return nil
}

func (ag *AptGet) Run() error {
	// Check if we could use APT's native API instead of calling an external command
	// It seems we can use C++ with CGO by using Swig.
	//
	// Using the API directly would be much more efficient and less risky that
	// calling a system command which is kind of ugly.
	//
	// Interesting links:
	// http://golang.org/doc/faq#Do_Go_programs_link_with_Cpp_programs
	// http://www.swig.org/Doc2.0/Go.html

	// XXX : crap
	args := append(defaultOptions, ag.Method)
	args = stringutils.RemoveDuplicates(append(args, ag.ExtraOptions...))
	args = stringutils.RemoveDuplicates(append(args, ag.Packages...))

	sc := osutils.SystemCommand{
		Path:                 APT_GET,
		Args:                 args,
		EnvVars:              envVariables,
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	output := sc.Run()
	if output.HasError() {
		return output
	}

	return nil
}

func (ag *AptGet) setMethod(args haiconf.CommandArgs) error {
	m, err := haiconf.CheckStringChoice("Method", args, availableMethods)

	if err != nil {
		return err
	}

	ag.Method = m

	return nil
}

func (ag *AptGet) setPackages(args haiconf.CommandArgs) error {
	pl, _ := haiconf.CheckStringList("Packages", args)

	if len(pl) > 0 {
		ag.Packages = stringutils.RemoveDuplicates(pl)
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
		ag.Packages = stringutils.RemoveDuplicates(pkgs)

		return nil
	}

	msg := "You must provide a value for either Packages or PackagesFromSource"
	return haiconf.NewArgError(msg, args)
}

func (ag *AptGet) setExtraOptions(args haiconf.CommandArgs) error {
	xtraOpts, _ := haiconf.CheckStringList("ExtraOptions", args)
	l := len(xtraOpts)

	if l > 0 {
		ag.ExtraOptions = stringutils.RemoveDuplicates(xtraOpts)
	}

	return nil
}
