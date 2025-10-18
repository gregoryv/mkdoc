package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
)

// replaceSections converts (§...) to a link to  #section-(\d+) link.
func replaceSections(stderr, w io.Writer, r io.Reader) {
	next := openSection
	in := bufio.NewReader(r)
	var err error
	for {
		next, err = next(w, in)
		if err != nil || next == nil {
			break
		}
	}
	if !errors.Is(err, io.EOF) {
		fmt.Fprint(stderr, err)
	}
}

func openSection(w io.Writer, r *bufio.Reader) (pipeFn, error) {
	// look for open character
	head, err := r.ReadBytes('(')
	if err != nil {
		w.Write(head)
		return nil, err
	}
	i := bytes.LastIndex(head, []byte("\n"))
	if i > 0 {
		w.Write(head[:i])
		head = head[i:]
	}
	if indented.Match(head) {
		w.Write(head)
		return openSection, nil
	}
	w.Write(head)
	return sectionRune, nil
}

// our indentation is 3, so one more is considered indented
var indented = regexp.MustCompile(`\n\s{4,}\"`)

func sectionRune(w io.Writer, r *bufio.Reader) (pipeFn, error) {
	v, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	}

	if v == '§' {
		return closeSection, nil
	}
	w.Write([]byte(string(v)))
	return openSection, nil
}

func closeSection(w io.Writer, r *bufio.Reader) (pipeFn, error) {
	// look for open character
	text, err := r.ReadBytes(')')
	if err != nil {
		w.Write(text)
		return nil, fmt.Errorf("missing right square parenthesis")
	}
	key := string(text[:len(text)-1])
	fmt.Fprintf(w, `<a href="#section-%s">section %s</a>)`, key, key)

	return openSection, nil
}
