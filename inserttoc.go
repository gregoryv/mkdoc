package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func inserttoc(w io.Writer, r, toc io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.ToLower(line) == "table of contents" {
			fmt.Fprintln(w, line)
			fmt.Fprintln(w)
			io.Copy(w, toc)
		} else {
			fmt.Fprintln(w, line)
		}
	}
}
