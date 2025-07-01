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

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
	"github.com/aws-research-wizard/tutorial-guard/pkg/interpreter"
)

func main() {
	fmt.Println("🔬 Testing Tutorial Guard on Real Spack-Manager-Go Tutorial")
	fmt.Println(strings.Repeat("=", 65))

	// Get API key from environment variable
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set ANTHROPIC_API_KEY environment variable")
	}
	client, err := ai.NewClaudeClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create Claude client: %v", err)
	}

	// Create tutorial interpreter
	config := interpreter.InterpreterConfig{
		MaxSteps:            20,
		StrictValidation:    false, // Allow some ambiguity for real tutorials
		AllowErrorRecovery:  true,
		ContextCompression:  true,
		ValidationThreshold: 0.7, // Lower threshold for real-world content
	}

	tutorialInterpreter := interpreter.NewTutorialInterpreter(client, config)

	// Real tutorial from Spack-Manager-Go README
	tutorial := &interpreter.Tutorial{
		Title:       "Installing Spack Manager Go",
		Description: "A step-by-step guide to install and set up Spack Manager Go for package management",
		Sections: []interpreter.TutorialSection{
			{
				Number:      1,
				Title:       "Prerequisites",
				Description: "Ensure you have the required software installed",
				Instructions: []interpreter.RawInstruction{
					{
						Text:    "Ensure you have Go 1.21 or later installed on your system",
						Context: "prerequisite-check",
					},
					{
						Text:    "Ensure you have a Spack installation (see Spack Documentation)",
						Context: "prerequisite-check",
					},
				},
			},
			{
				Number:      2,
				Title:       "Install from Source",
				Description: "Clone and build the project from source code",
				Instructions: []interpreter.RawInstruction{
					{
						Text:    "Clone the repository using: git clone https://github.com/spack-go/spack-manager.git",
						Context: "source-installation",
					},
					{
						Text:    "Navigate to the project directory: cd spack-manager",
						Context: "source-installation",
					},
					{
						Text:    "Build the project: go build -o spack-manager ./cmd/spack-manager",
						Context: "source-installation",
					},
					{
						Text:    "Install the binary to system path: sudo mv spack-manager /usr/local/bin/",
						Context: "source-installation",
					},
				},
			},
			{
				Number:      3,
				Title:       "Quick Start Setup",
				Description: "Set up environment and test the installation",
				Instructions: []interpreter.RawInstruction{
					{
						Text:    "Set the SPACK_ROOT environment variable: export SPACK_ROOT=\"/opt/spack\"",
						Context: "environment-setup",
					},
					{
						Text:    "Optionally set the binary cache for faster installations: export SPACK_BINARY_CACHE=\"https://cache.spack.io\"",
						Context: "environment-setup",
					},
					{
						Text:    "Test the installation by running: spack-manager --help",
						Context: "verification",
					},
				},
			},
		},
		Metadata: map[string]string{
			"source":     "spack-manager-go README.md",
			"language":   "bash",
			"platform":   "linux/macOS",
			"difficulty": "beginner",
			"duration":   "10-15 minutes",
		},
	}

	fmt.Printf("📖 Tutorial: %s\n", tutorial.Title)
	fmt.Printf("📋 Description: %s\n", tutorial.Description)
	fmt.Printf("🔢 Total Sections: %d\n", len(tutorial.Sections))
	fmt.Println()

	// Interpret the tutorial using AI
	ctx := context.Background()
	plan, err := tutorialInterpreter.InterpretTutorial(ctx, tutorial)
	if err != nil {
		log.Fatalf("Failed to interpret tutorial: %v", err)
	}

	// Display the AI's understanding
	fmt.Println("🧠 AI Tutorial Interpretation Results:")
	fmt.Println(strings.Repeat("=", 45))

	for _, step := range plan.Steps {
		fmt.Printf("\n📍 Section %d: %s\n", step.SectionNumber, step.Title)
		fmt.Printf("   📝 Description: %s\n", step.Description)
		fmt.Printf("   🔧 Instructions parsed: %d\n", len(step.Instructions))

		for j, instruction := range step.Instructions {
			fmt.Printf("\n   💡 Instruction %d:\n", j+1)
			fmt.Printf("      📄 Original: %s\n", instruction.Text)
			fmt.Printf("      🎯 AI Intent: %s\n", instruction.Intent)
			fmt.Printf("      📊 Confidence: %.2f\n", instruction.Confidence)
			fmt.Printf("      🔧 Actions: %d\n", len(instruction.Actions))

			for k, action := range instruction.Actions {
				fmt.Printf("         Action %d: %s\n", k+1, action.Description)
				fmt.Printf("         Command: %s\n", action.Command)
				if action.Validation.Expected != nil {
					fmt.Printf("         Validation: %s (expects: %v)\n",
						action.Validation.Type, action.Validation.Expected)
				}
			}

			if len(instruction.Prerequisites) > 0 {
				fmt.Printf("      ⚠️  Prerequisites: %v\n", instruction.Prerequisites)
			}

			if len(instruction.ExpectedOutcomes) > 0 {
				fmt.Printf("      ✅ Expected outcomes: %v\n", instruction.ExpectedOutcomes)
			}
		}
	}

	// Test a specific complex instruction parsing
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🔍 Deep Dive: Testing Complex Instruction Parsing")
	fmt.Println(strings.Repeat("=", 70))

	complexInstruction := "Clone the repository using: git clone https://github.com/spack-go/spack-manager.git"

	testContext := ai.TutorialContext{
		WorkingDirectory: "/tmp/tutorial-test",
		CurrentStep:      1,
		TotalSteps:       3,
		CreatedFiles:     []string{},
		ExecutedCommands: []string{},
		EnvironmentVars:  make(map[string]string),
	}

	parsed, err := client.ParseInstruction(ctx, complexInstruction, testContext)
	if err != nil {
		log.Printf("Failed to parse complex instruction: %v", err)
	} else {
		fmt.Printf("📝 Complex Instruction: %s\n", complexInstruction)
		fmt.Printf("🎯 AI Understanding: %s\n", parsed.Intent)
		fmt.Printf("📊 Confidence: %.2f\n", parsed.Confidence)
		fmt.Printf("💭 AI Reasoning: %s\n", parsed.Reasoning)

		fmt.Println("\n🔧 Parsed Actions:")
		for i, action := range parsed.Actions {
			fmt.Printf("   %d. %s\n", i+1, action.Description)
			fmt.Printf("      Command: %s\n", action.Command)
			fmt.Printf("      Type: %s\n", action.Type)
		}
	}

	// Test error handling with a problematic scenario
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🚨 Testing Error Interpretation")
	fmt.Println(strings.Repeat("=", 70))

	errorScenario := "fatal: repository 'https://github.com/spack-go/spack-manager.git' not found"

	errorInterpretation, err := client.InterpretError(ctx, errorScenario, testContext)
	if err != nil {
		log.Printf("Failed to interpret error: %v", err)
	} else {
		fmt.Printf("⚠️  Error: %s\n", errorScenario)
		fmt.Printf("🏷️  Type: %s\n", errorInterpretation.ErrorType)
		fmt.Printf("📖 Explanation: %s\n", errorInterpretation.Explanation)

		fmt.Println("\n💡 AI-Suggested Solutions:")
		for i, solution := range errorInterpretation.Solutions {
			fmt.Printf("   %d. (%.0f%% confidence) %s\n",
				i+1, solution.Probability*100, solution.Description)
			for _, cmd := range solution.Commands {
				fmt.Printf("      → %s\n", cmd)
			}
		}
	}

	// Display final statistics
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("📈 Tutorial Guard Performance Summary")
	fmt.Println(strings.Repeat("=", 70))

	stats := client.GetUsageStats()
	capabilities := client.GetCapabilities()

	fmt.Printf("🤖 AI Provider: %s (%s)\n", capabilities.Name, capabilities.Version)
	fmt.Printf("📊 Quality Score: %.2f/1.0\n", capabilities.QualityScore)
	fmt.Printf("🔢 Total API Requests: %d\n", stats.RequestsTotal)
	fmt.Printf("🎯 Total Tokens Used: %d\n", stats.TokensUsedTotal)
	fmt.Printf("💰 Total Cost: $%.4f\n", stats.CostTotal)
	fmt.Printf("📋 Instructions Processed: %d\n", countTotalInstructions(tutorial))
	fmt.Printf("🧠 Sections Interpreted: %d/%d\n", len(plan.Steps), len(tutorial.Sections))

	fmt.Println("\n🎉 Real Tutorial Test Complete!")
	fmt.Println("💡 Tutorial Guard successfully understood and parsed a real-world tutorial!")
}

func countTotalInstructions(tutorial *interpreter.Tutorial) int {
	total := 0
	for _, section := range tutorial.Sections {
		total += len(section.Instructions)
	}
	return total
}
