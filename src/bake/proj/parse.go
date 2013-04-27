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
	lDelim    = '{' // Denotes the start of a template directive
	rDelim    = '}' // Denotes the end of a template directive
	condStart = '?' // Denotes the start of a conditional insert
	condElsif = ':' // Denotes the else of a conditional insert
	condEnd   = '!' // Denotes the end of a conditional insert
)

type parseStatus byte

const (
	ignoreOutput parseStatus = 0
	genOutput    parseStatus = 1
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

	err := p.parseText(&in, out, genOutput)
	if err != nil {
		return err
	}

	if in.Peek() != scanner.EOF {
		return parseErr(&in, "Unexpected closing statement: '%c'",
			in.Peek())
	}

	return out.Flush()
}

func parseErr(stream *scanner.Scanner, text string, a ...interface{}) error {
	p := stream.Pos()
	text = fmt.Sprintf(text, a...)
	return fmt.Errorf("%s[%d:%d] %s", p.Filename, p.Line, p.Column, text)
}

// Parses text from in to out until it reaches EOF, reads '{:' or '{!'
func (p *Project) parseText(in *scanner.Scanner, out *bufio.Writer, s parseStatus) error {
	for finished := isEOF(in); !finished; finished = finished || isEOF(in) {
		readChar := false

		switch in.Peek() {
		case lDelim:
			in.Next()

			switch in.Peek() {
			case scanner.EOF:
				return parseErr(in,
					"Expected directive, got EOF")
			case lDelim:
				readChar = true
			case condEnd, condElsif:
				finished = true
			default:
				if err := p.parseDirect(in, out, s); err != nil {
					return err
				}
			}
		case rDelim:
			in.Next()

			switch in.Peek() {
			case scanner.EOF:
				return parseErr(in, "Escape '%c', got EOF",
					rDelim)
			case rDelim:
				readChar = true
			default:
				return parseErr(in, "Escape '%c', got '%c'",
					rDelim, in.Peek())
			}
		default:
			readChar = true
		}

		if readChar {
			c := in.Next()
			if s == genOutput {
				if n, err := out.WriteRune(c); err == nil && n < 1 {
					return fmt.Errorf("Couldn't write: %c", c)
				} else if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isEOF(in *scanner.Scanner) bool {
	return in.Peek() == scanner.EOF
}

func (p *Project) parseDirect(in *scanner.Scanner, out *bufio.Writer, s parseStatus) error {
	var err error

	if in.Peek() == condStart {
		err = p.parseCond(in, out, s)
	} else {
		err = p.parseVar(in, out, s)
	}

	return err
}

func (p *Project) parseVar(in *scanner.Scanner, out *bufio.Writer, s parseStatus) error {
	name, err := parseName(in)
	if err != nil {
		return err
	}

	if err = match(in, rDelim); err != nil {
		return err
	}

	if s == genOutput {
		val, ok := p.vars[name]
		if !ok {
			return parseErr(in, "Unknown variable '%s'", name)
		}

		if n, err := out.WriteString(val); n != len(val) {
			return parseErr(in, "Only wrote %d characters of '%s'(%d)",
				n, val, len(val))
		} else if err != nil {
			return err
		}
	}

	return out.Flush()
}

func (p *Project) parseCond(in *scanner.Scanner, out *bufio.Writer, s parseStatus) error {
	hasIncluded := false

	for in.Peek() != condEnd && in.Peek() != scanner.EOF {
		include, err := p.evalCondBool(in)
		if err != nil {
			return err
		}
		skipIfNewline(in)

		if !hasIncluded && include {
			if err := p.parseText(in, out, s); err != nil {
				return err
			}
			hasIncluded = true
		} else {
			if err := p.parseText(in, out, ignoreOutput); err != nil {
				return err
			}
		}
	}

	if err := match(in, condEnd); err != nil {
		return err
	}
	if err := match(in, rDelim); err != nil {
		return err
	}
	skipIfNewline(in)

	return nil
}

// TODO add to custom scanner
func skipIfNewline(in *scanner.Scanner) {
	if in.Peek() == '\n' {
		in.Next()
	}
}

// TODO add to custom scanner
func match(in *scanner.Scanner, r rune) error {
	var err error

	if c := in.Next(); c == scanner.EOF {
		err = parseErr(in, "Expected '%c', got EOF", r)
	} else if c != r {
		err = parseErr(in, "Expected '%c', got '%c'", r, c)
	}

	return err
}

func (p *Project) evalCondBool(in *scanner.Scanner) (eval bool, err error) {
	if c := in.Next(); c == scanner.EOF {
		err = parseErr(in, "Expected '%c' or '%c', got EOF",
			condStart, condElsif)
		return
	} else if c == condStart {
		if in.Peek() == rDelim {
			err = parseErr(in, "Start conditional can't be empty")
			return
		}
	} else if c != condElsif {
		err = parseErr(in, "Expected '%c' or '%c', got '%c'",
			condStart, condElsif, c)
		return
	}

	eval, err = p.evalCondName(in)
	if err != nil {
		return
	}
	err = match(in, rDelim)

	return
}

func (p *Project) evalCondName(in *scanner.Scanner) (eval bool, err error) {
	name, err := parseName(in)
	if err != nil {
		return
	}

	if 'a' <= name[0] && name[0] <= 'z' {
		eval = p.IsOfType(name)
	} else {
		eval = p.hasVar(name)
	}

	return
}

func (p *Project) hasVar(name string) bool {
	_, exists := p.vars[name]
	return exists
}

func parseName(in *scanner.Scanner) (s string, err error) {
	var buf bytes.Buffer

	for isAlpha(in.Peek()) {
		buf.WriteRune(in.Next())
	}

	s = buf.String()
	if len(s) == 0 {
		err = parseErr(in, "Expected name")
	}

	return
}

func isAlpha(r rune) bool {
	return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z'
}
