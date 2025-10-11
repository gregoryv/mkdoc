package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// replacerefs converts lines starting with `[\d+] ...` with named
// anchor.
func replacerefs(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "[") {
			i := strings.Index(line, "] ")
			if i > 0 {
				key := line[1:i]
				rest := line[i+2:]
				if _, err := strconv.Atoi(key); err == nil {
					fmt.Fprintf(w, `[<a id="ref-%s" href="#ref-%s">%s</a>] `, key, key, key)
					fmt.Fprintln(w, rest)
					continue
				}
			}
		}
		fmt.Fprintln(w, line)
	}
}
