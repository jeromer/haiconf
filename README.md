What is Haiconf ?
=================

Haiconf is an experimental configuration management system based on the power of Lua.

How does this stuff work ?
--------------------------

The idea is pretty simple. Haiconf will provide you with a set of primitives inspired by Puppet and others.
You simply have to write a program which uses theses primitives and this program will be applied sequentially.
No dependecy graph from hell, no DSL, no YAML.

Enough bullshit, show me some code
----------------------------------

Well, haiconf is under heavy development but if you are curious you can have a look at what an ideal configuration file may look like in the future:

-   <https://github.com/jeromer/haiconf/blob/master/prototypes/python/python.lua>
-   <https://github.com/jeromer/haiconf/blob/master/prototypes/ssh/ssh.lua>

You can also have a look at <https://github.com/jeromer/haiconf/blob/master/haiconf.lua> .

I want to test it
-----------------

Currently there is no package available.

Applying the commands below should be enough:

1.  Install lua **5.1** (no 5.2 please)
2.  git clone <git@github.com>:jeromer/haiconf.git
3.  cd ./haiconf
4.  ln -sv pwd $GOPATH/src/github.com/jeromer/
5.  make installdependencies tests
6.  change whatever you want in haiconf.lua
7.  go run main.go (haiconf.lua will automatically be applied)

FAQ
---

### Where is the master server ?

There is none. Haiconf is aimed at configuring a local node for the moment. Distribution will come later.
