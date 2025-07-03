package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Manager handles monitoring, SLA tracking, and alerting
type Manager struct {
	slaDefinitions  map[string]*SLADefinition
	slaMetrics      map[string]*SLAMetric
	alerts          map[string]*Alert
	dashboards      map[string]*Dashboard
	metrics         []Metric
	config          *MonitoringConfig
	dataDir         string
	mutex           sync.RWMutex
	collectors      map[string]MetricCollector
	evaluators      map[string]SLAEvaluator
}

// MetricCollector interface for collecting metrics from various sources
type MetricCollector interface {
	CollectMetrics() ([]Metric, error)
	GetName() string
}

// SLAEvaluator interface for evaluating SLA compliance
type SLAEvaluator interface {
	EvaluateSLA(sla *SLADefinition, metrics []Metric) (*SLAMetric, error)
	GetSupportedSLATypes() []SLAType
}

// NewManager creates a new monitoring manager
func NewManager(dataDir string) *Manager {
	return &Manager{
		slaDefinitions: make(map[string]*SLADefinition),
		slaMetrics:     make(map[string]*SLAMetric),
		alerts:         make(map[string]*Alert),
		dashboards:     make(map[string]*Dashboard),
		metrics:        make([]Metric, 0),
		config: &MonitoringConfig{
			MetricsRetention:     30 * 24 * time.Hour, // 30 days
			SampleRate:           5 * time.Second,
			AlertingEnabled:      true,
			NotificationChannels: make([]NotificationChannel, 0),
			SLADefinitions:       make([]SLADefinition, 0),
			Dashboards:          make([]Dashboard, 0),
		},
		dataDir:    dataDir,
		collectors: make(map[string]MetricCollector),
		evaluators: make(map[string]SLAEvaluator),
	}
}

// CreateSLA creates a new SLA definition
func (m *Manager) CreateSLA(sla *SLADefinition) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if sla.ID == "" {
		sla.ID = fmt.Sprintf("sla-%s-%d", sla.Type, time.Now().Unix())
	}

	// Set timestamps
	now := time.Now()
	sla.CreatedAt = now
	sla.UpdatedAt = now

	// Store SLA definition
	m.slaDefinitions[sla.ID] = sla

	// Initialize SLA metric
	slaMetric := &SLAMetric{
		SLADefinitionID: sla.ID,
		TargetValue:     sla.TargetValue,
		Status:          SLAStatusUnknown,
		LastEvaluated:   time.Now(),
	}
	m.slaMetrics[sla.ID] = slaMetric

	// Persist to disk
	return m.saveSLA(sla)
}

// GetSLA retrieves an SLA definition by ID
func (m *Manager) GetSLA(slaID string) (*SLADefinition, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	sla, exists := m.slaDefinitions[slaID]
	if !exists {
		return nil, fmt.Errorf("SLA %s not found", slaID)
	}

	return sla, nil
}

// ListSLAs returns all SLA definitions for a tenant
func (m *Manager) ListSLAs(tenantID string) []*SLADefinition {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var slas []*SLADefinition
	for _, sla := range m.slaDefinitions {
		if tenantID == "" || sla.TenantID == tenantID {
			slas = append(slas, sla)
		}
	}

	return slas
}

// UpdateSLA updates an existing SLA definition
func (m *Manager) UpdateSLA(slaID string, updates *SLADefinition) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	sla, exists := m.slaDefinitions[slaID]
	if !exists {
		return fmt.Errorf("SLA %s not found", slaID)
	}

	// Update fields
	if updates.Name != "" {
		sla.Name = updates.Name
	}
	if updates.Description != "" {
		sla.Description = updates.Description
	}
	if updates.TargetValue != 0 {
		sla.TargetValue = updates.TargetValue
	}
	if updates.Threshold != 0 {
		sla.Threshold = updates.Threshold
	}

	// Update timestamp
	sla.UpdatedAt = time.Now()

	// Persist to disk
	return m.saveSLA(sla)
}

// RecordMetric records a new metric
func (m *Manager) RecordMetric(metric Metric) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Set timestamp if not provided
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}

	// Add to metrics collection
	m.metrics = append(m.metrics, metric)

	// Clean up old metrics based on retention policy
	m.cleanupOldMetrics()

	return nil
}

// GetMetrics retrieves metrics based on query parameters
func (m *Manager) GetMetrics(query MetricQuery) ([]Metric, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var filteredMetrics []Metric

	for _, metric := range m.metrics {
		// Apply filters
		if query.MetricName != "" && metric.Name != query.MetricName {
			continue
		}

		if !query.StartTime.IsZero() && metric.Timestamp.Before(query.StartTime) {
			continue
		}

		if !query.EndTime.IsZero() && metric.Timestamp.After(query.EndTime) {
			continue
		}

		if len(query.Labels) > 0 {
			match := true
			for key, value := range query.Labels {
				if metric.Labels[key] != value {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		filteredMetrics = append(filteredMetrics, metric)
	}

	return filteredMetrics, nil
}

// EvaluateSLAs evaluates all active SLAs
func (m *Manager) EvaluateSLAs() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, sla := range m.slaDefinitions {
		if !sla.IsActive {
			continue
		}

		// Get relevant metrics for evaluation
		endTime := time.Now()
		startTime := endTime.Add(-sla.EvaluationWindow)

		query := MetricQuery{
			StartTime: startTime,
			EndTime:   endTime,
		}

		metrics, err := m.getMetricsUnsafe(query)
		if err != nil {
			continue
		}

		// Evaluate SLA compliance
		slaMetric, err := m.evaluateSLACompliance(sla, metrics)
		if err != nil {
			continue
		}

		// Update SLA metric
		m.slaMetrics[sla.ID] = slaMetric

		// Check for SLA breaches and generate alerts
		if err := m.checkSLABreaches(sla, slaMetric); err != nil {
			fmt.Printf("Error checking SLA breaches for %s: %v\n", sla.ID, err)
		}
	}

	return nil
}

// GetSLAMetrics returns current SLA metrics for a tenant
func (m *Manager) GetSLAMetrics(tenantID string) []*SLAMetric {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var metrics []*SLAMetric
	for slaID, slaMetric := range m.slaMetrics {
		if sla, exists := m.slaDefinitions[slaID]; exists {
			if tenantID == "" || sla.TenantID == tenantID {
				metrics = append(metrics, slaMetric)
			}
		}
	}

	return metrics
}

// CreateAlert creates a new alert
func (m *Manager) CreateAlert(alert *Alert) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if alert.ID == "" {
		alert.ID = fmt.Sprintf("alert-%d", time.Now().UnixNano())
	}

	alert.CreatedAt = time.Now()
	alert.Status = AlertStatusOpen

	m.alerts[alert.ID] = alert

	// Send notifications
	if err := m.sendAlertNotifications(alert); err != nil {
		fmt.Printf("Error sending alert notifications: %v\n", err)
	}

	return m.saveAlert(alert)
}

// GetAlerts returns alerts for a tenant
func (m *Manager) GetAlerts(tenantID string, status AlertStatus) []*Alert {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var alerts []*Alert
	for _, alert := range m.alerts {
		if tenantID != "" && alert.TenantID != tenantID {
			continue
		}

		if status != "" && alert.Status != status {
			continue
		}

		alerts = append(alerts, alert)
	}

	return alerts
}

// CreateDashboard creates a new monitoring dashboard
func (m *Manager) CreateDashboard(dashboard *Dashboard) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if dashboard.ID == "" {
		dashboard.ID = fmt.Sprintf("dashboard-%d", time.Now().Unix())
	}

	now := time.Now()
	dashboard.CreatedAt = now
	dashboard.UpdatedAt = now

	m.dashboards[dashboard.ID] = dashboard

	return m.saveDashboard(dashboard)
}

// GetDashboard retrieves a dashboard by ID
func (m *Manager) GetDashboard(dashboardID string) (*Dashboard, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	dashboard, exists := m.dashboards[dashboardID]
	if !exists {
		return nil, fmt.Errorf("dashboard %s not found", dashboardID)
	}

	return dashboard, nil
}

// ListDashboards returns all dashboards for a tenant
func (m *Manager) ListDashboards(tenantID string) []*Dashboard {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var dashboards []*Dashboard
	for _, dashboard := range m.dashboards {
		if tenantID == "" || dashboard.TenantID == tenantID {
			dashboards = append(dashboards, dashboard)
		}
	}

	return dashboards
}

// GenerateComplianceReport generates a compliance report for a tenant
func (m *Manager) GenerateComplianceReport(tenantID string, period ReportPeriod) (*ComplianceReport, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	report := &ComplianceReport{
		ID:           fmt.Sprintf("report-%s-%d", tenantID, time.Now().Unix()),
		TenantID:     tenantID,
		ReportPeriod: period,
		GeneratedAt:  time.Now(),
		SLACompliance: make([]SLAComplianceResult, 0),
	}

	var totalCompliance float64
	var compliantCount, breachedCount int

	// Evaluate each SLA for the report period
	for _, sla := range m.slaDefinitions {
		if sla.TenantID != tenantID {
			continue
		}

		// Get metrics for the report period
		query := MetricQuery{
			StartTime: period.StartTime,
			EndTime:   period.EndTime,
		}

		metrics, err := m.getMetricsUnsafe(query)
		if err != nil {
			continue
		}

		// Calculate compliance for this SLA
		complianceResult := m.calculateSLACompliance(sla, metrics, period)
		report.SLACompliance = append(report.SLACompliance, complianceResult)

		totalCompliance += complianceResult.ComplianceRate
		if complianceResult.Status == SLAStatusCompliant {
			compliantCount++
		} else if complianceResult.Status == SLAStatusBreached {
			breachedCount++
		}
	}

	// Calculate summary
	totalSLAs := len(report.SLACompliance)
	if totalSLAs > 0 {
		report.OverallScore = totalCompliance / float64(totalSLAs)
	}

	report.Summary = ComplianceSummary{
		TotalSLAs:         totalSLAs,
		CompliantSLAs:     compliantCount,
		BreachedSLAs:      breachedCount,
		AverageCompliance: report.OverallScore,
	}

	// Generate recommendations
	report.Recommendations = m.generateRecommendations(report)

	return report, nil
}

// Helper methods

// MetricQuery represents a query for metrics
type MetricQuery struct {
	MetricName string
	StartTime  time.Time
	EndTime    time.Time
	Labels     map[string]string
}

// getMetricsUnsafe retrieves metrics without locking (internal use)
func (m *Manager) getMetricsUnsafe(query MetricQuery) ([]Metric, error) {
	var filteredMetrics []Metric

	for _, metric := range m.metrics {
		// Apply filters
		if query.MetricName != "" && metric.Name != query.MetricName {
			continue
		}

		if !query.StartTime.IsZero() && metric.Timestamp.Before(query.StartTime) {
			continue
		}

		if !query.EndTime.IsZero() && metric.Timestamp.After(query.EndTime) {
			continue
		}

		if len(query.Labels) > 0 {
			match := true
			for key, value := range query.Labels {
				if metric.Labels[key] != value {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		filteredMetrics = append(filteredMetrics, metric)
	}

	return filteredMetrics, nil
}

// evaluateSLACompliance evaluates SLA compliance based on metrics
func (m *Manager) evaluateSLACompliance(sla *SLADefinition, metrics []Metric) (*SLAMetric, error) {
	// This is a simplified implementation - in practice, this would use
	// more sophisticated evaluation logic based on the SLA type and metric query
	
	slaMetric := &SLAMetric{
		SLADefinitionID: sla.ID,
		TargetValue:     sla.TargetValue,
		LastEvaluated:   time.Now(),
	}

	if len(metrics) == 0 {
		slaMetric.Status = SLAStatusUnknown
		return slaMetric, nil
	}

	// Calculate average value from metrics
	var total float64
	for _, metric := range metrics {
		total += metric.Value
	}
	averageValue := total / float64(len(metrics))
	slaMetric.CurrentValue = averageValue

	// Determine compliance based on alert condition
	var isCompliant bool
	switch sla.AlertCondition {
	case AlertConditionLessThan:
		isCompliant = averageValue < sla.TargetValue
	case AlertConditionGreaterThan:
		isCompliant = averageValue > sla.TargetValue
	default:
		isCompliant = averageValue == sla.TargetValue
	}

	if isCompliant {
		slaMetric.Status = SLAStatusCompliant
		slaMetric.ComplianceRate = 100.0
	} else {
		// Check if it's a warning or breach
		thresholdBreached := false
		switch sla.AlertCondition {
		case AlertConditionLessThan:
			thresholdBreached = averageValue > sla.Threshold
		case AlertConditionGreaterThan:
			thresholdBreached = averageValue < sla.Threshold
		}

		if thresholdBreached {
			slaMetric.Status = SLAStatusBreached
		} else {
			slaMetric.Status = SLAStatusWarning
		}

		// Calculate compliance rate based on how close we are to target
		distance := abs(averageValue - sla.TargetValue)
		maxDistance := abs(sla.Threshold - sla.TargetValue)
		if maxDistance > 0 {
			slaMetric.ComplianceRate = max(0, (1-(distance/maxDistance))*100)
		}
	}

	return slaMetric, nil
}

// checkSLABreaches checks for SLA breaches and creates alerts if needed
func (m *Manager) checkSLABreaches(sla *SLADefinition, slaMetric *SLAMetric) error {
	if slaMetric.Status == SLAStatusBreached || slaMetric.Status == SLAStatusWarning {
		// Create alert for SLA breach
		alert := &Alert{
			TenantID:     sla.TenantID,
			SLAID:        sla.ID,
			Type:         AlertTypeSLA,
			Title:        fmt.Sprintf("SLA Breach: %s", sla.Name),
			Description:  fmt.Sprintf("SLA %s is not meeting target of %v. Current value: %v", sla.Name, sla.TargetValue, slaMetric.CurrentValue),
			TriggerValue: slaMetric.CurrentValue,
			Threshold:    sla.Threshold,
			Metadata: map[string]string{
				"sla_id":   sla.ID,
				"sla_type": string(sla.Type),
			},
		}

		if slaMetric.Status == SLAStatusBreached {
			alert.Severity = AlertSeverityCritical
		} else {
			alert.Severity = AlertSeverityWarning
		}

		return m.CreateAlert(alert)
	}

	return nil
}

// sendAlertNotifications sends notifications for an alert
func (m *Manager) sendAlertNotifications(alert *Alert) error {
	// This would integrate with actual notification systems (email, Slack, etc.)
	fmt.Printf("🚨 ALERT: %s - %s\n", alert.Title, alert.Description)
	return nil
}

// calculateSLACompliance calculates compliance for a specific period
func (m *Manager) calculateSLACompliance(sla *SLADefinition, metrics []Metric, period ReportPeriod) SLAComplianceResult {
	result := SLAComplianceResult{
		SLAID:       sla.ID,
		SLAName:     sla.Name,
		TargetValue: sla.TargetValue,
		Status:      SLAStatusUnknown,
	}

	if len(metrics) == 0 {
		return result
	}

	// Calculate average value
	var total float64
	for _, metric := range metrics {
		total += metric.Value
	}
	result.AverageValue = total / float64(len(metrics))

	// Determine compliance
	var isCompliant bool
	switch sla.AlertCondition {
	case AlertConditionLessThan:
		isCompliant = result.AverageValue < sla.TargetValue
	case AlertConditionGreaterThan:
		isCompliant = result.AverageValue > sla.TargetValue
	default:
		isCompliant = result.AverageValue == sla.TargetValue
	}

	if isCompliant {
		result.Status = SLAStatusCompliant
		result.ComplianceRate = 100.0
	} else {
		result.Status = SLAStatusBreached
		// Calculate compliance rate
		distance := abs(result.AverageValue - sla.TargetValue)
		maxDistance := abs(sla.Threshold - sla.TargetValue)
		if maxDistance > 0 {
			result.ComplianceRate = max(0, (1-(distance/maxDistance))*100)
		}
	}

	return result
}

// generateRecommendations generates recommendations based on compliance report
func (m *Manager) generateRecommendations(report *ComplianceReport) []string {
	var recommendations []string

	if report.OverallScore < 80 {
		recommendations = append(recommendations, "Overall SLA compliance is below 80%. Consider reviewing SLA targets and system performance.")
	}

	if report.Summary.BreachedSLAs > 0 {
		recommendations = append(recommendations, fmt.Sprintf("%d SLAs are currently breached. Immediate attention required.", report.Summary.BreachedSLAs))
	}

	if report.Summary.TotalSLAs == 0 {
		recommendations = append(recommendations, "No SLAs defined. Consider creating SLAs to track service quality.")
	}

	return recommendations
}

// cleanupOldMetrics removes metrics older than the retention period
func (m *Manager) cleanupOldMetrics() {
	cutoff := time.Now().Add(-m.config.MetricsRetention)
	
	filteredMetrics := make([]Metric, 0, len(m.metrics))
	for _, metric := range m.metrics {
		if metric.Timestamp.After(cutoff) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}
	
	m.metrics = filteredMetrics
}

// Persistence methods

func (m *Manager) saveSLA(sla *SLADefinition) error {
	if m.dataDir == "" {
		return nil
	}

	slaDir := filepath.Join(m.dataDir, "slas")
	if err := os.MkdirAll(slaDir, 0755); err != nil {
		return err
	}

	slaFile := filepath.Join(slaDir, sla.ID+".json")
	data, err := json.MarshalIndent(sla, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(slaFile, data, 0644)
}

func (m *Manager) saveAlert(alert *Alert) error {
	if m.dataDir == "" {
		return nil
	}

	alertDir := filepath.Join(m.dataDir, "alerts")
	if err := os.MkdirAll(alertDir, 0755); err != nil {
		return err
	}

	alertFile := filepath.Join(alertDir, alert.ID+".json")
	data, err := json.MarshalIndent(alert, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(alertFile, data, 0644)
}

func (m *Manager) saveDashboard(dashboard *Dashboard) error {
	if m.dataDir == "" {
		return nil
	}

	dashboardDir := filepath.Join(m.dataDir, "dashboards")
	if err := os.MkdirAll(dashboardDir, 0755); err != nil {
		return err
	}

	dashboardFile := filepath.Join(dashboardDir, dashboard.ID+".json")
	data, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dashboardFile, data, 0644)
}

// Utility functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}