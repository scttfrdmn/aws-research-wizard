package tui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aws-research-wizard/go/internal/aws"
	"github.com/aws-research-wizard/go/internal/config"
)

// CostCalculatorModel represents the cost calculation interface
type CostCalculatorModel struct {
	table        table.Model
	domain       *config.DomainPack
	calculator   *aws.PricingCalculator
	estimates    map[string]*aws.CostEstimate
	selectedInstance string
	quitting     bool
}

// NewCostCalculator creates a new cost calculator
func NewCostCalculator(domain *config.DomainPack, region string) (*CostCalculatorModel, error) {
	calculator, err := aws.NewPricingCalculator(region)
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing calculator: %w", err)
	}
	
	// Create table columns
	columns := []table.Column{
		{Title: "Instance Type", Width: 15},
		{Title: "vCPUs", Width: 6},
		{Title: "Memory", Width: 10},
		{Title: "Hourly", Width: 8},
		{Title: "Monthly", Width: 10},
		{Title: "Annual", Width: 12},
		{Title: "Spot Savings", Width: 12},
	}
	
	// Calculate cost estimates for all recommended instances
	estimates := make(map[string]*aws.CostEstimate)
	var rows []table.Row
	
	for _, rec := range domain.AWSInstanceRecommendations {
		estimate, err := calculator.CalculateCost(rec.InstanceType)
		if err != nil {
			// Skip instances we can't calculate costs for
			continue
		}
		
		estimates[rec.InstanceType] = estimate
		
		rows = append(rows, table.Row{
			rec.InstanceType,
			fmt.Sprintf("%d", rec.VCPUs),
			fmt.Sprintf("%d GB", rec.MemoryGB),
			fmt.Sprintf("$%.3f", rec.CostPerHour),
			fmt.Sprintf("$%.0f", rec.CostPerHour*24*30.44),
			fmt.Sprintf("$%.0f", rec.CostPerHour*24*365),
			fmt.Sprintf("$%.0f (70%%)", rec.CostPerHour*24*30.44*0.7),
		})
	}
	
	// Sort rows by monthly cost
	sort.Slice(rows, func(i, j int) bool {
		return estimates[rows[i][0]].MonthlyCost < estimates[rows[j][0]].MonthlyCost
	})
	
	// Create and configure table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	
	// Style the table
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
	t.SetStyles(s)
	
	return &CostCalculatorModel{
		table:      t,
		domain:     domain,
		calculator: calculator,
		estimates:  estimates,
	}, nil
}

// Init initializes the model
func (m *CostCalculatorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *CostCalculatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// Get selected instance
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) > 0 {
				m.selectedInstance = selectedRow[0]
				return m, tea.Quit
			}
		}
	}
	
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the model
func (m *CostCalculatorModel) View() string {
	if m.quitting {
		return ""
	}
	
	title := titleStyle.Render(fmt.Sprintf("💰 Cost Calculator - %s", m.domain.Name))
	
	// Domain info section
	domainInfo := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Render(fmt.Sprintf(
			"Domain: %s\nDescription: %s\nTarget Users: %s",
			m.domain.Name,
			m.domain.Description,
			m.domain.TargetUsers,
		))
	
	// Cost optimization tips
	tips := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")).
		Padding(1).
		Render("💡 Optimization Tips:\n" +
			"• Use Spot instances for 70% savings\n" +
			"• Reserved instances save 30-60%\n" +
			"• Consider S3 Intelligent Tiering\n" +
			"• Enable detailed monitoring")
	
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: navigate • enter: select instance • q: quit")
	
	// Get current selection info
	selectedRow := m.table.SelectedRow()
	var selectionInfo string
	if len(selectedRow) > 0 {
		instanceType := selectedRow[0]
		if estimate, exists := m.estimates[instanceType]; exists {
			selectionInfo = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("84")).
				Padding(1).
				Render(fmt.Sprintf(
					"Selected: %s\n"+
						"Specs: %d vCPUs, %s RAM\n"+
						"Cost: $%.3f/hour, $%.0f/month\n"+
						"Spot Savings: $%.0f/month (70%%)\n"+
						"Reserved Savings: $%.0f/month (40%%)",
					instanceType,
					estimate.VCPUs,
					estimate.Memory,
					estimate.HourlyCost,
					estimate.MonthlyCost,
					estimate.SpotSavings*24*30.44,
					estimate.ReservedSavings*24*30.44,
				))
		}
	}
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			domainInfo,
			"  ",
			tips,
		),
		"",
		baseStyle.Render(m.table.View()),
		"",
		selectionInfo,
		"",
		help,
	)
	
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// GetSelectedInstance returns the selected instance type
func (m *CostCalculatorModel) GetSelectedInstance() string {
	return m.selectedInstance
}

// GetEstimate returns the cost estimate for an instance type
func (m *CostCalculatorModel) GetEstimate(instanceType string) *aws.CostEstimate {
	return m.estimates[instanceType]
}

// RunCostCalculator runs the cost calculation TUI
func RunCostCalculator(domain *config.DomainPack, region string) (string, *aws.CostEstimate, error) {
	model, err := NewCostCalculator(domain, region)
	if err != nil {
		return "", nil, err
	}
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", nil, fmt.Errorf("failed to run cost calculator: %w", err)
	}
	
	if calcModel, ok := finalModel.(*CostCalculatorModel); ok {
		selectedInstance := calcModel.GetSelectedInstance()
		estimate := calcModel.GetEstimate(selectedInstance)
		return selectedInstance, estimate, nil
	}
	
	return "", nil, fmt.Errorf("unexpected model type")
}