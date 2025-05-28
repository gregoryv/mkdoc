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

var maxLineWidth = 69 // 72 - rfc indent of 3

func includeReq(w io.Writer, r io.Reader, requirements []string) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.ToLower(line) == "<list of requirements>" {
			fmt.Fprintln(w)
			for _, req := range requirements {
				// find first space, ie. after "R101 ..."
				i := strings.Index(req, " ")

				// link the requirement id
				id := req[:i]
				fmt.Fprintf(w, `<a href="#%s">%s</a>`, id, id)

				if len(req) < maxLineWidth {
					fmt.Fprintln(w, req[i:])
					fmt.Fprintln(w)
					continue
				}

				// find last space
				j := strings.LastIndex(req[:maxLineWidth], " ")
				fmt.Fprintln(w, req[i:j])
				// indentation
				fmt.Fprint(w, strings.Repeat(" ", i))
				// the rest
				fmt.Fprintln(w, req[j:])
				fmt.Fprintln(w)
			}
		} else {
			fmt.Fprintln(w, line)
		}
	}
}

func parsereq(w io.Writer, r io.Reader) []string {
	// use a pipe to parse sentences and just copy the data
	r1, w1 := io.Pipe()
	wboth := io.MultiWriter(w, w1)
	// start the copy
	go func() {
		io.Copy(wboth, r)
		w1.Close()
	}()

	// parse sentences
	var buf bytes.Buffer
	sentences(&buf, r1)

	// filter sentences that look like requirements
	s := bufio.NewScanner(&buf)
	res := make([]string, 0)
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "(#R") {
			// move the (#R...) to front of line
			res = append(res, moveTagToFront(line))
		}
	}
	return res
}

var re = regexp.MustCompile(`^(.*?\b)\(#(R\d+)\)(.*)$`)

func moveTagToFront(input string) string {
	return re.ReplaceAllString(input, "$2 $1$3")
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
