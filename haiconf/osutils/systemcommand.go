// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package osutils

import (
	"bytes"
	"fmt"
	"os"
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
		msgFmt = append(msgFmt, "Error message was \"%s\".")
		fmtArgs = append(fmtArgs, sce.ExitMessage)
	}

	if sce.Stdout != "" {
		msgFmt = append(msgFmt, "StdOut was \"%s\".")
		fmtArgs = append(fmtArgs, sce.Stdout)
	}

	if sce.Stderr != "" {
		msgFmt = append(msgFmt, "StdErr was \"%s\".")
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
		path = "/bin/sh"

		fullCmd := fmt.Sprintf("%s %s", sc.Path, strings.Join(sc.Args, " "))
		args = []string{"sh", "-c", fullCmd}
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
	sc.addPathIfMissing()

	envv := make([]string, len(sc.EnvVars))

	i := 0
	for name, value := range sc.EnvVars {
		envv[i] = name + "=" + value
		i++
	}

	return envv
}

func (sc *SystemCommand) addPathIfMissing() {
	if sc.EnvVars == nil {
		return
	}

	_, present := sc.EnvVars["PATH"]

	if !present {
		sc.EnvVars["PATH"] = os.Getenv("PATH")
	}
}
