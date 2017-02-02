package main

import "strings"

type source struct {
	codes    []string
	filename string
}

func newSrc(filename string, src ...string) source {
	if src == nil {
		src = make([]string, 0)
	}

	return source{
		filename: filename,
		codes:    src,
	}
}

func (s *source) addLine(line string) {
	s.codes = append(s.codes, line)
}

func (s source) String() string {
	return strings.Join(s.codes, "\n")
}
