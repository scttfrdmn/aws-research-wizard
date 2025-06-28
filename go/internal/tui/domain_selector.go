package tui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aws-research-wizard/go/internal/config"
)

var (
	// Styling
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1)
	
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(true)
)

// DomainSelectorModel represents the domain selection interface
type DomainSelectorModel struct {
	table    table.Model
	domains  map[string]*config.DomainPack
	selected *config.DomainPack
	quitting bool
}

// NewDomainSelector creates a new domain selector
func NewDomainSelector(domains map[string]*config.DomainPack) *DomainSelectorModel {
	// Create table columns
	columns := []table.Column{
		{Title: "Domain", Width: 25},
		{Title: "Description", Width: 50},
		{Title: "Users", Width: 15},
		{Title: "Monthly Cost", Width: 12},
	}
	
	// Create table rows from domains
	var rows []table.Row
	var domainNames []string
	
	// Sort domains by name for consistent display
	for name := range domains {
		domainNames = append(domainNames, name)
	}
	sort.Strings(domainNames)
	
	for _, name := range domainNames {
		domain := domains[name]
		
		// Format users - already a string in YAML
		users := domain.TargetUsers
		if len(users) > 12 {
			users = users[:12] + "..."
		}
		
		// Format cost
		cost := fmt.Sprintf("$%.0f", domain.EstimatedCost.Total)
		
		// Truncate description if too long
		description := domain.Description
		if len(description) > 47 {
			description = description[:47] + "..."
		}
		
		rows = append(rows, table.Row{
			name,
			description,
			users,
			cost,
		})
	}
	
	// Create and configure table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
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
	
	return &DomainSelectorModel{
		table:   t,
		domains: domains,
	}
}

// Init initializes the model
func (m *DomainSelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *DomainSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// Get selected domain
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) > 0 {
				domainName := selectedRow[0]
				m.selected = m.domains[domainName]
				return m, tea.Quit
			}
		}
	}
	
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the model
func (m *DomainSelectorModel) View() string {
	if m.quitting {
		return ""
	}
	
	title := titleStyle.Render("🔬 AWS Research Wizard - Domain Selection")
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: navigate • enter: select • q: quit")
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		baseStyle.Render(m.table.View()),
		"",
		help,
	)
	
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// GetSelected returns the selected domain
func (m *DomainSelectorModel) GetSelected() *config.DomainPack {
	return m.selected
}

// RunDomainSelector runs the domain selection TUI
func RunDomainSelector(domains map[string]*config.DomainPack) (*config.DomainPack, error) {
	model := NewDomainSelector(domains)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run domain selector: %w", err)
	}
	
	if selectorModel, ok := finalModel.(*DomainSelectorModel); ok {
		return selectorModel.GetSelected(), nil
	}
	
	return nil, fmt.Errorf("unexpected model type")
}