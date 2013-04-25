// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"strio"
)

const (
	exitStatusErr = "exit status "
)

type testCmd interface {
	Run(expDescr string) (*result, error)
}

type testCmd_ struct {
	cmd_ string
	vars map[string]string
}

func newTestCmd(cmd string) *testCmd_ {
	return &testCmd_{cmd, map[string]string{}}
}

func (t *testCmd_) Run(expDescr string) (*result, error) {
	cmdLine := t.cmd()
	parts := strings.Split(cmdLine, " ")
	cmd := exec.Command(parts[0], parts[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("couldn't get stdout for '%s': %v",
			cmdLine, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("couldn't get stderr for '%s': %v",
			cmdLine, err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("couldn't start '%s': %v", cmdLine, err)
	}

	output, err := strio.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("error reading stdout for '%s': %v",
			cmdLine, err)
	}

	errput, err := strio.ReadAll(stderr)
	if err != nil {
		return nil, fmt.Errorf("error reading stderr for '%s': %v",
			cmdLine, err)
	}

	// err is set if exit status is not 0
	err = cmd.Wait()

	exitStatus := 0
	if err != nil {
		if strings.HasPrefix(err.Error(), exitStatusErr) {
			// returning non-zero is not an error in this context,
			// just a possible outcome of running an executable, so
			// we are free to overwrite the value of err in the
			// following if statement

			stat := err.Error()[len(exitStatusErr):]
			if exitStatus, err = strconv.Atoi(stat); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("error waiting for '%s': %v",
				cmdLine, err)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &result{
		success_: exitStatus == 0,
		descr_:   cmdLine,
		detail_: fmt.Sprintf(
			"dir\t%s\nstdout\t|%s\nstderr\t|%s\nexit\t%d\nwant %s",
			cwd,
			strings.Replace(output, "\n", "\n\t|", -1),
			strings.Replace(errput, "\n", "\n\t|", -1),
			exitStatus,
			expDescr,
		),
	}, nil
}

func (t *testCmd_) cmd() string {
	s := t.cmd_
	for name, val := range t.vars {
		s = strings.Replace(s, "{"+name+"}", val, -1)
	}
	return s
}

func (t *testCmd_) AddVars(vars map[string]string) {
	for name, val := range vars {
		t.vars[name] = val
	}
}

type testAction interface {
	AddVars(vars map[string]string)
	beforeBake() (*result, error)
	afterBake() (*result, error)
}

type result struct {
	success_ bool
	descr_   string
	detail_  string
}

func (r *result) asExpected() bool {
	return r.success()
}

func (r *result) success() bool {
	return r.success_
}

func (r *result) descr() string {
	return r.descr_
}

func (r *result) detail() string {
	return r.detail_
}

// A command is a support command for tests that should run successfully before
// and after bake is run.
type command struct {
	cmd *testCmd_
}

func newCommand(cmd string) testAction {
	return &command{newTestCmd(cmd)}
}

func (c *command) AddVars(vars map[string]string) {
	c.cmd.AddVars(vars)
}

func (c *command) beforeBake() (*result, error) {
	return c.cmd.Run("exit status = 0 (command test before bake)")
}

func (c *command) afterBake() (*result, error) {
	return c.cmd.Run("exit status = 0 (command test after bake)")
}

type pass struct {
	cmd *testCmd_
}

func newPass(cmd string) testAction {
	return &pass{newTestCmd(cmd)}
}

func (p *pass) AddVars(vars map[string]string) {
	p.cmd.AddVars(vars)
}

func (p *pass) beforeBake() (*result, error) {
	r, err := p.cmd.Run("exit status != 0 (pass test before bake)")

	if err != nil {
		return nil, err
	}
	return &result{!r.success(), r.descr(), r.detail()}, nil
}

func (p *pass) afterBake() (*result, error) {
	return p.cmd.Run("exit status = 0 (pass test after bake)")
}

type buildPass struct {
	cmd *testCmd_
}

func newBuildPass(cmd string) testAction {
	return &buildPass{newTestCmd(cmd)}
}

func (p *buildPass) AddVars(vars map[string]string) {
	p.cmd.AddVars(vars)
}

func (p *buildPass) beforeBake() (*result, error) {
	r, err := p.cmd.Run("error (build pass test before bake)")

	if err != nil {
		return &result{true, p.cmd.cmd(), err.Error()}, nil
	}
	return &result{false, r.descr(), r.detail()}, nil
}

func (p *buildPass) afterBake() (*result, error) {
	return p.cmd.Run("exit status = 0 (build pass test after bake)")
}
