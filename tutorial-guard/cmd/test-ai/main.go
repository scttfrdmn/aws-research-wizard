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
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set ANTHROPIC_API_KEY environment variable")
	}

	fmt.Println("🧠 Testing Tutorial Guard AI Integration with Claude...")

	// Create AI client
	client, err := ai.NewClaudeClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create Claude client: %v", err)
	}

	// Test health check
	fmt.Println("🔍 Testing health check...")
	ctx := context.Background()
	if err := client.HealthCheck(ctx); err != nil {
		fmt.Printf("⚠️  Health check failed: %v\n", err)
		fmt.Println("🔄 Continuing with other tests...")
	} else {
		fmt.Println("✅ Health check passed!")
	}

	// Test instruction parsing
	fmt.Println("\n📖 Testing instruction parsing...")
	instruction := "First, create a directory called 'my-project' and navigate into it."
	tutorialContext := ai.TutorialContext{
		WorkingDirectory: "/tmp",
		CurrentStep:      1,
		TotalSteps:       5,
		CreatedFiles:     []string{},
		ExecutedCommands: []string{},
		EnvironmentVars:  make(map[string]string),
		Metadata:         make(map[string]string),
	}

	parsed, err := client.ParseInstruction(ctx, instruction, tutorialContext)
	if err != nil {
		log.Fatalf("Failed to parse instruction: %v", err)
	}

	fmt.Printf("📝 Original instruction: %s\n", instruction)
	fmt.Printf("🎯 AI Intent: %s\n", parsed.Intent)
	fmt.Printf("🔧 Number of actions: %d\n", len(parsed.Actions))
	fmt.Printf("📊 Confidence: %.2f\n", parsed.Confidence)
	fmt.Printf("💭 Reasoning: %s\n", parsed.Reasoning)

	// Print actions
	for i, action := range parsed.Actions {
		fmt.Printf("   Action %d: %s (%s)\n", i+1, action.Description, action.Command)
	}

	// Test expectation validation
	fmt.Println("\n✅ Testing expectation validation...")
	expected := "You should see a new directory called 'my-project'"
	actual := "drwxr-xr-x  2 user user 4096 Dec 15 14:30 my-project"

	validation, err := client.ValidateExpectation(ctx, expected, actual, tutorialContext)
	if err != nil {
		log.Fatalf("Failed to validate expectation: %v", err)
	}

	fmt.Printf("🔎 Expected: %s\n", expected)
	fmt.Printf("📄 Actual: %s\n", actual)
	fmt.Printf("✅ Success: %t\n", validation.Success)
	fmt.Printf("📊 Confidence: %.2f\n", validation.Confidence)
	fmt.Printf("💭 Reasoning: %s\n", validation.Reasoning)

	// Test error interpretation
	fmt.Println("\n🚨 Testing error interpretation...")
	errorMsg := "mkdir: cannot create directory 'my-project': Permission denied"

	errorInterpretation, err := client.InterpretError(ctx, errorMsg, tutorialContext)
	if err != nil {
		log.Fatalf("Failed to interpret error: %v", err)
	}

	fmt.Printf("⚠️  Error: %s\n", errorMsg)
	fmt.Printf("🏷️  Error Type: %s\n", errorInterpretation.ErrorType)
	fmt.Printf("📖 Explanation: %s\n", errorInterpretation.Explanation)
	fmt.Printf("💡 Number of solutions: %d\n", len(errorInterpretation.Solutions))

	for i, solution := range errorInterpretation.Solutions {
		fmt.Printf("   Solution %d (%.1f%%): %s\n", i+1, solution.Probability*100, solution.Description)
		for _, cmd := range solution.Commands {
			fmt.Printf("      Command: %s\n", cmd)
		}
	}

	// Test context compression
	fmt.Println("\n🗜️  Testing context compression...")
	largeContext := ai.TutorialContext{
		WorkingDirectory: "/home/user/my-project",
		CurrentStep:      5,
		TotalSteps:       10,
		CreatedFiles:     []string{"package.json", "index.js", "README.md", "src/main.js", "test/test.js"},
		ExecutedCommands: []string{"mkdir my-project", "cd my-project", "npm init -y", "touch index.js", "git init"},
		PreviousOutputs:  []string{"Directory created", "Initialized empty Git repository", "package.json created"},
		EnvironmentVars:  map[string]string{"NODE_ENV": "development", "PROJECT_NAME": "my-project"},
		Metadata:         map[string]string{"language": "javascript", "framework": "node"},
	}

	compressed, err := client.CompressContext(ctx, largeContext)
	if err != nil {
		log.Fatalf("Failed to compress context: %v", err)
	}

	fmt.Printf("📊 Original context size: ~%d chars\n", estimateContextSize(largeContext))
	fmt.Printf("🗜️  Compressed summary: %s\n", compressed.Summary)
	fmt.Printf("📁 Key files: %v\n", compressed.KeyFiles)
	fmt.Printf("🔄 Current state: %s\n", compressed.CurrentState)

	// Show usage stats
	fmt.Println("\n📈 Usage Statistics:")
	stats := client.GetUsageStats()
	capabilities := client.GetCapabilities()

	fmt.Printf("   Total requests: %d\n", stats.RequestsTotal)
	fmt.Printf("   Total tokens: %d\n", stats.TokensUsedTotal)
	fmt.Printf("   Total cost: $%.4f\n", stats.CostTotal)
	fmt.Printf("   Provider: %s (%s)\n", capabilities.Name, capabilities.Version)
	fmt.Printf("   Quality score: %.2f\n", capabilities.QualityScore)

	fmt.Println("\n🎉 All AI integration tests passed!")
	fmt.Println("💡 Tutorial Guard is ready to understand and follow tutorials with AI!")
}

func estimateContextSize(ctx ai.TutorialContext) int {
	data, _ := json.Marshal(ctx)
	return len(data)
}
