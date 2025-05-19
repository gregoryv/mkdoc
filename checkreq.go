package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// checkreq checks that all requirement level words as specified in
// RFC 2119 are followed by (#R\d+)
func checkreq(e, w io.Writer, r io.Reader) {
	s := bufio.NewScanner(r)
	var checkNOT bool

	var prev string
	var line string
	var lineno int
	ok := true

	keywords := []string{
		"MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
		"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", "OPTIONAL",
	}

	index := make(map[string]int)

	warn := func(v string) {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, `"`) {
			// quoted like in explaining the words
			return
		}
		if !strings.HasPrefix(v, "(#R") {
			fmt.Fprintln(e, prev, line, "line:", lineno, "WARNING! untagged requirement")
			ok = false
		}
		i := strings.Index(v, ")")
		if i > 1 {
			key := v[1:i]
			if prevline, found := index[key]; found {
				fmt.Fprintln(e, prev, line, "line:", lineno, "WARNING! duplicate", key, "defined at line:", prevline)
				return
			}
			index[key] = lineno
		}
	}

loop:
	for s.Scan() {
		lineno++
		prev = line // save for context when we warn
		line = s.Text()
		fmt.Fprintln(w, line)

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
				if j == len(line) {
					// end of line
					if words == "MUST" || words == "SHALL" || words == "SHOULD" {
						// can be followed by NOT
						checkNOT = true
						continue loop
					}
				}
				warn(line[j:])
			}
		}
	}

	if !ok {
		fmt.Fprintln(e, `

Each keyword as defined in RFC 2119 SHOULD be tagged with a
requirement, ie.

  This sentence MUST(#R1) have a tag.`)
	}
}
