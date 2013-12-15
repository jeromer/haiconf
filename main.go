// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	lua "github.com/aarzilli/golua/lua"
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/fs"
	"github.com/jeromer/haiconf/haiconf/pkg"
	"github.com/jeromer/haiconf/haiconf/utils/httpget"
	"github.com/jeromer/haiconf/haiconf/utils/targz"
	"github.com/stevedonovan/luar"
	"log"
	"os"
)

var (
	flagConfigFile = flag.String("config", "./haiconf.lua", "Path to config file")
	flagVerbose    = flag.Bool("verbose", true, "Verbose mode")
)

func main() {
	flag.Parse()

	conf := NewConf()
	defer conf.Close()

	err := conf.DoFile(*flagConfigFile)
	if err != nil {
		panic(err)
	}

	conf.RunMain()
}

// -------------------

type Conf struct {
	Inputs luar.Map
	l      *lua.State
}

func NewConf() *Conf {
	c := Conf{
		l: luar.Init(),
	}

	c.registerCommands()

	return &c
}

func (c *Conf) registerCommands() {
	luar.Register(c.l, "", luar.Map{
		"Directory": Directory,
		"File":      File,
		"AptGet":    AptGet,
		"HttpGet":   HttpGet,
		"TarGz":     TarGz,
		"UnTarGz":   UnTarGz,
	})
}

func (c *Conf) DoFile(f string) error {
	return c.l.DoFile(f)
}

func (c *Conf) Close() {
	c.l.Close()
}

func (c *Conf) RunMain() {
	fun := luar.NewLuaObjectFromName(c.l, "Main")
	_, err := fun.Call()

	if err != nil {
		panic(err)
	}
}

// -------------------

func Directory(args haiconf.CommandArgs) {
	runCommand(new(fs.Directory), args)
}

func File(args haiconf.CommandArgs) {
	runCommand(new(fs.File), args)
}

func AptGet(args haiconf.CommandArgs) {
	runCommand(new(pkg.AptGet), args)
}

func HttpGet(args haiconf.CommandArgs) {
	runCommand(new(httpget.HttpGet), args)
}

func TarGz(args haiconf.CommandArgs) {
	runCommand(new(targz.TarGz), args)
}

func UnTarGz(args haiconf.CommandArgs) {
	runCommand(new(targz.UnTarGz), args)
}

func runCommand(c haiconf.Commander, args haiconf.CommandArgs) {
	rc := haiconf.RuntimeConfig{
		Verbose: *flagVerbose,
		Output:  os.Stdout,
	}

	err := c.SetDefault(&rc)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = c.SetUserConfig(args)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = c.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
