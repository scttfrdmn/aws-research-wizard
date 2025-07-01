/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package interpreter

import (
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
)

// Tutorial represents a complete tutorial document
type Tutorial struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Sections    []TutorialSection `json:"sections"`
	Metadata    map[string]string `json:"metadata"`
}

// TutorialSection represents a section of a tutorial
type TutorialSection struct {
	Number       int                 `json:"number"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Instructions []RawInstruction    `json:"instructions"`
	CodeExamples []extractor.Example `json:"code_examples"`
	Metadata     map[string]string   `json:"metadata"`
}

// RawInstruction represents an unprocessed instruction from the tutorial
type RawInstruction struct {
	Text     string            `json:"text"`
	Context  string            `json:"context"`
	Metadata map[string]string `json:"metadata"`
}

// TutorialPlan represents the AI's understanding and execution plan for a tutorial
type TutorialPlan struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Steps       []TutorialStep     `json:"steps"`
	Context     ai.TutorialContext `json:"context"`
	Metadata    map[string]string  `json:"metadata"`
}

// TutorialStep represents a single step in the tutorial execution plan
type TutorialStep struct {
	SectionNumber int                `json:"section_number"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Instructions  []Instruction      `json:"instructions"`
	Context       ai.TutorialContext `json:"context"`
	Metadata      map[string]string  `json:"metadata"`
}

// Instruction represents an AI-interpreted instruction
type Instruction struct {
	Text             string            `json:"text"`
	Type             InstructionType   `json:"type"`
	Intent           string            `json:"intent"`
	Actions          []Action          `json:"actions"`
	Prerequisites    []string          `json:"prerequisites"`
	ExpectedOutcomes []string          `json:"expected_outcomes"`
	Confidence       float64           `json:"confidence"`
	Reasoning        string            `json:"reasoning"`
	Metadata         map[string]string `json:"metadata"`
}

// InstructionType defines the type of instruction
type InstructionType string

const (
	InstructionTypeAction      InstructionType = "action"      // Execute an action
	InstructionTypeValidation  InstructionType = "validation"  // Validate something
	InstructionTypeInformation InstructionType = "information" // Informational text
	InstructionTypeCode        InstructionType = "code"        // Execute code
	InstructionTypeConditional InstructionType = "conditional" // Conditional logic
)

// Action represents a specific action to be executed
type Action struct {
	Type        ActionType        `json:"type"`
	Command     string            `json:"command"`
	Description string            `json:"description"`
	Language    string            `json:"language"` // For code actions
	Validation  ValidationRule    `json:"validation"`
	Timeout     time.Duration     `json:"timeout"`
	Metadata    map[string]string `json:"metadata"`
}

// ActionType defines the types of actions
type ActionType string

const (
	ActionTypeCommand     ActionType = "command"     // Execute shell command
	ActionTypeValidate    ActionType = "validate"    // Validate a condition
	ActionTypeWait        ActionType = "wait"        // Wait for a condition
	ActionTypeDownload    ActionType = "download"    // Download a file
	ActionTypeExtract     ActionType = "extract"     // Extract an archive
	ActionTypeNavigate    ActionType = "navigate"    // Change directory
	ActionTypeCheck       ActionType = "check"       // Check if something exists
	ActionTypeConditional ActionType = "conditional" // Conditional logic
)

// ValidationRule defines how to validate an action's success
type ValidationRule struct {
	Type      ValidationType    `json:"type"`
	Condition string            `json:"condition"`
	Expected  interface{}       `json:"expected"` // Can be string, bool, number
	Operator  string            `json:"operator"`
	Metadata  map[string]string `json:"metadata"`
}

// ValidationType defines types of validation
type ValidationType string

const (
	ValidationTypeExitCode   ValidationType = "exit_code"
	ValidationTypeFileExists ValidationType = "file_exists"
	ValidationTypeOutput     ValidationType = "output"
	ValidationTypeContains   ValidationType = "contains"
	ValidationTypeRegex      ValidationType = "regex"
	ValidationTypeCustom     ValidationType = "custom"
)

// ExecutionResult represents the result of executing an action
type ExecutionResult struct {
	Action           Action            `json:"action"`
	Success          bool              `json:"success"`
	ExitCode         int               `json:"exit_code"`
	Output           string            `json:"output"`
	ErrorOutput      string            `json:"error_output"`
	Duration         time.Duration     `json:"duration"`
	Command          string            `json:"command"`
	WorkingDirectory string            `json:"working_directory"`
	CreatedFiles     []string          `json:"created_files"`
	EnvironmentVars  map[string]string `json:"environment_vars"`
	Metadata         map[string]string `json:"metadata"`
}
