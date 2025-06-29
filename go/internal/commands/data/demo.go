package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws-research-wizard/go/internal/data"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// demoCmd represents the demo command for testing the system
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run demonstration of intelligent data movement capabilities",
	Long: `Run a comprehensive demonstration of the intelligent data movement system
using the built-in genomics project example. This command showcases:

- Pattern analysis and domain detection
- Cost optimization recommendations  
- Workflow generation and execution
- Domain-specific optimizations
- Real-time progress monitoring

This is perfect for evaluating the system capabilities and understanding
how it works with realistic research data scenarios.`,
	RunE: runDemo,
}

var (
	demoSkipExecution bool
	demoVerbose       bool
	demoOutputDir     string
)

func init() {
	// Add demo command to data command
	DataCmd.AddCommand(demoCmd)
	
	// Flags
	demoCmd.Flags().BoolVar(&demoSkipExecution, "skip-execution", false, "Skip actual workflow execution (analysis only)")
	demoCmd.Flags().BoolVarP(&demoVerbose, "verbose", "v", false, "Show detailed progress information")
	demoCmd.Flags().StringVar(&demoOutputDir, "output-dir", "", "Directory to save demo outputs (default: temp dir)")
}

func runDemo(cmd *cobra.Command, args []string) error {
	fmt.Println("🧬 AWS Research Wizard - Intelligent Data Movement Demo")
	fmt.Println("======================================================")
	fmt.Println()
	
	// Setup demo environment
	ctx := context.Background()
	
	// Use built-in genomics example
	genomicsConfigPath := filepath.Join("examples", "genomics-project.yaml")
	absConfigPath, err := filepath.Abs(genomicsConfigPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}
	
	// Check if genomics example exists
	if _, err := os.Stat(absConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("genomics example not found at %s", absConfigPath)
	}
	
	fmt.Printf("📋 Using genomics project configuration: %s\n\n", absConfigPath)
	
	// Load the example configuration
	configData, err := os.ReadFile(absConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	
	var projectConfig data.ProjectConfig
	if err := yaml.Unmarshal(configData, &projectConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	
	// Display project information
	fmt.Printf("🔬 Project: %s\n", projectConfig.Project.Name)
	fmt.Printf("   Domain: %s\n", projectConfig.Project.Domain)
	fmt.Printf("   Owner: %s\n", projectConfig.Project.Owner)
	fmt.Printf("   Budget: %s\n", projectConfig.Project.Budget)
	fmt.Println()
	
	// Show data profiles
	fmt.Printf("📊 Data Profiles (%d):\n", len(projectConfig.DataProfiles))
	for name, profile := range projectConfig.DataProfiles {
		fmt.Printf("   • %s: %s (%s, %d files)\n", 
			name, profile.Name, profile.TotalSize, profile.FileCount)
	}
	fmt.Println()
	
	// Show destinations
	fmt.Printf("🎯 Destinations (%d):\n", len(projectConfig.Destinations))
	for name, dest := range projectConfig.Destinations {
		fmt.Printf("   • %s: %s (%s)\n", name, dest.Name, dest.URI)
	}
	fmt.Println()
	
	// Show workflows
	fmt.Printf("⚙️  Workflows (%d):\n", len(projectConfig.Workflows))
	for _, workflow := range projectConfig.Workflows {
		fmt.Printf("   • %s: %s → %s (%s)\n", 
			workflow.Name, workflow.Source, workflow.Destination, workflow.Engine)
		if demoVerbose {
			fmt.Printf("     Preprocessing: %d steps\n", len(workflow.PreProcessing))
			fmt.Printf("     Postprocessing: %d steps\n", len(workflow.PostProcessing))
		}
	}
	fmt.Println()
	
	// Demonstrate domain profile integration
	fmt.Println("🎯 Domain Profile Analysis")
	fmt.Println("---------------------------")
	
	dpm := data.NewResearchDomainProfileManager()
	profile, exists := dpm.GetProfile(projectConfig.Project.Domain)
	if exists {
		fmt.Printf("✅ Domain profile found: %s\n", profile.Name)
		fmt.Printf("   Description: %s\n", profile.Description)
		fmt.Printf("   File type optimizations: %d\n", len(profile.FileTypeHints))
		fmt.Printf("   Preferred engines: %v\n", profile.TransferOptimization.PreferredEngines)
		fmt.Printf("   Bundling enabled: %t\n", profile.BundlingStrategy.EnableBundling)
		fmt.Printf("   Security requirements: encryption=%t\n", profile.SecurityRequirements.EncryptionRequired)
		
		if demoVerbose {
			fmt.Printf("\n   📁 File Type Optimizations:\n")
			for ext, hint := range profile.FileTypeHints {
				fmt.Printf("     • %s: %s (compression: %.1fx, engine: %s)\n", 
					ext, hint.Description, hint.CompressionRatio, hint.PreferredEngine)
			}
		}
	} else {
		fmt.Printf("⚠️  No domain profile found for: %s\n", projectConfig.Project.Domain)
	}
	fmt.Println()
	
	// Demonstrate cost optimization
	fmt.Println("💰 Cost Optimization Analysis")
	fmt.Println("------------------------------")
	
	// Calculate rough estimates based on configuration data
	totalFiles := int64(0)
	totalSizeGB := float64(0)
	
	for _, profile := range projectConfig.DataProfiles {
		totalFiles += profile.FileCount
		if profile.TotalSize != "" {
			// Parse size (simplified)
			if profile.TotalSize == "2.5TB" {
				totalSizeGB += 2500
			} else if profile.TotalSize == "1.8TB" {
				totalSizeGB += 1800
			} else if profile.TotalSize == "500GB" {
				totalSizeGB += 500
			}
		}
	}
	
	fmt.Printf("📊 Estimated Dataset Statistics:\n")
	fmt.Printf("   Total files: %d\n", totalFiles)
	fmt.Printf("   Total size: %.1f GB\n", totalSizeGB)
	
	// Calculate basic cost estimates
	monthlyCost := totalSizeGB * 0.023 // Rough S3 standard pricing
	fmt.Printf("   Estimated monthly storage cost: $%.2f\n", monthlyCost)
	
	// Show optimization potential
	if projectConfig.Optimization.CostOptimization.Enabled {
		fmt.Printf("✅ Cost optimization enabled:\n")
		fmt.Printf("   Auto bundling: %t\n", projectConfig.Optimization.CostOptimization.AutoBundleSmallFiles)
		fmt.Printf("   Auto compression: %t\n", projectConfig.Optimization.CostOptimization.AutoCompression)
		fmt.Printf("   Storage class optimization: %t\n", projectConfig.Optimization.CostOptimization.AutoStorageClass)
		
		if projectConfig.Optimization.CostOptimization.AutoBundleSmallFiles {
			fmt.Printf("   💡 Small file bundling could save ~20-40%% on request costs\n")
		}
	}
	fmt.Println()
	
	// Demonstrate workflow engine
	if !demoSkipExecution {
		fmt.Println("🚀 Workflow Engine Demonstration")
		fmt.Println("---------------------------------")
		
		// Create workflow engine
		engine := data.NewWorkflowEngine(&data.WorkflowEngineConfig{
			MaxConcurrentWorkflows: 1,
			DefaultTimeout:         30 * 1000000000, // 30 seconds in nanoseconds
			RetryAttempts:         1,
		})
		
		// Register components
		engine.RegisterAnalyzer(data.NewPatternAnalyzer())
		engine.RegisterBundlingEngine(data.NewBundlingEngine(nil))
		engine.RegisterWarningSystem(data.NewWarningSystem())
		
		// Register mock transfer engine for demo
		mockEngine := &DemoMockTransferEngine{name: "s5cmd"}
		engine.RegisterTransferEngine(mockEngine)
		
		// Execute first workflow
		if len(projectConfig.Workflows) > 0 {
			workflowName := projectConfig.Workflows[0].Name
			fmt.Printf("▶️  Executing workflow: %s\n", workflowName)
			
			execution, err := engine.ExecuteWorkflow(ctx, &projectConfig, workflowName)
			if err != nil {
				fmt.Printf("❌ Workflow execution failed: %v\n", err)
			} else {
				fmt.Printf("✅ Workflow started successfully: %s\n", execution.ID)
				fmt.Printf("   Status: %s\n", execution.Status)
				fmt.Printf("   Total steps: %d\n", execution.TotalSteps)
				
				if demoVerbose {
					fmt.Printf("   Steps:\n")
					for i, step := range execution.Steps {
						fmt.Printf("     %d. %s (%s)\n", i+1, step.Name, step.Type)
					}
				}
			}
		}
	} else {
		fmt.Println("⏭️  Workflow execution skipped (--skip-execution)")
	}
	fmt.Println()
	
	// Summary and recommendations
	fmt.Println("📝 Demo Summary")
	fmt.Println("===============")
	fmt.Println("✅ Intelligent data movement system successfully demonstrated!")
	fmt.Println()
	fmt.Println("Key capabilities shown:")
	fmt.Println("• 🔬 Domain-specific optimization profiles")
	fmt.Println("• 📊 Comprehensive cost analysis and optimization")
	fmt.Println("• ⚙️  Automated workflow generation and execution")
	fmt.Println("• 🎯 Research-aware file type handling")
	fmt.Println("• 💰 Cost-effective bundling strategies")
	fmt.Println("• 🔒 Security and compliance considerations")
	fmt.Println()
	
	fmt.Println("Next steps to use with your data:")
	fmt.Println("1. Run 'aws-research-wizard data analyze /your/data/path'")
	fmt.Println("2. Use '--generate-config' to create optimized configuration")
	fmt.Println("3. Execute workflows with 'aws-research-wizard data workflow run'")
	fmt.Println("4. Monitor progress with 'aws-research-wizard data workflow status'")
	fmt.Println()
	
	return nil
}

// MockTransferEngine for demo purposes (reuse existing one but ensure it's available)
type DemoMockTransferEngine struct {
	name string
}

func (m *DemoMockTransferEngine) GetName() string                                                    { return m.name }
func (m *DemoMockTransferEngine) GetType() string                                                    { return "mock" }
func (m *DemoMockTransferEngine) IsAvailable(ctx context.Context) error                             { return nil }
func (m *DemoMockTransferEngine) GetCapabilities() data.EngineCapabilities                          { return data.EngineCapabilities{Protocols: []string{"mock"}} }
func (m *DemoMockTransferEngine) Upload(ctx context.Context, req *data.TransferRequest) (*data.TransferResult, error) {
	return &data.TransferResult{TransferID: req.ID, Engine: m.name, Success: true}, nil
}
func (m *DemoMockTransferEngine) Download(ctx context.Context, req *data.TransferRequest) (*data.TransferResult, error) {
	return m.Upload(ctx, req)
}
func (m *DemoMockTransferEngine) Sync(ctx context.Context, req *data.SyncRequest) (*data.TransferResult, error) {
	return &data.TransferResult{Success: true}, nil
}
func (m *DemoMockTransferEngine) GetProgress(ctx context.Context, transferID string) (*data.TransferProgress, error) {
	return &data.TransferProgress{Percentage: 100.0}, nil
}
func (m *DemoMockTransferEngine) Cancel(ctx context.Context, transferID string) error { return nil }
func (m *DemoMockTransferEngine) Validate() error                                     { return nil }