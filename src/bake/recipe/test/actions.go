// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"os/exec"
	"strings"
	"strio"
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

	err = cmd.Wait()

	return &result{
		success_: err == nil,
		descr_:   cmdLine,
		detail_:  fmt.Sprintf("stdout{%s}\nstderr{%s}", output, errput),
	}, err
}

func (c *command) afterBake() (*result, error) {
	return runCmd(c.cmd)
}
