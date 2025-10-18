package stp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseTOC(stderr, w, toc io.Writer, r io.Reader, width int) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		s, h, ok := parseSection(stderr, line)
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		// print section link to stderr
		// three spaces separation as found in example RFC's
		fmt.Fprintf(toc,
			`<a href="#section-%s">%s</a>   %s`, s, s, h,
		)
		// fill remainig width using dots
		n := width - len(s) - len(h) - 5
		dots := strings.Repeat(".", n)
		fmt.Fprintln(toc, " ", dots)

		// just forward the line untouched, see linksections
		fmt.Fprintln(w, line)
	}
}

func LinkSections(stderr, w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		s, h, ok := parseSection(stderr, line)
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		// print toc link to toc
		format := `<a id="section-%s" href="#section-%s">%s</a> %s`
		fmt.Fprintf(w, format, s, s, s, h)
		fmt.Fprintln(w)
	}
}

func parseSection(stderr io.Writer, line string) (s, h string, ok bool) {
	if !strings.HasPrefix(line, "§") {
		return
	}
	// find section identifier
	i := strings.Index(line, " ")
	if i == -1 {
		fmt.Fprint(stderr, "WARNING! section has no identifier")
		return
	}
	s = strings.TrimLeft(line[:i], "§")
	h = strings.TrimSpace(line[i:])
	ok = true
	return
}

func InsertTOC(w io.Writer, r, toc io.Reader) {
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
