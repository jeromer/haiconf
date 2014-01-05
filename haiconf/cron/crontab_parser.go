// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"bytes"
	"unicode"
)

type CrontabParser struct {
	Buff []byte
}

func (cp *CrontabParser) Parse() ([]Cronjob, error) {
	var cronjobs []Cronjob

	lines := bytes.Split(cp.Buff, []byte("\n"))

	cj := Cronjob{
		// XXX : arbitrary but should be enough
		Env: make(map[string]string, 10),
	}

	found := false
	for _, line := range lines {
		if found {
			cj = Cronjob{
				Env: make(map[string]string, 10),
			}

			found = false
		}

		l := NewCrontabLine(line)

		if l.IsBlank() || l.IsComment() {
			continue
		}

		if l.IsEnvVar() {
			k, v := l.SplitEnvVar()
			cj.Env[k] = v
		}

		if l.IsCrontabLine() {
			cj.Schedule, cj.Command = l.SplitCrontabLine()
			cronjobs = append(cronjobs, cj)
			found = true
		}
	}

	return cronjobs, nil
}

type CrontabLine struct {
	Buff []byte
}

func NewCrontabLine(l []byte) *CrontabLine {
	return &CrontabLine{
		Buff: bytes.TrimSpace(l),
	}
}

func (cl *CrontabLine) IsComment() bool {
	return len(cl.Buff) > 0 && cl.Buff[0] == '#'
}

func (cl *CrontabLine) IsBlank() bool {
	return len(cl.Buff) <= 0
}

func (cl *CrontabLine) String() string {
	return string(cl.Buff)
}

func (cl *CrontabLine) IsEnvVar() bool {
	if cl.IsBlank() {
		return false
	}

	if !unicode.IsUpper(rune(cl.Buff[0])) {
		return false
	}

	parts := bytes.Split(cl.Buff, []byte("="))
	return len(parts) == 2
}

func (cl *CrontabLine) IsCrontabLine() bool {
	if cl.IsBlank() {
		return false
	}

	b := cl.Buff[0]

	if b == '*' || b == '@' || unicode.IsDigit(rune(b)) {
		return true
	}

	return false
}

func (cl *CrontabLine) SplitEnvVar() (string, string) {
	parts := bytes.Split(cl.Buff, []byte("="))

	k := ""
	v := ""

	if len(parts) >= 1 {
		k = string(parts[0])
	}

	if len(parts) == 2 {
		v = string(parts[1])
	}

	return k, v
}

func (cl *CrontabLine) SplitCrontabLine() ([]string, string) {
	schedule := make([]string, 5)
	command := ""

	parts := bytes.Split(cl.Buff, []byte(" "))

	max := 5

	if parts[0][0] == '@' {
		max = 1
		predefined := string(parts[0][1:])
		schedule = PREDEFINED_SCHEDULES[predefined]
	} else {
		i := 0
		for i < max {
			schedule[i] = string(parts[i])
			i++
		}
	}

	command = string(bytes.Join(parts[max:], []byte(" ")))

	return schedule, command
}
