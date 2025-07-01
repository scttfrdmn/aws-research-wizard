package extractor

import (
	"regexp"
	"strings"
)

var (
	// Regular expressions for HTML code blocks
	preCodeRegex = regexp.MustCompile(`(?s)<pre[^>]*(?:class="[^"]*language-([^"]*)[^"]*"[^>]*)?[^>]*>(.*?)</pre>`)
	codeRegex    = regexp.MustCompile(`(?s)<code[^>]*(?:class="[^"]*language-([^"]*)[^"]*"[^>]*)?[^>]*>(.*?)</code>`)
)

// ExtractFromContent extracts code examples from HTML content
func (h *HTMLExtractor) ExtractFromContent(content []byte, source string) ([]Example, error) {
	var examples []Example
	contentStr := string(content)

	// Extract from <pre> tags (more likely to be code blocks)
	preMatches := preCodeRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range preMatches {
		language := match[1] // May be empty
		code := h.cleanHTMLCode(match[2])

		if h.isValidCode(code, language) {
			example := h.createExample(code, language, source, contentStr)
			examples = append(examples, example)
		}
	}

	// Extract from <code> tags (but filter more strictly)
	codeMatches := codeRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range codeMatches {
		language := match[1]
		code := h.cleanHTMLCode(match[2])

		// For <code> tags, be more selective - require minimum length and language
		if len(strings.TrimSpace(code)) > 20 && h.isValidCode(code, language) {
			example := h.createExample(code, language, source, contentStr)
			examples = append(examples, example)
		}
	}

	return h.deduplicateExamples(examples), nil
}

// cleanHTMLCode removes HTML entities and cleans up code content
func (h *HTMLExtractor) cleanHTMLCode(rawCode string) string {
	// Decode common HTML entities
	code := strings.ReplaceAll(rawCode, "&lt;", "<")
	code = strings.ReplaceAll(code, "&gt;", ">")
	code = strings.ReplaceAll(code, "&amp;", "&")
	code = strings.ReplaceAll(code, "&quot;", "\"")
	code = strings.ReplaceAll(code, "&#x27;", "'")
	code = strings.ReplaceAll(code, "&#39;", "'")

	// Remove HTML tags that might be nested inside
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)
	code = htmlTagRegex.ReplaceAllString(code, "")

	// Clean up whitespace
	code = strings.TrimSpace(code)

	return code
}

// isValidCode checks if the extracted content is likely to be actual code
func (h *HTMLExtractor) isValidCode(code, language string) bool {
	cleanCode := strings.TrimSpace(code)

	// Skip empty or very short content
	if len(cleanCode) < 3 {
		return false
	}

	// Skip if it's just a single word (likely inline code)
	if !strings.Contains(cleanCode, " ") && !strings.Contains(cleanCode, "\n") && len(cleanCode) < 20 {
		return false
	}

	// Check if language is supported
	if language != "" && !h.isLanguageSupported(language) {
		return false
	}

	// Skip if it looks like regular text (has many common English words)
	englishWords := []string{"the", "and", "for", "are", "but", "not", "you", "all", "can", "had", "her", "was", "one", "our", "out", "day", "get", "has", "him", "his", "how", "man", "new", "now", "old", "see", "two", "way", "who", "boy", "did", "its", "let", "put", "say", "she", "too", "use"}
	wordCount := 0
	englishWordCount := 0
	words := strings.Fields(strings.ToLower(cleanCode))

	for _, word := range words {
		wordCount++
		for _, englishWord := range englishWords {
			if word == englishWord {
				englishWordCount++
				break
			}
		}
	}

	// If more than 50% are common English words, probably not code
	if wordCount > 5 && float64(englishWordCount)/float64(wordCount) > 0.5 {
		return false
	}

	return true
}

// createExample creates an Example from extracted HTML code
func (h *HTMLExtractor) createExample(code, language, source, fullContent string) Example {
	// Try to extract context by finding surrounding text
	context := h.extractHTMLContext(code, fullContent)

	// Normalize language
	normalizedLanguage := h.normalizeLanguage(language)

	// Detect tags (reuse logic from markdown extractor)
	tags := h.detectTags(code, context, normalizedLanguage)

	// Generate line numbers (approximate)
	lineStart := strings.Count(fullContent[:strings.Index(fullContent, code)], "\n") + 1
	lineEnd := lineStart + strings.Count(code, "\n")

	return Example{
		ID:            generateExampleID(source, code, lineStart),
		Source:        source,
		Language:      normalizedLanguage,
		Code:          code,
		Context:       context,
		Tags:          tags,
		Prerequisites: h.detectPrerequisites(code, context),
		Metadata:      h.extractMetadata(code, context),
		LineStart:     lineStart,
		LineEnd:       lineEnd,
	}
}

// extractHTMLContext tries to extract meaningful context around code blocks
func (h *HTMLExtractor) extractHTMLContext(code, fullContent string) string {
	// Find the position of the code in the full content
	codeIndex := strings.Index(fullContent, code)
	if codeIndex == -1 {
		return ""
	}

	// Look for surrounding text within reasonable distance
	contextStart := maxInt(0, codeIndex-500)
	contextEnd := minInt(len(fullContent), codeIndex+len(code)+500)

	contextArea := fullContent[contextStart:contextEnd]

	// Extract text from HTML tags in the context area
	textRegex := regexp.MustCompile(`>([^<]+)<`)
	matches := textRegex.FindAllStringSubmatch(contextArea, -1)

	var contextParts []string
	for _, match := range matches {
		text := strings.TrimSpace(match[1])
		if len(text) > 10 && !strings.Contains(text, code) {
			contextParts = append(contextParts, text)
		}
	}

	if len(contextParts) > 3 {
		contextParts = contextParts[:3] // Limit context length
	}

	return strings.Join(contextParts, " ")
}

// deduplicateExamples removes duplicate examples based on code content
func (h *HTMLExtractor) deduplicateExamples(examples []Example) []Example {
	seen := make(map[string]bool)
	var unique []Example

	for _, example := range examples {
		codeHash := generateExampleID("", example.Code, 0)
		if !seen[codeHash] {
			seen[codeHash] = true
			unique = append(unique, example)
		}
	}

	return unique
}

// Helper methods that delegate to the shared markdown extractor logic
func (h *HTMLExtractor) isLanguageSupported(language string) bool {
	normalized := h.normalizeLanguage(language)
	for _, supported := range h.config.Languages {
		if supported == normalized {
			return true
		}
	}
	return false
}

func (h *HTMLExtractor) normalizeLanguage(language string) string {
	// Reuse the same logic as markdown extractor
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

func (h *HTMLExtractor) detectTags(code, context, language string) []string {
	// Simple tag detection - could be enhanced
	var tags []string

	if language != "" {
		tags = append(tags, language)
	}

	codeUpper := strings.ToUpper(code)
	if strings.Contains(codeUpper, "INSTALL") {
		tags = append(tags, "install")
	}
	if strings.Contains(codeUpper, "EXAMPLE") {
		tags = append(tags, "example")
	}

	return tags
}

func (h *HTMLExtractor) detectPrerequisites(code, context string) []string {
	// Simple prerequisite detection
	return []string{}
}

func (h *HTMLExtractor) extractMetadata(code, context string) map[string]string {
	metadata := make(map[string]string)
	metadata["source_format"] = "html"
	return metadata
}

// Helper functions
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
