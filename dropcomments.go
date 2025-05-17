package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func dropcomments(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)
	var inside bool
	for s.Scan() {
		line := s.Text()

		if inside {
			j := strings.Index(line, "-->")
			if j >= 0 {
				fmt.Fprintln(w, line[j+3:])
				inside = false
				continue
			}
			continue // skip line
		}

		i := strings.Index(line, "<!--")
		if i >= 0 {
			fmt.Fprint(w, line[:i])
			j := strings.Index(line, "-->")
			if j >= 0 {
				fmt.Fprintln(w, line[j+3:])
				continue
			}
			inside = true
			continue // skip rest of line
		}
		fmt.Fprintln(w, line) // not a comment
	}
}
