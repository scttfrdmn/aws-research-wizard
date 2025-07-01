/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package interpreter

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// TutorialInterpreter understands and follows tutorials using AI
type TutorialInterpreter struct {
	aiClient *ai.Client
	context  *ai.TutorialContext
	config   InterpreterConfig
}

// InterpreterConfig holds configuration for the interpreter
type InterpreterConfig struct {
	MaxSteps            int     `json:"max_steps"`
	StrictValidation    bool    `json:"strict_validation"`
	AllowErrorRecovery  bool    `json:"allow_error_recovery"`
	ContextCompression  bool    `json:"context_compression"`
	ValidationThreshold float64 `json:"validation_threshold"`
}

// NewTutorialInterpreter creates a new tutorial interpreter
func NewTutorialInterpreter(aiClient *ai.Client, config InterpreterConfig) *TutorialInterpreter {
	if config.ValidationThreshold == 0 {
		config.ValidationThreshold = 0.8
	}

	return &TutorialInterpreter{
		aiClient: aiClient,
		context: &ai.TutorialContext{
			WorkingDirectory: "/tmp/tutorial-guard",
			EnvironmentVars:  make(map[string]string),
			CreatedFiles:     []string{},
			ExecutedCommands: []string{},
			PreviousOutputs:  []string{},
			CurrentStep:      0,
			TotalSteps:       0,
			Metadata:         make(map[string]string),
		},
		config: config,
	}
}

// InterpretTutorial processes an entire tutorial and generates execution plans
func (t *TutorialInterpreter) InterpretTutorial(ctx context.Context, tutorial *Tutorial) (*TutorialPlan, error) {
	plan := &TutorialPlan{
		Title:       tutorial.Title,
		Description: tutorial.Description,
		Steps:       []TutorialStep{},
		Context:     *t.context,
	}

	t.context.TotalSteps = len(tutorial.Sections)

	for i, section := range tutorial.Sections {
		t.context.CurrentStep = i + 1

		step, err := t.interpretSection(ctx, section)
		if err != nil {
			return nil, fmt.Errorf("failed to interpret section %d: %w", i+1, err)
		}

		plan.Steps = append(plan.Steps, *step)

		// Compress context if it's getting too large
		if t.config.ContextCompression && t.shouldCompressContext() {
			compressed, err := t.aiClient.CompressContext(ctx, *t.context)
			if err == nil {
				t.updateContextFromCompressed(compressed)
			}
		}
	}

	return plan, nil
}

// interpretSection processes a single tutorial section
func (t *TutorialInterpreter) interpretSection(ctx context.Context, section TutorialSection) (*TutorialStep, error) {
	step := &TutorialStep{
		SectionNumber: section.Number,
		Title:         section.Title,
		Description:   section.Description,
		Instructions:  []Instruction{},
		Context:       *t.context,
	}

	// Process each instruction in the section
	for _, instruction := range section.Instructions {
		interpreted, err := t.interpretInstruction(ctx, instruction)
		if err != nil {
			return nil, fmt.Errorf("failed to interpret instruction '%s': %w", instruction.Text, err)
		}

		step.Instructions = append(step.Instructions, *interpreted)
	}

	// Process any code examples
	for _, example := range section.CodeExamples {
		codeInstruction, err := t.interpretCodeExample(ctx, example)
		if err != nil {
			return nil, fmt.Errorf("failed to interpret code example: %w", err)
		}

		step.Instructions = append(step.Instructions, *codeInstruction)
	}

	return step, nil
}

// interpretInstruction processes a single instruction using AI
func (t *TutorialInterpreter) interpretInstruction(ctx context.Context, instruction RawInstruction) (*Instruction, error) {
	// Use AI to parse the instruction
	parsed, err := t.aiClient.ParseInstruction(ctx, instruction.Text, *t.context)
	if err != nil {
		return nil, fmt.Errorf("AI parsing failed: %w", err)
	}

	// Convert AI response to our instruction format
	aiInstruction := &Instruction{
		Text:             instruction.Text,
		Type:             InstructionTypeAction,
		Intent:           parsed.Intent,
		Actions:          convertAIActions(parsed.Actions),
		Prerequisites:    parsed.Prerequisites,
		ExpectedOutcomes: parsed.ExpectedOutcomes,
		Confidence:       parsed.Confidence,
		Reasoning:        parsed.Reasoning,
		Metadata:         parsed.Metadata,
	}

	// Validate confidence threshold
	if t.config.StrictValidation && parsed.Confidence < t.config.ValidationThreshold {
		return nil, fmt.Errorf("AI confidence %f below threshold %f", parsed.Confidence, t.config.ValidationThreshold)
	}

	return aiInstruction, nil
}

// interpretCodeExample processes code examples
func (t *TutorialInterpreter) interpretCodeExample(ctx context.Context, example extractor.Example) (*Instruction, error) {
	// Create instruction from code example
	instructionText := fmt.Sprintf("Execute the following %s code: %s", example.Language, example.Code)

	instruction := &Instruction{
		Text: instructionText,
		Type: InstructionTypeCode,
		Actions: []Action{
			{
				Type:        ActionTypeCommand,
				Command:     example.Code,
				Description: fmt.Sprintf("Execute %s code", example.Language),
				Language:    example.Language,
				Validation: ValidationRule{
					Type:      ValidationTypeExitCode,
					Expected:  0,
					Condition: "exit_code == 0",
				},
			},
		},
		ExpectedOutcomes: []string{"Code executes successfully"},
		Confidence:       0.9, // High confidence for direct code execution
		Metadata: map[string]string{
			"source":     example.Source,
			"language":   example.Language,
			"line_start": fmt.Sprintf("%d", example.LineStart),
			"line_end":   fmt.Sprintf("%d", example.LineEnd),
		},
	}

	return instruction, nil
}

// ValidateOutcome validates if an execution result matches expectations
func (t *TutorialInterpreter) ValidateOutcome(ctx context.Context, expected string, actual ExecutionResult) (*ai.ValidationResult, error) {
	actualOutput := actual.Output
	if actual.ErrorOutput != "" {
		actualOutput += "\nErrors: " + actual.ErrorOutput
	}

	return t.aiClient.ValidateExpectation(ctx, expected, actualOutput, *t.context)
}

// HandleError interprets errors and suggests recovery actions
func (t *TutorialInterpreter) HandleError(ctx context.Context, errorMsg string) (*ai.ErrorInterpretation, error) {
	return t.aiClient.InterpretError(ctx, errorMsg, *t.context)
}

// UpdateContext updates the tutorial context with execution results
func (t *TutorialInterpreter) UpdateContext(result ExecutionResult) {
	t.context.ExecutedCommands = append(t.context.ExecutedCommands, result.Command)
	t.context.PreviousOutputs = append(t.context.PreviousOutputs, result.Output)

	// Update working directory if changed
	if result.WorkingDirectory != "" {
		t.context.WorkingDirectory = result.WorkingDirectory
	}

	// Track created files
	if result.CreatedFiles != nil {
		t.context.CreatedFiles = append(t.context.CreatedFiles, result.CreatedFiles...)
	}

	// Update environment variables
	if result.EnvironmentVars != nil {
		for key, value := range result.EnvironmentVars {
			t.context.EnvironmentVars[key] = value
		}
	}
}

// shouldCompressContext determines if context should be compressed
func (t *TutorialInterpreter) shouldCompressContext() bool {
	// Compress if context is getting large
	totalContextSize := len(strings.Join(t.context.ExecutedCommands, " ")) +
		len(strings.Join(t.context.PreviousOutputs, " ")) +
		len(strings.Join(t.context.CreatedFiles, " "))

	return totalContextSize > 5000 // 5KB threshold
}

// updateContextFromCompressed updates context with compressed version
func (t *TutorialInterpreter) updateContextFromCompressed(compressed *ai.CompressedContext) {
	// Keep essential information, compress the rest
	t.context.Metadata["compressed_summary"] = compressed.Summary
	t.context.CreatedFiles = compressed.KeyFiles

	// Truncate long arrays to save space
	if len(t.context.ExecutedCommands) > 10 {
		t.context.ExecutedCommands = t.context.ExecutedCommands[len(t.context.ExecutedCommands)-5:]
	}
	if len(t.context.PreviousOutputs) > 10 {
		t.context.PreviousOutputs = t.context.PreviousOutputs[len(t.context.PreviousOutputs)-5:]
	}
}

// convertAIActions converts AI actions to our action format
func convertAIActions(aiActions []ai.Action) []Action {
	actions := make([]Action, len(aiActions))

	for i, aiAction := range aiActions {
		actions[i] = Action{
			Type:        ActionType(aiAction.Type),
			Command:     aiAction.Command,
			Description: aiAction.Description,
			Validation:  convertValidationRule(aiAction.Validation),
			Timeout:     aiAction.Timeout,
			Metadata:    aiAction.Metadata,
		}
	}

	return actions
}

// convertValidationRule converts AI validation rule to our format
func convertValidationRule(aiRule ai.ValidationRule) ValidationRule {
	return ValidationRule{
		Type:      ValidationType(aiRule.Type),
		Condition: aiRule.Condition,
		Expected:  aiRule.Expected,
		Operator:  aiRule.Operator,
		Metadata:  aiRule.Metadata,
	}
}
