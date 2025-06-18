package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

// sentenceSpace warns if two adjecent sentences are not separated by
// double spaces.
func sentenceSpace(w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)

	var lineno int
	for s.Scan() {
		lineno++
		line := s.Text()
		fmt.Fprintln(w, line)

		if !strings.Contains(line, ". ") || strings.HasPrefix(line, "§") {
			continue
		}
		ends := strings.Split(line, ". ")
		for _, end := range ends[1:] {
			if len(end) == 0 {
				continue
			}
			if end[0] == ' ' {
				continue
			}
			r, _ := utf8.DecodeRuneInString(end)
			if unicode.IsUpper(r) {
				log.Printf("line %v: %s\nmissing double space between sentences", lineno, line)
			}
		}
	}
}
