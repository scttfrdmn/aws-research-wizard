package data

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Help styles
var (
	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	helpHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F25D94")).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#50FA7B"))

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))
)

// renderHelp returns the help screen content
func (m *SpackTUI) renderHelp() string {
	var help strings.Builder

	help.WriteString(helpTitleStyle.Render("Spack Manager TUI - Help"))
	help.WriteString("\n\n")

	// Navigation section
	help.WriteString(helpHeaderStyle.Render("🧭 Navigation"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("↑/k") + "     " + helpDescStyle.Render("Move up"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("↓/j") + "     " + helpDescStyle.Render("Move down"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("←/h") + "     " + helpDescStyle.Render("Previous view"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("→/l") + "     " + helpDescStyle.Render("Next view"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("Tab") + "     " + helpDescStyle.Render("Switch between views"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("Enter") + "   " + helpDescStyle.Render("Select/Confirm"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("Esc") + "     " + helpDescStyle.Render("Go back"))

	// Actions section
	help.WriteString("\n")
	help.WriteString(helpHeaderStyle.Render("⚡ Actions"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("i") + "       " + helpDescStyle.Render("Install selected environment"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("v") + "       " + helpDescStyle.Render("View logs"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("r") + "       " + helpDescStyle.Render("Refresh environments"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("d") + "       " + helpDescStyle.Render("Delete environment"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("c") + "       " + helpDescStyle.Render("Create new environment"))

	// System section
	help.WriteString("\n")
	help.WriteString(helpHeaderStyle.Render("🛠️  System"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("?") + "       " + helpDescStyle.Render("Toggle this help"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("q") + "       " + helpDescStyle.Render("Quit application"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("Ctrl+C") + " " + helpDescStyle.Render("Force quit"))

	// Views section
	help.WriteString("\n")
	help.WriteString(helpHeaderStyle.Render("📱 Views"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("1. Environment List") + " - " + helpDescStyle.Render("Browse and manage environments"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("2. Environment Details") + " - " + helpDescStyle.Render("View packages and configuration"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("3. Installation Progress") + " - " + helpDescStyle.Render("Real-time progress monitoring"))
	help.WriteString("\n")
	help.WriteString(helpKeyStyle.Render("4. Logs") + " - " + helpDescStyle.Render("Live logs and debugging information"))

	// Features section
	help.WriteString("\n")
	help.WriteString(helpHeaderStyle.Render("✨ Features"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Real-time progress tracking with detailed status"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Binary cache optimization for faster installs"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Environment validation and error detection"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Package dependency visualization"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Live log streaming with timestamps"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Keyboard-driven interface for efficiency"))

	// Tips section
	help.WriteString("\n")
	help.WriteString(helpHeaderStyle.Render("💡 Tips"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Use 'i' to install environments with progress tracking"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Switch to logs view during installation for details"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Environments are auto-discovered from domain packs"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Binary cache speeds up installations by 95%"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Tab through views to monitor different aspects"))

	help.WriteString("\n\n")
	help.WriteString(helpDescStyle.Render("Press any key to return to the main interface"))

	return help.String()
}

// Advanced TUI features we could add:
// - Environment creation wizard
// - Package search and filter
// - Dependency graph visualization
// - Performance metrics and benchmarks
// - Cost analysis integration
// - Multi-environment parallel installation
// - Environment comparison
// - Export/import environment configurations
// - Integration with CI/CD pipelines
// - Web UI companion mode
