// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package proj provides project generation functionality.
package proj

type Project struct {
	lang    string
	types   []string
	verbose bool
	vars    map[string]string
}

func New(lang string, types []string, v bool, vars map[string]string) Project {
	return Project{lang, types, v, vars}
}

func (p *Project) IsOfType(t string) bool {
	for _, val := range p.types {
		if val == t {
			return true
		}
	}
	return false
}
