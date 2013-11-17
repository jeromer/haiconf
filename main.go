package main

import (
	"flag"
	"fmt"
	lua "github.com/aarzilli/golua/lua"
	"github.com/jeromer/haiconf/haiconf/fs"
	"github.com/stevedonovan/luar"
)

var (
	flagConfigFile = flag.String("config", "./haiconf.lua", "Path to config file")
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

	luar.Register(c.l, "", luar.Map{
		"Directory": Directory,
		"File":      File,
	})

	return &c
}

func (c *Conf) DoFile(f string) error {
	return c.l.DoFile(f)
}

func (c *Conf) Close() {
	c.l.Close()
}

// -------------------

func (c *Conf) RunMain() {
	fun := luar.NewLuaObjectFromName(c.l, "Main")
	_, err := fun.Call()

	if err != nil {
		panic(err)
	}
}

func Directory(m map[string]interface{}) {
	err := fs.ApplyDirectory(m)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func File(m map[string]interface{}) {
	err := fs.ApplyFile(m)
	if err != nil {
		fmt.Println(err.Error())
	}
}
