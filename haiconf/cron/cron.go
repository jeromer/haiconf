// Copyright 2013 Jérôme Renard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Usage in lua configuration file
//  Cron({
//      Command = "/path/to/ntpdate",
//      Ensure = "present",
//
//      Env = {
//          PATH     = "$PATH:/usr/bin/foo",
//          ENV_VAR2 = "foo-bar",
//      },
//
//      Schedule = {
//          -- yearly / monthly / weekly / daily / hourly
//          Predefined = "yearly",
//
//          -- alternatively :
//          Minute   = "*",
//          Hour     = "*",
//          MonthDay = "*",
//          Month    = "*",
//          WeekDay  = "*",
//      },
//
//      Owner="root",
//  })

package cron

import (
	"github.com/jeromer/haiconf/haiconf"
	"github.com/jeromer/haiconf/haiconf/utils"
	"os/user"
)

var (
	PREDEFINED_SCHEDULES = map[string][]string{
		//         []string{minute, hour, monthday, month, weekday}
		"yearly":  []string{"0", "0", "1", "1", "*"},
		"monthly": []string{"0", "0", "1", "*", "*"},
		"weekly":  []string{"0", "0", "*", "*", "0"},
		"daily":   []string{"0", "0", "*", "*", "*"},
		"hourly":  []string{"0", "*", "*", "*", "*"},
	}

	NON_PREDEFINED_SCHEDULES = []string{
		"Minute",
		"Hour",
		"MonthDay",
		"Month",
		"WeekDay",
	}
)

type Cron struct {
	Command  string
	Ensure   string
	Env      map[string]string
	Schedule []string
	Owner    *user.User

	rc *haiconf.RuntimeConfig
}

func (c *Cron) SetDefault(rc *haiconf.RuntimeConfig) error {
	*c = Cron{
		Command:  "",
		Ensure:   haiconf.ENSURE_PRESENT,
		Env:      map[string]string{},
		Schedule: make([]string, 5),
		Owner:    new(user.User),
		rc:       rc,
	}

	return nil
}

func (c *Cron) SetUserConfig(args haiconf.CommandArgs) error {
	var err error
	type setter func(haiconf.CommandArgs) error

	setters := []setter{
		c.setCommand,
		c.setOwner,
		c.setEnv,
		c.setSchedule,
		c.setEnsure,
	}

	for _, s := range setters {
		err = s(args)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cron) Run() error {
	cj := Cronjob{
		Schedule: c.Schedule,
		Command:  c.Command,
		Env:      c.Env,
	}

	ct := NewCrontab(c.Owner)

	if c.Ensure == haiconf.ENSURE_PRESENT {
		haiconf.Output(c.rc, "Adding cronjob %s for user %s", cj.Command, c.Owner.Username)
		return ct.Add(cj)
	}

	haiconf.Output(c.rc, "Removing cronjob %s for user %s", cj.Command, c.Owner.Username)
	return ct.Remove(cj)
}

func (c *Cron) setCommand(args haiconf.CommandArgs) error {
	// XXX : check command really exists ?
	cmd, err := haiconf.CheckString("Command", args)
	if err != nil {
		return err
	}

	c.Command = cmd

	return nil
}

func (c *Cron) setEnsure(args haiconf.CommandArgs) error {
	e, err := haiconf.CheckEnsure(args)
	if err != nil {
		return err
	}

	c.Ensure = e

	return nil
}

func (c *Cron) setOwner(args haiconf.CommandArgs) error {
	u, err := haiconf.CheckSystemUser("Owner", args)
	if err != nil {
		return err
	}

	c.Owner = u
	return nil
}

func (c *Cron) setEnv(args haiconf.CommandArgs) error {
	_, present := args["Env"]
	if present {
		env, err := utils.ToStringMap(args["Env"].(map[string]interface{}))
		if err != nil {
			return err
		}

		c.Env = env
	}

	return nil
}

func (c *Cron) setSchedule(args haiconf.CommandArgs) error {
	_, present := args["Schedule"]
	if !present {
		return haiconf.NewArgError("Schedule must be provided", args)
	}

	schedule, err := utils.ToStringMap(args["Schedule"].(map[string]interface{}))
	if err != nil {
		return err
	}

	predefined, err := c.checkPredefined(schedule)
	if err != nil {
		return err
	}

	if len(predefined) > 0 {
		c.Schedule = PREDEFINED_SCHEDULES[predefined]
		return nil
	}

	nonPredefined, err := c.checkNonPredefined(schedule)
	if err != nil {
		return err
	}
	c.Schedule = nonPredefined

	return nil
}

func (c *Cron) checkPredefined(schedule map[string]string) (string, error) {
	var chosen string

	predefined, present := schedule["Predefined"]
	if !present {
		return chosen, nil
	}

	predef := haiconf.CommandArgs{
		"Predefined": predefined,
	}

	var keywords []string
	for k, _ := range PREDEFINED_SCHEDULES {
		keywords = append(keywords, k)
	}

	chosen, err := haiconf.CheckStringChoice("Predefined", predef, keywords)
	if err != nil {
		return "", err
	}

	return chosen, nil
}

func (c *Cron) checkNonPredefined(schedule map[string]string) ([]string, error) {
	sched := make([]string, len(c.Schedule))

	tmp := make(haiconf.CommandArgs, len(schedule))
	for k, v := range schedule {
		tmp[k] = v
	}

	for i, nps := range NON_PREDEFINED_SCHEDULES {
		s, err := haiconf.CheckString(nps, tmp)
		if err != nil {
			return []string{}, err
		}
		sched[i] = s
	}

	return sched, nil
}
