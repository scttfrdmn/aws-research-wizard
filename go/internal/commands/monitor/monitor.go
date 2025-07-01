package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/aws-research-wizard/go/internal/aws"
	"github.com/aws-research-wizard/go/internal/tui"
)

// NewMonitorCommand creates the monitor subcommand
func NewMonitorCommand() *cobra.Command {
	var refreshRate int
	var stackName string
	var instanceID string
	var showCosts bool
	var showAlerts bool
	var autoRefresh bool
	var outputFormat string

	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Real-time resource monitoring and dashboards",
		Long: `Real-time monitoring dashboard for AWS research environments with:

- Live EC2 instance metrics (CPU, memory, network, disk I/O)
- Real-time cost tracking with daily/monthly projections
- CloudWatch alarms and alert management
- Performance optimization recommendations
- Interactive terminal dashboards
- Cost breakdown by service and project`,
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveMonitor(cmd, refreshRate, stackName, instanceID, showCosts, showAlerts, autoRefresh, outputFormat)
		},
	}

	// Add flags
	monitorCmd.PersistentFlags().IntVar(&refreshRate, "refresh", 30, "Refresh interval in seconds")
	monitorCmd.PersistentFlags().StringVar(&stackName, "stack", "", "CloudFormation stack name to monitor")
	monitorCmd.PersistentFlags().StringVar(&instanceID, "instance", "", "EC2 instance ID to monitor")
	monitorCmd.PersistentFlags().BoolVar(&showCosts, "costs", true, "Show cost tracking")
	monitorCmd.PersistentFlags().BoolVar(&showAlerts, "alerts", true, "Show alert status")
	monitorCmd.PersistentFlags().BoolVar(&autoRefresh, "auto-refresh", true, "Enable auto-refresh")
	monitorCmd.PersistentFlags().StringVar(&outputFormat, "format", "dashboard", "Output format: dashboard, json, table")

	// Add subcommands
	monitorCmd.AddCommand(
		createDashboardCommand(&refreshRate, &stackName, &instanceID, &showCosts, &showAlerts, &autoRefresh, &outputFormat),
		createCostCommand(),
		createAlertsCommand(),
		createInstancesCommand(&instanceID),
		createStacksCommand(&stackName),
	)

	return monitorCmd
}

func runInteractiveMonitor(cmd *cobra.Command, refreshRate int, stackName, instanceID string, showCosts, showAlerts, autoRefresh bool, outputFormat string) {
	ctx := context.Background()

	region, _ := cmd.Flags().GetString("region")

	fmt.Printf("📊 AWS Research Wizard - Real-time Monitoring\n")
	fmt.Printf("Region: %s | Refresh: %ds | Auto-refresh: %v\n\n", region, refreshRate, autoRefresh)

	// Initialize AWS client
	awsClient, err := aws.NewClient(ctx, region)
	if err != nil {
		log.Fatalf("Failed to initialize AWS client: %v", err)
	}

	// Validate AWS credentials
	if err := awsClient.ValidateCredentials(ctx); err != nil {
		log.Fatalf("AWS credentials validation failed: %v", err)
	}

	// Create monitoring manager
	monitoringManager := aws.NewMonitoringManager(awsClient)
	infraManager := aws.NewInfrastructureManager(awsClient)

	// Configuration for the monitoring dashboard
	config := &tui.MonitoringConfig{
		Region:       region,
		RefreshRate:  time.Duration(refreshRate) * time.Second,
		StackName:    stackName,
		InstanceID:   instanceID,
		ShowCosts:    showCosts,
		ShowAlerts:   showAlerts,
		AutoRefresh:  autoRefresh,
		OutputFormat: outputFormat,
	}

	// Start the monitoring dashboard
	dashboard := tui.NewMonitoringDashboard(awsClient, monitoringManager, infraManager, config)

	if err := dashboard.Run(ctx); err != nil {
		log.Fatalf("Monitoring dashboard failed: %v", err)
	}
}

func createDashboardCommand(refreshRate *int, stackName, instanceID *string, showCosts, showAlerts, autoRefresh *bool, outputFormat *string) *cobra.Command {
	return &cobra.Command{
		Use:   "dashboard",
		Short: "Launch interactive monitoring dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveMonitor(cmd, *refreshRate, *stackName, *instanceID, *showCosts, *showAlerts, *autoRefresh, *outputFormat)
		},
	}
}

func createCostCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "costs",
		Short: "Show cost analysis and tracking",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			region, _ := cmd.Flags().GetString("region")
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			monitoringManager := aws.NewMonitoringManager(awsClient)

			// Get cost data for the last 30 days
			endDate := time.Now().Format("2006-01-02")
			startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

			costs, err := monitoringManager.GetCostData(ctx, startDate, endDate, []string{"SERVICE"})
			if err != nil {
				log.Fatalf("Failed to get cost data: %v", err)
			}

			fmt.Printf("💰 Cost Analysis - Last 30 Days\n")
			fmt.Printf("Period: %s to %s\n\n", startDate, endDate)

			if len(costs) == 0 {
				fmt.Println("No cost data available for the specified period.")
				return
			}

			// Group costs by service
			serviceCosts := make(map[string]float64)
			totalCost := 0.0

			for _, cost := range costs {
				serviceCosts[cost.Service] += cost.Amount
				totalCost += cost.Amount
			}

			fmt.Printf("Service Breakdown:\n")
			for service, amount := range serviceCosts {
				percentage := (amount / totalCost) * 100
				fmt.Printf("  %-20s $%8.2f (%5.1f%%)\n", service, amount, percentage)
			}

			fmt.Printf("\nTotal Cost: $%.2f\n", totalCost)
			fmt.Printf("Daily Average: $%.2f\n", totalCost/30)
			fmt.Printf("Monthly Projection: $%.2f\n", totalCost)
		},
	}
}

func createAlertsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "alerts",
		Short: "Show CloudWatch alerts and alarms",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			region, _ := cmd.Flags().GetString("region")
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			monitoringManager := aws.NewMonitoringManager(awsClient)

			alarms, err := monitoringManager.ListAlarms(ctx, "")
			if err != nil {
				log.Fatalf("Failed to list alarms: %v", err)
			}

			fmt.Printf("🚨 CloudWatch Alarms (%d total)\n\n", len(alarms))

			if len(alarms) == 0 {
				fmt.Println("No CloudWatch alarms found.")
				return
			}

			for _, alarm := range alarms {
				status := "🟢"
				if alarm.State == "ALARM" {
					status = "🔴"
				} else if alarm.State == "INSUFFICIENT_DATA" {
					status = "🟡"
				}

				fmt.Printf("%s %s\n", status, alarm.AlarmName)
				fmt.Printf("   Metric: %s/%s\n", alarm.Namespace, alarm.MetricName)
				fmt.Printf("   Threshold: %s %.2f\n", alarm.ComparisonOp, alarm.Threshold)
				fmt.Printf("   State: %s\n", alarm.State)
				if alarm.StateReason != "" {
					fmt.Printf("   Reason: %s\n", alarm.StateReason)
				}
				fmt.Printf("\n")
			}
		},
	}
}

func createInstancesCommand(instanceID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "instances",
		Short: "Show EC2 instance status and metrics",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			region, _ := cmd.Flags().GetString("region")
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			infraManager := aws.NewInfrastructureManager(awsClient)

			// List instances
			filters := make(map[string][]string)
			if *instanceID != "" {
				filters["instance-id"] = []string{*instanceID}
			}

			instances, err := infraManager.ListInstances(ctx, filters)
			if err != nil {
				log.Fatalf("Failed to list instances: %v", err)
			}

			fmt.Printf("🖥️  EC2 Instances (%d total)\n\n", len(instances))

			if len(instances) == 0 {
				fmt.Println("No EC2 instances found.")
				return
			}

			for _, instance := range instances {
				status := "🟢"
				if instance.State == "stopped" {
					status = "🔴"
				} else if instance.State == "pending" || instance.State == "stopping" {
					status = "🟡"
				}

				fmt.Printf("%s %s (%s)\n", status, instance.InstanceID, instance.InstanceType)
				fmt.Printf("   State: %s\n", instance.State)
				fmt.Printf("   Public IP: %s\n", instance.PublicIP)
				fmt.Printf("   Private IP: %s\n", instance.PrivateIP)
				fmt.Printf("   AZ: %s\n", instance.AvailabilityZone)
				fmt.Printf("   Launch Time: %s\n", instance.LaunchTime.Format(time.RFC3339))

				// Show tags
				if len(instance.Tags) > 0 {
					fmt.Printf("   Tags: ")
					first := true
					for key, value := range instance.Tags {
						if !first {
							fmt.Printf(", ")
						}
						fmt.Printf("%s=%s", key, value)
						first = false
					}
					fmt.Printf("\n")
				}
				fmt.Printf("\n")
			}
		},
	}
}

func createStacksCommand(stackName *string) *cobra.Command {
	return &cobra.Command{
		Use:   "stacks",
		Short: "Show CloudFormation stack status",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			region, _ := cmd.Flags().GetString("region")
			awsClient, err := aws.NewClient(ctx, region)
			if err != nil {
				log.Fatalf("Failed to initialize AWS client: %v", err)
			}

			infraManager := aws.NewInfrastructureManager(awsClient)

			if *stackName != "" {
				// Show specific stack
				stackInfo, err := infraManager.GetStackInfo(ctx, *stackName)
				if err != nil {
					log.Fatalf("Failed to get stack info: %v", err)
				}

				fmt.Printf("📚 Stack: %s\n\n", stackInfo.StackName)
				fmt.Printf("Status: %s\n", stackInfo.Status)
				fmt.Printf("Created: %s\n", stackInfo.CreatedTime.Format(time.RFC3339))

				if stackInfo.UpdatedTime != nil {
					fmt.Printf("Updated: %s\n", stackInfo.UpdatedTime.Format(time.RFC3339))
				}

				if len(stackInfo.Outputs) > 0 {
					fmt.Printf("\nOutputs:\n")
					for key, value := range stackInfo.Outputs {
						fmt.Printf("  %s: %s\n", key, value)
					}
				}

				if len(stackInfo.Parameters) > 0 {
					fmt.Printf("\nParameters:\n")
					for key, value := range stackInfo.Parameters {
						fmt.Printf("  %s: %s\n", key, value)
					}
				}
			} else {
				fmt.Printf("📚 CloudFormation Stacks\n\n")
				fmt.Printf("Use --stack flag to view details of a specific stack.\n")
				fmt.Printf("Example: aws-research-wizard monitor stacks --stack my-research-stack\n")
			}
		},
	}
}
