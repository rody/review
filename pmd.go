package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type pmdReport struct {
	Version string `json:"pmdVersion"`
	Files   []struct {
		Filename   string `json:"filename"`
		Violations []struct {
			Line        int    `json:"beginLine"`
			Description string `json:"description"`
			Rule        string `json:"rule"`
			Priority    int    `json:"priority"`
		} `json:"violations"`
	} `json:"files"`
}

func readPMDFile(path string) (*pmdReport, error) {
	var report pmdReport

	f, err := os.Open(path)
	if err != nil {
		return &report, fmt.Errorf("could not open pmd report: %w", err)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&report); err != nil {
		return &report, fmt.Errorf("could not parse pmd report: %w", err)
	}

	return &report, err
}

func (r *pmdReport) violations() ([]violation, error) {
	var result []violation
	for _, f := range r.Files {
		for _, v := range f.Violations {
			result = append(result, violation{
				rule:     v.Rule,
				line:     v.Line,
				filename: f.Filename,
				message:  v.Description,
				severity: toSeverity(v.Priority),
			})
		}
	}

	return result, nil
}

func toSeverity(pmdPriority int) severity {
	switch pmdPriority {
	case 1:
		return severityError
	case 2, 3:
		return severityWarning
	default:
		return severityInfo
	}
}
