package extractor

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// New creates a new file extractor with the given configuration
func New(config Config) Extractor {
	if config.MaxSize == 0 {
		config.MaxSize = 10 * 1024 * 1024 // 10MB default
	}
	if len(config.Languages) == 0 {
		config.Languages = []string{"bash", "sh", "shell", "go", "python", "yaml", "json"}
	}
	if len(config.Include) == 0 {
		config.Include = []string{"*.md", "*.html", "*.rst"}
	}

	return &FileExtractor{
		config:   config,
		markdown: &MarkdownExtractor{config: config},
		html:     &HTMLExtractor{config: config},
	}
}

// ExtractFromPath extracts examples from a file or directory path
func (e *FileExtractor) ExtractFromPath(path string) (*ExampleSet, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	var allExamples []Example

	if info.IsDir() {
		err = e.walkDirectory(path, &allExamples)
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory %s: %w", path, err)
		}
	} else {
		examples, err := e.ExtractFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to extract from file %s: %w", path, err)
		}
		allExamples = examples
	}

	return &ExampleSet{
		Name:        filepath.Base(path),
		Description: fmt.Sprintf("Examples extracted from %s", path),
		Examples:    allExamples,
		Metadata: map[string]string{
			"source_path":    path,
			"total_examples": fmt.Sprintf("%d", len(allExamples)),
		},
		ExtractedAt: time.Now(),
	}, nil
}

// ExtractFromFile extracts examples from a single file
func (e *FileExtractor) ExtractFromFile(filename string) ([]Example, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", filename, err)
	}

	if info.Size() > e.config.MaxSize {
		return nil, fmt.Errorf("file %s is too large (%d bytes)", filename, info.Size())
	}

	if !e.shouldProcessFile(filename) {
		return nil, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return e.ExtractFromContent(content, filename)
}

// ExtractFromContent extracts examples from content bytes
func (e *FileExtractor) ExtractFromContent(content []byte, source string) ([]Example, error) {
	format := e.detectFormat(content, source)

	switch format {
	case "markdown":
		return e.markdown.ExtractFromContent(content, source)
	case "html":
		return e.html.ExtractFromContent(content, source)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// walkDirectory recursively walks a directory and extracts examples
func (e *FileExtractor) walkDirectory(root string, allExamples *[]Example) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if !e.config.Recursive && path != root {
				return filepath.SkipDir
			}
			return nil
		}

		if !e.shouldProcessFile(path) {
			return nil
		}

		examples, err := e.ExtractFromFile(path)
		if err != nil {
			// Log error but continue processing other files
			fmt.Fprintf(os.Stderr, "Warning: failed to extract from %s: %v\n", path, err)
			return nil
		}

		*allExamples = append(*allExamples, examples...)
		return nil
	})
}

// shouldProcessFile checks if a file should be processed based on include/exclude patterns
func (e *FileExtractor) shouldProcessFile(filename string) bool {
	base := filepath.Base(filename)

	// Check exclude patterns first
	for _, pattern := range e.config.Exclude {
		if matched, _ := filepath.Match(pattern, base); matched {
			return false
		}
	}

	// Check include patterns
	for _, pattern := range e.config.Include {
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
	}

	return false
}

// detectFormat detects the format of the content
func (e *FileExtractor) detectFormat(content []byte, source string) string {
	if e.config.Format != "auto" {
		return e.config.Format
	}

	ext := strings.ToLower(filepath.Ext(source))
	switch ext {
	case ".md", ".markdown":
		return "markdown"
	case ".html", ".htm":
		return "html"
	case ".rst":
		return "rst"
	default:
		// Try to detect from content
		contentStr := string(content)
		if strings.Contains(contentStr, "```") || strings.Contains(contentStr, "~~~") {
			return "markdown"
		}
		if strings.Contains(contentStr, "<html") || strings.Contains(contentStr, "<pre>") {
			return "html"
		}
		return "unknown"
	}
}

// generateExampleID generates a unique ID for an example
func generateExampleID(source, code string, lineStart int) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%d:%s", source, lineStart, code)))
	return fmt.Sprintf("%x", hash)[:12]
}
