function Main()
    Directory({
        Path    = "/private/tmp/haiconf/testdirectory",
        Mode    = "0755",
        Owner   = "jerome",
        Group   = "wheel",
		-- XXX : lua boolean or string ?
        Recurse = true,
        Ensure  = "present",
		-- Ensure = "absent",
    })
end
