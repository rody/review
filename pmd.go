package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type pmdReport struct {
	Runs []struct {
		Results []struct {
			Level   string `json:"level"`
			RuleId  string `json:"ruleId"`
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
			Locations []struct {
				PhysicalLocation struct {
					ArtifactLocation struct {
						Uri string `json:"uri"`
					} `json:"artifactLocation"`
					Region struct {
						StartLine int `json:"startLine"`
						EndLine   int `json:"endLine"`
					} `json:"region"`
				} `json:"physicalLocation"`
			} `json:"locations"`
		} `json:"results"`
	} `json:"runs"`
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

	for _, f := range r.Runs {
		for _, x := range f.Results {
			for _, l := range x.Locations {
				result = append(result, violation{
					rule:     x.RuleId,
					line:     l.PhysicalLocation.Region.StartLine,
					filename: stripPath(l.PhysicalLocation.ArtifactLocation.Uri, "force-app"),
					message:  x.Message.Text,
					severity: x.Level,
				})
			}

		}
	}

	return result, nil
}

func stripPath(input, keyword string) string {
	index := strings.Index(input, keyword)
	if index != -1 {
		return input[index:]
	}
	return input
}
