--
-- Ideal public usage:
--
-- ssh = SSH.New()
-- ssh:InstallServer()
-- ssh:InstallClient()
--

-- class declaration
local SSH = {}
SSH.__index = SSH

-- ----------------
-- private methods
-- ----------------

local function installPkg(what)
    pkgName = "server"

    if what == "client" then
        pkgName = what
    end

    AptGet({
        Method = "install",
        Packages = {"openssh-client"},
    })
end

local function checkService()
    Service({
        Name   = "ssh-server",
        Ensure = "running",
    })
end

local function configureServer()
    File({
        Path     = "/etc/ssh/sshd_config",
        Mode     = 0640,
        Owner    = "root",
        Group    = "root",
        Ensure   = "present",
        Template = "/absolute/path/to/templates/etc/sshd_config",
        BindVariables = {
            Hostname      = Haiconf.Hostname,
            PortNumer     = 1234,
            AllowGroups   = true,
            GroupsToAllow = {"root", "video", "sullivan", "googly-bear"},
        }
    })
end

local function configureClient()
    File({
        Path     = "/etc/ssh/ssh_config",
        Mode     = 0644,
        Owner    = "root",
        Group    = "root",
        Ensure   = "present",
        Template = "/absolute/path/to/templates/etc/ssh_config",
        BindVariables = {
            -- let's pretend some variables have been bound and the template
            -- actually exists
        }
    })
end
-- -----------------
-- public stuff
-- -----------------
function SSH.New()
  local self = setmetatable({}, SSH)
  return self
end

function SSH:InstallServer()
    installPkg("server")
    checkService()
    configureServer()
end

function SSH:InstallClient()
    installPkg("client")
    configureClient()
end
