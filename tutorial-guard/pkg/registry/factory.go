/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package registry

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
)

// ProviderFactory creates and configures AI providers
type ProviderFactory struct {
	defaultConfigs map[string]ProviderConfig
}

// ProviderType defines supported provider types
type ProviderType string

const (
	ProviderTypeClaude ProviderType = "claude"
	ProviderTypeGPT4   ProviderType = "gpt4"
	ProviderTypeGemini ProviderType = "gemini"
	ProviderTypeLlama  ProviderType = "llama"
	ProviderTypeLocal  ProviderType = "local"
)

// CreateProviderOptions defines options for creating providers
type CreateProviderOptions struct {
	Type     ProviderType      `json:"type"`
	Name     string            `json:"name"`
	APIKey   string            `json:"api_key,omitempty"`
	Endpoint string            `json:"endpoint,omitempty"`
	Model    string            `json:"model,omitempty"`
	Region   string            `json:"region,omitempty"`
	Config   ProviderConfig    `json:"config"`
	Override map[string]string `json:"override,omitempty"`
}

// NewProviderFactory creates a new provider factory with default configurations
func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		defaultConfigs: make(map[string]ProviderConfig),
	}

	// Set up default configurations for each provider type
	factory.setupDefaultConfigs()

	return factory
}

// CreateProvider creates and configures a new AI provider
func (f *ProviderFactory) CreateProvider(options CreateProviderOptions) (ai.Provider, error) {
	// Merge with default config
	config := f.mergeWithDefaults(options)

	// Validate configuration
	if err := f.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create provider based on type
	switch options.Type {
	case ProviderTypeClaude:
		return f.createClaudeProvider(options, config)

	case ProviderTypeGPT4:
		return f.createGPT4Provider(options, config)

	case ProviderTypeGemini:
		return f.createGeminiProvider(options, config)

	case ProviderTypeLlama:
		return f.createLlamaProvider(options, config)

	case ProviderTypeLocal:
		return f.createLocalProvider(options, config)

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", options.Type)
	}
}

// CreateDefaultRegistry creates a registry with standard provider configurations
func (f *ProviderFactory) CreateDefaultRegistry() (*ProviderRegistry, error) {
	config := RoutingConfig{
		DefaultStrategy:  StrategyIntelligent,
		FallbackChain:    []string{"claude", "gpt4", "gemini"},
		LoadBalancing:    LoadBalancingWeighted,
		QualityThreshold: 0.7,
		CostOptimization: true,
		LatencyThreshold: 30 * time.Second,
		MaxRetries:       3,
		CircuitBreaker: CircuitBreakerConfig{
			Enabled:          true,
			FailureThreshold: 5,
			RecoveryTimeout:  60 * time.Second,
			SuccessThreshold: 3,
		},
	}

	registry := NewProviderRegistry(config)

	// Register Claude provider if API key available
	if claudeKey := os.Getenv("ANTHROPIC_API_KEY"); claudeKey != "" {
		claudeProvider, err := f.CreateProvider(CreateProviderOptions{
			Type:   ProviderTypeClaude,
			Name:   "claude",
			APIKey: claudeKey,
			Model:  "claude-3-5-sonnet-20241022",
		})
		if err == nil {
			registry.Register("claude", claudeProvider, f.defaultConfigs["claude"])
		}
	}

	// Register OpenAI provider if API key available
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		gpt4Provider, err := f.CreateProvider(CreateProviderOptions{
			Type:   ProviderTypeGPT4,
			Name:   "gpt4",
			APIKey: openaiKey,
			Model:  "gpt-4-turbo-preview",
		})
		if err == nil {
			registry.Register("gpt4", gpt4Provider, f.defaultConfigs["gpt4"])
		}
	}

	// Register Google provider if API key available
	if googleKey := os.Getenv("GOOGLE_AI_API_KEY"); googleKey != "" {
		geminiProvider, err := f.CreateProvider(CreateProviderOptions{
			Type:   ProviderTypeGemini,
			Name:   "gemini",
			APIKey: googleKey,
			Model:  "gemini-pro",
		})
		if err == nil {
			registry.Register("gemini", geminiProvider, f.defaultConfigs["gemini"])
		}
	}

	return registry, nil
}

// setupDefaultConfigs initializes default configurations for each provider type
func (f *ProviderFactory) setupDefaultConfigs() {
	// Claude (Anthropic) default configuration
	f.defaultConfigs["claude"] = ProviderConfig{
		Name:          "claude",
		Priority:      100, // Highest priority (proven quality)
		Weight:        0.4,
		MaxConcurrent: 10,
		Timeout:       30 * time.Second,
		RetryPolicy: RetryPolicy{
			MaxRetries:      3,
			InitialDelay:    1 * time.Second,
			MaxDelay:        30 * time.Second,
			BackoffFactor:   2.0,
			RetryableErrors: []string{"rate_limit", "timeout", "server_error"},
		},
		CostLimit: CostLimit{
			DailyCostLimit:   50.0,   // $50/day
			MonthlyCostLimit: 1000.0, // $1000/month
			PerRequestLimit:  0.50,   // $0.50 per request
			AlertThreshold:   0.8,    // Alert at 80%
		},
		Capabilities: []string{
			"instruction_parsing",
			"expectation_validation",
			"context_compression",
			"error_analysis",
			"natural_language",
			"code_understanding",
		},
		Regions: []string{"global"},
		Metadata: map[string]string{
			"model_family":    "claude-3.5",
			"context_length":  "200000",
			"output_tokens":   "4096",
			"multimodal":      "false",
			"fine_tunable":    "false",
			"certified_level": "gold",
		},
	}

	// GPT-4 (OpenAI) default configuration
	f.defaultConfigs["gpt4"] = ProviderConfig{
		Name:          "gpt4",
		Priority:      90, // High priority
		Weight:        0.3,
		MaxConcurrent: 8,
		Timeout:       45 * time.Second,
		RetryPolicy: RetryPolicy{
			MaxRetries:      3,
			InitialDelay:    2 * time.Second,
			MaxDelay:        60 * time.Second,
			BackoffFactor:   2.0,
			RetryableErrors: []string{"rate_limit", "timeout", "server_error"},
		},
		CostLimit: CostLimit{
			DailyCostLimit:   75.0,   // $75/day (typically more expensive)
			MonthlyCostLimit: 1500.0, // $1500/month
			PerRequestLimit:  1.00,   // $1.00 per request
			AlertThreshold:   0.8,
		},
		Capabilities: []string{
			"instruction_parsing",
			"expectation_validation",
			"context_compression",
			"error_analysis",
			"natural_language",
			"code_understanding",
			"multimodal",
		},
		Regions: []string{"global"},
		Metadata: map[string]string{
			"model_family":    "gpt-4",
			"context_length":  "128000",
			"output_tokens":   "4096",
			"multimodal":      "true",
			"fine_tunable":    "true",
			"certified_level": "silver",
		},
	}

	// Gemini (Google) default configuration
	f.defaultConfigs["gemini"] = ProviderConfig{
		Name:          "gemini",
		Priority:      80, // Medium-high priority
		Weight:        0.2,
		MaxConcurrent: 6,
		Timeout:       35 * time.Second,
		RetryPolicy: RetryPolicy{
			MaxRetries:      3,
			InitialDelay:    1 * time.Second,
			MaxDelay:        45 * time.Second,
			BackoffFactor:   2.0,
			RetryableErrors: []string{"rate_limit", "timeout", "server_error"},
		},
		CostLimit: CostLimit{
			DailyCostLimit:   40.0,  // $40/day (competitive pricing)
			MonthlyCostLimit: 800.0, // $800/month
			PerRequestLimit:  0.40,  // $0.40 per request
			AlertThreshold:   0.8,
		},
		Capabilities: []string{
			"instruction_parsing",
			"expectation_validation",
			"context_compression",
			"error_analysis",
			"natural_language",
			"code_understanding",
			"multimodal",
		},
		Regions: []string{"global"},
		Metadata: map[string]string{
			"model_family":    "gemini-pro",
			"context_length":  "32000",
			"output_tokens":   "2048",
			"multimodal":      "true",
			"fine_tunable":    "false",
			"certified_level": "silver",
		},
	}

	// Llama (Meta/Local) default configuration
	f.defaultConfigs["llama"] = ProviderConfig{
		Name:          "llama",
		Priority:      70, // Lower priority (experimental)
		Weight:        0.1,
		MaxConcurrent: 4,
		Timeout:       60 * time.Second, // Local models may be slower
		RetryPolicy: RetryPolicy{
			MaxRetries:      2,
			InitialDelay:    3 * time.Second,
			MaxDelay:        30 * time.Second,
			BackoffFactor:   1.5,
			RetryableErrors: []string{"timeout", "server_error"},
		},
		CostLimit: CostLimit{
			DailyCostLimit:   10.0,  // $10/day (lower cost)
			MonthlyCostLimit: 200.0, // $200/month
			PerRequestLimit:  0.10,  // $0.10 per request
			AlertThreshold:   0.9,
		},
		Capabilities: []string{
			"instruction_parsing",
			"natural_language",
			"code_understanding",
		},
		Regions: []string{"local", "us-east-1"},
		Metadata: map[string]string{
			"model_family":    "llama-2",
			"context_length":  "4096",
			"output_tokens":   "2048",
			"multimodal":      "false",
			"fine_tunable":    "true",
			"certified_level": "bronze",
		},
	}

	// Local provider default configuration
	f.defaultConfigs["local"] = ProviderConfig{
		Name:          "local",
		Priority:      50, // Lowest priority (fallback)
		Weight:        0.05,
		MaxConcurrent: 2,
		Timeout:       120 * time.Second, // Very slow
		RetryPolicy: RetryPolicy{
			MaxRetries:      1,
			InitialDelay:    5 * time.Second,
			MaxDelay:        30 * time.Second,
			BackoffFactor:   1.0,
			RetryableErrors: []string{"timeout"},
		},
		CostLimit: CostLimit{
			DailyCostLimit:   1.0,  // $1/day (mostly compute cost)
			MonthlyCostLimit: 30.0, // $30/month
			PerRequestLimit:  0.01, // $0.01 per request
			AlertThreshold:   0.95,
		},
		Capabilities: []string{
			"instruction_parsing",
			"natural_language",
		},
		Regions: []string{"local"},
		Metadata: map[string]string{
			"model_family":    "local",
			"context_length":  "2048",
			"output_tokens":   "1024",
			"multimodal":      "false",
			"fine_tunable":    "true",
			"certified_level": "unverified",
		},
	}
}

// mergeWithDefaults merges user options with default configuration
func (f *ProviderFactory) mergeWithDefaults(options CreateProviderOptions) ProviderConfig {
	// Start with default config for this provider type
	config := f.defaultConfigs[string(options.Type)]

	// Override with user-provided config
	if options.Config.Name != "" {
		config.Name = options.Config.Name
	}
	if options.Config.Priority != 0 {
		config.Priority = options.Config.Priority
	}
	if options.Config.Weight != 0 {
		config.Weight = options.Config.Weight
	}
	if options.Config.MaxConcurrent != 0 {
		config.MaxConcurrent = options.Config.MaxConcurrent
	}
	if options.Config.Timeout != 0 {
		config.Timeout = options.Config.Timeout
	}

	// Apply any override values
	if options.Override != nil {
		if config.Metadata == nil {
			config.Metadata = make(map[string]string)
		}
		for key, value := range options.Override {
			config.Metadata[key] = value
		}
	}

	return config
}

// validateConfig validates provider configuration
func (f *ProviderFactory) validateConfig(config ProviderConfig) error {
	if config.Name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}
	if config.Weight < 0 || config.Weight > 1 {
		return fmt.Errorf("weight must be between 0 and 1")
	}
	if config.MaxConcurrent <= 0 {
		return fmt.Errorf("max concurrent must be positive")
	}
	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	return nil
}

// createClaudeProvider creates a Claude/Anthropic provider
func (f *ProviderFactory) createClaudeProvider(options CreateProviderOptions, config ProviderConfig) (ai.Provider, error) {
	if options.APIKey == "" {
		return nil, fmt.Errorf("Claude provider requires API key")
	}

	// Use existing Claude provider creation logic
	claudeOptions := []ai.ClaudeOption{}

	if options.Model != "" {
		claudeOptions = append(claudeOptions, ai.WithModel(options.Model))
	}

	provider := ai.NewClaudeProvider(options.APIKey, claudeOptions...)
	return provider, nil
}

// Placeholder implementations for other providers
// In a production system, these would integrate with respective SDKs

func (f *ProviderFactory) createGPT4Provider(options CreateProviderOptions, config ProviderConfig) (ai.Provider, error) {
	if options.APIKey == "" {
		return nil, fmt.Errorf("GPT-4 provider requires API key")
	}

	// Create GPT-4 provider options
	gpt4Options := []ai.GPT4Option{}

	if options.Model != "" {
		gpt4Options = append(gpt4Options, ai.WithGPT4Model(options.Model))
	}

	if options.Endpoint != "" {
		gpt4Options = append(gpt4Options, ai.WithGPT4BaseURL(options.Endpoint))
	}

	// Set max tokens based on config
	if config.Metadata["output_tokens"] != "" {
		if maxTokens := parseInt(config.Metadata["output_tokens"]); maxTokens > 0 {
			gpt4Options = append(gpt4Options, ai.WithGPT4MaxTokens(maxTokens))
		}
	}

	provider := ai.NewGPT4Provider(options.APIKey, gpt4Options...)
	return provider, nil
}

func (f *ProviderFactory) createGeminiProvider(options CreateProviderOptions, config ProviderConfig) (ai.Provider, error) {
	if options.APIKey == "" {
		return nil, fmt.Errorf("Gemini provider requires API key")
	}

	// Create Gemini provider options
	geminiOptions := []ai.GeminiOption{}

	if options.Model != "" {
		geminiOptions = append(geminiOptions, ai.WithGeminiModel(options.Model))
	}

	if options.Endpoint != "" {
		geminiOptions = append(geminiOptions, ai.WithGeminiBaseURL(options.Endpoint))
	}

	// Set max tokens based on config
	if config.Metadata["output_tokens"] != "" {
		if maxTokens := parseInt(config.Metadata["output_tokens"]); maxTokens > 0 {
			geminiOptions = append(geminiOptions, ai.WithGeminiMaxTokens(maxTokens))
		}
	}

	provider := ai.NewGeminiProvider(options.APIKey, geminiOptions...)
	return provider, nil
}

func (f *ProviderFactory) createLlamaProvider(options CreateProviderOptions, config ProviderConfig) (ai.Provider, error) {
	// TODO: Implement Llama provider (possibly via Ollama or similar)
	// This would integrate with local model serving infrastructure
	return nil, fmt.Errorf("Llama provider not yet implemented")
}

func (f *ProviderFactory) createLocalProvider(options CreateProviderOptions, config ProviderConfig) (ai.Provider, error) {
	// TODO: Implement local model provider
	// This would integrate with local inference engines
	return nil, fmt.Errorf("Local provider not yet implemented")
}

// GetSupportedProviders returns a list of supported provider types
func (f *ProviderFactory) GetSupportedProviders() []ProviderType {
	return []ProviderType{
		ProviderTypeClaude,
		ProviderTypeGPT4,
		ProviderTypeGemini,
		ProviderTypeLlama,
		ProviderTypeLocal,
	}
}

// GetProviderCapabilities returns the capabilities of a specific provider type
func (f *ProviderFactory) GetProviderCapabilities(providerType ProviderType) ([]string, error) {
	config, exists := f.defaultConfigs[string(providerType)]
	if !exists {
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
	return config.Capabilities, nil
}

// EstimateProviderCost estimates the cost for a provider type and usage pattern
func (f *ProviderFactory) EstimateProviderCost(providerType ProviderType, dailyRequests int) (float64, error) {
	config, exists := f.defaultConfigs[string(providerType)]
	if !exists {
		return 0, fmt.Errorf("unsupported provider type: %s", providerType)
	}

	// Simple cost estimation based on average request cost
	avgRequestCost := config.CostLimit.PerRequestLimit * 0.1 // Assume 10% of limit is typical
	return float64(dailyRequests) * avgRequestCost, nil
}

// parseInt parses a string to int, returning 0 if invalid
func parseInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}
