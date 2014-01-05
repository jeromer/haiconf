// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"fmt"
	"github.com/jeromer/haiconf/haiconf/osutils"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

type Crontab struct {
	Path string
	User *user.User
}

func NewCrontab(u *user.User) *Crontab {
	return &Crontab{
		Path: "/usr/bin/crontab",
		User: u,
	}
}

func (c *Crontab) Add(cj Cronjob) error {
	cronjobs, err := c.Read()
	if err != nil {
		return err
	}

	cronjobs = c.RemoveDuplicates(append(cronjobs, cj))
	return c.Save(cronjobs)
}

func (c *Crontab) Remove(cj Cronjob) error {
	cronjobs, err := c.Read()
	if err != nil {
		return err
	}

	idx := c.buildCronIndex(cronjobs)
	delete(idx, cj.Hash())

	cronjobs = []Cronjob{}
	for _, cj := range idx {
		cronjobs = append(cronjobs, cj)
	}

	return c.Save(cronjobs)
}

func (c *Crontab) Read() ([]Cronjob, error) {
	sc := osutils.SystemCommand{
		Path:                 c.Path,
		Args:                 []string{"-u", c.User.Username, "-l"},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	output := sc.Run()
	if output.HasError() {
		nilCronjobs := []Cronjob{}
		noCrontab := fmt.Sprintf("no crontab for %s", c.User.Username)
		stdErr := strings.Trim(strings.TrimSpace(output.Stderr), "\n")

		if stdErr == noCrontab {
			return nilCronjobs, nil
		}

		return nilCronjobs, output
	}

	cp := CrontabParser{
		Buff: []byte(output.Stdout),
	}

	cronjobs, err := cp.Parse()
	if err != nil {
		return []Cronjob{}, err
	}

	return c.RemoveDuplicates(cronjobs), nil
}

func (c *Crontab) Save(cronjobs []Cronjob) error {
	fileName, err := c.writeTmpCrontab(cronjobs)
	if err != nil {
		return err
	}
	defer os.Remove(fileName)

	// -----

	sc := osutils.SystemCommand{
		Path:                 c.Path,
		Args:                 []string{"-u", c.User.Username, fileName},
		ExecDir:              os.TempDir(),
		EnableShellExpansion: true,
	}

	output := sc.Run()
	if output.HasError() {
		return output
	}

	return nil
}

func (c *Crontab) RemoveDuplicates(cronjobs []Cronjob) []Cronjob {
	uniques := []Cronjob{}
	for _, cronjob := range c.buildCronIndex(cronjobs) {
		uniques = append(uniques, cronjob)
	}

	return uniques
}

func (c *Crontab) buildCronIndex(cronjobs []Cronjob) map[string]Cronjob {
	index := make(map[string]Cronjob, len(cronjobs))

	for _, cronjob := range cronjobs {
		h := cronjob.Hash()
		_, exists := index[h]
		if exists {
			continue
		}

		index[h] = cronjob
	}

	return index
}

func (c *Crontab) writeTmpCrontab(cronjobs []Cronjob) (string, error) {
	content := c.buildCrontabContent(cronjobs)
	fd, err := ioutil.TempFile("", "cron")
	if err != nil {
		return "", err
	}

	_, err = fd.WriteString(content)
	if err != nil {
		return "", err
	}

	err = fd.Close()
	if err != nil {
		return "", err
	}

	return fd.Name(), nil
}

func (c *Crontab) buildCrontabContent(cronjobs []Cronjob) string {
	buff := make([]string, len(cronjobs))

	for i, c := range cronjobs {
		buff[i] = c.BuildBlock()
	}

	return strings.Join(buff, "\n")
}
