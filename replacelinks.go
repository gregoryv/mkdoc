package txtfmt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

func replacelinks(w io.Writer, r io.Reader, links map[string]string) {
	next := openTag
	in := bufio.NewReader(r)
	var err error
	for {
		next, err = next(w, in, links)
		if err != nil || next == nil {
			break
		}
	}
	if !errors.Is(err, io.EOF) {
		log.Print(err)
	}
}

func openTag(w io.Writer, r *bufio.Reader, links map[string]string) (parseFn, error) {
	// look for open character
	head, err := r.ReadBytes('[')
	if err != nil {
		w.Write(head)
		return nil, err
	}
	if n := len(head); n > 0 && head[n-1] == '[' {
		// skip [
		head = head[:n-1]
	}
	w.Write(head)
	return closeTag, nil
}

func closeTag(w io.Writer, r *bufio.Reader, links map[string]string) (parseFn, error) {
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

	url, found := links[key]
	if !found {
		w.Write([]byte{'['})
		w.Write(text)
		log.Printf("missing reference [%s]\n", key)
		return openTag, nil
	}

	fmt.Fprintf(w, `<a href="%s">%s</a>`, url, key)

	return openTag, nil
}

type parseFn func(io.Writer, *bufio.Reader, map[string]string) (parseFn, error)
