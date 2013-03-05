// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package proj provides project generation functionality.
package proj

type Project struct {
	lang    string
	types   []string
	verbose bool
	resolve bool
	vars    map[string]string
}

func New(lg string, ts []string, v bool, r bool, vs map[string]string) Project {
	return Project{lg, ts, v, r, vs}
}

func (p *Project) IsOfType(t string) bool {
	for _, val := range p.types {
		if val == t {
			return true
		}
	}
	return false
}
