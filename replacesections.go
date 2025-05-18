package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
)

// replaceSections converts (§...) to a link to  #section-(\d+) link.
func replaceSections(w io.Writer, r io.Reader) {
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
		log.Print(err)
	}
}

func openSection(w io.Writer, r *bufio.Reader) (parseSectionFn, error) {
	// look for open character
	head, err := r.ReadBytes('(')
	if err != nil {
		w.Write(head)
		return nil, err
	}
	w.Write(head)
	return sectionRune, nil
}

func sectionRune(w io.Writer, r *bufio.Reader) (parseSectionFn, error) {
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

func closeSection(w io.Writer, r *bufio.Reader) (parseSectionFn, error) {
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

type parseSectionFn func(io.Writer, *bufio.Reader) (parseSectionFn, error)
