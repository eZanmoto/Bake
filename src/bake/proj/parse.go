// Copyright 2012 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package proj

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"text/scanner"
)

const (
	lDelim        = '{' // Denotes the start of a template directive
	rDelim        = '}' // Denotes the end of a template directive
	varDepInc     = '?' // Denotes a variable-dependant include
	depsListEnd   = ':' // Denotes the end of a list of dependencies
	depsListDelim = '&' // Delimits elements in a dependency list
)

type NullWriter bool

func (w *NullWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

var (
	stdnil    = new(NullWriter)
	stdnilbuf = bufio.NewWriter(stdnil)
)

func (p *Project) parseStr(src string) (string, error) {
	var out bytes.Buffer
	in := bytes.NewBufferString(src)
	err := p.parse(in, &out)

	if err != nil {
		return "", fmt.Errorf("'%s'%v", src, err)
	}

	return out.String(), nil
}

func (p *Project) parse(reader io.Reader, writer io.Writer) error {
	out := bufio.NewWriter(writer)
	var in scanner.Scanner
	in.Init(reader)

	for in.Peek() != scanner.EOF {
		assignedChar := false
		var char rune

		switch in.Peek() {
		case lDelim:
			in.Next()
			if in.Peek() == lDelim {
				in.Next()
				assignedChar, char = true, lDelim
			} else {
				err := p.parseDirective(&in, out)

				if err != nil {
					return err
				}
			}
		case rDelim:
			in.Next()
			if c := in.Peek(); c == rDelim {
				in.Next()
				assignedChar, char = true, rDelim
			} else {
				return parseErr(&in, "Expected '%c', got '%c'",
					rDelim, c)
			}
		case scanner.EOF:
			break
		default:
			assignedChar, char = true, in.Next()
		}

		if assignedChar {
			if n, err := out.WriteRune(char); n < 1 {
				return fmt.Errorf("Couldn't write: %v", err)
			} else if err != nil {
				return err
			}
		}

		if out.Available() < 1 {
			if err := out.Flush(); err != nil {
				return err
			}
		}
	}

	return out.Flush()
}

func (p *Project) parseDirective(in *scanner.Scanner, out *bufio.Writer) error {
	var err error

	switch in.Peek() {
	case varDepInc:
		in.Next()
		err = p.parseVarInc(in, out)
	default:
		err = p.parseInsert(in, out)
	}

	return err
}

func (p *Project) parseVarInc(in *scanner.Scanner, out *bufio.Writer) error {
	varDeps, err := readDepsList(in)
	if err != nil {
		return err
	}

	allRecognized := true
	for _, name := range varDeps {
		if _, ok := p.vars[name]; !ok {
			allRecognized = false
			break
		}
	}

	if !allRecognized {
		return p.exitDirective(in)
	}

	for in.Peek() != scanner.EOF {
		assignedChar := false
		var char rune

		switch in.Peek() {
		case lDelim:
			in.Next()
			if in.Peek() == lDelim {
				in.Next()
				assignedChar, char = true, lDelim
			} else {
				err := p.parseInsert(in, out)

				if err != nil {
					return err
				}
			}
		case rDelim:
			in.Next()
			if c := in.Peek(); c == rDelim {
				in.Next()
				assignedChar, char = true, rDelim
			} else {
				break
			}
		case scanner.EOF:
			break
		default:
			assignedChar, char = true, in.Next()
		}

		if assignedChar {
			if n, err := out.WriteRune(char); n < 1 {
				return fmt.Errorf("Couldn't write: %v", err)
			} else if err != nil {
				return err
			}
		}

		if out.Available() < 1 {
			if err := out.Flush(); err != nil {
				return err
			}
		}
	}

	return out.Flush()
}

func readDepsList(in *scanner.Scanner) ([]string, error) {
	var buf bytes.Buffer
	deps := make([]string, 0, 1)

	for in.Peek() != scanner.EOF && in.Peek() != depsListEnd {
		c := in.Next()
		switch {
		case c == depsListDelim:
			str := buf.String()
			if len(str) == 0 {
				return nil, parseErr(in, "Empty variable name")
			}
			deps = append(deps, str)
			buf.Reset()
		case 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z':
			buf.WriteRune(c)
		default:
			return nil, parseErr(in, "Unexpected character, '%c'",
				c)
		}
	}

	if in.Peek() != depsListEnd {
		return nil, parseErr(in, "Expected '%c', encountered EOF",
			depsListEnd)
	}
	in.Next() // Consume depsListEnd

	str := buf.String()
	if len(str) == 0 {
		return nil, parseErr(in, "Empty variable name")
	}
	deps = append(deps, str)

	return deps, nil
}

func (p *Project) exitDirective(in *scanner.Scanner) error {
	finished := false

	for in.Peek() != scanner.EOF && !finished {
		next := in.Next()
		if next == lDelim {
			if in.Peek() == lDelim {
				in.Next()
			}
			p.parseInsert(in, stdnilbuf)
		} else if next == rDelim {
			if in.Peek() == rDelim {
				in.Next()
			} else {
				finished = true
			}
		}
	}

	if !finished {
		return parseErr(in, "Expected '%c', encountered EOF", rDelim)
	}

	return nil
}

func (p *Project) parseInsert(in *scanner.Scanner, out *bufio.Writer) error {
	var buf bytes.Buffer

	for in.Peek() != scanner.EOF && in.Peek() != rDelim {
		c := in.Next()
		if !('A' <= c && c <= 'Z' || 'a' <= c && c <= 'z') {
			return parseErr(in, "'%c' in include name", c)
		}
		buf.WriteRune(c)
	}

	if in.Peek() != rDelim {
		return parseErr(in, "Expected '%c', encountered EOF", rDelim)
	}

	in.Next() // Consume rDelim

	name := buf.String()
	val, ok := p.vars[name]
	if !ok {
		return parseErr(in, "Unknown variable '"+name+"'")
	}

	if n, err := out.WriteString(val); n != len(val) {
		return parseErr(in, "Only wrote %d characters of '%s'(%d)",
			n, val, len(val))
	} else if err != nil {
		return nil
	}

	return out.Flush()
}

func parseErr(stream *scanner.Scanner, text string, a ...interface{}) error {
	p := stream.Pos()
	text = fmt.Sprintf(text, a...)
	return fmt.Errorf("%s[%d:%d] %s", p.Filename, p.Line, p.Column, text)
}
