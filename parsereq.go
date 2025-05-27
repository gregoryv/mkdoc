package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func includeReq(w io.Writer, r io.Reader, requirements []string) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.ToLower(line) == "<list of requirements>" {
			fmt.Fprintln(w)
			for _, req := range requirements {
				fmt.Fprintln(w, req)
			}
		} else {
			fmt.Fprintln(w, line)
		}
	}
}

func parsereq(w io.Writer, r io.Reader) []string {
	s := bufio.NewScanner(r)
	res := make([]string, 0)
	for s.Scan() {
		line := s.Text()
		fmt.Fprintln(w, line)
		if strings.Contains(line, "(#R") {
			res = append(res, moveTagToFront(line))
		}
	}
	return res
}

var re = regexp.MustCompile(`^(.*?\b)\(#(R\d+)\)(.*)$`)

func moveTagToFront(input string) string {
	return re.ReplaceAllString(input, "<a href=\"#$2\">$2</a> $1$3")
}

func sentences(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	// Set the split function for the scanning operation.
	scanner.Split(ScanSentences)

	var i int

	for scanner.Scan() {
		i++
		line := scanner.Bytes()
		if i := bytes.LastIndex(line, doubleNL); i > -1 {
			// found an empty line, this is normal after headings
			// ignore heading
			line = line[i+2:]
		}
		// for some reason
		// line = bytes.ReplaceAll(line, oneNL, oneSpace)
		// results in many allocations during benchmarks
		for j := range line {
			if line[j] == nl {
				line[j] = ' '
			}
		}
		line = bytes.TrimSpace(line)
		if len(line) > 1 { // one character followed by ., ? or !
			w.Write(line)
			w.Write(oneNL)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

var (
	nl       = byte('\n')
	oneNL    = []byte{nl}
	oneSpace = []byte{' '}
	doubleNL = []byte{nl, nl}
)

// ScanSentence is a split function for a Scanner that returns
// sentence.
func ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip until first capital letter is found
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if unicode.IsUpper(r) {
			break
		}
	} // capital letter found

	// find what looks like end of sentence.
	var width int
	for i := start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		switch r {
		case '.', '?', '!':
			return i + width, data[start : i+width], nil
			break
		}
	}
	// Request more data.
	return start, nil, nil
}
