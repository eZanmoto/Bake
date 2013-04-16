// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"strio"
)

const (
	exitStatusErr = "exit status "
)

type testAction interface {
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
	cmd string
}

func (c *command) beforeBake() (*result, error) {
	return runCmd(c.cmd)
}

func runCmd(cmdLine string) (*result, error) {
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
			// we are free to overwrite the value in the following
			// if statement

			stat := err.Error()[len(exitStatusErr):]
			if exitStatus, err = strconv.Atoi(stat); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("error waiting for '%s': %v",
				cmdLine, err)
		}
	}

	return &result{
		success_: exitStatus == 0,
		descr_:   cmdLine,
		detail_: fmt.Sprintf("stdout{%s}\nstderr{%s}\nexit status %d",
			output, errput, exitStatus),
	}, nil
}

func (c *command) afterBake() (*result, error) {
	return runCmd(c.cmd)
}

type pass struct {
	cmd string
}

func (p *pass) beforeBake() (*result, error) {
	r, err := runCmd(p.cmd)

	if err != nil {
		return nil, err
	}
	return &result{!r.success(), r.descr(), r.detail()}, nil
}

func (p *pass) afterBake() (*result, error) {
	return runCmd(p.cmd)
}

type buildPass struct {
	cmd string
}

func (p *buildPass) beforeBake() (*result, error) {
	r, err := runCmd(p.cmd)

	if err != nil {
		return &result{true, p.cmd, err.Error()}, nil
	}
	return &result{false, r.descr(), r.detail()}, nil
}

func (p *buildPass) afterBake() (*result, error) {
	return runCmd(p.cmd)
}
