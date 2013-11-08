--
-- Ideal public usage:
--
-- p = Python.New()
-- p:Install({
--  Pip        = true,
--  PythonDev  = true,
--  VirtualEnv = true,
--  Gunicorn   = true,
-- })
--

local Python = {}
Python.__index = Python

-- ----------------
-- private methods
-- ----------------

local function installPythonDev()
    AptGet({
        Method = "install",
        Packages = {"python-dev"},
    })
end

local function installVirtualEnv()
    AptGet({
        Method = "install",
        Packages = {"python-virtualen"},
    })
end

local function installGunicorn()
    AptGet({
        Method = "install",
        Packages = {"python-virtualen"},
    })

    Service({
        Name   = "gunicorn",
        Ensure = "running",
    })
end

local function installPip()
    AptGet({
        Method = "install",
        Packages = {"python-pip"},
    })
end

-- -----------------
-- public stuff
-- -----------------
function Python.New()
  local self = setmetatable({}, Python)
  return self
end

function Python:Install(conf)
    if conf["Pip"] then
        installPip()
    end

    if conf["PythonDev"] then
        installPythonDev()
    end

    if conf["VirtualEnv"] then
        installVirtualEnv()
    end

    if conf["Gunicorn"] then
        installGunicorn()
    end
end

function Python:NewVirtualEnv(name)
    -- we can obviously handle way more things
    Exec({
        Command = "virtualenv --distribute" .. name,
        ExecutionDir = "/some/execution/dir",
        RunAs = "root",
    })
end

---

p = Python.New()
p:Install({
    Pip        = true,
    PythonDev  = false,
    VirtualEnv = true,
    Gunicorn   = true,
})
p:NewVirtualEnv("foo")
