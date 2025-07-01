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
	"fmt"
	"os"
)

// Client provides a high-level interface for AI operations
type Client struct {
	provider Provider
	config   ClientConfig
}

// ClientConfig holds configuration for the AI client
type ClientConfig struct {
	DefaultTimeout   string `json:"default_timeout"`
	MaxRetries       int    `json:"max_retries"`
	CacheEnabled     bool   `json:"cache_enabled"`
	CostOptimization bool   `json:"cost_optimization"`
}

// NewClient creates a new AI client with the specified provider
func NewClient(provider Provider, config ClientConfig) *Client {
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	return &Client{
		provider: provider,
		config:   config,
	}
}

// NewClaudeClient creates a new AI client with Claude as the provider
func NewClaudeClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		// Try to get from environment
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("no Anthropic API key provided")
		}
	}

	claudeProvider := NewClaudeProvider(apiKey)

	config := ClientConfig{
		DefaultTimeout:   "30s",
		MaxRetries:       3,
		CacheEnabled:     true,
		CostOptimization: true,
	}

	return NewClient(claudeProvider, config), nil
}

// ParseInstruction parses a natural language instruction into structured actions
func (c *Client) ParseInstruction(ctx context.Context, instruction string, tutorialContext TutorialContext) (*ParsedInstruction, error) {
	return c.provider.ParseInstruction(ctx, instruction, tutorialContext)
}

// ValidateExpectation validates if actual output matches expected outcome
func (c *Client) ValidateExpectation(ctx context.Context, expected, actual string, tutorialContext TutorialContext) (*ValidationResult, error) {
	return c.provider.ValidateExpectation(ctx, expected, actual, tutorialContext)
}

// CompressContext compresses tutorial context for efficiency
func (c *Client) CompressContext(ctx context.Context, fullContext TutorialContext) (*CompressedContext, error) {
	return c.provider.CompressContext(ctx, fullContext)
}

// InterpretError interprets an error and suggests solutions
func (c *Client) InterpretError(ctx context.Context, errorMsg string, tutorialContext TutorialContext) (*ErrorInterpretation, error) {
	return c.provider.InterpretError(ctx, errorMsg, tutorialContext)
}

// GetCapabilities returns the provider's capabilities
func (c *Client) GetCapabilities() ProviderCapabilities {
	return c.provider.GetCapabilities()
}

// HealthCheck verifies the AI provider is working
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.provider.HealthCheck(ctx)
}

// GetUsageStats returns current usage statistics
func (c *Client) GetUsageStats() UsageStats {
	return c.provider.GetUsageStats()
}

// GetPerformanceMetrics returns performance metrics
func (c *Client) GetPerformanceMetrics() PerformanceMetrics {
	return c.provider.GetPerformanceMetrics()
}
