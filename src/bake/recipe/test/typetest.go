// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package test

import (
	"errors"
	"fmt"
	"io"
	"strio"
)

const (
	testDirectiveIndex = 0
)

func readTypeTests(in strio.LineReader) ([]*typeTest, error) {
	var err error

	tests := make([]*typeTest, 0, 1)
	for {
		var descr string
		descr, err = in.ChompLine()
		if err != nil {
			if err == io.EOF {
				err = errors.New("expected commands, got EOF")
			}
			break
		}

		if len(descr) == 0 {
			err = errors.New("description cannot be empty")
			break
		}

		var actions []testAction
		actions, err = readTypeTestCommands(in)
		tests = append(tests, &typeTest{descr, actions})

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}

	if err == nil && len(tests) == 0 {
		err = fmt.Errorf("didn't read any tests")
	}
	return tests, err
}

type typeTest struct {
	descr_   string
	actions_ []testAction
}

func (t *typeTest) descr() string {
	return t.descr_
}

func (t *typeTest) actions() []testAction {
	return t.actions_
}

func readTypeTestCommands(in strio.LineReader) ([]testAction, error) {
	var err error
	actions := make([]testAction, 0, 1)

	for {
		var line string
		line, err = in.ChompLine()

		if len(line) == 0 {
			if err == io.EOF {
				err = errors.New("superfluous newline at EOF")
			}
			break
		}

		if err != nil && err != io.EOF {
			break
		}
		// inv: len(line) > 0 && (err == nil || err == io.EOF)

		cmd := line[testDirectiveIndex+1:]

		action := parseTestAction(rune(line[testDirectiveIndex]), cmd)
		if action == nil {
			err = fmt.Errorf("'%c' is not a valid test directive",
				line[testDirectiveIndex])
			break
		}
		actions = append(actions, action)

		if err == io.EOF {
			break
		}
	}

	if err == nil && len(actions) == 0 {
		err = fmt.Errorf("each test must have commands")
	}

	return actions, err
}

func parseTestAction(actionSpecifier rune, cmd string) testAction {
	var action testAction

	// this is the only place the test directive constants are used, so they
	// are defined using magic constants since it is easier to link the
	// symbols to the actions they should perform
	switch actionSpecifier {
	case ' ':
		action = &command{cmd}
	case '+':
		action = &pass{cmd}
	case '=':
		action = &buildPass{cmd}
	default:
		action = nil
	}

	return action
}
