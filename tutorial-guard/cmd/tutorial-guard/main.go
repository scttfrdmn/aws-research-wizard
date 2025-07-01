/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package main

import (
	"fmt"
	"os"

	"github.com/aws-research-wizard/tutorial-guard/pkg/extractor"
	"github.com/aws-research-wizard/tutorial-guard/pkg/reporter"
	"github.com/aws-research-wizard/tutorial-guard/pkg/runner"
	"github.com/aws-research-wizard/tutorial-guard/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tutorial-guard",
		Short: "Automated testing for documentation tutorials",
		Long: `Tutorial Guard automatically extracts and tests code examples from documentation
to ensure they remain accurate and functional.

Examples:
  tutorial-guard scan ./docs/          # Extract code examples
  tutorial-guard test --env docker     # Test in Docker environment
  tutorial-guard validate --format md  # Validate Markdown files`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add subcommands
	rootCmd.AddCommand(
		newScanCommand(),
		newTestCommand(),
		newValidateCommand(),
		newVersionCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newScanCommand() *cobra.Command {
	var (
		format    string
		output    string
		recursive bool
		include   []string
		exclude   []string
	)

	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Extract code examples from documentation",
		Long: `Scan documentation files and extract testable code examples.

Supports Markdown, HTML, and other documentation formats.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := extractor.Config{
				Format:    format,
				Recursive: recursive,
				Include:   include,
				Exclude:   exclude,
			}

			ext := extractor.New(config)
			examples, err := ext.ExtractFromPath(args[0])
			if err != nil {
				return fmt.Errorf("failed to extract examples: %w", err)
			}

			rep := reporter.New(output)
			return rep.WriteExamples(examples)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "auto", "Documentation format (markdown, html, auto)")
	cmd.Flags().StringVarP(&output, "output", "o", "examples.json", "Output file for extracted examples")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Scan directories recursively")
	cmd.Flags().StringSliceVar(&include, "include", []string{"*.md", "*.html"}, "File patterns to include")
	cmd.Flags().StringSliceVar(&exclude, "exclude", []string{}, "File patterns to exclude")

	return cmd
}

func newTestCommand() *cobra.Command {
	var (
		environment string
		timeout     string
		parallel    int
		cleanup     bool
		verbose     bool
	)

	cmd := &cobra.Command{
		Use:   "test [examples-file]",
		Short: "Execute extracted code examples",
		Long: `Run extracted code examples in isolated environments and validate results.

Supports Docker, local shell, and cloud environments.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			examplesFile := "examples.json"
			if len(args) > 0 {
				examplesFile = args[0]
			}

			config := runner.Config{
				Environment: environment,
				Timeout:     timeout,
				Parallel:    parallel,
				Cleanup:     cleanup,
				Verbose:     verbose,
			}

			r := runner.New(config)
			results, err := r.RunExamplesFromFile(examplesFile)
			if err != nil {
				return fmt.Errorf("failed to run tests: %w", err)
			}

			val := validator.New()
			return val.ValidateResults(results)
		},
	}

	cmd.Flags().StringVarP(&environment, "environment", "e", "docker", "Test environment (docker, local, aws)")
	cmd.Flags().StringVar(&timeout, "timeout", "5m", "Timeout for individual tests")
	cmd.Flags().IntVarP(&parallel, "parallel", "p", 1, "Number of parallel test runners")
	cmd.Flags().BoolVar(&cleanup, "cleanup", true, "Cleanup resources after tests")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func newValidateCommand() *cobra.Command {
	var (
		format string
		strict bool
	)

	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate documentation format and examples",
		Long: `Validate that documentation follows best practices and
contains well-formed, testable examples.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := validator.Config{
				Format: format,
				Strict: strict,
			}

			val := validator.New()
			return val.ValidatePath(args[0], config)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "auto", "Documentation format")
	cmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation rules")

	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tutorial-guard %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Built: %s\n", date)
		},
	}
}
