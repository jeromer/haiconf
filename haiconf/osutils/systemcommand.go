package osutils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type SystemCommand struct {
	Path                 string
	Args                 []string
	EnvVars              map[string]string
	ExecDir              string
	EnableShellExpansion bool
	cmd                  *exec.Cmd
}

type SystemCommandError struct {
	FullCommand string
	ExitMessage string
	Stdout      string
	Stderr      string
}

func (sce SystemCommandError) Error() string {
	msgFmt := []string{"Error with command \"%s\"."}
	fmtArgs := []interface{}{sce.FullCommand}

	if sce.ExitMessage != "" {
		msgFmt = append(msgFmt, "Error message : \"%s\".")
		fmtArgs = append(fmtArgs, sce.ExitMessage)
	}

	if sce.Stdout != "" {
		msgFmt = append(msgFmt, "StdOut : \"%s\".")
		fmtArgs = append(fmtArgs, sce.Stdout)
	}

	if sce.Stderr != "" {
		msgFmt = append(msgFmt, "StdErr : \"%s\".")
		fmtArgs = append(fmtArgs, sce.Stderr)
	}

	return fmt.Sprintf(strings.Join(msgFmt, " "), fmtArgs...)
}

func (sc *SystemCommand) Run() error {
	cmd := sc.buildCmd()

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if err != nil {
		e := SystemCommandError{
			Stdout:      stdOut.String(),
			Stderr:      stdErr.String(),
			FullCommand: cmd.Path + " " + strings.Join(cmd.Args, " "),
			ExitMessage: err.Error(),
		}

		return e
	}

	return nil
}

func (sc *SystemCommand) buildCmd() *exec.Cmd {
	path := sc.Path
	args := sc.Args

	if sc.EnableShellExpansion {
		shellCommand := []string{"sh", "-c", sc.Path}
		path = "/bin/sh"
		args = append(shellCommand, sc.Args...)
	}

	cmd := &exec.Cmd{
		Path: path,
		Args: args,
		Env:  sc.buildEnvVars(),
		Dir:  sc.ExecDir,
	}

	return cmd
}

func (sc *SystemCommand) buildEnvVars() []string {
	envv := make([]string, len(sc.EnvVars))

	i := 0
	for name, value := range sc.EnvVars {
		envv[i] = name + "=" + value
		i++
	}

	return envv
}
