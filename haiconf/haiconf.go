package haiconf

import (
	"fmt"
)

type CommandArgs map[string]interface{}

type Commander interface {
	SetDefault() error
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
