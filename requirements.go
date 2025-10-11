package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

// replaceRequirements converts (#R\d+) to a link to  #requirement-(\d+) link.
func replaceRequirements(w io.Writer, r io.Reader) {
	next := openRequirement
	in := bufio.NewReader(r)
	var err error
	for {
		next, err = next(w, in)
		if err != nil || next == nil {
			break
		}
	}
	if !errors.Is(err, io.EOF) {
		log.Print(err)
	}
}

func openRequirement(w io.Writer, r *bufio.Reader) (pipeFn, error) {
	// look for open character
	head, err := r.ReadBytes('(')
	if err != nil {
		w.Write(head)
		return nil, err
	}
	n := len(head)
	if n > 0 {
		w.Write(head[:n-1]) // skip '('
	}
	return requirementRune, nil
}

func requirementRune(w io.Writer, r *bufio.Reader) (pipeFn, error) {
	v, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	}

	if v == '#' {
		return closeRequirement, nil
	}
	w.Write([]byte("(" + string(v)))
	return openRequirement, nil
}

func closeRequirement(w io.Writer, r *bufio.Reader) (pipeFn, error) {
	// look for open character
	text, err := r.ReadBytes(')')
	if err != nil {
		w.Write(text)
		return nil, fmt.Errorf("missing right square parenthesis")
	}
	key := string(text[:len(text)-1])
	no := strings.TrimPrefix(key, "R")
	fmt.Fprintf(w, `<sub><a id="%s" href="#%s">%s</a></sub>`, key, key, no)

	return openRequirement, nil
}

type pipeFn func(io.Writer, *bufio.Reader) (pipeFn, error)
