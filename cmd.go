package main

import "fmt"

type violation struct {
	rule     string
	line     int
	filename string
	message  string
	severity string
}

type linter interface {
	violations() ([]violation, error)
}

type reporter interface {
	report(violations []violation) error
}

type contentChecker interface {
	contains(filename string, line int) bool
}

type command struct {
	linter   linter
	content  contentChecker
	reporter reporter
}

func (cmd *command) run() (int, error) {
	var toReport []violation
	violations, err := cmd.linter.violations()
	if err != nil {
		return 0, fmt.Errorf("could not get rule violations: %w", err)
	}

	violationCount := 0
	for _, v := range violations {
		if cmd.content.contains(v.filename, v.line) {
			violationCount++
			toReport = append(toReport, v)
		}
	}
	return violationCount, cmd.reporter.report(toReport)
}
