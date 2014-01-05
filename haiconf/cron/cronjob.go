// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Cronjob struct {
	Env      map[string]string
	Schedule []string
	Command  string
}

func (c *Cronjob) BuildBlock() string {
	block := ""
	if len(c.Env) > 0 {
		block = c.buildEnvVarsBlock() + "\n"
	}

	return block + c.buildJobLine() + "\n"
}

func (c *Cronjob) Hash() string {
	envVarNames := []string{}
	for name, _ := range c.Env {
		envVarNames = append(envVarNames, name)
	}
	sort.Strings(envVarNames)

	var envVarBuff string
	for _, name := range envVarNames {
		envVarBuff += name + "=" + c.Env[name]
	}

	h := sha1.New()
	io.WriteString(h, envVarBuff+c.buildJobLine())

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *Cronjob) buildEnvVarsBlock() string {
	var buff []string

	for k, v := range c.Env {
		buff = append(buff, k+"="+v)
	}

	return strings.Join(buff, "\n")
}

func (c *Cronjob) buildJobLine() string {
	tmp := append(c.Schedule, c.Command)
	return strings.Join(tmp, " ")
}
