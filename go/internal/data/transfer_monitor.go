package data

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TransferMonitor provides real-time monitoring of data transfers
type TransferMonitor struct {
	s3Manager       *S3Manager
	pipelineManager *PipelineManager

	// Monitoring state
	mu        sync.RWMutex
	transfers map[string]*MonitoredTransfer
	pipelines map[string]*MonitoredPipeline

	// UI state
	refreshInterval time.Duration
	lastUpdate      time.Time
}

// MonitoredTransfer represents a transfer being monitored
type MonitoredTransfer struct {
	ID          string
	Type        string // "upload" or "download"
	Source      string
	Destination string
	Progress    TransferProgress
	Status      TransferStatus
	StartTime   time.Time
}

// MonitoredPipeline represents a pipeline being monitored
type MonitoredPipeline struct {
	ID                  string
	Name                string
	Status              PipelineStatus
	TotalJobs           int
	CompletedJobs       int
	FailedJobs          int
	Progress            float64
	StartTime           time.Time
	EstimatedCompletion time.Time
}

// TransferStatus represents the status of a transfer
type TransferStatus string

const (
	TransferStatusQueued    TransferStatus = "queued"
	TransferStatusActive    TransferStatus = "active"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
	TransferStatusCancelled TransferStatus = "cancelled"
)

// MonitorModel represents the TUI model for transfer monitoring
type MonitorModel struct {
	monitor     *TransferMonitor
	table       table.Model
	progress    progress.Model
	width       int
	height      int
	selectedTab int
	tabs        []string
	quitting    bool
}

// NewTransferMonitor creates a new transfer monitor
func NewTransferMonitor(s3Manager *S3Manager, pipelineManager *PipelineManager) *TransferMonitor {
	return &TransferMonitor{
		s3Manager:       s3Manager,
		pipelineManager: pipelineManager,
		transfers:       make(map[string]*MonitoredTransfer),
		pipelines:       make(map[string]*MonitoredPipeline),
		refreshInterval: 1 * time.Second,
		lastUpdate:      time.Now(),
	}
}

// StartMonitoring begins monitoring transfers and pipelines
func (tm *TransferMonitor) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(tm.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tm.updateMonitoringData()
		}
	}
}

// updateMonitoringData refreshes monitoring data from managers
func (tm *TransferMonitor) updateMonitoringData() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Update S3 transfers
	activeTransfers := tm.s3Manager.GetActiveTransfers()
	for id, progress := range activeTransfers {
		if transfer, exists := tm.transfers[id]; exists {
			transfer.Progress = progress
		} else {
			// Parse transfer ID to get type and details
			parts := parseTransferID(id)
			tm.transfers[id] = &MonitoredTransfer{
				ID:          id,
				Type:        parts.Type,
				Source:      parts.Source,
				Destination: parts.Destination,
				Progress:    progress,
				Status:      TransferStatusActive,
				StartTime:   progress.StartTime,
			}
		}
	}

	// Update pipeline status
	pipelines := tm.pipelineManager.ListPipelines()
	for _, pipeline := range pipelines {
		monitored := &MonitoredPipeline{
			ID:        pipeline.ID,
			Name:      pipeline.Name,
			Status:    pipeline.Status,
			TotalJobs: len(pipeline.Jobs),
			StartTime: pipeline.CreatedAt,
		}

		// Calculate progress
		completed := 0
		failed := 0
		totalProgress := 0.0

		for _, job := range pipeline.Jobs {
			switch job.Status {
			case JobStatusCompleted:
				completed++
				totalProgress += 100.0
			case JobStatusFailed:
				failed++
			case JobStatusRunning:
				totalProgress += job.Progress
			}
		}

		monitored.CompletedJobs = completed
		monitored.FailedJobs = failed

		if monitored.TotalJobs > 0 {
			monitored.Progress = totalProgress / float64(monitored.TotalJobs)
		}

		tm.pipelines[pipeline.ID] = monitored
	}

	tm.lastUpdate = time.Now()
}

// GetTransferSummary returns a summary of active transfers
func (tm *TransferMonitor) GetTransferSummary() TransferSummary {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	summary := TransferSummary{
		ActiveTransfers:       0,
		CompletedTransfers:    0,
		FailedTransfers:       0,
		TotalBytesTransferred: 0,
		AverageSpeed:          0,
	}

	var totalSpeed int64
	var speedCount int

	for _, transfer := range tm.transfers {
		switch transfer.Status {
		case TransferStatusActive:
			summary.ActiveTransfers++
			if transfer.Progress.Speed > 0 {
				totalSpeed += transfer.Progress.Speed
				speedCount++
			}
		case TransferStatusCompleted:
			summary.CompletedTransfers++
		case TransferStatusFailed:
			summary.FailedTransfers++
		}

		summary.TotalBytesTransferred += transfer.Progress.BytesTransferred
	}

	if speedCount > 0 {
		summary.AverageSpeed = totalSpeed / int64(speedCount)
	}

	return summary
}

// GetPipelineSummary returns a summary of active pipelines
func (tm *TransferMonitor) GetPipelineSummary() PipelineSummary {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	summary := PipelineSummary{
		ActivePipelines:    0,
		CompletedPipelines: 0,
		FailedPipelines:    0,
		TotalJobs:          0,
		CompletedJobs:      0,
	}

	for _, pipeline := range tm.pipelines {
		switch pipeline.Status {
		case PipelineStatusRunning:
			summary.ActivePipelines++
		case PipelineStatusCompleted:
			summary.CompletedPipelines++
		case PipelineStatusFailed:
			summary.FailedPipelines++
		}

		summary.TotalJobs += pipeline.TotalJobs
		summary.CompletedJobs += pipeline.CompletedJobs
	}

	return summary
}

// TransferSummary provides aggregated transfer statistics
type TransferSummary struct {
	ActiveTransfers       int
	CompletedTransfers    int
	FailedTransfers       int
	TotalBytesTransferred int64
	AverageSpeed          int64
}

// PipelineSummary provides aggregated pipeline statistics
type PipelineSummary struct {
	ActivePipelines    int
	CompletedPipelines int
	FailedPipelines    int
	TotalJobs          int
	CompletedJobs      int
}

// NewMonitorModel creates a new TUI model for monitoring
func NewMonitorModel(monitor *TransferMonitor) MonitorModel {
	tabs := []string{"Transfers", "Pipelines", "Summary"}

	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Type", Width: 10},
		{Title: "Progress", Width: 15},
		{Title: "Speed", Width: 12},
		{Title: "ETA", Width: 10},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

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

	p := progress.New(progress.WithDefaultGradient())

	return MonitorModel{
		monitor:     monitor,
		table:       t,
		progress:    p,
		tabs:        tabs,
		selectedTab: 0,
	}
}

// Init initializes the model
func (m MonitorModel) Init() tea.Cmd {
	return tea.Batch(
		m.updateData(),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

// Update handles messages
func (m MonitorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(msg.Width - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.selectedTab = (m.selectedTab + 1) % len(m.tabs)
			return m, m.updateData()
		case "shift+tab":
			m.selectedTab = (m.selectedTab - 1 + len(m.tabs)) % len(m.tabs)
			return m, m.updateData()
		}

	case tickMsg:
		return m, tea.Batch(
			m.updateData(),
			tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			}),
		)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the model
func (m MonitorModel) View() string {
	if m.quitting {
		return "Exiting monitor...\n"
	}

	doc := strings.Builder{}

	// Render tabs
	var renderedTabs []string
	for i, tab := range m.tabs {
		var style lipgloss.Style
		isActive := i == m.selectedTab
		if isActive {
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("36")).
				Background(lipgloss.Color("57"))
		} else {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))
		}
		renderedTabs = append(renderedTabs, style.Render(tab))
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(tabsRow)
	doc.WriteString("\n\n")

	// Render content based on selected tab
	switch m.selectedTab {
	case 0: // Transfers
		doc.WriteString(m.renderTransfersView())
	case 1: // Pipelines
		doc.WriteString(m.renderPipelinesView())
	case 2: // Summary
		doc.WriteString(m.renderSummaryView())
	}

	return doc.String()
}

// renderTransfersView renders the transfers tab
func (m MonitorModel) renderTransfersView() string {
	m.monitor.mu.RLock()
	defer m.monitor.mu.RUnlock()

	rows := make([]table.Row, 0, len(m.monitor.transfers))
	for _, transfer := range m.monitor.transfers {
		speed := formatBytes(transfer.Progress.Speed) + "/s"
		eta := transfer.Progress.ETA.Round(time.Second).String()
		progress := fmt.Sprintf("%.1f%%", transfer.Progress.Percentage)

		rows = append(rows, table.Row{
			transfer.ID[:20], // Truncate ID
			transfer.Type,
			progress,
			speed,
			eta,
			string(transfer.Status),
		})
	}

	m.table.SetRows(rows)
	return m.table.View()
}

// renderPipelinesView renders the pipelines tab
func (m MonitorModel) renderPipelinesView() string {
	m.monitor.mu.RLock()
	defer m.monitor.mu.RUnlock()

	doc := strings.Builder{}

	for _, pipeline := range m.monitor.pipelines {
		status := string(pipeline.Status)
		progress := fmt.Sprintf("%.1f%%", pipeline.Progress)
		jobs := fmt.Sprintf("%d/%d", pipeline.CompletedJobs, pipeline.TotalJobs)

		doc.WriteString(fmt.Sprintf("Pipeline: %s\n", pipeline.Name))
		doc.WriteString(fmt.Sprintf("Status: %s | Progress: %s | Jobs: %s\n", status, progress, jobs))
		doc.WriteString(m.progress.ViewAs(pipeline.Progress / 100.0))
		doc.WriteString("\n\n")
	}

	return doc.String()
}

// renderSummaryView renders the summary tab
func (m MonitorModel) renderSummaryView() string {
	transferSummary := m.monitor.GetTransferSummary()
	pipelineSummary := m.monitor.GetPipelineSummary()

	doc := strings.Builder{}

	doc.WriteString("Transfer Summary:\n")
	doc.WriteString(fmt.Sprintf("  Active: %d | Completed: %d | Failed: %d\n",
		transferSummary.ActiveTransfers, transferSummary.CompletedTransfers, transferSummary.FailedTransfers))
	doc.WriteString(fmt.Sprintf("  Total Transferred: %s\n", formatBytes(transferSummary.TotalBytesTransferred)))
	doc.WriteString(fmt.Sprintf("  Average Speed: %s/s\n", formatBytes(transferSummary.AverageSpeed)))

	doc.WriteString("\nPipeline Summary:\n")
	doc.WriteString(fmt.Sprintf("  Active: %d | Completed: %d | Failed: %d\n",
		pipelineSummary.ActivePipelines, pipelineSummary.CompletedPipelines, pipelineSummary.FailedPipelines))
	doc.WriteString(fmt.Sprintf("  Jobs: %d/%d completed\n",
		pipelineSummary.CompletedJobs, pipelineSummary.TotalJobs))

	return doc.String()
}

// updateData returns a command to update the monitoring data
func (m MonitorModel) updateData() tea.Cmd {
	return func() tea.Msg {
		return updateDataMsg{}
	}
}

// Message types for the TUI
type tickMsg time.Time
type updateDataMsg struct{}

// parseTransferID parses a transfer ID to extract components
func parseTransferID(id string) transferIDParts {
	// ID format: "type:bucket:key"
	parts := strings.Split(id, ":")
	if len(parts) < 3 {
		return transferIDParts{Type: "unknown", Source: id, Destination: "unknown"}
	}

	transferType := parts[0]
	bucket := parts[1]
	key := parts[2]

	if transferType == "upload" {
		return transferIDParts{
			Type:        transferType,
			Source:      "local",
			Destination: fmt.Sprintf("s3://%s/%s", bucket, key),
		}
	} else {
		return transferIDParts{
			Type:        transferType,
			Source:      fmt.Sprintf("s3://%s/%s", bucket, key),
			Destination: "local",
		}
	}
}

type transferIDParts struct {
	Type        string
	Source      string
	Destination string
}

// formatBytes is already defined in pattern_analyzer.go
