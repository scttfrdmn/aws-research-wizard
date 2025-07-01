package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aws-research-wizard/go/internal/aws"
)

// MonitoringConfig contains configuration for the monitoring dashboard
type MonitoringConfig struct {
	Region       string
	RefreshRate  time.Duration
	StackName    string
	InstanceID   string
	ShowCosts    bool
	ShowAlerts   bool
	AutoRefresh  bool
	OutputFormat string
}

// MonitoringDashboard provides real-time AWS resource monitoring
type MonitoringDashboard struct {
	client            *aws.Client
	monitoringManager *aws.MonitoringManager
	infraManager      *aws.InfrastructureManager
	config            *MonitoringConfig
}

// NewMonitoringDashboard creates a new monitoring dashboard
func NewMonitoringDashboard(client *aws.Client, monitoring *aws.MonitoringManager, infra *aws.InfrastructureManager, config *MonitoringConfig) *MonitoringDashboard {
	return &MonitoringDashboard{
		client:            client,
		monitoringManager: monitoring,
		infraManager:      infra,
		config:            config,
	}
}

// DashboardModel represents the monitoring dashboard state
type DashboardModel struct {
	dashboard     *MonitoringDashboard
	instances     []aws.InstanceInfo
	costs         []aws.CostData
	alarms        []aws.AlarmInfo
	metrics       map[string]*aws.InstanceMetrics
	lastUpdate    time.Time
	refreshTicker *time.Ticker
	selectedTab   int
	instanceTable table.Model
	costTable     table.Model
	alarmTable    table.Model
	quitting      bool
	loading       bool
	error         string
}

// Tab definitions
const (
	TabInstances = iota
	TabCosts
	TabAlarms
	TabMetrics
)

var tabNames = []string{"Instances", "Costs", "Alarms", "Metrics"}

// Run starts the monitoring dashboard
func (md *MonitoringDashboard) Run(ctx context.Context) error {
	// Initialize the model
	model := &DashboardModel{
		dashboard:     md,
		selectedTab:   TabInstances,
		metrics:       make(map[string]*aws.InstanceMetrics),
		refreshTicker: time.NewTicker(md.config.RefreshRate),
	}

	// Initialize tables
	model.initializeTables()

	// Load initial data
	model.refreshData(ctx)

	// Start the TUI
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Handle refresh ticker in a goroutine
	go func() {
		for range model.refreshTicker.C {
			if model.quitting {
				return
			}
			p.Send(refreshMsg{})
		}
	}()

	_, err := p.Run()
	model.refreshTicker.Stop()
	return err
}

// Initialize table models
func (dm *DashboardModel) initializeTables() {
	// Instance table
	instanceColumns := []table.Column{
		{Title: "Instance ID", Width: 19},
		{Title: "Type", Width: 12},
		{Title: "State", Width: 10},
		{Title: "Public IP", Width: 15},
		{Title: "AZ", Width: 12},
		{Title: "Launch Time", Width: 16},
	}
	dm.instanceTable = table.New(
		table.WithColumns(instanceColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Cost table
	costColumns := []table.Column{
		{Title: "Service", Width: 20},
		{Title: "Amount", Width: 12},
		{Title: "Currency", Width: 8},
		{Title: "Period", Width: 20},
	}
	dm.costTable = table.New(
		table.WithColumns(costColumns),
		table.WithHeight(10),
	)

	// Alarm table
	alarmColumns := []table.Column{
		{Title: "Alarm Name", Width: 25},
		{Title: "State", Width: 12},
		{Title: "Metric", Width: 20},
		{Title: "Threshold", Width: 15},
	}
	dm.alarmTable = table.New(
		table.WithColumns(alarmColumns),
		table.WithHeight(10),
	)

	// Style tables
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	dm.instanceTable.SetStyles(s)
	dm.costTable.SetStyles(s)
	dm.alarmTable.SetStyles(s)
}

// Refresh data from AWS
func (dm *DashboardModel) refreshData(ctx context.Context) {
	dm.loading = true
	dm.error = ""

	// Refresh instances
	if err := dm.refreshInstances(ctx); err != nil {
		dm.error = fmt.Sprintf("Failed to refresh instances: %v", err)
	}

	// Refresh costs if enabled
	if dm.dashboard.config.ShowCosts {
		if err := dm.refreshCosts(ctx); err != nil {
			dm.error = fmt.Sprintf("Failed to refresh costs: %v", err)
		}
	}

	// Refresh alarms if enabled
	if dm.dashboard.config.ShowAlerts {
		if err := dm.refreshAlarms(ctx); err != nil {
			dm.error = fmt.Sprintf("Failed to refresh alarms: %v", err)
		}
	}

	// Refresh metrics for running instances
	dm.refreshMetrics(ctx)

	dm.lastUpdate = time.Now()
	dm.loading = false
}

func (dm *DashboardModel) refreshInstances(ctx context.Context) error {
	filters := make(map[string][]string)
	if dm.dashboard.config.InstanceID != "" {
		filters["instance-id"] = []string{dm.dashboard.config.InstanceID}
	}

	instances, err := dm.dashboard.infraManager.ListInstances(ctx, filters)
	if err != nil {
		return err
	}

	dm.instances = instances

	// Update instance table
	var rows []table.Row
	for _, instance := range instances {
		state := instance.State
		if state == "running" {
			state = "🟢 " + state
		} else if state == "stopped" {
			state = "🔴 " + state
		} else {
			state = "🟡 " + state
		}

		rows = append(rows, table.Row{
			instance.InstanceID,
			instance.InstanceType,
			state,
			instance.PublicIP,
			instance.AvailabilityZone,
			instance.LaunchTime.Format("15:04:05"),
		})
	}

	dm.instanceTable.SetRows(rows)
	return nil
}

func (dm *DashboardModel) refreshCosts(ctx context.Context) error {
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // Last 7 days

	costs, err := dm.dashboard.monitoringManager.GetCostData(ctx, startDate, endDate, []string{"SERVICE"})
	if err != nil {
		return err
	}

	dm.costs = costs

	// Update cost table
	var rows []table.Row
	serviceCosts := make(map[string]float64)

	for _, cost := range costs {
		serviceCosts[cost.Service] += cost.Amount
	}

	// Sort by cost amount
	type serviceCost struct {
		service string
		amount  float64
	}
	var sortedCosts []serviceCost
	for service, amount := range serviceCosts {
		sortedCosts = append(sortedCosts, serviceCost{service, amount})
	}
	sort.Slice(sortedCosts, func(i, j int) bool {
		return sortedCosts[i].amount > sortedCosts[j].amount
	})

	for _, sc := range sortedCosts {
		rows = append(rows, table.Row{
			sc.service,
			fmt.Sprintf("$%.2f", sc.amount),
			"USD",
			fmt.Sprintf("%s to %s", startDate, endDate),
		})
	}

	dm.costTable.SetRows(rows)
	return nil
}

func (dm *DashboardModel) refreshAlarms(ctx context.Context) error {
	alarms, err := dm.dashboard.monitoringManager.ListAlarms(ctx, "")
	if err != nil {
		return err
	}

	dm.alarms = alarms

	// Update alarm table
	var rows []table.Row
	for _, alarm := range alarms {
		state := alarm.State
		if state == "OK" {
			state = "🟢 " + state
		} else if state == "ALARM" {
			state = "🔴 " + state
		} else {
			state = "🟡 " + state
		}

		threshold := fmt.Sprintf("%s %.1f", alarm.ComparisonOp, alarm.Threshold)
		metric := fmt.Sprintf("%s/%s", alarm.Namespace, alarm.MetricName)

		rows = append(rows, table.Row{
			alarm.AlarmName,
			state,
			metric,
			threshold,
		})
	}

	dm.alarmTable.SetRows(rows)
	return nil
}

func (dm *DashboardModel) refreshMetrics(ctx context.Context) {
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Hour) // Last hour

	for _, instance := range dm.instances {
		if instance.State == "running" {
			metrics, err := dm.dashboard.monitoringManager.GetInstanceMetrics(
				ctx, instance.InstanceID, startTime, endTime,
			)
			if err == nil {
				dm.metrics[instance.InstanceID] = metrics
			}
		}
	}
}

// Message types
type refreshMsg struct{}

// Init initializes the model
func (dm *DashboardModel) Init() tea.Cmd {
	return tea.Tick(dm.dashboard.config.RefreshRate, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

// Update handles messages and updates the model
func (dm *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			dm.quitting = true
			return dm, tea.Quit
		case "r":
			// Manual refresh
			ctx := context.Background()
			dm.refreshData(ctx)
		case "tab":
			dm.selectedTab = (dm.selectedTab + 1) % len(tabNames)
		case "shift+tab":
			dm.selectedTab = (dm.selectedTab - 1 + len(tabNames)) % len(tabNames)
		}

	case refreshMsg:
		if dm.dashboard.config.AutoRefresh && !dm.quitting {
			ctx := context.Background()
			dm.refreshData(ctx)
			return dm, tea.Tick(dm.dashboard.config.RefreshRate, func(t time.Time) tea.Msg {
				return refreshMsg{}
			})
		}
	}

	// Update active table
	switch dm.selectedTab {
	case TabInstances:
		dm.instanceTable, cmd = dm.instanceTable.Update(msg)
	case TabCosts:
		dm.costTable, cmd = dm.costTable.Update(msg)
	case TabAlarms:
		dm.alarmTable, cmd = dm.alarmTable.Update(msg)
	}

	return dm, cmd
}

// View renders the model
func (dm *DashboardModel) View() string {
	if dm.quitting {
		return ""
	}

	// Header
	header := titleStyle.Render("📊 AWS Research Wizard - Real-time Monitoring")
	status := fmt.Sprintf("Region: %s | Last Update: %s",
		dm.dashboard.config.Region,
		dm.lastUpdate.Format("15:04:05"))

	// Loading indicator
	if dm.loading {
		status += " | 🔄 Loading..."
	}

	// Error display
	errorMsg := ""
	if dm.error != "" {
		errorMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(fmt.Sprintf("Error: %s", dm.error))
	}

	// Tab navigation
	var tabs []string
	for i, name := range tabNames {
		if i == dm.selectedTab {
			tabs = append(tabs, selectedStyle.Render(name))
		} else {
			tabs = append(tabs, name)
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Content based on selected tab
	var content string
	switch dm.selectedTab {
	case TabInstances:
		content = dm.renderInstancesTab()
	case TabCosts:
		content = dm.renderCostsTab()
	case TabAlarms:
		content = dm.renderAlarmsTab()
	case TabMetrics:
		content = dm.renderMetricsTab()
	}

	// Help text
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("tab: switch tabs • r: refresh • q: quit")

	// Combine all sections
	sections := []string{header, status, tabBar, content}
	if errorMsg != "" {
		sections = append(sections, errorMsg)
	}
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (dm *DashboardModel) renderInstancesTab() string {
	summary := fmt.Sprintf("EC2 Instances (%d total)", len(dm.instances))

	// Instance state summary
	states := make(map[string]int)
	for _, instance := range dm.instances {
		states[instance.State]++
	}

	var stateSummary []string
	for state, count := range states {
		stateSummary = append(stateSummary, fmt.Sprintf("%s: %d", state, count))
	}

	summary += " | " + strings.Join(stateSummary, ", ")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		summary,
		"",
		baseStyle.Render(dm.instanceTable.View()),
	)
}

func (dm *DashboardModel) renderCostsTab() string {
	if !dm.dashboard.config.ShowCosts {
		return "Cost tracking disabled. Enable with --costs flag."
	}

	// Calculate total cost
	serviceCosts := make(map[string]float64)
	totalCost := 0.0
	for _, cost := range dm.costs {
		serviceCosts[cost.Service] += cost.Amount
		totalCost += cost.Amount
	}

	summary := fmt.Sprintf("Cost Summary - Last 7 Days | Total: $%.2f | Daily Avg: $%.2f",
		totalCost, totalCost/7)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		summary,
		"",
		baseStyle.Render(dm.costTable.View()),
	)
}

func (dm *DashboardModel) renderAlarmsTab() string {
	if !dm.dashboard.config.ShowAlerts {
		return "Alert monitoring disabled. Enable with --alerts flag."
	}

	// Alarm state summary
	states := make(map[string]int)
	for _, alarm := range dm.alarms {
		states[alarm.State]++
	}

	var stateSummary []string
	for state, count := range states {
		stateSummary = append(stateSummary, fmt.Sprintf("%s: %d", state, count))
	}

	summary := fmt.Sprintf("CloudWatch Alarms (%d total)", len(dm.alarms))
	if len(stateSummary) > 0 {
		summary += " | " + strings.Join(stateSummary, ", ")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		summary,
		"",
		baseStyle.Render(dm.alarmTable.View()),
	)
}

func (dm *DashboardModel) renderMetricsTab() string {
	if len(dm.metrics) == 0 {
		return "No metrics available. Metrics are collected for running instances."
	}

	var content []string
	content = append(content, fmt.Sprintf("Instance Metrics (%d instances)", len(dm.metrics)))
	content = append(content, "")

	for instanceID, metrics := range dm.metrics {
		content = append(content, fmt.Sprintf("Instance: %s", instanceID))

		// Latest CPU utilization
		if len(metrics.CPUUtilization) > 0 {
			latest := metrics.CPUUtilization[len(metrics.CPUUtilization)-1]
			content = append(content, fmt.Sprintf("  CPU: %.1f%% at %s",
				latest.Value, latest.Timestamp.Format("15:04:05")))
		}

		// Latest network metrics
		if len(metrics.NetworkIn) > 0 && len(metrics.NetworkOut) > 0 {
			netIn := metrics.NetworkIn[len(metrics.NetworkIn)-1]
			netOut := metrics.NetworkOut[len(metrics.NetworkOut)-1]
			content = append(content, fmt.Sprintf("  Network: In %.1f MB, Out %.1f MB",
				netIn.Value/1024/1024, netOut.Value/1024/1024))
		}

		content = append(content, "")
	}

	return strings.Join(content, "\n")
}

// RunCostAnalysis provides a simplified cost analysis view
func RunCostAnalysis(client *aws.Client, region string, days int) error {
	ctx := context.Background()
	monitoringManager := aws.NewMonitoringManager(client)

	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	costs, err := monitoringManager.GetCostData(ctx, startDate, endDate, []string{"SERVICE"})
	if err != nil {
		return fmt.Errorf("failed to get cost data: %w", err)
	}

	// Display cost analysis
	fmt.Printf("💰 Cost Analysis - Last %d Days\n", days)
	fmt.Printf("Period: %s to %s\n\n", startDate, endDate)

	if len(costs) == 0 {
		fmt.Println("No cost data available for the specified period.")
		return nil
	}

	serviceCosts := make(map[string]float64)
	totalCost := 0.0

	for _, cost := range costs {
		serviceCosts[cost.Service] += cost.Amount
		totalCost += cost.Amount
	}

	// Sort services by cost
	type serviceCost struct {
		service string
		amount  float64
	}
	var sortedCosts []serviceCost
	for service, amount := range serviceCosts {
		sortedCosts = append(sortedCosts, serviceCost{service, amount})
	}
	sort.Slice(sortedCosts, func(i, j int) bool {
		return sortedCosts[i].amount > sortedCosts[j].amount
	})

	fmt.Printf("Service Breakdown:\n")
	for _, sc := range sortedCosts {
		percentage := (sc.amount / totalCost) * 100
		fmt.Printf("  %-20s $%8.2f (%5.1f%%)\n", sc.service, sc.amount, percentage)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total Cost: $%.2f\n", totalCost)
	fmt.Printf("  Daily Average: $%.2f\n", totalCost/float64(days))
	fmt.Printf("  Monthly Projection: $%.2f\n", totalCost*30/float64(days))

	return nil
}
