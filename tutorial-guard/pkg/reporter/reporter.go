/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// Reporter handles output formatting for extracted examples
type Reporter struct {
	outputFile string
}

// New creates a new reporter with the specified output file
func New(outputFile string) *Reporter {
	return &Reporter{
		outputFile: outputFile,
	}
}

// WriteExamples writes extracted examples to the output file
func (r *Reporter) WriteExamples(examples *extractor.ExampleSet) error {
	data, err := json.MarshalIndent(examples, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal examples: %w", err)
	}

	if err := os.WriteFile(r.outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", r.outputFile, err)
	}

	fmt.Printf("Extracted %d examples to %s\n", len(examples.Examples), r.outputFile)
	return nil
}
