package main

import (
	"fmt"
	"io"
	"os"

	"github.com/waigani/diffparser"
)

type diffChecker struct {
	changes map[string][]int
}

func readDiffFile(path string) (*diffChecker, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readDiff(f)
}

func readDiff(r io.Reader) (*diffChecker, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read diff: %w", err)
	}

	return newDiffChecker(string(content))
}

func newDiffChecker(content string) (*diffChecker, error) {
	d, err := diffparser.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("could not parse diff: %w", err)
	}

	return &diffChecker{
		changes: d.Changed(),
	}, nil
}

func (d *diffChecker) contains(filename string, line int) bool {
	lines, found := d.changes[filename]
	if !found {
		return false
	}

	for _, l := range lines {
		if l == line {
			return true
		}
	}

	return false
}
