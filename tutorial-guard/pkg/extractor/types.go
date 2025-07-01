/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package extractor

import (
	"time"
)

// Example represents a single code example extracted from documentation
type Example struct {
	ID            string            `json:"id"`
	Source        string            `json:"source"`        // File path or URL
	Language      string            `json:"language"`      // bash, go, python, etc.
	Code          string            `json:"code"`          // The actual code content
	Context       string            `json:"context"`       // Surrounding documentation
	Tags          []string          `json:"tags"`          // Categories: setup, demo, cleanup
	Prerequisites []string          `json:"prerequisites"` // Dependencies on other examples
	Metadata      map[string]string `json:"metadata"`      // Additional metadata
	LineStart     int               `json:"line_start"`    // Starting line in source file
	LineEnd       int               `json:"line_end"`      // Ending line in source file
}

// ExampleSet represents a collection of related examples
type ExampleSet struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Examples    []Example         `json:"examples"`
	Metadata    map[string]string `json:"metadata"`
	ExtractedAt time.Time         `json:"extracted_at"`
}

// Config holds configuration for the extractor
type Config struct {
	Format    string   // markdown, html, auto
	Recursive bool     // scan directories recursively
	Include   []string // file patterns to include
	Exclude   []string // file patterns to exclude
	MaxSize   int64    // maximum file size to process (bytes)
	Languages []string // code languages to extract
}

// Extractor interface defines the contract for extracting examples
type Extractor interface {
	ExtractFromFile(filename string) ([]Example, error)
	ExtractFromPath(path string) (*ExampleSet, error)
	ExtractFromContent(content []byte, source string) ([]Example, error)
}

// MarkdownExtractor handles Markdown files
type MarkdownExtractor struct {
	config Config
}

// HTMLExtractor handles HTML files
type HTMLExtractor struct {
	config Config
}

// FileExtractor is the main extractor that delegates to format-specific extractors
type FileExtractor struct {
	config   Config
	markdown *MarkdownExtractor
	html     *HTMLExtractor
}
