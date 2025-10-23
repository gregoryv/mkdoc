package stp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// CheckRequirements checks that all requirement level words as specified in
// RFC 2119 are followed by (#R\d+)
func CheckRequirements(stderr, w io.Writer, r io.Reader) {
	in := bufio.NewReader(r)
	var checkNOT bool

	var prev string
	var line string
	var lineno int
	ok := true

	keywords := []string{
		"MUST NOT", "SHALL NOT", "SHOULD NOT", "MUST", "REQUIRED",
		"SHALL", "SHOULD", "RECOMMENDED", "MAY", "OPTIONAL",
	}

	index := make(map[string]int)

	warn := func(v string) {
		if !strings.HasPrefix(v, "(#R") {
			fmt.Fprintln(stderr, prev, line, "line:", lineno, "WARNING! untagged requirement")
			ok = false
		}
		i := strings.Index(v, ")")
		if i > 1 {
			key := v[1:i]
			if prevline, found := index[key]; found {
				fmt.Fprintln(stderr, prev, line, "line:", lineno, "WARNING! duplicate", key, "defined at line:", prevline)
				return
			}
			index[key] = lineno
		}
	}
	indented := regexp.MustCompile(`\s{1,}\"`)
loop:
	for {
		lineno++
		prev = line // save for context when we warn
		line, err := in.ReadString('\n')
		if len(line) == 0 && errors.Is(err, io.EOF) {
			break loop
		}
		fmt.Fprint(w, line)

		// quoted like in explaining the words
		if indented.Match([]byte(line)) {
			continue
		}
		// i.e. MUST\nNOT
		if checkNOT {
			if strings.HasPrefix(line, "NOT") {
				warn(line[3:])
			} else {
				warn(line)
			}
			checkNOT = false
		}
		for _, words := range keywords {
			if i := strings.Index(line, words); i > -1 {
				j := i + len(words)
				if err == nil {
					j++ // +1 for the newline
				}
				if j == len(line) {
					// end of line
					if words == "MUST" || words == "SHALL" || words == "SHOULD" {
						// can be followed by NOT
						checkNOT = true
						continue loop
					}
				}
				warn(line[j:])
				// we assume only one keywords is on each line
				continue loop
			}
		}
	}

	if !ok {
		fmt.Fprintln(stderr, `

Each keyword as defined in RFC 2119 SHOULD be tagged with a
requirement, ie.

  This sentence MUST(#R1) have a tag.`)
	}
}
