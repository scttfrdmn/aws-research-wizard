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
	"sort"
	"sync"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
)

// ProviderRegistry manages multiple AI providers and intelligent routing
type ProviderRegistry struct {
	providers      map[string]*ProviderEntry
	routingConfig  RoutingConfig
	qualityTracker *QualityTracker
	mutex          sync.RWMutex
}

// ProviderEntry represents a registered AI provider with metadata
type ProviderEntry struct {
	Provider       ai.Provider
	Config         ProviderConfig
	Status         ProviderStatus
	QualityMetrics QualityMetrics
	LastHealthy    time.Time
	CreatedAt      time.Time
}

// ProviderConfig defines provider-specific configuration
type ProviderConfig struct {
	Name          string            `json:"name"`
	Priority      int               `json:"priority"`       // Higher number = higher priority
	Weight        float64           `json:"weight"`         // For load balancing (0.0-1.0)
	MaxConcurrent int               `json:"max_concurrent"` // Max concurrent requests
	Timeout       time.Duration     `json:"timeout"`        // Request timeout
	RetryPolicy   RetryPolicy       `json:"retry_policy"`   // Retry configuration
	CostLimit     CostLimit         `json:"cost_limit"`     // Cost controls
	Capabilities  []string          `json:"capabilities"`   // Supported features
	Regions       []string          `json:"regions"`        // Supported regions
	Metadata      map[string]string `json:"metadata"`
}

// ProviderStatus tracks the current state of a provider
type ProviderStatus struct {
	Available         bool      `json:"available"`
	HealthStatus      string    `json:"health_status"` // healthy, degraded, unhealthy
	CurrentLoad       int       `json:"current_load"`  // Current concurrent requests
	LastHealthCheck   time.Time `json:"last_health_check"`
	ConsecutiveErrors int       `json:"consecutive_errors"`
	RateLimited       bool      `json:"rate_limited"`
	RateLimitReset    time.Time `json:"rate_limit_reset"`
}

// QualityMetrics tracks provider performance and quality
type QualityMetrics struct {
	RequestCount    int64         `json:"request_count"`
	SuccessRate     float64       `json:"success_rate"`
	AverageLatency  time.Duration `json:"average_latency"`
	AverageCost     float64       `json:"average_cost"`
	ConfidenceScore float64       `json:"confidence_score"` // Average AI confidence
	AccuracyScore   float64       `json:"accuracy_score"`   // Validation accuracy
	TokenEfficiency float64       `json:"token_efficiency"` // Output quality per token
	ErrorRate       float64       `json:"error_rate"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// RoutingConfig defines how requests are routed to providers
type RoutingConfig struct {
	DefaultStrategy  RoutingStrategy      `json:"default_strategy"`
	FallbackChain    []string             `json:"fallback_chain"` // Provider names in fallback order
	LoadBalancing    LoadBalancingMode    `json:"load_balancing"`
	QualityThreshold float64              `json:"quality_threshold"` // Minimum quality score
	CostOptimization bool                 `json:"cost_optimization"` // Prefer lower cost providers
	LatencyThreshold time.Duration        `json:"latency_threshold"` // Maximum acceptable latency
	MaxRetries       int                  `json:"max_retries"`
	CircuitBreaker   CircuitBreakerConfig `json:"circuit_breaker"`
}

// RoutingStrategy defines how to select providers
type RoutingStrategy string

const (
	StrategyPriority     RoutingStrategy = "priority"      // Use highest priority available
	StrategyRoundRobin   RoutingStrategy = "round_robin"   // Distribute evenly
	StrategyWeighted     RoutingStrategy = "weighted"      // Use provider weights
	StrategyCostOptimal  RoutingStrategy = "cost_optimal"  // Minimize cost
	StrategyQualityFirst RoutingStrategy = "quality_first" // Maximize quality
	StrategyLatencyFirst RoutingStrategy = "latency_first" // Minimize latency
	StrategyIntelligent  RoutingStrategy = "intelligent"   // AI-driven routing decisions
)

// LoadBalancingMode defines load balancing behavior
type LoadBalancingMode string

const (
	LoadBalancingRoundRobin LoadBalancingMode = "round_robin"
	LoadBalancingWeighted   LoadBalancingMode = "weighted"
	LoadBalancingLeastLoad  LoadBalancingMode = "least_load"
	LoadBalancingResponse   LoadBalancingMode = "response_time"
)

// RetryPolicy defines retry behavior for failed requests
type RetryPolicy struct {
	MaxRetries      int           `json:"max_retries"`
	InitialDelay    time.Duration `json:"initial_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	BackoffFactor   float64       `json:"backoff_factor"`
	RetryableErrors []string      `json:"retryable_errors"`
}

// CostLimit defines cost controls for a provider
type CostLimit struct {
	DailyCostLimit   float64 `json:"daily_cost_limit"`
	MonthlyCostLimit float64 `json:"monthly_cost_limit"`
	PerRequestLimit  float64 `json:"per_request_limit"`
	AlertThreshold   float64 `json:"alert_threshold"`
}

// CircuitBreakerConfig defines circuit breaker behavior
type CircuitBreakerConfig struct {
	Enabled          bool          `json:"enabled"`
	FailureThreshold int           `json:"failure_threshold"` // Failures before opening
	RecoveryTimeout  time.Duration `json:"recovery_timeout"`  // Time before trying again
	SuccessThreshold int           `json:"success_threshold"` // Successes to close circuit
}

// QualityTracker monitors and learns from provider performance
type QualityTracker struct {
	metrics map[string]*QualityMetrics
	window  time.Duration // Time window for metrics
	decay   float64       // Decay factor for older metrics
	mutex   sync.RWMutex
}

// RoutingRequest represents a request for AI provider routing
type RoutingRequest struct {
	Type         ai.RequestType     `json:"type"`
	Priority     ai.Priority        `json:"priority"`
	MaxCost      float64            `json:"max_cost"`
	MaxLatency   time.Duration      `json:"max_latency"`
	RequiredCaps []string           `json:"required_capabilities"`
	Region       string             `json:"region"`
	Context      ai.TutorialContext `json:"context"`
	Metadata     map[string]string  `json:"metadata"`
}

// RoutingResult contains the selected provider and routing metadata
type RoutingResult struct {
	Provider         ai.Provider       `json:"-"` // Selected provider
	ProviderName     string            `json:"provider_name"`
	Reason           string            `json:"reason"`       // Why this provider was selected
	Alternatives     []string          `json:"alternatives"` // Other viable providers
	EstimatedCost    float64           `json:"estimated_cost"`
	EstimatedLatency time.Duration     `json:"estimated_latency"`
	QualityScore     float64           `json:"quality_score"`
	RoutingTime      time.Duration     `json:"routing_time"`
	Metadata         map[string]string `json:"metadata"`
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry(config RoutingConfig) *ProviderRegistry {
	return &ProviderRegistry{
		providers:      make(map[string]*ProviderEntry),
		routingConfig:  config,
		qualityTracker: NewQualityTracker(24*time.Hour, 0.95), // 24h window, 5% decay
		mutex:          sync.RWMutex{},
	}
}

// Register adds a new AI provider to the registry
func (r *ProviderRegistry) Register(name string, provider ai.Provider, config ProviderConfig) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	entry := &ProviderEntry{
		Provider: provider,
		Config:   config,
		Status: ProviderStatus{
			Available:         true,
			HealthStatus:      "unknown",
			CurrentLoad:       0,
			LastHealthCheck:   time.Time{},
			ConsecutiveErrors: 0,
			RateLimited:       false,
		},
		QualityMetrics: QualityMetrics{
			RequestCount:    0,
			SuccessRate:     1.0, // Start optimistic
			ConfidenceScore: 0.8, // Default confidence
			AccuracyScore:   0.8, // Default accuracy
			TokenEfficiency: 0.8, // Default efficiency
			ErrorRate:       0.0,
			LastUpdated:     time.Now(),
		},
		LastHealthy: time.Now(),
		CreatedAt:   time.Now(),
	}

	r.providers[name] = entry
	return nil
}

// Route selects the best provider for a given request
func (r *ProviderRegistry) Route(ctx context.Context, request RoutingRequest) (*RoutingResult, error) {
	start := time.Now()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get available providers that meet requirements
	candidates := r.getCandidateProviders(request)
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no available providers meet requirements")
	}

	// Apply routing strategy
	selected, reason := r.applyRoutingStrategy(candidates, request)
	if selected == nil {
		return nil, fmt.Errorf("routing strategy failed to select provider")
	}

	// Build routing result
	result := &RoutingResult{
		Provider:         selected.Provider,
		ProviderName:     selected.Config.Name,
		Reason:           reason,
		Alternatives:     r.getAlternativeNames(candidates, selected.Config.Name),
		EstimatedCost:    r.estimateCost(selected, request),
		EstimatedLatency: r.estimateLatency(selected, request),
		QualityScore:     selected.QualityMetrics.AccuracyScore,
		RoutingTime:      time.Since(start),
		Metadata:         make(map[string]string),
	}

	// Add routing metadata
	result.Metadata["strategy"] = string(r.routingConfig.DefaultStrategy)
	result.Metadata["candidates_count"] = fmt.Sprintf("%d", len(candidates))
	result.Metadata["selection_time_ms"] = fmt.Sprintf("%.2f", float64(result.RoutingTime.Nanoseconds())/1e6)

	return result, nil
}

// getCandidateProviders filters providers based on request requirements
func (r *ProviderRegistry) getCandidateProviders(request RoutingRequest) []*ProviderEntry {
	var candidates []*ProviderEntry

	for _, entry := range r.providers {
		if !r.isProviderViable(entry, request) {
			continue
		}
		candidates = append(candidates, entry)
	}

	return candidates
}

// isProviderViable checks if a provider meets basic requirements
func (r *ProviderRegistry) isProviderViable(entry *ProviderEntry, request RoutingRequest) bool {
	// Check availability
	if !entry.Status.Available || entry.Status.HealthStatus == "unhealthy" {
		return false
	}

	// Check rate limiting
	if entry.Status.RateLimited && time.Now().Before(entry.Status.RateLimitReset) {
		return false
	}

	// Check concurrent load
	if entry.Status.CurrentLoad >= entry.Config.MaxConcurrent {
		return false
	}

	// Check quality threshold
	if entry.QualityMetrics.AccuracyScore < r.routingConfig.QualityThreshold {
		return false
	}

	// Check cost limits
	if request.MaxCost > 0 {
		estimatedCost := r.estimateCost(entry, request)
		if estimatedCost > request.MaxCost {
			return false
		}
	}

	// Check latency requirements
	if request.MaxLatency > 0 {
		estimatedLatency := r.estimateLatency(entry, request)
		if estimatedLatency > request.MaxLatency {
			return false
		}
	}

	// Check required capabilities
	if len(request.RequiredCaps) > 0 {
		providerCaps := make(map[string]bool)
		for _, cap := range entry.Config.Capabilities {
			providerCaps[cap] = true
		}

		for _, requiredCap := range request.RequiredCaps {
			if !providerCaps[requiredCap] {
				return false
			}
		}
	}

	// Check region support
	if request.Region != "" {
		regionSupported := false
		for _, region := range entry.Config.Regions {
			if region == request.Region || region == "global" {
				regionSupported = true
				break
			}
		}
		if !regionSupported {
			return false
		}
	}

	return true
}

// applyRoutingStrategy selects a provider based on the configured strategy
func (r *ProviderRegistry) applyRoutingStrategy(candidates []*ProviderEntry, request RoutingRequest) (*ProviderEntry, string) {
	if len(candidates) == 0 {
		return nil, "no candidates available"
	}

	switch r.routingConfig.DefaultStrategy {
	case StrategyPriority:
		return r.selectByPriority(candidates), "highest priority provider"

	case StrategyCostOptimal:
		return r.selectByCost(candidates, request), "lowest cost provider"

	case StrategyQualityFirst:
		return r.selectByQuality(candidates), "highest quality provider"

	case StrategyLatencyFirst:
		return r.selectByLatency(candidates, request), "lowest latency provider"

	case StrategyWeighted:
		return r.selectByWeight(candidates), "weighted random selection"

	case StrategyIntelligent:
		return r.selectIntelligently(candidates, request), "AI-optimized selection"

	default:
		// Fallback to priority
		return r.selectByPriority(candidates), "fallback to priority"
	}
}

// selectByPriority selects the provider with highest priority
func (r *ProviderRegistry) selectByPriority(candidates []*ProviderEntry) *ProviderEntry {
	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Config.Priority > candidates[j].Config.Priority
	})

	return candidates[0]
}

// selectByCost selects the provider with lowest estimated cost
func (r *ProviderRegistry) selectByCost(candidates []*ProviderEntry, request RoutingRequest) *ProviderEntry {
	if len(candidates) == 0 {
		return nil
	}

	var best *ProviderEntry
	var bestCost float64

	for _, candidate := range candidates {
		cost := r.estimateCost(candidate, request)
		if best == nil || cost < bestCost {
			best = candidate
			bestCost = cost
		}
	}

	return best
}

// selectByQuality selects the provider with highest quality score
func (r *ProviderRegistry) selectByQuality(candidates []*ProviderEntry) *ProviderEntry {
	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].QualityMetrics.AccuracyScore > candidates[j].QualityMetrics.AccuracyScore
	})

	return candidates[0]
}

// selectByLatency selects the provider with lowest estimated latency
func (r *ProviderRegistry) selectByLatency(candidates []*ProviderEntry, request RoutingRequest) *ProviderEntry {
	if len(candidates) == 0 {
		return nil
	}

	var best *ProviderEntry
	var bestLatency time.Duration

	for _, candidate := range candidates {
		latency := r.estimateLatency(candidate, request)
		if best == nil || latency < bestLatency {
			best = candidate
			bestLatency = latency
		}
	}

	return best
}

// selectByWeight performs weighted random selection
func (r *ProviderRegistry) selectByWeight(candidates []*ProviderEntry) *ProviderEntry {
	if len(candidates) == 0 {
		return nil
	}

	// For now, use simple round-robin as weighted selection requires randomization
	// In a production system, this would implement proper weighted random selection
	return candidates[0]
}

// selectIntelligently uses AI-driven selection logic
func (r *ProviderRegistry) selectIntelligently(candidates []*ProviderEntry, request RoutingRequest) *ProviderEntry {
	if len(candidates) == 0 {
		return nil
	}

	// Intelligent selection considers multiple factors:
	// - Request type and complexity
	// - Provider strengths for specific tasks
	// - Current load and performance
	// - Cost-quality trade-offs

	var best *ProviderEntry
	var bestScore float64

	for _, candidate := range candidates {
		score := r.calculateIntelligentScore(candidate, request)
		if best == nil || score > bestScore {
			best = candidate
			bestScore = score
		}
	}

	return best
}

// calculateIntelligentScore computes a composite score for intelligent routing
func (r *ProviderRegistry) calculateIntelligentScore(entry *ProviderEntry, request RoutingRequest) float64 {
	// Base quality score (40% weight)
	qualityScore := entry.QualityMetrics.AccuracyScore * 0.4

	// Performance score (30% weight)
	latencyScore := 1.0 - (float64(entry.QualityMetrics.AverageLatency.Milliseconds()) / 10000.0) // Normalize to 0-1
	if latencyScore < 0 {
		latencyScore = 0
	}
	performanceScore := (entry.QualityMetrics.SuccessRate*0.7 + latencyScore*0.3) * 0.3

	// Cost efficiency score (20% weight)
	// Lower cost = higher score
	costScore := (1.0 - entry.QualityMetrics.AverageCost/0.01) * 0.2 // Assume max reasonable cost of $0.01
	if costScore < 0 {
		costScore = 0
	}

	// Load balancing score (10% weight)
	loadScore := (1.0 - float64(entry.Status.CurrentLoad)/float64(entry.Config.MaxConcurrent)) * 0.1

	return qualityScore + performanceScore + costScore + loadScore
}

// Helper functions for cost and latency estimation
func (r *ProviderRegistry) estimateCost(entry *ProviderEntry, request RoutingRequest) float64 {
	// Simple estimation based on request type and provider metrics
	baseCost := entry.QualityMetrics.AverageCost

	// Adjust based on request type
	switch request.Type {
	case ai.RequestParseInstruction:
		return baseCost * 1.5 // Instructions are more complex
	case ai.RequestValidateExpectation:
		return baseCost * 0.8 // Validation is simpler
	case ai.RequestCompressContext:
		return baseCost * 0.5 // Compression is lightweight
	case ai.RequestInterpretError:
		return baseCost * 1.2 // Error interpretation is moderately complex
	default:
		return baseCost
	}
}

func (r *ProviderRegistry) estimateLatency(entry *ProviderEntry, request RoutingRequest) time.Duration {
	// Base latency from metrics
	baseLatency := entry.QualityMetrics.AverageLatency

	// Adjust for current load
	loadFactor := 1.0 + (float64(entry.Status.CurrentLoad)/float64(entry.Config.MaxConcurrent))*0.5

	return time.Duration(float64(baseLatency.Nanoseconds()) * loadFactor)
}

// getAlternativeNames returns names of alternative providers
func (r *ProviderRegistry) getAlternativeNames(candidates []*ProviderEntry, selected string) []string {
	var alternatives []string
	for _, candidate := range candidates {
		if candidate.Config.Name != selected {
			alternatives = append(alternatives, candidate.Config.Name)
		}
	}
	return alternatives
}

// NewQualityTracker creates a new quality tracking system
func NewQualityTracker(window time.Duration, decay float64) *QualityTracker {
	return &QualityTracker{
		metrics: make(map[string]*QualityMetrics),
		window:  window,
		decay:   decay,
		mutex:   sync.RWMutex{},
	}
}
