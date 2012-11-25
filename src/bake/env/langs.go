// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package env

import "io/ioutil"

// SupportedLangs returns the languages supported by bake
func SupportedLangs() ([]string, error) {
	templatesPath, err := TemplatesPath()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(templatesPath)
	if err != nil {
		return nil, err
	}

	langs := make([]string, len(files))
	numLangs := 0
	for _, file := range files {
		if file.IsDir() {
			langs[numLangs] = file.Name()
			numLangs++
		}
	}

	return langs[:numLangs], nil
}
