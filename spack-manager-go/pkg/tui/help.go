package tui

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
			Foreground(lipgloss.Color("#04B575"))

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

// RenderHelp renders the help documentation
func RenderHelp() string {
	var help strings.Builder

	help.WriteString(helpTitleStyle.Render("🚀 Spack Manager - Help Guide"))
	help.WriteString("\n\n")

	// Navigation section
	help.WriteString(helpHeaderStyle.Render("🧭 Navigation"))
	help.WriteString("\n")
	help.WriteString(formatHelpItem("↑/k, ↓/j", "Navigate up/down in lists"))
	help.WriteString(formatHelpItem("←/h, →/l", "Navigate left/right in tables"))
	help.WriteString(formatHelpItem("enter", "Select item or view details"))
	help.WriteString(formatHelpItem("esc", "Go back to previous view"))
	help.WriteString(formatHelpItem("tab", "Switch between UI elements"))

	// Environment Management
	help.WriteString(helpHeaderStyle.Render("📦 Environment Management"))
	help.WriteString("\n")
	help.WriteString(formatHelpItem("r", "Refresh environment list"))
	help.WriteString(formatHelpItem("i", "Install selected environment"))
	help.WriteString(formatHelpItem("enter", "View environment details"))
	help.WriteString(formatHelpItem("d", "Delete environment (with confirmation)"))

	// Installation
	help.WriteString(helpHeaderStyle.Render("⚡ Installation"))
	help.WriteString("\n")
	help.WriteString(formatHelpItem("i", "Start installation of selected environment"))
	help.WriteString(formatHelpItem("l", "View installation logs"))
	help.WriteString(formatHelpItem("p", "Show detailed progress"))
	help.WriteString(formatHelpItem("ctrl+c", "Cancel installation (if supported)"))

	// Views and Interface
	help.WriteString(helpHeaderStyle.Render("👁️  Views"))
	help.WriteString("\n")
	help.WriteString(formatHelpItem("1", "Environment list view"))
	help.WriteString(formatHelpItem("2", "Environment detail view"))
	help.WriteString(formatHelpItem("3", "Installation progress view"))
	help.WriteString(formatHelpItem("4", "Logs view"))
	help.WriteString(formatHelpItem("?", "Toggle this help"))

	// General Controls
	help.WriteString(helpHeaderStyle.Render("⚙️  General"))
	help.WriteString("\n")
	help.WriteString(formatHelpItem("q", "Quit application"))
	help.WriteString(formatHelpItem("ctrl+c", "Force quit"))
	help.WriteString(formatHelpItem("ctrl+l", "Clear screen"))
	help.WriteString(formatHelpItem("F5", "Force refresh all"))

	// Tips and Tricks
	help.WriteString(helpHeaderStyle.Render("💡 Tips & Tricks"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Use 'r' frequently to refresh the environment list"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Installation progress is shown in real-time"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Check logs view for detailed installation information"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Binary cache speeds up installations significantly"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Environment details show all installed packages"))

	// Spack Information
	help.WriteString(helpHeaderStyle.Render("📚 About Spack"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Spack is a package manager for scientific computing"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Environments allow isolated package collections"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Binary caches provide pre-built packages"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Visit https://spack.io for more information"))

	// Status Information
	help.WriteString(helpHeaderStyle.Render("📊 Status Indicators"))
	help.WriteString("\n")
	help.WriteString("• " + helpKeyStyle.Render("✓") + " " + helpDescStyle.Render("Package installed successfully"))
	help.WriteString("\n")
	help.WriteString("• " + helpKeyStyle.Render("⚡") + " " + helpDescStyle.Render("Installation in progress"))
	help.WriteString("\n")
	help.WriteString("• " + helpKeyStyle.Render("❌") + " " + helpDescStyle.Render("Installation failed"))
	help.WriteString("\n")
	help.WriteString("• " + helpKeyStyle.Render("⏳") + " " + helpDescStyle.Render("Waiting/queued"))
	help.WriteString("\n")
	help.WriteString("• " + helpKeyStyle.Render("🔄") + " " + helpDescStyle.Render("Refreshing data"))

	// Troubleshooting
	help.WriteString(helpHeaderStyle.Render("🔧 Troubleshooting"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("If environments don't appear, check Spack installation"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Installation failures are logged in the logs view"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Use 'r' to refresh if data seems stale"))
	help.WriteString("\n")
	help.WriteString("• " + helpDescStyle.Render("Check terminal size if UI appears corrupted"))

	help.WriteString("\n")
	help.WriteString(helpHeaderStyle.Render("🚀 Happy Packaging!"))
	help.WriteString("\n")
	help.WriteString(helpDescStyle.Render("Press '?' again to close this help"))

	return help.String()
}

// formatHelpItem formats a help item with key and description
func formatHelpItem(key, description string) string {
	return "  " + helpKeyStyle.Render(key) + "  " + helpDescStyle.Render(description) + "\n"
}
