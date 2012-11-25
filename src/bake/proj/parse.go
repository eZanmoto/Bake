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
	lDelim = '{'
	rDelim = '}'
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
	var buf bytes.Buffer

	for in.Peek() != scanner.EOF && in.Peek() != rDelim {
		buf.WriteRune(in.Next())
	}

	if in.Peek() == scanner.EOF {
		return parseErr(in, "Expected '%c', encountered EOF", rDelim)
	}

	in.Next() // Consume rDelim

	str := buf.String()
	val, ok := p.vars[str]
	if !ok {
		return parseErr(in, "Unknown variable '"+str+"'")
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
