/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeProvider implements the Provider interface using Anthropic's Claude
type ClaudeProvider struct {
	client      *anthropic.Client
	model       string
	maxTokens   int
	temperature float64
	cache       *ResponseCache
	metrics     *ProviderMetrics
	usageStats  *UsageStats
}

// NewClaudeProvider creates a new Claude AI provider
func NewClaudeProvider(apiKey string, options ...ClaudeOption) *ClaudeProvider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	provider := &ClaudeProvider{
		client:      &client,
		model:       "claude-3-5-sonnet-20241022",
		maxTokens:   4000,
		temperature: 0.0, // Deterministic for better caching
		cache:       NewResponseCache(),
		metrics:     NewProviderMetrics(),
		usageStats:  &UsageStats{},
	}

	// Apply options
	for _, opt := range options {
		opt(provider)
	}

	return provider
}

// ClaudeOption allows customization of Claude provider
type ClaudeOption func(*ClaudeProvider)

// WithModel sets the Claude model to use
func WithModel(model string) ClaudeOption {
	return func(p *ClaudeProvider) {
		p.model = model
	}
}

// WithMaxTokens sets the maximum tokens for responses
func WithMaxTokens(maxTokens int) ClaudeOption {
	return func(p *ClaudeProvider) {
		p.maxTokens = maxTokens
	}
}

// WithTemperature sets the temperature for responses
func WithTemperature(temperature float64) ClaudeOption {
	return func(p *ClaudeProvider) {
		p.temperature = temperature
	}
}

// ParseInstruction uses Claude to parse natural language instructions
func (c *ClaudeProvider) ParseInstruction(ctx context.Context, instruction string, tutorialContext TutorialContext) (*ParsedInstruction, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("parse:%s:%s", hashString(instruction), hashContext(tutorialContext))
	if cached := c.cache.Get(cacheKey); cached != nil {
		if parsed, ok := cached.(*ParsedInstruction); ok {
			return parsed, nil
		}
	}

	prompt := c.buildInstructionPrompt(instruction, tutorialContext)

	params := anthropic.MessageNewParams{}
	params.Model = anthropic.Model(c.model)
	params.MaxTokens = int64(c.maxTokens)
	params.Temperature = anthropic.Float(c.temperature)
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}
	response, err := c.client.Messages.New(ctx, params)

	if err != nil {
		c.metrics.RecordError()
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	c.metrics.RecordRequest(int(response.Usage.InputTokens), int(response.Usage.OutputTokens))
	c.updateUsageStats(int(response.Usage.InputTokens), int(response.Usage.OutputTokens))

	// Parse Claude's response into structured format
	responseText := ""
	if len(response.Content) > 0 {
		textBlock := response.Content[0].AsText()
		responseText = textBlock.Text
	}
	parsed, err := c.parseInstructionResponse(responseText, instruction)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, parsed, 24*time.Hour)

	return parsed, nil
}

// ValidateExpectation uses Claude to validate if actual output matches expected outcome
func (c *ClaudeProvider) ValidateExpectation(ctx context.Context, expected, actual string, tutorialContext TutorialContext) (*ValidationResult, error) {
	cacheKey := fmt.Sprintf("validate:%s:%s", hashString(expected+actual), hashContext(tutorialContext))
	if cached := c.cache.Get(cacheKey); cached != nil {
		if result, ok := cached.(*ValidationResult); ok {
			return result, nil
		}
	}

	prompt := c.buildValidationPrompt(expected, actual, tutorialContext)

	params := anthropic.MessageNewParams{}
	params.Model = anthropic.Model(c.model)
	params.MaxTokens = int64(2000)
	params.Temperature = anthropic.Float(c.temperature)
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}
	response, err := c.client.Messages.New(ctx, params)

	if err != nil {
		c.metrics.RecordError()
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	c.metrics.RecordRequest(int(response.Usage.InputTokens), int(response.Usage.OutputTokens))
	c.updateUsageStats(int(response.Usage.InputTokens), int(response.Usage.OutputTokens))

	responseText := ""
	if len(response.Content) > 0 {
		textBlock := response.Content[0].AsText()
		responseText = textBlock.Text
	}
	result, err := c.parseValidationResponse(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validation response: %w", err)
	}

	c.cache.Set(cacheKey, result, 1*time.Hour)
	return result, nil
}

// CompressContext uses Claude to compress tutorial context for efficiency
func (c *ClaudeProvider) CompressContext(ctx context.Context, fullContext TutorialContext) (*CompressedContext, error) {
	prompt := c.buildCompressionPrompt(fullContext)

	// Use faster model for compression
	params := anthropic.MessageNewParams{}
	params.Model = anthropic.Model("claude-3-5-haiku-20241022")
	params.MaxTokens = int64(500)
	params.Temperature = anthropic.Float(0.0)
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}
	response, err := c.client.Messages.New(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	c.updateUsageStats(int(response.Usage.InputTokens), int(response.Usage.OutputTokens))

	responseText := ""
	if len(response.Content) > 0 {
		textBlock := response.Content[0].AsText()
		responseText = textBlock.Text
	}
	compressed, err := c.parseCompressionResponse(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse compression response: %w", err)
	}

	return compressed, nil
}

// InterpretError uses Claude to interpret and suggest solutions for errors
func (c *ClaudeProvider) InterpretError(ctx context.Context, errorMsg string, tutorialContext TutorialContext) (*ErrorInterpretation, error) {
	prompt := c.buildErrorPrompt(errorMsg, tutorialContext)

	params := anthropic.MessageNewParams{}
	params.Model = anthropic.Model(c.model)
	params.MaxTokens = int64(3000)
	params.Temperature = anthropic.Float(0.1) // Slightly creative for solution suggestions
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}
	response, err := c.client.Messages.New(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	c.updateUsageStats(int(response.Usage.InputTokens), int(response.Usage.OutputTokens))

	responseText := ""
	if len(response.Content) > 0 {
		textBlock := response.Content[0].AsText()
		responseText = textBlock.Text
	}
	interpretation, err := c.parseErrorResponse(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse error response: %w", err)
	}

	return interpretation, nil
}

// buildInstructionPrompt creates a prompt for instruction parsing
func (c *ClaudeProvider) buildInstructionPrompt(instruction string, context TutorialContext) string {
	return fmt.Sprintf(`You are an expert at understanding technical tutorial instructions. Your job is to parse natural language instructions into structured, executable actions.

Context:
- Working Directory: %s
- Current Step: %d of %d
- Files Created: %v
- Recent Commands: %v

Instruction to Parse:
"%s"

Please analyze this instruction and respond with a JSON object containing:
{
  "intent": "Brief description of what this instruction wants to accomplish",
  "actions": [
    {
      "type": "command|validate|download|navigate|check",
      "command": "specific shell command to execute (if applicable)",
      "description": "human-readable description",
      "validation": {
        "type": "exit_code|file_exists|output|contains",
        "condition": "what to check",
        "expected": "expected result"
      }
    }
  ],
  "prerequisites": ["list of things that must be done first"],
  "expected_outcomes": ["what should happen after executing these actions"],
  "confidence": 0.95,
  "reasoning": "explanation of your interpretation"
}

Focus on:
1. Understanding the intent behind the instruction
2. Breaking down complex instructions into specific actions
3. Identifying what needs to be validated
4. Being specific about commands and file paths
5. Considering the current context and state

Respond only with the JSON object, no additional text.`,
		context.WorkingDirectory,
		context.CurrentStep,
		context.TotalSteps,
		context.CreatedFiles,
		context.ExecutedCommands,
		instruction)
}

// buildValidationPrompt creates a prompt for expectation validation
func (c *ClaudeProvider) buildValidationPrompt(expected, actual string, context TutorialContext) string {
	return fmt.Sprintf(`You are validating whether actual output matches expected outcomes in a tutorial.

Expected Outcome:
"%s"

Actual Output:
"%s"

Context:
- Working Directory: %s
- Step: %d of %d

Please analyze whether the actual output satisfies the expected outcome and respond with JSON:
{
  "success": true/false,
  "confidence": 0.95,
  "reasoning": "detailed explanation of why it succeeded or failed",
  "differences": ["list of specific differences if any"],
  "suggestions": ["suggestions for improvement if it failed"]
}

Consider:
1. The expected outcome might be described in natural language, not exact text
2. Minor formatting differences are usually okay
3. Focus on whether the functional outcome was achieved
4. Be helpful in explaining differences

Respond only with the JSON object.`,
		expected, actual,
		context.WorkingDirectory,
		context.CurrentStep,
		context.TotalSteps)
}

// buildCompressionPrompt creates a prompt for context compression
func (c *ClaudeProvider) buildCompressionPrompt(context TutorialContext) string {
	return fmt.Sprintf(`Compress this tutorial context into a brief summary while preserving essential information:

Current Context:
- Working Directory: %s
- Step: %d of %d
- Files Created: %v
- Commands Executed: %v
- Environment Variables: %v

Provide a JSON response:
{
  "summary": "2-3 sentence summary of what has been accomplished",
  "key_files": ["list of important files that were created"],
  "current_state": "brief description of current working state"
}

Focus on preserving information that would be needed for subsequent tutorial steps.
Respond only with the JSON object.`,
		context.WorkingDirectory,
		context.CurrentStep,
		context.TotalSteps,
		context.CreatedFiles,
		context.ExecutedCommands,
		context.EnvironmentVars)
}

// buildErrorPrompt creates a prompt for error interpretation
func (c *ClaudeProvider) buildErrorPrompt(errorMsg string, context TutorialContext) string {
	return fmt.Sprintf(`You are an expert at diagnosing technical errors. Analyze this error and provide solutions.

Error Message:
"%s"

Context:
- Working Directory: %s
- Recent Commands: %v
- Environment: %v

Provide a JSON response:
{
  "error_type": "category of error (e.g., 'permission_denied', 'command_not_found', 'network_error')",
  "explanation": "clear explanation of what went wrong",
  "solutions": [
    {
      "description": "what to do",
      "commands": ["list of commands to run"],
      "probability": 0.8
    }
  ],
  "confidence": 0.9
}

Focus on:
1. Identifying the root cause
2. Providing actionable solutions
3. Ordering solutions by likelihood of success
4. Being specific about commands to run

Respond only with the JSON object.`,
		errorMsg,
		context.WorkingDirectory,
		context.ExecutedCommands,
		context.EnvironmentVars)
}

// parseInstructionResponse parses Claude's instruction parsing response
func (c *ClaudeProvider) parseInstructionResponse(response, originalText string) (*ParsedInstruction, error) {
	// Extract JSON from response (Claude sometimes adds extra text)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]

	var parsed ParsedInstruction
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	parsed.OriginalText = originalText
	return &parsed, nil
}

// parseValidationResponse parses Claude's validation response
func (c *ClaudeProvider) parseValidationResponse(response string) (*ValidationResult, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]

	var result ValidationResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &result, nil
}

// parseCompressionResponse parses Claude's context compression response
func (c *ClaudeProvider) parseCompressionResponse(response string) (*CompressedContext, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]

	var compressed CompressedContext
	if err := json.Unmarshal([]byte(jsonStr), &compressed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &compressed, nil
}

// parseErrorResponse parses Claude's error interpretation response
func (c *ClaudeProvider) parseErrorResponse(response string) (*ErrorInterpretation, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart:jsonEnd]

	var interpretation ErrorInterpretation
	if err := json.Unmarshal([]byte(jsonStr), &interpretation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &interpretation, nil
}

// GetCapabilities returns Claude's capabilities
func (c *ClaudeProvider) GetCapabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Name:               "claude",
		Version:            c.model,
		MaxContextLength:   200000,
		SupportedLanguages: []string{"en", "es", "fr", "de", "ja", "zh"},
		Features:           []string{"instruction_parsing", "expectation_validation", "context_compression", "error_analysis"},
		CostModel:          CostModel{InputTokenPrice: 0.000003, OutputTokenPrice: 0.000015, Currency: "USD"},
		LatencyProfile:     LatencyProfile{P50: "2s", P95: "8s", P99: "15s"},
		QualityScore:       0.94,
		CertificationLevel: CertifiedGold,
	}
}

// GetCostEstimate estimates the cost of a request
func (c *ClaudeProvider) GetCostEstimate(request AIRequest) CostEstimate {
	// Rough estimation based on content length
	estimatedTokens := len(request.Content)/4 + request.MaxTokens
	capabilities := c.GetCapabilities()
	cost := float64(estimatedTokens) * capabilities.CostModel.InputTokenPrice

	return CostEstimate{
		EstimatedTokens: estimatedTokens,
		EstimatedCost:   cost,
		Currency:        capabilities.CostModel.Currency,
	}
}

// GetPerformanceMetrics returns current performance metrics
func (c *ClaudeProvider) GetPerformanceMetrics() PerformanceMetrics {
	return c.metrics.GetMetrics()
}

// HealthCheck verifies the provider is working
func (c *ClaudeProvider) HealthCheck(ctx context.Context) error {
	// Simple health check with minimal tokens
	params := anthropic.MessageNewParams{}
	params.Model = anthropic.Model("claude-3-5-haiku-20241022")
	params.MaxTokens = int64(10)
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("Respond with 'OK'")),
	}
	_, err := c.client.Messages.New(ctx, params)

	return err
}

// GetUsageStats returns current usage statistics
func (c *ClaudeProvider) GetUsageStats() UsageStats {
	return *c.usageStats
}

// updateUsageStats updates internal usage tracking
func (c *ClaudeProvider) updateUsageStats(inputTokens, outputTokens int) {
	capabilities := c.GetCapabilities()
	cost := float64(inputTokens)*capabilities.CostModel.InputTokenPrice +
		float64(outputTokens)*capabilities.CostModel.OutputTokenPrice

	c.usageStats.RequestsToday++
	c.usageStats.RequestsTotal++
	c.usageStats.TokensUsedToday += inputTokens + outputTokens
	c.usageStats.TokensUsedTotal += inputTokens + outputTokens
	c.usageStats.CostToday += cost
	c.usageStats.CostTotal += cost
}

// Helper functions
func hashString(s string) string {
	// Simple hash function - could be improved with crypto/sha256
	hash := 0
	for _, char := range s {
		hash = hash*31 + int(char)
	}
	return fmt.Sprintf("%x", hash)
}

func hashContext(ctx TutorialContext) string {
	return hashString(fmt.Sprintf("%s:%d:%v", ctx.WorkingDirectory, ctx.CurrentStep, ctx.CreatedFiles))
}
