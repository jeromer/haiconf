function Main()
    RuntimeConfig({
        RollbackOnError = true,
    })

    --
    -- Core commands
    --
    Directory({
        Path    = "/path/to/directory",
        Mode    = 0755,
        Owner   = "root",
        Group   = "root",
        Recurse = true,
        Ensure  = "present",
    })

    File({
        Path     = "/path/to/file.ext",
        Mode     = 0644,
        Owner    = "root",
        Group    = "root",
        Ensure   = "present",
        Template = "/path/to/template.ext"
    })

    -- related to RollbackOnError
    -- if execution of commands below fails for some reason
    -- the runtime will rollback only up to this point
    Checkpoint({
        Id = "Checkpoint 1",
    })

    AptGet({
        Method = "install",

        -- defined here:
        Packages = {"vim", "mutt", "cowsay"},

        -- or alternatively:
        PackageFromSource = "/path/to/packages.to.install.txt",

        -- automatically added to the apt-get call
        ExtraOptions = {
            "--download-only",
            "--simulate",
            "--fix-broken",
        }
    })

    Cron({
        Command = "/path/to/ntpdate",
        Ensure = "present",

        Env = {
            PATH     = "$PATH:/usr/bin/foo",
            ENV_VAR2 = "foo-bar",
        },

        Schedule = {
            -- yearly / monthly / weekly / daily / hourly / reboot etc
            Predefined = "yearly",

            -- alternatively :
            WeekDay  = "*",
            Month    = "*",
            MonthDay = "*",
            Hour     = "*",
            Minute   = "*",
        },

        RunAs = "root",
    })

    Service({
        Ensure = "running",
        Name = "foo",
        -- XXX : I need to figure out the rest
    })

    Exec({
        Command = "/usr/bin/foo",
        ExecutionDir = "/home/someuser",
        RunAs = "root",

        Env = {
            PATH     = "$PATH:/usr/bin/foo",
            ENV_VAR2 = "foo-bar",
        }
    })

    --
    -- small utilities
    --
    Curl({
        Url = "http://some.url/path.to.pkg",
    })

    TarGz({
        From = "/path/to/from/dir",
        Source = "/path/to/dest/dir.tar.gz",
    })

    UnTarGz({
        From = "/path/to/from/dir.tar.gz",
        Source = "/path/to/dest/dir",
    })

    Touch({
        File = "/path/to/file",
    })
end
