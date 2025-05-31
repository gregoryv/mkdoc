package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/gregoryv/sentences"
)

var maxLineWidth = 69 // 72 - rfc indent of 3

func includeReq(w io.Writer, r io.Reader, requirements []string) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.ToLower(line) == "<list of requirements>" {
			fmt.Fprintln(w)
			for _, req := range requirements {
				// find first space, ie. after "R101 ..."
				i := strings.Index(req, " ")

				// link the requirement id
				id := req[:i]
				fmt.Fprintf(w, `<a href="#%s">%s</a>`, id, id)

				if len(req) < maxLineWidth {
					fmt.Fprintln(w, req[i:])
					fmt.Fprintln(w)
					continue
				}

				// find last space
				j := strings.LastIndex(req[:maxLineWidth], " ")
				fmt.Fprintln(w, req[i:j])
				// indentation
				fmt.Fprint(w, strings.Repeat(" ", i))
				// the rest
				fmt.Fprintln(w, req[j:])
				fmt.Fprintln(w)
			}
		} else {
			fmt.Fprintln(w, line)
		}
	}
}

func parsereq(w io.Writer, r io.Reader) []string {
	// use a pipe to parse sentences and just copy the data
	r1, w1 := io.Pipe()
	wboth := io.MultiWriter(w, w1)
	// start the copy
	go func() {
		io.Copy(wboth, r)
		w1.Close()
	}()

	// parse sentences
	var buf bytes.Buffer
	cmd := sentences.Command{
		In:  bufio.NewReader(r1),
		Out: &buf,
	}
	cmd.Run()

	// filter sentences that look like requirements
	s := bufio.NewScanner(&buf)
	res := make([]string, 0)
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "(#R") {
			// move the (#R...) to front of line
			res = append(res, moveTagToFront(line))
		}
	}
	return res
}

var re = regexp.MustCompile(`^(.*?\b)\(#(R\d+)\)(.*)$`)

func moveTagToFront(input string) string {
	return re.ReplaceAllString(input, "$2 $1$3")
}
