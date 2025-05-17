package txtfmt

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
		if strings.HasPrefix(line, "§") {
			fmt.Fprintln(w, line)
			continue
		} else {
			lc := strings.ToLower(line)
			if nonSectionHeader[lc] {
				fmt.Fprintln(w, line)
				continue
			}
		}
		fmt.Fprintln(w, "  ", line)
	}
}

var nonSectionHeader = map[string]bool{
	"status of this memo": true,
	"copyright notice":    true,
	"abstract":            true,
	"table of contents":   true,
}
