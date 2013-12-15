// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package haiconf

import (
	"fmt"
	"io"
)

type CommandArgs map[string]interface{}
type RuntimeConfig struct {
	Verbose bool
	Output  io.Writer
}

type Commander interface {
	SetDefault(*RuntimeConfig) error
	SetUserConfig(CommandArgs) error
	Run() error
}

type HaiconfError struct {
	Msg  string
	Args CommandArgs
}

func NewArgError(m string, args CommandArgs) *HaiconfError {
	return &HaiconfError{Msg: m, Args: args}
}

func (err *HaiconfError) Error() string {
	return fmt.Sprintf("%s. Received args : %+v", err.Msg, err.Args)
}

func Output(rc *RuntimeConfig, msgFmt string, msgArgs ...interface{}) {
	if !rc.Verbose {
		return
	}

	msg := fmt.Sprintf(msgFmt+"\n", msgArgs...)
	io.WriteString(rc.Output, msg)
}
