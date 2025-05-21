package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// parsetoc does two things; writes named sections to w and the table
// of contents to toc.
func parsetoc(w, toc io.Writer, r io.Reader, width int) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "§") {
			fmt.Fprintln(w, line)
			continue
		}
		// find section identifier
		i := strings.Index(line, " ")
		if i == -1 {
			log.Print("WARNING! section has no identifier")
			continue
		}
		s := strings.TrimLeft(line[:i], "§")

		h := strings.TrimSpace(line[i:])
		// print section link to stderr
		// three spaces separation as found in example RFC's
		fmt.Fprintf(toc,
			`<a href="#section-%s">%s</a>   %s`, s, s, h,
		)
		// fill remainig width using dots
		n := width - len(s) - len(h) - 5
		dots := strings.Repeat(".", n)
		fmt.Fprintln(toc, " ", dots)

		// print toc link to toc
		format := `<a name="section-%s" href="#section-%s">%s</a> %s`
		fmt.Fprintf(w, format, s, s, s, h)
		fmt.Fprintln(w)
	}
}

func linksections(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "§") {
			fmt.Fprintln(w, line)
			continue
		}
		// find section identifier
		i := strings.Index(line, " ")
		if i == -1 {
			log.Print("WARNING! section has no identifier")
			continue
		}
		s := strings.TrimLeft(line[:i], "§")

		h := strings.TrimSpace(line[i:])

		// print toc link to toc
		format := `<a name="section-%s" href="#section-%s">%s</a> %s`
		fmt.Fprintf(w, format, s, s, s, h)
		fmt.Fprintln(w)
	}
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
