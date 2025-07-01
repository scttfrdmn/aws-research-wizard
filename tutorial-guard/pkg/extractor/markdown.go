package extractor

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	// Regular expressions for different code block formats
	fencedCodeBlockRegex = regexp.MustCompile("^```([a-zA-Z0-9_+-]*)")
	tildeCodeBlockRegex  = regexp.MustCompile("^~~~([a-zA-Z0-9_+-]*)")
	indentedCodeRegex    = regexp.MustCompile("^    (.*)$")
)

// ExtractFromContent extracts code examples from Markdown content
func (m *MarkdownExtractor) ExtractFromContent(content []byte, source string) ([]Example, error) {
	var examples []Example
	lines := strings.Split(string(content), "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check for fenced code blocks (```)
		if match := fencedCodeBlockRegex.FindStringSubmatch(line); match != nil {
			example, endLine, err := m.extractFencedBlock(lines, i, "```", match[1], source)
			if err != nil {
				continue // Skip malformed blocks
			}
			if example != nil {
				examples = append(examples, *example)
			}
			i = endLine
			continue
		}

		// Check for tilde code blocks (~~~)
		if match := tildeCodeBlockRegex.FindStringSubmatch(line); match != nil {
			example, endLine, err := m.extractFencedBlock(lines, i, "~~~", match[1], source)
			if err != nil {
				continue // Skip malformed blocks
			}
			if example != nil {
				examples = append(examples, *example)
			}
			i = endLine
			continue
		}

		// Check for indented code blocks (4 spaces)
		if indentedCodeRegex.MatchString(line) {
			example, endLine := m.extractIndentedBlock(lines, i, source)
			if example != nil {
				examples = append(examples, *example)
			}
			i = endLine
			continue
		}
	}

	return examples, nil
}

// extractFencedBlock extracts a fenced code block (``` or ~~~)
func (m *MarkdownExtractor) extractFencedBlock(lines []string, startLine int, fence, language, source string) (*Example, int, error) {
	var codeLines []string
	endLine := startLine

	// Find the closing fence
	for i := startLine + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], fence) {
			endLine = i
			break
		}
		codeLines = append(codeLines, lines[i])
	}

	if endLine == startLine {
		return nil, startLine, fmt.Errorf("unclosed code block at line %d", startLine+1)
	}

	// Skip if language is not in our supported list
	if language != "" && !m.isLanguageSupported(language) {
		return nil, endLine, nil
	}

	code := strings.Join(codeLines, "\n")

	// Skip empty code blocks
	if strings.TrimSpace(code) == "" {
		return nil, endLine, nil
	}

	// Extract context (previous and next few lines)
	context := m.extractContext(lines, startLine, endLine)

	// Detect tags and metadata
	tags := m.detectTags(code, context, language)
	metadata := m.extractMetadata(code, context)

	example := &Example{
		ID:            generateExampleID(source, code, startLine+1),
		Source:        source,
		Language:      m.normalizeLanguage(language),
		Code:          code,
		Context:       context,
		Tags:          tags,
		Prerequisites: m.detectPrerequisites(code, context),
		Metadata:      metadata,
		LineStart:     startLine + 1,
		LineEnd:       endLine + 1,
	}

	return example, endLine, nil
}

// extractIndentedBlock extracts an indented code block (4 spaces)
func (m *MarkdownExtractor) extractIndentedBlock(lines []string, startLine int, source string) (*Example, int) {
	var codeLines []string
	endLine := startLine

	// Collect all consecutive indented lines
	for i := startLine; i < len(lines); i++ {
		if match := indentedCodeRegex.FindStringSubmatch(lines[i]); match != nil {
			codeLines = append(codeLines, match[1])
			endLine = i
		} else if strings.TrimSpace(lines[i]) == "" {
			// Allow empty lines within indented blocks
			codeLines = append(codeLines, "")
			endLine = i
		} else {
			break
		}
	}

	code := strings.Join(codeLines, "\n")
	code = strings.TrimRight(code, "\n")

	// Skip empty or very short code blocks
	if len(strings.TrimSpace(code)) < 3 {
		return nil, endLine
	}

	// Try to detect language from content
	language := m.detectLanguageFromContent(code)

	// Extract context
	context := m.extractContext(lines, startLine, endLine)

	example := &Example{
		ID:            generateExampleID(source, code, startLine+1),
		Source:        source,
		Language:      language,
		Code:          code,
		Context:       context,
		Tags:          m.detectTags(code, context, language),
		Prerequisites: m.detectPrerequisites(code, context),
		Metadata:      m.extractMetadata(code, context),
		LineStart:     startLine + 1,
		LineEnd:       endLine + 1,
	}

	return example, endLine
}

// extractContext extracts surrounding context for better understanding
func (m *MarkdownExtractor) extractContext(lines []string, startLine, endLine int) string {
	contextLines := 3 // Lines before and after to include

	start := max(0, startLine-contextLines)
	end := min(len(lines), endLine+contextLines+1)

	var contextParts []string

	// Add lines before the code block
	if start < startLine {
		before := strings.Join(lines[start:startLine], "\n")
		if strings.TrimSpace(before) != "" {
			contextParts = append(contextParts, "Before:\n"+before)
		}
	}

	// Add lines after the code block
	if end > endLine+1 {
		after := strings.Join(lines[endLine+1:end], "\n")
		if strings.TrimSpace(after) != "" {
			contextParts = append(contextParts, "After:\n"+after)
		}
	}

	return strings.Join(contextParts, "\n\n")
}

// detectTags identifies the purpose and category of the code example
func (m *MarkdownExtractor) detectTags(code, context, language string) []string {
	var tags []string

	codeUpper := strings.ToUpper(code)
	contextUpper := strings.ToUpper(context)

	// Language-specific tags
	if language != "" {
		tags = append(tags, language)
	}

	// Purpose tags
	if strings.Contains(codeUpper, "INSTALL") || strings.Contains(contextUpper, "INSTALL") {
		tags = append(tags, "install")
	}
	if strings.Contains(codeUpper, "SETUP") || strings.Contains(contextUpper, "SETUP") {
		tags = append(tags, "setup")
	}
	if strings.Contains(codeUpper, "EXAMPLE") || strings.Contains(contextUpper, "EXAMPLE") {
		tags = append(tags, "example")
	}
	if strings.Contains(codeUpper, "DEMO") || strings.Contains(contextUpper, "DEMO") {
		tags = append(tags, "demo")
	}
	if strings.Contains(codeUpper, "CLEANUP") || strings.Contains(contextUpper, "CLEANUP") {
		tags = append(tags, "cleanup")
	}
	if strings.Contains(codeUpper, "TEST") || strings.Contains(contextUpper, "TEST") {
		tags = append(tags, "test")
	}

	// Command type tags for shell scripts
	if language == "bash" || language == "shell" || language == "sh" {
		if strings.Contains(code, "curl") || strings.Contains(code, "wget") {
			tags = append(tags, "download")
		}
		if strings.Contains(code, "docker") {
			tags = append(tags, "docker")
		}
		if strings.Contains(code, "aws ") {
			tags = append(tags, "aws")
		}
		if strings.Contains(code, "git ") {
			tags = append(tags, "git")
		}
		if strings.Contains(code, "npm ") || strings.Contains(code, "yarn ") {
			tags = append(tags, "npm")
		}
		if strings.Contains(code, "go ") || strings.Contains(code, "go.mod") {
			tags = append(tags, "golang")
		}
	}

	// Complexity tags
	lineCount := len(strings.Split(code, "\n"))
	if lineCount <= 3 {
		tags = append(tags, "simple")
	} else if lineCount > 10 {
		tags = append(tags, "complex")
	}

	return tags
}

// detectPrerequisites identifies dependencies on other examples
func (m *MarkdownExtractor) detectPrerequisites(code, context string) []string {
	var prereqs []string

	// Look for references to previous steps
	contextLower := strings.ToLower(context)
	if strings.Contains(contextLower, "after") || strings.Contains(contextLower, "following") {
		prereqs = append(prereqs, "previous-step")
	}

	// Look for specific dependencies in comments
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			if strings.Contains(strings.ToLower(line), "requires") {
				prereqs = append(prereqs, "explicit-requirement")
			}
		}
	}

	return prereqs
}

// extractMetadata extracts additional metadata from code and context
func (m *MarkdownExtractor) extractMetadata(code, context string) map[string]string {
	metadata := make(map[string]string)

	// Estimate execution time based on content
	if strings.Contains(code, "sleep") {
		metadata["estimated_duration"] = "long"
	} else if len(strings.Split(code, "\n")) > 10 {
		metadata["estimated_duration"] = "medium"
	} else {
		metadata["estimated_duration"] = "short"
	}

	// Check if requires user input
	if strings.Contains(code, "read ") || strings.Contains(code, "input(") {
		metadata["requires_input"] = "true"
	}

	// Check if modifies system
	if strings.Contains(code, "sudo ") || strings.Contains(code, "rm -rf") {
		metadata["system_modifying"] = "true"
	}

	// Check if creates files
	if strings.Contains(code, ">") || strings.Contains(code, "touch ") || strings.Contains(code, "mkdir") {
		metadata["creates_files"] = "true"
	}

	return metadata
}

// isLanguageSupported checks if a language is in our supported list
func (m *MarkdownExtractor) isLanguageSupported(language string) bool {
	normalized := m.normalizeLanguage(language)
	for _, supported := range m.config.Languages {
		if supported == normalized {
			return true
		}
	}
	return false
}

// normalizeLanguage normalizes language names
func (m *MarkdownExtractor) normalizeLanguage(language string) string {
	switch strings.ToLower(language) {
	case "bash", "sh", "shell", "zsh":
		return "bash"
	case "golang", "go":
		return "go"
	case "python", "py", "python3":
		return "python"
	case "yaml", "yml":
		return "yaml"
	case "json":
		return "json"
	case "javascript", "js":
		return "javascript"
	case "typescript", "ts":
		return "typescript"
	case "dockerfile", "docker":
		return "dockerfile"
	default:
		return strings.ToLower(language)
	}
}

// detectLanguageFromContent tries to detect language from code content
func (m *MarkdownExtractor) detectLanguageFromContent(code string) string {
	codeUpper := strings.ToUpper(code)

	// Shell/bash indicators
	if strings.HasPrefix(code, "#!") && (strings.Contains(code, "bash") || strings.Contains(code, "sh")) {
		return "bash"
	}
	if strings.Contains(codeUpper, "#!/BIN/") {
		return "bash"
	}

	// Common shell commands
	shellCommands := []string{"LS ", "CD ", "PWD", "ECHO ", "CAT ", "GREP ", "AWK ", "SED ", "FIND "}
	for _, cmd := range shellCommands {
		if strings.Contains(codeUpper, cmd) {
			return "bash"
		}
	}

	// Go indicators
	if strings.Contains(code, "package main") || strings.Contains(code, "func main()") {
		return "go"
	}

	// Python indicators
	if strings.Contains(code, "def ") || strings.Contains(code, "import ") || strings.Contains(code, "print(") {
		return "python"
	}

	// YAML indicators
	if strings.Contains(code, ":") && (strings.Contains(code, "  ") || strings.Contains(code, "\t")) {
		// Simple heuristic for YAML
		lines := strings.Split(code, "\n")
		yamlLike := 0
		for _, line := range lines {
			if strings.Contains(line, ":") && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				yamlLike++
			}
		}
		if yamlLike > len(lines)/2 {
			return "yaml"
		}
	}

	return "text" // Default to text if can't detect
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
