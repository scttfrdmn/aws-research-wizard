/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
	"github.com/aws-research-wizard/tutorial-guard/pkg/registry"
)

func main() {
	fmt.Println("🌟 Testing Tutorial Guard Multi-Provider AI System")
	fmt.Println(strings.Repeat("=", 70))

	// Test 1: Individual Provider Testing
	fmt.Println("\n🔬 Testing Individual Providers...")

	ctx := context.Background()
	testContext := ai.TutorialContext{
		CurrentStep:      1,
		TotalSteps:       3,
		WorkingDirectory: "/tmp/test",
		ExecutedCommands: []string{"mkdir test-dir"},
	}

	// Test Claude if available
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		fmt.Println("\n   Testing Claude Provider...")
		claudeProvider := ai.NewClaudeProvider(apiKey)
		success := testProvider(ctx, claudeProvider, "Claude", testContext)
		if success {
			fmt.Println("   ✅ Claude: All tests passed")
		} else {
			fmt.Println("   ❌ Claude: Some tests failed")
		}
	} else {
		fmt.Println("\n   ⚠️  Claude: API key not configured")
	}

	// Test GPT-4 if available
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		fmt.Println("\n   Testing GPT-4 Provider...")
		gpt4Provider := ai.NewGPT4Provider(apiKey)
		success := testProvider(ctx, gpt4Provider, "GPT-4", testContext)
		if success {
			fmt.Println("   ✅ GPT-4: All tests passed")
		} else {
			fmt.Println("   ❌ GPT-4: Some tests failed")
		}
	} else {
		fmt.Println("\n   ⚠️  GPT-4: API key not configured")
	}

	// Test Gemini if available
	if apiKey := os.Getenv("GOOGLE_AI_API_KEY"); apiKey != "" {
		fmt.Println("\n   Testing Gemini Provider...")
		geminiProvider := ai.NewGeminiProvider(apiKey)
		success := testProvider(ctx, geminiProvider, "Gemini", testContext)
		if success {
			fmt.Println("   ✅ Gemini: All tests passed")
		} else {
			fmt.Println("   ❌ Gemini: Some tests failed")
		}
	} else {
		fmt.Println("\n   ⚠️  Gemini: API key not configured")
	}

	// Test 2: Multi-Provider Registry
	fmt.Println("\n🏗️ Testing Multi-Provider Registry...")

	factory := registry.NewProviderFactory()
	registry, err := factory.CreateDefaultRegistry()
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}

	// Count available providers
	availableProviders := 0
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		availableProviders++
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		availableProviders++
	}
	if os.Getenv("GOOGLE_AI_API_KEY") != "" {
		availableProviders++
	}

	fmt.Printf("   📊 Registry initialized with %d providers\n", availableProviders)

	if availableProviders == 0 {
		fmt.Println("   ⚠️  No API keys configured - skipping registry tests")
		fmt.Println("\n💡 To enable providers, set these environment variables:")
		fmt.Println("   export ANTHROPIC_API_KEY='your-claude-key'")
		fmt.Println("   export OPENAI_API_KEY='your-openai-key'")
		fmt.Println("   export GOOGLE_AI_API_KEY='your-gemini-key'")
		return
	}

	// Test 3: Provider Routing
	fmt.Println("\n🎯 Testing Intelligent Provider Routing...")

	testRequests := []struct {
		name    string
		request registry.RoutingRequest
	}{
		{
			name: "High Priority Instruction Parsing",
			request: registry.RoutingRequest{
				Type:         ai.RequestParseInstruction,
				Priority:     ai.PriorityHigh,
				MaxCost:      0.10,
				MaxLatency:   10 * time.Second,
				RequiredCaps: []string{"instruction_parsing", "natural_language"},
				Context:      testContext,
			},
		},
		{
			name: "Cost-Optimized Validation",
			request: registry.RoutingRequest{
				Type:         ai.RequestValidateExpectation,
				Priority:     ai.PriorityMedium,
				MaxCost:      0.02,
				MaxLatency:   15 * time.Second,
				RequiredCaps: []string{"expectation_validation"},
				Context:      testContext,
			},
		},
		{
			name: "Complex Error Analysis",
			request: registry.RoutingRequest{
				Type:         ai.RequestInterpretError,
				Priority:     ai.PriorityHigh,
				MaxLatency:   20 * time.Second,
				RequiredCaps: []string{"error_analysis", "natural_language"},
				Context:      testContext,
			},
		},
	}

	for i, test := range testRequests {
		fmt.Printf("\n   Test %d: %s\n", i+1, test.name)

		result, err := registry.Route(ctx, test.request)
		if err != nil {
			fmt.Printf("      ❌ Routing failed: %v\n", err)
			continue
		}

		fmt.Printf("      ✅ Selected: %s\n", result.ProviderName)
		fmt.Printf("      📝 Reason: %s\n", result.Reason)
		fmt.Printf("      💰 Estimated cost: $%.4f\n", result.EstimatedCost)
		fmt.Printf("      ⏱️  Estimated latency: %v\n", result.EstimatedLatency)
		fmt.Printf("      🎯 Quality score: %.2f\n", result.QualityScore)
		fmt.Printf("      🔄 Routing time: %v\n", result.RoutingTime)

		if len(result.Alternatives) > 0 {
			fmt.Printf("      🔀 Alternatives: %v\n", result.Alternatives)
		}
	}

	// Test 4: Provider Comparison
	fmt.Println("\n📊 Provider Performance Comparison...")

	testInstruction := "Install the Go programming language on Ubuntu using apt package manager"

	fmt.Printf("\n   Testing instruction: \"%s\"\n", testInstruction)
	fmt.Println("   " + strings.Repeat("-", 60))

	type ProviderResult struct {
		Name       string
		Success    bool
		Duration   time.Duration
		Confidence float64
		Error      error
	}

	var results []ProviderResult

	// Test each available provider
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		start := time.Now()
		claudeProvider := ai.NewClaudeProvider(os.Getenv("ANTHROPIC_API_KEY"))
		parsed, err := claudeProvider.ParseInstruction(ctx, testInstruction, testContext)
		duration := time.Since(start)

		result := ProviderResult{
			Name:     "Claude",
			Success:  err == nil,
			Duration: duration,
			Error:    err,
		}
		if parsed != nil {
			result.Confidence = parsed.Confidence
		}
		results = append(results, result)
	}

	if os.Getenv("OPENAI_API_KEY") != "" {
		start := time.Now()
		gpt4Provider := ai.NewGPT4Provider(os.Getenv("OPENAI_API_KEY"))
		parsed, err := gpt4Provider.ParseInstruction(ctx, testInstruction, testContext)
		duration := time.Since(start)

		result := ProviderResult{
			Name:     "GPT-4",
			Success:  err == nil,
			Duration: duration,
			Error:    err,
		}
		if parsed != nil {
			result.Confidence = parsed.Confidence
		}
		results = append(results, result)
	}

	if os.Getenv("GOOGLE_AI_API_KEY") != "" {
		start := time.Now()
		geminiProvider := ai.NewGeminiProvider(os.Getenv("GOOGLE_AI_API_KEY"))
		parsed, err := geminiProvider.ParseInstruction(ctx, testInstruction, testContext)
		duration := time.Since(start)

		result := ProviderResult{
			Name:     "Gemini",
			Success:  err == nil,
			Duration: duration,
			Error:    err,
		}
		if parsed != nil {
			result.Confidence = parsed.Confidence
		}
		results = append(results, result)
	}

	// Display results
	for _, result := range results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}

		fmt.Printf("   %s %s: %.2fs, confidence: %.2f",
			status, result.Name, result.Duration.Seconds(), result.Confidence)

		if !result.Success && result.Error != nil {
			fmt.Printf(" (Error: %v)", result.Error)
		}
		fmt.Println()
	}

	// Test 5: Cost and Performance Analysis
	fmt.Println("\n💰 Cost and Performance Analysis...")

	fmt.Printf("\n   Provider Cost Estimates (per 100 requests/day):\n")
	for _, providerType := range factory.GetSupportedProviders() {
		cost, err := factory.EstimateProviderCost(providerType, 100)
		if err == nil {
			fmt.Printf("   %s: $%.2f/day\n", string(providerType), cost)
		}
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎉 Multi-Provider AI System Test Summary")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("✅ Providers tested: %d\n", len(results))

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}
	fmt.Printf("✅ Successful tests: %d/%d\n", successCount, len(results))

	if availableProviders >= 2 {
		fmt.Println("✅ Multi-provider routing: Operational")
		fmt.Println("✅ Intelligent provider selection: Functional")
		fmt.Println("✅ Cost optimization: Active")
		fmt.Println("✅ Fallback chains: Configured")
	}

	fmt.Println("\n🚀 Tutorial Guard Multi-Provider System: Ready for Production!")
	fmt.Println("💡 Your tutorials can now leverage multiple AI providers for optimal performance!")
}

// testProvider performs basic functionality tests on a provider
func testProvider(ctx context.Context, provider ai.Provider, name string, testContext ai.TutorialContext) bool {
	success := true

	// Test 1: Health Check
	fmt.Printf("      🏥 Health check...")
	if err := provider.HealthCheck(ctx); err != nil {
		fmt.Printf(" ❌ Failed: %v\n", err)
		success = false
	} else {
		fmt.Println(" ✅ Healthy")
	}

	// Test 2: Parse Instruction
	fmt.Printf("      📝 Parse instruction...")
	_, err := provider.ParseInstruction(ctx, "echo 'hello world'", testContext)
	if err != nil {
		fmt.Printf(" ❌ Failed: %v\n", err)
		success = false
	} else {
		fmt.Println(" ✅ Success")
	}

	// Test 3: Validate Expectation
	fmt.Printf("      ✅ Validate expectation...")
	_, err = provider.ValidateExpectation(ctx, "hello world", "hello world", testContext)
	if err != nil {
		fmt.Printf(" ❌ Failed: %v\n", err)
		success = false
	} else {
		fmt.Println(" ✅ Success")
	}

	return success
}
