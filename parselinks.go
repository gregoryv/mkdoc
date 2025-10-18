package stp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseLinks(w io.Writer, r io.Reader) map[string]string {
	s := bufio.NewScanner(r)
	res := make(map[string]string)

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "[") {
			i := strings.Index(line, "]: ")
			if i > 0 {
				key := line[1:i]
				url := line[i+3:]
				res[key] = url
				continue // don't write the link line back
			}
		}
		fmt.Fprintln(w, line)
	}
	return res
}
