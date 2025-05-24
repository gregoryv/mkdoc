package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// emptylines warns if empty lines contains only spaces or tabs
func emptyLines(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)

	var lineno int
	var lastLine string
	for s.Scan() {
		lineno++
		line := s.Text()
		fmt.Fprintln(w, line)

		tmp := strings.TrimSpace(line)
		if len(tmp) == 0 && line != tmp {
			log.Printf("line %v: %s\n%q only space", lineno, lastLine, line)
		}
		lastLine = line
	}
}
