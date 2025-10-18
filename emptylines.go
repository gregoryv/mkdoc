package stp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// EmptyLines warns if empty lines only contain spaces or tabs.
func EmptyLines(stderr, w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)

	var lineno int
	var lastLine string
	for s.Scan() {
		lineno++
		line := s.Text()
		fmt.Fprintln(w, line)

		tmp := strings.TrimSpace(line)
		if len(tmp) == 0 && line != tmp {
			fmt.Fprintf(
				stderr,
				"line %v: %s\n%q only space",
				lineno, lastLine, line,
			)
		}
		lastLine = line
	}
}
