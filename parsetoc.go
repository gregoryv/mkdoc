package txtfmt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func parsetoc(w, toc io.Writer, r io.Reader, width int) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "§") {
			i := strings.Index(line, " ")
			if i > 0 {
				s := strings.TrimLeft(line[:i], "§")
				h := strings.TrimSpace(line[i:])
				// print section link to stderr
				// three spaces separation as found in example RFC's
				fmt.Fprintf(toc,
					`<a href="#section-%s">%s</a>   %s`, s, s, h,
				)
				n := width - len(s) - len(h) - 5
				dots := strings.Repeat(".", n)
				fmt.Fprintln(toc, " ", dots)

				// print toc link to stdout
				fmt.Fprintf(w,
					`<a name="section-%s" href="#section-%s">%s</a> %s`, s, s, s, h,
				)
				fmt.Fprintln(w)
			}
		} else {
			fmt.Fprintln(w, line)
		}
	}
}
