// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package env

import (
	"errors"
	"os"
	"path"
)

const (
	bakeVar      = "BAKE"      // The name of the bake environment variable
	templatesDir = "templates" // The directory containing bake templates
)

func templatesPath() (string, error) {
	bakeDir := os.Getenv(bakeVar)
	if len(bakeDir) == 0 {
		return "", errors.New("bake environment variable not set")
	}

	templatesPath := path.Join(bakeDir, templatesDir)
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		return "", errors.New("bake root doesn't contain templates")
	}

	return templatesPath, nil
}
