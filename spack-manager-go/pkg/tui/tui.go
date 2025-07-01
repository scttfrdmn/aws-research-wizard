package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/spack-go/spack-manager/pkg/manager"
)

// TUI represents the main TUI application for Spack management
type TUI struct {
	spackManager    manager.Manager
	environments    []manager.Environment
	currentView     viewType
	selectedEnv     int
	installProgress map[string]float64
	logs            []string

	// Bubble Tea components
	list     list.Model
	table    table.Model
	progress progress.Model
	spinner  spinner.Model
	viewport viewport.Model

	// UI state
	width       int
	height      int
	installing  bool
	quitting    bool
	showingHelp bool
}

type viewType int

const (
	envListView viewType = iota
	envDetailView
	installView
	logsView
	helpView
)

// Colors and styles
var (
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#F25D94")
	successColor   = lipgloss.Color("#04B575")
	warningColor   = lipgloss.Color("#FF8700")
	errorColor     = lipgloss.Color("#FF5F87")
	mutedColor     = lipgloss.Color("#626262")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			BorderBottom(true).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(primaryColor).
			Bold(true)
)

// New creates a new TUI instance
func New(spackManager manager.Manager) *TUI {
	// Initialize list
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Spack Environments"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Initialize table
	columns := []table.Column{
		{Title: "Package", Width: 30},
		{Title: "Version", Width: 15},
		{Title: "Status", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	// Initialize progress bar
	p := progress.New(progress.WithDefaultGradient())

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	// Initialize viewport
	vp := viewport.New(0, 0)

	return &TUI{
		spackManager:    spackManager,
		environments:    []manager.Environment{},
		currentView:     envListView,
		installProgress: make(map[string]float64),
		list:            l,
		table:           t,
		progress:        p,
		spinner:         s,
		viewport:        vp,
	}
}

// environmentItem represents a list item for environments
type environmentItem struct {
	env manager.Environment
}

func (i environmentItem) FilterValue() string { return i.env.Name }
func (i environmentItem) Title() string       { return i.env.Name }
func (i environmentItem) Description() string {
	return fmt.Sprintf("%d packages", len(i.env.Packages))
}

// Init implements tea.Model
func (m *TUI) Init() tea.Cmd {
	return tea.Batch(
		m.loadEnvironments(),
		m.spinner.Tick,
	)
}

// Update implements tea.Model
func (m *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateComponentSizes()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Quit):
			if m.installing {
				return m, nil // Don't quit during installation
			}
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keyMap.Help):
			m.showingHelp = !m.showingHelp
			if m.showingHelp {
				m.currentView = helpView
			} else {
				m.currentView = envListView
			}

		case key.Matches(msg, keyMap.Back):
			if m.currentView != envListView && !m.installing {
				m.currentView = envListView
			}

		case key.Matches(msg, keyMap.Install):
			if m.currentView == envListView && len(m.environments) > 0 {
				return m, m.startInstallation()
			}

		case key.Matches(msg, keyMap.Refresh):
			return m, m.loadEnvironments()

		case key.Matches(msg, keyMap.Select):
			if m.currentView == envListView && len(m.environments) > 0 {
				m.currentView = envDetailView
			}
		}

	case environmentsLoadedMsg:
		m.environments = msg.environments
		m.updateEnvironmentList()

	case installationStartedMsg:
		m.installing = true
		m.currentView = installView

	case installationCompleteMsg:
		m.installing = false
		if msg.success {
			m.logs = append(m.logs, "Installation completed successfully!")
		} else {
			m.logs = append(m.logs, fmt.Sprintf("Installation failed: %s", msg.error))
		}

	case progressUpdateMsg:
		m.installProgress[msg.update.Package] = msg.update.Progress
		m.logs = append(m.logs, msg.update.Message)
		if len(m.logs) > 100 { // Keep log size manageable
			m.logs = m.logs[1:]
		}
	}

	// Update components based on current view
	switch m.currentView {
	case envListView:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
		m.selectedEnv = m.list.Index()

	case envDetailView:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)

	case installView:
		var updatedProgress tea.Model
		updatedProgress, cmd = m.progress.Update(msg)
		if p, ok := updatedProgress.(progress.Model); ok {
			m.progress = p
		}
		cmds = append(cmds, cmd)
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case logsView:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *TUI) View() string {
	if m.quitting {
		return "Thanks for using Spack Manager! 🚀\n"
	}

	var content string

	switch m.currentView {
	case envListView:
		content = m.renderEnvironmentList()
	case envDetailView:
		content = m.renderEnvironmentDetail()
	case installView:
		content = m.renderInstallation()
	case logsView:
		content = m.renderLogs()
	case helpView:
		content = m.renderHelp()
	}

	// Add header
	header := titleStyle.Render("🚀 Spack Manager")

	// Add status bar
	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		statusBar,
	)
}

// renderEnvironmentList renders the environment list view
func (m *TUI) renderEnvironmentList() string {
	if len(m.environments) == 0 {
		return "No Spack environments found.\n\nPress 'r' to refresh or 'q' to quit."
	}

	return m.list.View()
}

// renderEnvironmentDetail renders the environment detail view
func (m *TUI) renderEnvironmentDetail() string {
	if m.selectedEnv >= len(m.environments) {
		return "No environment selected"
	}

	env := m.environments[m.selectedEnv]

	title := headerStyle.Render(fmt.Sprintf("Environment: %s", env.Name))

	var rows []table.Row
	for _, pkg := range env.Packages {
		parts := strings.Split(pkg, "@")
		name := parts[0]
		version := ""
		if len(parts) > 1 {
			version = parts[1]
		}

		status := "✓ Installed"
		rows = append(rows, table.Row{name, version, status})
	}

	m.table.SetRows(rows)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.table.View(),
	)
}

// renderInstallation renders the installation view
func (m *TUI) renderInstallation() string {
	if !m.installing {
		return "No installation in progress"
	}

	var content strings.Builder

	content.WriteString(headerStyle.Render("Installing Environment"))
	content.WriteString("\n\n")

	// Show spinner
	content.WriteString(m.spinner.View())
	content.WriteString(" Installing packages...\n\n")

	// Show progress for each package
	for pkg, progress := range m.installProgress {
		content.WriteString(fmt.Sprintf("%s: %s\n", pkg, m.progress.ViewAs(progress)))
	}

	// Show recent logs
	content.WriteString("\n" + headerStyle.Render("Recent Activity:") + "\n")
	logCount := len(m.logs)
	start := 0
	if logCount > 5 {
		start = logCount - 5
	}

	for i := start; i < logCount; i++ {
		content.WriteString(fmt.Sprintf("• %s\n", m.logs[i]))
	}

	return content.String()
}

// renderLogs renders the logs view
func (m *TUI) renderLogs() string {
	m.viewport.SetContent(strings.Join(m.logs, "\n"))
	return headerStyle.Render("Installation Logs") + "\n\n" + m.viewport.View()
}

// renderHelp renders the help view
func (m *TUI) renderHelp() string {
	return RenderHelp()
}

// renderStatusBar renders the bottom status bar
func (m *TUI) renderStatusBar() string {
	var status string

	switch m.currentView {
	case envListView:
		status = "Press 'enter' to view details, 'i' to install, 'r' to refresh, '?' for help"
	case envDetailView:
		status = "Press 'i' to install, 'esc' to go back"
	case installView:
		status = "Installation in progress... Press 'l' to view logs"
	case helpView:
		status = "Press '?' to close help"
	}

	return lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(status)
}

// updateComponentSizes updates component sizes based on terminal size
func (m *TUI) updateComponentSizes() {
	headerHeight := 3
	statusHeight := 1
	contentHeight := m.height - headerHeight - statusHeight

	m.list.SetSize(m.width, contentHeight)
	m.table.SetWidth(m.width)
	m.table.SetHeight(contentHeight)
	m.viewport.Width = m.width
	m.viewport.Height = contentHeight
}

// updateEnvironmentList updates the list with current environments
func (m *TUI) updateEnvironmentList() {
	items := make([]list.Item, len(m.environments))
	for i, env := range m.environments {
		items[i] = environmentItem{env: env}
	}
	m.list.SetItems(items)
}

// Message types for tea.Model communication
type environmentsLoadedMsg struct {
	environments []manager.Environment
}

type installationStartedMsg struct{}

type installationCompleteMsg struct {
	success bool
	error   string
}

type progressUpdateMsg struct {
	update manager.ProgressUpdate
}

// Commands for tea.Model
func (m *TUI) loadEnvironments() tea.Cmd {
	return func() tea.Msg {
		environments, err := m.spackManager.ListEnvironments()
		if err != nil {
			// TODO: Handle error better
			return environmentsLoadedMsg{environments: []manager.Environment{}}
		}

		// Load detailed info for each environment
		for i, env := range environments {
			if detailedEnv, err := m.spackManager.GetEnvironmentInfo(env.Name); err == nil {
				environments[i] = *detailedEnv
			}
		}

		return environmentsLoadedMsg{environments: environments}
	}
}

func (m *TUI) startInstallation() tea.Cmd {
	return func() tea.Msg {
		if m.selectedEnv >= len(m.environments) {
			return installationCompleteMsg{success: false, error: "No environment selected"}
		}

		env := m.environments[m.selectedEnv]

		// Start installation in background
		go func() {
			progressChan := make(chan manager.ProgressUpdate, 100)

			// Monitor progress
			go func() {
				for update := range progressChan {
					// Send progress update to TUI
					// Note: In a real implementation, you'd need a way to send this back to the TUI
					_ = update
				}
			}()

			err := m.spackManager.InstallEnvironment(env.Name, progressChan)
			_ = err // TODO: Handle error
		}()

		return installationStartedMsg{}
	}
}

// Key mappings
type keyMapType struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Select  key.Binding
	Back    key.Binding
	Install key.Binding
	Refresh key.Binding
	Help    key.Binding
	Quit    key.Binding
}

var keyMap = keyMapType{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Run starts the TUI application
func Run(spackManager manager.Manager) error {
	tui := New(spackManager)

	p := tea.NewProgram(tui, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
