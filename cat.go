package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func cat(stderr, w io.Writer, r io.Reader, delim string) {
	s := bufio.NewScanner(r)
	prefix := "<cat "
	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, prefix) {
			f := line[len(prefix):]
			f = strings.Trim(f, " >")
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprint(stderr, err)
				continue
			}
			io.Copy(w, fh)
			fh.Close()
		} else {
			fmt.Fprintln(w, line)
		}
	}
}
