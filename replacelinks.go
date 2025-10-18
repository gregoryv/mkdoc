package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// replacelinks converts [...] to a link using the given index of
// links. If it's [\d+] it's converted to a #ref-(\d+) link.
func replacelinks(stderr, w io.Writer, r io.Reader, links map[string]string) {
	next := openTag
	in := bufio.NewReader(r)
	var err error
	for {
		next, err = next(stderr, w, in, links)
		if err != nil || next == nil {
			break
		}
	}
	if !errors.Is(err, io.EOF) {
		fmt.Fprint(stderr, err)
	}
}

func openTag(stderr, w io.Writer, r *bufio.Reader, links map[string]string) (parseFn, error) {
	// look for open character
	head, err := r.ReadBytes('[')
	if err != nil {
		w.Write(head)
		return nil, err
	}
	// skip indented text
	i := bytes.LastIndex(head, []byte("\n"))
	if i > 0 {
		w.Write(head[:i])
		head = head[i:]
	}
	if indented.Match(head) {
		w.Write(head)
		return openTag, nil
	}
	if n := len(head); n > 0 && head[n-1] == '[' {
		// skip [
		head = head[:n-1]
	}
	w.Write(head)
	return closeTag, nil
}

func closeTag(stderr, w io.Writer, r *bufio.Reader, links map[string]string) (parseFn, error) {
	// look for open character
	text, err := r.ReadBytes(']')
	if err != nil {
		w.Write(text)
		return nil, fmt.Errorf("missing right square bracket")
	}
	key := string(text[:len(text)-1])
	if strings.HasPrefix(key, "<a") {
		// looks like an anchor already, skip it
		w.Write([]byte{'['})
		w.Write(text)
		return openTag, nil
	}
	if _, err := strconv.Atoi(key); err == nil {
		// looks like an index, use local ref link
		fmt.Fprintf(w, `[<a href="#ref-%s">%s</a>]`, key, key)
		return openTag, nil
	}
	if strings.HasPrefix(key, "#R") {
		fmt.Fprintf(w, `<a href="%s">%s</a>`, key, key[1:])
		return openTag, nil
	}
	url, found := links[key]
	if !found {
		w.Write([]byte{'['})
		w.Write(text)
		fmt.Fprintf(stderr, "missing reference [%s]\n", key)
		return openTag, nil
	}

	fmt.Fprintf(w, `<a href="%s">%s</a>`, url, key)

	return openTag, nil
}

type parseFn func(e, w io.Writer, r *bufio.Reader, _ map[string]string) (parseFn, error)
