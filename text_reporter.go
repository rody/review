package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type textReporter struct {
	w io.Writer
}

func newTextReporter(w io.Writer) *textReporter {
	return &textReporter{
		w,
	}
}

func (tr textReporter) report(violations []violation) error {
	if len(violations) == 0 {
		// nothing to print
		return nil
	}

	slices.SortFunc(violations, func(a, b violation) int {
		return int(a.severity) - int(b.severity)
	})

	for _, v := range violations {
		fmt.Fprintf(tr.w, "%s\t%s\t%s:%d\t%s\n", strings.ToUpper(v.severity.String()), v.rule, v.filename, v.line, v.message)
	}

	return nil
}
