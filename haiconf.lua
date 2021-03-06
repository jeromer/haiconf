function Main()
-- 	Directory({
-- 		Path    = "/private/tmp/haiconf/testdirectory",
-- 		Mode    = "0755",
-- 		Owner   = "jerome",
-- 		Group   = "wheel",
-- 	 -- XXX : lua boolean or string ?
-- 		Recurse = true,
-- 		Ensure  = "present",
-- 	 -- Ensure = "absent",
-- 	})
--
--
-- 	File({
-- 		Path   = "/tmp/testtpl.txt",
-- 		Mode   = "0644",
-- 		Owner  = "jerome",
-- 		Group  = "wheel",
-- 		Ensure = "present",
-- 		Source = "/tmp/sometemplate.tpl",
-- 		TemplateVariables = {
-- 			VarString = "some string",
-- 			VarBoolean = false,
-- 			VarInt = 1234,
-- 			VarFloat = 3.14,
-- 			VarTable = {"one", "two", "three"},
-- 			VarMap = {a="1", b="2"},
-- 		},
-- 	})

    AptGet({
        Method = "install",
        -- defined here:
        Packages = {"vim", "mutt", "cowsay"},
        -- or alternatively:
        -- PackageFromSource = "/path/to/packages.to.install.txt",
        -- automatically added to the apt-get call
        -- ExtraOptions = {
        --     "--download-only",
        --     "--simulate",
        --     "--fix-broken",
        -- }
    })

    HttpGet({
        From = "http://example.com/",
        To = "/tmp/example.html",
    })

    TarGz({
       Source = "/etc/",
       Dest = "/tmp/etc.tar.gz",
    })

    UnTarGz({
       Source = "/tmp/etc.tar.gz",
       Dest = "/tmp",
    })

    Cron({
        Command = "/path/to/ntpdate",
        Ensure = "present",

        Env = {
            PATH     = "$PATH:/usr/bin/foo",
            ENV_VAR2 = "foo-bar",
        },

        Schedule = {

            -- yearly / monthly / weekly / daily / hourly
            Predefined = "yearly",

            -- alternatively :
            WeekDay  = "*",
            Month    = "*",
            MonthDay = "*",
            Hour     = "*",
            Minute   = "*",
        },

        Owner="vagrant",
    })

    Group({
        Name="testgroup",
        Ensure="present",
    })
end
