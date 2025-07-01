package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpackTUI represents the main TUI application for Spack management
type SpackTUI struct {
	spackManager    *SpackManager
	environments    []SpackEnvironment
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

// Key bindings
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Tab     key.Binding
	Escape  key.Binding
	Install key.Binding
	Logs    key.Binding
	Quit    key.Binding
	Help    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous view"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next view"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/confirm"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch view"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install environment"),
	),
	Logs: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "view logs"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#F25D94")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#F25D94"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	progressBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#50FA7B"))
)

// Progress update message
type progressMsg struct {
	env    string
	update ProgressUpdate
}

// Log message
type logMsg struct {
	message string
}

// Environment list item
type envItem struct {
	env SpackEnvironment
}

func (i envItem) FilterValue() string { return i.env.Name }
func (i envItem) Title() string       { return i.env.Name }
func (i envItem) Description() string {
	status := "Ready"
	if i.env.Concretized {
		status = "Installed"
	}
	return fmt.Sprintf("%s • %d packages", status, len(i.env.Packages))
}

// NewSpackTUI creates a new TUI instance
func NewSpackTUI(spackManager *SpackManager, environments []SpackEnvironment) *SpackTUI {
	// Create list items
	items := make([]list.Item, len(environments))
	for i, env := range environments {
		items[i] = envItem{env: env}
	}

	// Initialize list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Spack Environments"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	// Initialize table for environment details
	columns := []table.Column{
		{Title: "Package", Width: 30},
		{Title: "Version", Width: 15},
		{Title: "Status", Width: 15},
		{Title: "Build Time", Width: 15},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Initialize progress bar
	prog := progress.New(progress.WithDefaultGradient())

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize viewport for logs
	vp := viewport.New(0, 0)
	vp.SetContent("Spack Manager Logs\n" + strings.Repeat("─", 50) + "\n")

	return &SpackTUI{
		spackManager:    spackManager,
		environments:    environments,
		currentView:     envListView,
		installProgress: make(map[string]float64),
		logs:            []string{},
		list:            l,
		table:           t,
		progress:        prog,
		spinner:         s,
		viewport:        vp,
	}
}

// Init implements tea.Model
func (m *SpackTUI) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadEnvironments(),
	)
}

// Update implements tea.Model
func (m *SpackTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update component sizes
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 8)

		m.table.SetWidth(msg.Width - 4)
		m.table.SetHeight(msg.Height - 12)

		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Tab):
			m.nextView()

		case key.Matches(msg, keys.Escape):
			if m.currentView != envListView {
				m.currentView = envListView
			}

		case key.Matches(msg, keys.Install):
			if m.currentView == envListView && !m.installing {
				return m, m.installSelected()
			}

		case key.Matches(msg, keys.Logs):
			m.currentView = logsView

		case key.Matches(msg, keys.Help):
			m.showingHelp = !m.showingHelp
			if m.showingHelp {
				m.currentView = helpView
			} else {
				m.currentView = envListView
			}

		case key.Matches(msg, keys.Enter):
			if m.currentView == envListView {
				m.currentView = envDetailView
				m.updateTable()
			}
		}

	case progressMsg:
		m.installProgress[msg.env] = msg.update.Progress
		m.addLog(fmt.Sprintf("[%s] %s: %s", msg.env, msg.update.Stage, msg.update.Message))

		if msg.update.IsComplete {
			m.installing = false
			m.addLog(fmt.Sprintf("✅ Environment %s installation complete", msg.env))
		}

		if msg.update.IsError {
			m.addLog(fmt.Sprintf("❌ Error in %s: %s", msg.env, msg.update.Message))
		}

	case logMsg:
		m.addLog(msg.message)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update active component
	switch m.currentView {
	case envListView:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case envDetailView:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	case logsView:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *SpackTUI) View() string {
	if m.quitting {
		return "Goodbye! 👋\n"
	}

	if m.width == 0 {
		return "Loading..."
	}

	// Header
	header := titleStyle.Render("Spack Manager TUI") + "\n\n"

	// Content based on current view
	var content string
	switch m.currentView {
	case envListView:
		content = m.renderEnvironmentList()
	case envDetailView:
		content = m.renderEnvironmentDetail()
	case installView:
		content = m.renderInstallProgress()
	case logsView:
		content = m.renderLogs()
	case helpView:
		content = m.renderHelp()
	}

	// Footer with help
	footer := m.renderFooter()

	return header + content + "\n" + footer
}

func (m *SpackTUI) renderEnvironmentList() string {
	return m.list.View()
}

func (m *SpackTUI) renderEnvironmentDetail() string {
	if m.selectedEnv >= len(m.environments) {
		return errorStyle.Render("No environment selected")
	}

	env := m.environments[m.selectedEnv]

	// Environment header
	header := headerStyle.Render(fmt.Sprintf("Environment: %s", env.Name)) + "\n\n"

	// Environment info
	info := fmt.Sprintf("Packages: %d\n", len(env.Packages))
	info += fmt.Sprintf("View: %t\n", env.View)
	info += fmt.Sprintf("Concretized: %t\n\n", env.Concretized)

	// Package table
	return header + info + m.table.View()
}

func (m *SpackTUI) renderInstallProgress() string {
	if !m.installing {
		return normalStyle.Render("No installation in progress")
	}

	content := headerStyle.Render("Installation Progress") + "\n\n"

	for env, progress := range m.installProgress {
		content += fmt.Sprintf("%s %s\n",
			env,
			m.progress.ViewAs(progress/100.0))
	}

	content += "\n" + m.spinner.View() + " Installing packages..."

	return content
}

func (m *SpackTUI) renderLogs() string {
	return headerStyle.Render("Logs") + "\n\n" + m.viewport.View()
}

func (m *SpackTUI) renderFooter() string {
	help := "Press ? for help • "

	switch m.currentView {
	case envListView:
		help += "i: install • enter: details • v: logs • q: quit"
	case envDetailView:
		help += "esc: back • tab: switch view • q: quit"
	case installView:
		help += "esc: back • v: logs • q: quit"
	case logsView:
		help += "esc: back • tab: switch view • q: quit"
	case helpView:
		help += "esc: back • any key: return • q: quit"
	}

	return infoStyle.Render(help)
}

func (m *SpackTUI) nextView() {
	switch m.currentView {
	case envListView:
		if m.installing {
			m.currentView = installView
		} else {
			m.currentView = envDetailView
		}
	case envDetailView:
		if m.installing {
			m.currentView = installView
		} else {
			m.currentView = logsView
		}
	case installView:
		m.currentView = logsView
	case logsView:
		m.currentView = envListView
	case helpView:
		m.currentView = envListView
	}
}

func (m *SpackTUI) updateTable() {
	if m.selectedEnv >= len(m.environments) {
		return
	}

	env := m.environments[m.selectedEnv]
	rows := make([]table.Row, len(env.Packages))

	for i, pkg := range env.Packages {
		status := "Ready"
		if env.Concretized {
			status = "Installed"
		}

		rows[i] = table.Row{
			pkg,
			"latest", // Would extract version from package spec
			status,
			"-", // Would show actual build time
		}
	}

	m.table.SetRows(rows)
}

func (m *SpackTUI) addLog(message string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, message)
	m.logs = append(m.logs, logEntry)

	// Update viewport content
	content := "Spack Manager Logs\n" + strings.Repeat("─", 50) + "\n"
	for _, log := range m.logs {
		content += log + "\n"
	}
	m.viewport.SetContent(content)

	// Auto-scroll to bottom
	m.viewport.GotoBottom()
}

func (m *SpackTUI) loadEnvironments() tea.Cmd {
	return func() tea.Msg {
		// In a real implementation, this would load environments from Spack
		return logMsg{message: "Environments loaded"}
	}
}

func (m *SpackTUI) installSelected() tea.Cmd {
	if m.list.SelectedItem() == nil {
		return func() tea.Msg {
			return logMsg{message: "No environment selected"}
		}
	}

	item := m.list.SelectedItem().(envItem)
	env := item.env

	m.installing = true
	m.currentView = installView
	m.addLog(fmt.Sprintf("Starting installation of environment: %s", env.Name))

	return func() tea.Msg {
		ctx := context.Background()

		// Start installation with progress monitoring
		progressChan, err := m.spackManager.InstallEnvironment(ctx, &env)
		if err != nil {
			return logMsg{message: fmt.Sprintf("Failed to start installation: %v", err)}
		}

		// Monitor progress in background
		go func() {
			for range *progressChan {
				// Send progress updates to TUI
				// In a real implementation, you'd use a proper channel or callback
			}
		}()

		return logMsg{message: "Installation started"}
	}
}

// RunSpackTUI starts the TUI application
func RunSpackTUI(spackManager *SpackManager, environments []SpackEnvironment) error {
	tui := NewSpackTUI(spackManager, environments)

	p := tea.NewProgram(tui, tea.WithAltScreen())
	_, err := p.Run()

	return err
}

// SpackTUICommand creates a new command for the TUI
func SpackTUICommand(spackManager *SpackManager) tea.Cmd {
	return func() tea.Msg {
		// Load available environments
		environments := []SpackEnvironment{
			{
				Name:     "genomics_lab",
				Packages: []string{"bwa", "samtools", "gatk", "star"},
				View:     true,
			},
			{
				Name:     "climate_modeling",
				Packages: []string{"netcdf", "hdf5", "openmpi", "fftw"},
				View:     true,
			},
			{
				Name:     "ai_research_studio",
				Packages: []string{"python", "pytorch", "tensorflow", "cuda"},
				View:     true,
			},
		}

		// Launch TUI
		if err := RunSpackTUI(spackManager, environments); err != nil {
			return logMsg{message: fmt.Sprintf("TUI error: %v", err)}
		}

		return logMsg{message: "TUI session ended"}
	}
}
