package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func rfcindent(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "§") || len(line) == 0 || line[0] == '\t' {
			// no indent
			fmt.Fprintln(w, line)
			continue
		}
		lc := strings.ToLower(line)
		if nonSectionHeader[lc] {
			fmt.Fprintln(w, line)
			continue
		}

		fmt.Fprintln(w, "  ", line)
	}
}

var nonSectionHeader = map[string]bool{
	"status of this memo": true,
	"copyright notice":    true,
	"abstract":            true,
	"table of contents":   true,

	// https://www.rfc-editor.org/old/instructions2authors.txt
	"ipr statement": true,
}
