package monitoring

import (
	"time"
)

// MetricType defines the type of metric being collected
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

// Metric represents a monitoring metric
type Metric struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Value       float64           `json:"value"`
	Labels      map[string]string `json:"labels"`
	Timestamp   time.Time         `json:"timestamp"`
	Description string            `json:"description"`
	Unit        string            `json:"unit"`
}

// SLAType defines the type of SLA being tracked
type SLAType string

const (
	SLATypeAvailability SLAType = "availability"
	SLATypePerformance  SLAType = "performance"
	SLATypeReliability  SLAType = "reliability"
	SLATypeCost         SLAType = "cost"
	SLATypeSecurity     SLAType = "security"
)

// SLADefinition defines an SLA with its targets and measurement criteria
type SLADefinition struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Type            SLAType           `json:"type"`
	TenantID        string            `json:"tenantId"`
	Description     string            `json:"description"`
	TargetValue     float64           `json:"targetValue"`
	Threshold       float64           `json:"threshold"`
	Unit            string            `json:"unit"`
	MetricQuery     string            `json:"metricQuery"`
	EvaluationWindow time.Duration    `json:"evaluationWindow"`
	AlertCondition  AlertCondition    `json:"alertCondition"`
	IsActive        bool              `json:"isActive"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	Tags            map[string]string `json:"tags"`
}

// AlertCondition defines when an alert should be triggered
type AlertCondition string

const (
	AlertConditionGreaterThan AlertCondition = "greater_than"
	AlertConditionLessThan    AlertCondition = "less_than"
	AlertConditionEquals      AlertCondition = "equals"
	AlertConditionNotEquals   AlertCondition = "not_equals"
)

// SLAMetric represents the current state of an SLA
type SLAMetric struct {
	SLADefinitionID string    `json:"slaDefinitionId"`
	CurrentValue    float64   `json:"currentValue"`
	TargetValue     float64   `json:"targetValue"`
	ComplianceRate  float64   `json:"complianceRate"`
	Status          SLAStatus `json:"status"`
	LastEvaluated   time.Time `json:"lastEvaluated"`
	LastAlert       time.Time `json:"lastAlert"`
	AlertCount      int       `json:"alertCount"`
	TrendDirection  string    `json:"trendDirection"`
}

// SLAStatus represents the current status of an SLA
type SLAStatus string

const (
	SLAStatusCompliant    SLAStatus = "compliant"
	SLAStatusWarning      SLAStatus = "warning"
	SLAStatusBreached     SLAStatus = "breached"
	SLAStatusUnknown      SLAStatus = "unknown"
	SLAStatusMaintenance  SLAStatus = "maintenance"
)

// Alert represents a monitoring alert
type Alert struct {
	ID           string            `json:"id"`
	TenantID     string            `json:"tenantId"`
	SLAID        string            `json:"slaId"`
	Type         AlertType         `json:"type"`
	Severity     AlertSeverity     `json:"severity"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Status       AlertStatus       `json:"status"`
	TriggerValue float64           `json:"triggerValue"`
	Threshold    float64           `json:"threshold"`
	Metadata     map[string]string `json:"metadata"`
	CreatedAt    time.Time         `json:"createdAt"`
	ResolvedAt   time.Time         `json:"resolvedAt,omitempty"`
	AcknowledgedAt time.Time       `json:"acknowledgedAt,omitempty"`
	AcknowledgedBy string          `json:"acknowledgedBy,omitempty"`
}

// AlertType defines the type of alert
type AlertType string

const (
	AlertTypeSLA         AlertType = "sla"
	AlertTypeSystem      AlertType = "system"
	AlertTypeSecurity    AlertType = "security"
	AlertTypePerformance AlertType = "performance"
	AlertTypeCost        AlertType = "cost"
)

// AlertSeverity defines the severity level of an alert
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
	AlertSeverityEmergency AlertSeverity = "emergency"
)

// AlertStatus defines the current status of an alert
type AlertStatus string

const (
	AlertStatusOpen         AlertStatus = "open"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
	AlertStatusSuppressed   AlertStatus = "suppressed"
)

// Dashboard represents a custom monitoring dashboard
type Dashboard struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenantId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Widgets     []DashboardWidget `json:"widgets"`
	Layout      DashboardLayout   `json:"layout"`
	IsShared    bool              `json:"isShared"`
	CreatedBy   string            `json:"createdBy"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Tags        map[string]string `json:"tags"`
}

// DashboardWidget represents a widget on a dashboard
type DashboardWidget struct {
	ID          string            `json:"id"`
	Type        WidgetType        `json:"type"`
	Title       string            `json:"title"`
	MetricQuery string            `json:"metricQuery"`
	Position    WidgetPosition    `json:"position"`
	Size        WidgetSize        `json:"size"`
	Config      map[string]interface{} `json:"config"`
	RefreshRate time.Duration     `json:"refreshRate"`
}

// WidgetType defines the type of dashboard widget
type WidgetType string

const (
	WidgetTypeLineChart   WidgetType = "line_chart"
	WidgetTypeBarChart    WidgetType = "bar_chart"
	WidgetTypePieChart    WidgetType = "pie_chart"
	WidgetTypeGauge       WidgetType = "gauge"
	WidgetTypeNumber      WidgetType = "number"
	WidgetTypeTable       WidgetType = "table"
	WidgetTypeHeatmap     WidgetType = "heatmap"
	WidgetTypeSLAStatus   WidgetType = "sla_status"
)

// WidgetPosition defines the position of a widget on the dashboard
type WidgetPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// WidgetSize defines the size of a widget
type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DashboardLayout defines the layout configuration for a dashboard
type DashboardLayout struct {
	Columns     int  `json:"columns"`
	AutoResize  bool `json:"autoResize"`
	Responsive  bool `json:"responsive"`
	GridSize    int  `json:"gridSize"`
}

// ComplianceReport represents a compliance report for SLA tracking
type ComplianceReport struct {
	ID               string                   `json:"id"`
	TenantID         string                   `json:"tenantId"`
	ReportPeriod     ReportPeriod            `json:"reportPeriod"`
	GeneratedAt      time.Time                `json:"generatedAt"`
	SLACompliance    []SLAComplianceResult   `json:"slaCompliance"`
	OverallScore     float64                  `json:"overallScore"`
	Summary          ComplianceSummary        `json:"summary"`
	Recommendations  []string                 `json:"recommendations"`
}

// ReportPeriod defines the time period for a report
type ReportPeriod struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Type      string    `json:"type"` // daily, weekly, monthly, quarterly
}

// SLAComplianceResult represents the compliance result for a specific SLA
type SLAComplianceResult struct {
	SLAID           string    `json:"slaId"`
	SLAName         string    `json:"slaName"`
	ComplianceRate  float64   `json:"complianceRate"`
	AverageValue    float64   `json:"averageValue"`
	TargetValue     float64   `json:"targetValue"`
	BreachCount     int       `json:"breachCount"`
	BreachDuration  time.Duration `json:"breachDuration"`
	Status          SLAStatus `json:"status"`
	TrendDirection  string    `json:"trendDirection"`
}

// ComplianceSummary provides an overall summary of compliance
type ComplianceSummary struct {
	TotalSLAs           int     `json:"totalSlas"`
	CompliantSLAs       int     `json:"compliantSlas"`
	BreachedSLAs        int     `json:"breachedSlas"`
	AverageCompliance   float64 `json:"averageCompliance"`
	CriticalBreaches    int     `json:"criticalBreaches"`
	ImprovementTrend    string  `json:"improvementTrend"`
}

// MonitoringConfig represents the monitoring system configuration
type MonitoringConfig struct {
	MetricsRetention    time.Duration     `json:"metricsRetention"`
	SampleRate          time.Duration     `json:"sampleRate"`
	AlertingEnabled     bool              `json:"alertingEnabled"`
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
	SLADefinitions      []SLADefinition   `json:"slaDefinitions"`
	Dashboards          []Dashboard       `json:"dashboards"`
}

// NotificationChannel defines how alerts are delivered
type NotificationChannel struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"` // email, slack, webhook, pagerduty
	Config   map[string]interface{} `json:"config"`
	IsActive bool                   `json:"isActive"`
}