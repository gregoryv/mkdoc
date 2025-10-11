package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func parsetoc(w, toc io.Writer, r io.Reader, width int) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		s, h, ok := parseSection(line)
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		// print section link to stderr
		// three spaces separation as found in example RFC's
		id := fmt.Sprintf("section-%s", s)
		fmt.Fprintf(toc,
			`<a href="#%s">%s</a>   <a href="#%s">%s</a>`, id, s, id, h,
		)
		// fill remainig width using dots
		n := width - len(s) - len(h) - 5
		dots := strings.Repeat(".", n)
		fmt.Fprintln(toc, " ", dots)

		// just forward the line untouched, see linksections
		fmt.Fprintln(w, line)
	}
}

func linksections(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		s, h, ok := parseSection(line)
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		// print toc link to toc
		format := `<a name="section-%s" href="#section-%s">%s</a> %s`
		fmt.Fprintf(w, format, s, s, s, h)
		fmt.Fprintln(w)
	}
}

func parseSection(line string) (s, h string, ok bool) {
	if !strings.HasPrefix(line, "§") {
		return
	}
	// find section identifier
	i := strings.Index(line, " ")
	if i == -1 {
		log.Print("WARNING! section has no identifier")
		return
	}
	s = strings.TrimLeft(line[:i], "§")
	h = strings.TrimSpace(line[i:])
	ok = true
	return
}

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
