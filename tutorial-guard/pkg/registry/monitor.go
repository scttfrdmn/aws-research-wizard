/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package registry

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProviderMonitor manages health checks and performance monitoring for providers
type ProviderMonitor struct {
	registry        *ProviderRegistry
	healthInterval  time.Duration
	metricsInterval time.Duration
	alertThresholds AlertThresholds
	stopChan        chan struct{}
	wg              sync.WaitGroup
	callbacks       []MonitorCallback
	mutex           sync.RWMutex
}

// AlertThresholds defines when to trigger alerts
type AlertThresholds struct {
	ErrorRate        float64       `json:"error_rate"`        // Alert if error rate exceeds this
	LatencyThreshold time.Duration `json:"latency_threshold"` // Alert if latency exceeds this
	CostThreshold    float64       `json:"cost_threshold"`    // Alert if cost exceeds this
	SuccessRate      float64       `json:"success_rate"`      // Alert if success rate drops below this
}

// MonitorCallback defines callback functions for monitoring events
type MonitorCallback func(event MonitorEvent)

// MonitorEvent represents a monitoring event
type MonitorEvent struct {
	Type         EventType         `json:"type"`
	ProviderName string            `json:"provider_name"`
	Timestamp    time.Time         `json:"timestamp"`
	Message      string            `json:"message"`
	Severity     Severity          `json:"severity"`
	Metrics      *QualityMetrics   `json:"metrics,omitempty"`
	Metadata     map[string]string `json:"metadata"`
}

// EventType defines types of monitoring events
type EventType string

const (
	EventHealthCheck    EventType = "health_check"
	EventProviderDown   EventType = "provider_down"
	EventProviderUp     EventType = "provider_up"
	EventHighErrorRate  EventType = "high_error_rate"
	EventHighLatency    EventType = "high_latency"
	EventHighCost       EventType = "high_cost"
	EventLowSuccessRate EventType = "low_success_rate"
	EventMetricsUpdated EventType = "metrics_updated"
	EventRateLimited    EventType = "rate_limited"
	EventCircuitOpen    EventType = "circuit_open"
	EventCircuitClosed  EventType = "circuit_closed"
)

// Severity defines alert severity levels
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

// NewProviderMonitor creates a new provider monitoring system
func NewProviderMonitor(registry *ProviderRegistry) *ProviderMonitor {
	return &ProviderMonitor{
		registry:        registry,
		healthInterval:  30 * time.Second,
		metricsInterval: 5 * time.Minute,
		alertThresholds: AlertThresholds{
			ErrorRate:        0.1, // 10% error rate
			LatencyThreshold: 30 * time.Second,
			CostThreshold:    100.0, // $100/day
			SuccessRate:      0.9,   // 90% success rate
		},
		stopChan:  make(chan struct{}),
		callbacks: make([]MonitorCallback, 0),
		mutex:     sync.RWMutex{},
	}
}

// Start begins monitoring all registered providers
func (m *ProviderMonitor) Start(ctx context.Context) error {
	m.wg.Add(2)

	// Start health check routine
	go m.healthCheckRoutine(ctx)

	// Start metrics collection routine
	go m.metricsRoutine(ctx)

	m.emitEvent(MonitorEvent{
		Type:      EventHealthCheck,
		Timestamp: time.Now(),
		Message:   "Provider monitoring started",
		Severity:  SeverityInfo,
		Metadata:  map[string]string{"action": "start"},
	})

	return nil
}

// Stop gracefully stops the monitoring system
func (m *ProviderMonitor) Stop() error {
	close(m.stopChan)
	m.wg.Wait()

	m.emitEvent(MonitorEvent{
		Type:      EventHealthCheck,
		Timestamp: time.Now(),
		Message:   "Provider monitoring stopped",
		Severity:  SeverityInfo,
		Metadata:  map[string]string{"action": "stop"},
	})

	return nil
}

// AddCallback registers a callback for monitoring events
func (m *ProviderMonitor) AddCallback(callback MonitorCallback) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.callbacks = append(m.callbacks, callback)
}

// SetHealthInterval configures health check frequency
func (m *ProviderMonitor) SetHealthInterval(interval time.Duration) {
	m.healthInterval = interval
}

// SetMetricsInterval configures metrics collection frequency
func (m *ProviderMonitor) SetMetricsInterval(interval time.Duration) {
	m.metricsInterval = interval
}

// SetAlertThresholds configures alert thresholds
func (m *ProviderMonitor) SetAlertThresholds(thresholds AlertThresholds) {
	m.alertThresholds = thresholds
}

// healthCheckRoutine performs periodic health checks on all providers
func (m *ProviderMonitor) healthCheckRoutine(ctx context.Context) {
	defer m.wg.Done()

	ticker := time.NewTicker(m.healthInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.performHealthChecks(ctx)
		}
	}
}

// metricsRoutine performs periodic metrics collection and analysis
func (m *ProviderMonitor) metricsRoutine(ctx context.Context) {
	defer m.wg.Done()

	ticker := time.NewTicker(m.metricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.collectMetrics(ctx)
		}
	}
}

// performHealthChecks checks the health of all registered providers
func (m *ProviderMonitor) performHealthChecks(ctx context.Context) {
	m.registry.mutex.RLock()
	providers := make(map[string]*ProviderEntry)
	for name, entry := range m.registry.providers {
		providers[name] = entry
	}
	m.registry.mutex.RUnlock()

	for name, entry := range providers {
		go m.checkProviderHealth(ctx, name, entry)
	}
}

// checkProviderHealth performs a health check on a specific provider
func (m *ProviderMonitor) checkProviderHealth(ctx context.Context, name string, entry *ProviderEntry) {
	start := time.Now()

	// Create timeout context for health check
	healthCtx, cancel := context.WithTimeout(ctx, entry.Config.Timeout)
	defer cancel()

	// Perform health check
	err := entry.Provider.HealthCheck(healthCtx)
	duration := time.Since(start)

	// Update provider status
	m.registry.mutex.Lock()
	previousStatus := entry.Status.HealthStatus

	if err != nil {
		entry.Status.ConsecutiveErrors++
		entry.Status.HealthStatus = m.determineHealthStatus(entry.Status.ConsecutiveErrors)
		entry.Status.Available = entry.Status.HealthStatus != "unhealthy"

		// Emit provider down event if status changed
		if previousStatus != "unhealthy" && entry.Status.HealthStatus == "unhealthy" {
			m.emitEvent(MonitorEvent{
				Type:         EventProviderDown,
				ProviderName: name,
				Timestamp:    time.Now(),
				Message:      fmt.Sprintf("Provider %s is unhealthy: %v", name, err),
				Severity:     SeverityError,
				Metadata: map[string]string{
					"error":              err.Error(),
					"consecutive_errors": fmt.Sprintf("%d", entry.Status.ConsecutiveErrors),
					"duration_ms":        fmt.Sprintf("%.2f", float64(duration.Nanoseconds())/1e6),
				},
			})
		}
	} else {
		entry.Status.ConsecutiveErrors = 0
		entry.Status.HealthStatus = "healthy"
		entry.Status.Available = true
		entry.LastHealthy = time.Now()

		// Emit provider up event if status improved
		if previousStatus == "unhealthy" {
			m.emitEvent(MonitorEvent{
				Type:         EventProviderUp,
				ProviderName: name,
				Timestamp:    time.Now(),
				Message:      fmt.Sprintf("Provider %s is now healthy", name),
				Severity:     SeverityInfo,
				Metadata: map[string]string{
					"duration_ms": fmt.Sprintf("%.2f", float64(duration.Nanoseconds())/1e6),
				},
			})
		}
	}

	entry.Status.LastHealthCheck = time.Now()
	m.registry.mutex.Unlock()
}

// determineHealthStatus determines health status based on consecutive errors
func (m *ProviderMonitor) determineHealthStatus(consecutiveErrors int) string {
	switch {
	case consecutiveErrors == 0:
		return "healthy"
	case consecutiveErrors < 3:
		return "degraded"
	default:
		return "unhealthy"
	}
}

// collectMetrics collects and analyzes provider metrics
func (m *ProviderMonitor) collectMetrics(ctx context.Context) {
	m.registry.mutex.RLock()
	providers := make(map[string]*ProviderEntry)
	for name, entry := range m.registry.providers {
		providers[name] = entry
	}
	m.registry.mutex.RUnlock()

	for name, entry := range providers {
		m.analyzeProviderMetrics(name, entry)
	}
}

// analyzeProviderMetrics analyzes provider metrics and triggers alerts if needed
func (m *ProviderMonitor) analyzeProviderMetrics(name string, entry *ProviderEntry) {
	metrics := entry.QualityMetrics

	// Check error rate
	if metrics.ErrorRate > m.alertThresholds.ErrorRate {
		m.emitEvent(MonitorEvent{
			Type:         EventHighErrorRate,
			ProviderName: name,
			Timestamp:    time.Now(),
			Message:      fmt.Sprintf("High error rate detected: %.2f%% (threshold: %.2f%%)", metrics.ErrorRate*100, m.alertThresholds.ErrorRate*100),
			Severity:     SeverityWarning,
			Metrics:      &metrics,
			Metadata: map[string]string{
				"current_rate": fmt.Sprintf("%.4f", metrics.ErrorRate),
				"threshold":    fmt.Sprintf("%.4f", m.alertThresholds.ErrorRate),
			},
		})
	}

	// Check latency
	if metrics.AverageLatency > m.alertThresholds.LatencyThreshold {
		m.emitEvent(MonitorEvent{
			Type:         EventHighLatency,
			ProviderName: name,
			Timestamp:    time.Now(),
			Message:      fmt.Sprintf("High latency detected: %v (threshold: %v)", metrics.AverageLatency, m.alertThresholds.LatencyThreshold),
			Severity:     SeverityWarning,
			Metrics:      &metrics,
			Metadata: map[string]string{
				"current_latency": metrics.AverageLatency.String(),
				"threshold":       m.alertThresholds.LatencyThreshold.String(),
			},
		})
	}

	// Check success rate
	if metrics.SuccessRate < m.alertThresholds.SuccessRate {
		m.emitEvent(MonitorEvent{
			Type:         EventLowSuccessRate,
			ProviderName: name,
			Timestamp:    time.Now(),
			Message:      fmt.Sprintf("Low success rate detected: %.2f%% (threshold: %.2f%%)", metrics.SuccessRate*100, m.alertThresholds.SuccessRate*100),
			Severity:     SeverityWarning,
			Metrics:      &metrics,
			Metadata: map[string]string{
				"current_rate": fmt.Sprintf("%.4f", metrics.SuccessRate),
				"threshold":    fmt.Sprintf("%.4f", m.alertThresholds.SuccessRate),
			},
		})
	}

	// Check cost
	// Note: This would need access to cost tracking data
	// For now, emit a metrics updated event
	m.emitEvent(MonitorEvent{
		Type:         EventMetricsUpdated,
		ProviderName: name,
		Timestamp:    time.Now(),
		Message:      "Provider metrics updated",
		Severity:     SeverityInfo,
		Metrics:      &metrics,
		Metadata: map[string]string{
			"requests":     fmt.Sprintf("%d", metrics.RequestCount),
			"success_rate": fmt.Sprintf("%.4f", metrics.SuccessRate),
			"avg_latency":  metrics.AverageLatency.String(),
			"confidence":   fmt.Sprintf("%.4f", metrics.ConfidenceScore),
		},
	})
}

// emitEvent sends a monitoring event to all registered callbacks
func (m *ProviderMonitor) emitEvent(event MonitorEvent) {
	m.mutex.RLock()
	callbacks := make([]MonitorCallback, len(m.callbacks))
	copy(callbacks, m.callbacks)
	m.mutex.RUnlock()

	for _, callback := range callbacks {
		go callback(event)
	}
}

// RecordRequest records the result of a provider request for metrics tracking
func (m *ProviderMonitor) RecordRequest(providerName string, result RequestResult) {
	m.registry.mutex.Lock()
	defer m.registry.mutex.Unlock()

	entry, exists := m.registry.providers[providerName]
	if !exists {
		return
	}

	// Update metrics
	metrics := &entry.QualityMetrics
	metrics.RequestCount++

	// Update success rate using exponential moving average
	alpha := 0.1 // Smoothing factor
	if result.Success {
		metrics.SuccessRate = metrics.SuccessRate*(1-alpha) + alpha
		entry.Status.ConsecutiveErrors = 0
	} else {
		metrics.SuccessRate = metrics.SuccessRate * (1 - alpha)
		entry.Status.ConsecutiveErrors++
		metrics.ErrorRate = metrics.ErrorRate*(1-alpha) + alpha
	}

	// Update latency using exponential moving average
	if result.Duration > 0 {
		if metrics.AverageLatency == 0 {
			metrics.AverageLatency = result.Duration
		} else {
			avgNanos := float64(metrics.AverageLatency.Nanoseconds())
			newNanos := float64(result.Duration.Nanoseconds())
			updatedNanos := avgNanos*(1-alpha) + newNanos*alpha
			metrics.AverageLatency = time.Duration(updatedNanos)
		}
	}

	// Update cost
	if result.Cost > 0 {
		if metrics.AverageCost == 0 {
			metrics.AverageCost = result.Cost
		} else {
			metrics.AverageCost = metrics.AverageCost*(1-alpha) + result.Cost*alpha
		}
	}

	// Update confidence score
	if result.Confidence > 0 {
		if metrics.ConfidenceScore == 0 {
			metrics.ConfidenceScore = result.Confidence
		} else {
			metrics.ConfidenceScore = metrics.ConfidenceScore*(1-alpha) + result.Confidence*alpha
		}
	}

	metrics.LastUpdated = time.Now()

	// Check for rate limiting
	if result.RateLimited {
		entry.Status.RateLimited = true
		entry.Status.RateLimitReset = result.RateLimitReset

		m.emitEvent(MonitorEvent{
			Type:         EventRateLimited,
			ProviderName: providerName,
			Timestamp:    time.Now(),
			Message:      fmt.Sprintf("Provider %s is rate limited until %v", providerName, result.RateLimitReset),
			Severity:     SeverityWarning,
			Metadata: map[string]string{
				"reset_time": result.RateLimitReset.Format(time.RFC3339),
			},
		})
	}
}

// RequestResult represents the result of an AI provider request
type RequestResult struct {
	Success        bool          `json:"success"`
	Duration       time.Duration `json:"duration"`
	Cost           float64       `json:"cost"`
	Confidence     float64       `json:"confidence"`
	TokensUsed     int           `json:"tokens_used"`
	RateLimited    bool          `json:"rate_limited"`
	RateLimitReset time.Time     `json:"rate_limit_reset"`
	Error          error         `json:"error,omitempty"`
}

// GetProviderStatus returns the current status of all providers
func (m *ProviderMonitor) GetProviderStatus() map[string]ProviderStatus {
	m.registry.mutex.RLock()
	defer m.registry.mutex.RUnlock()

	status := make(map[string]ProviderStatus)
	for name, entry := range m.registry.providers {
		status[name] = entry.Status
	}

	return status
}

// GetProviderMetrics returns the current metrics for all providers
func (m *ProviderMonitor) GetProviderMetrics() map[string]QualityMetrics {
	m.registry.mutex.RLock()
	defer m.registry.mutex.RUnlock()

	metrics := make(map[string]QualityMetrics)
	for name, entry := range m.registry.providers {
		metrics[name] = entry.QualityMetrics
	}

	return metrics
}

// GetHealthSummary returns a summary of provider health
func (m *ProviderMonitor) GetHealthSummary() HealthSummary {
	m.registry.mutex.RLock()
	defer m.registry.mutex.RUnlock()

	summary := HealthSummary{
		TotalProviders:     len(m.registry.providers),
		HealthyProviders:   0,
		DegradedProviders:  0,
		UnhealthyProviders: 0,
		LastUpdated:        time.Now(),
	}

	for _, entry := range m.registry.providers {
		switch entry.Status.HealthStatus {
		case "healthy":
			summary.HealthyProviders++
		case "degraded":
			summary.DegradedProviders++
		case "unhealthy":
			summary.UnhealthyProviders++
		}
	}

	return summary
}

// HealthSummary provides an overview of provider health
type HealthSummary struct {
	TotalProviders     int       `json:"total_providers"`
	HealthyProviders   int       `json:"healthy_providers"`
	DegradedProviders  int       `json:"degraded_providers"`
	UnhealthyProviders int       `json:"unhealthy_providers"`
	LastUpdated        time.Time `json:"last_updated"`
}
