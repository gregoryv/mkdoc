package stp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gregoryv/sentences"
)

var maxLineWidth = 69 // 72 - rfc indent of 3

func ListRequirements(w io.Writer, r io.Reader, requirements []string) {
	slices.SortFunc(requirements, func(a, b string) int {
		i := strings.Index(a, " ")
		ra, _ := strconv.ParseInt(a[1:i], 10, 64)
		j := strings.Index(b, " ")
		rb, _ := strconv.ParseInt(b[1:j], 10, 64)
		switch {
		case ra < rb:
			return -1
		case ra > rb:
			return 1
		default:
			return 0
		}
	})

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

// ParseRequirements finds all sentences with requirement identifiers
// starting with '(#R'. The identifier is moved to the front of the
// sentence.
func ParseRequirements(w io.Writer, r io.Reader) []string {
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
	cmd := sentences.Command{
		In:  bufio.NewReader(r1),
		Out: &buf,
	}
	cmd.Run()

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
