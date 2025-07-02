/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package certification

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
)

// CertificationLevel defines the quality certification levels for AI providers
type CertificationLevel string

const (
	CertificationGold      CertificationLevel = "gold"      // Premium tier: 95%+ accuracy, <2s latency, enterprise SLA
	CertificationSilver    CertificationLevel = "silver"    // Production tier: 90%+ accuracy, <5s latency, business SLA
	CertificationBronze    CertificationLevel = "bronze"    // Development tier: 80%+ accuracy, <10s latency, basic SLA
	CertificationUnverified CertificationLevel = "unverified" // Untested or failed certification
)

// QualityCertifier manages provider quality assessment and certification
type QualityCertifier struct {
	testSuites       map[string]*CertificationTestSuite
	certifications   map[string]*ProviderCertification
	benchmarkResults map[string]*BenchmarkResults
	config           CertificationConfig
	mutex            sync.RWMutex
}

// CertificationConfig defines configuration for the certification system
type CertificationConfig struct {
	MinTestCases          int           `json:"min_test_cases"`           // Minimum test cases for certification
	AccuracyThresholds    Thresholds    `json:"accuracy_thresholds"`      // Accuracy requirements per level
	LatencyThresholds     Thresholds    `json:"latency_thresholds"`       // Latency requirements per level
	ReliabilityThresholds Thresholds    `json:"reliability_thresholds"`   // Reliability requirements per level
	CertificationTimeout  time.Duration `json:"certification_timeout"`    // Maximum time for certification
	RecertificationPeriod time.Duration `json:"recertification_period"`   // How often to re-certify
	AutoFailover          bool          `json:"auto_failover"`             // Auto-demote on performance drop
}

// Thresholds defines quality thresholds for different certification levels
type Thresholds struct {
	Gold   float64 `json:"gold"`
	Silver float64 `json:"silver"`
	Bronze float64 `json:"bronze"`
}

// ProviderCertification represents a provider's certification status
type ProviderCertification struct {
	ProviderName      string                 `json:"provider_name"`
	Level             CertificationLevel     `json:"level"`
	Score             float64                `json:"score"`                // Overall quality score (0-100)
	AccuracyScore     float64                `json:"accuracy_score"`       // Accuracy percentage
	LatencyScore      float64                `json:"latency_score"`        // Latency score (lower is better)
	ReliabilityScore  float64                `json:"reliability_score"`    // Reliability percentage
	CostEfficiency    float64                `json:"cost_efficiency"`      // Cost per quality unit
	IssuedAt          time.Time              `json:"issued_at"`
	ExpiresAt         time.Time              `json:"expires_at"`
	TestResults       []TestResult           `json:"test_results"`
	Capabilities      []string               `json:"capabilities"`         // Certified capabilities
	Limitations       []string               `json:"limitations"`          // Known limitations
	SLACompliance     SLAMetrics             `json:"sla_compliance"`
	Metadata          map[string]string      `json:"metadata"`
}

// SLAMetrics tracks service level agreement compliance
type SLAMetrics struct {
	UptimePercentage     float64       `json:"uptime_percentage"`
	AverageResponseTime  time.Duration `json:"average_response_time"`
	ErrorRate            float64       `json:"error_rate"`
	RateLimitViolations  int           `json:"rate_limit_violations"`
	ServiceInterruptions int           `json:"service_interruptions"`
}

// CertificationTestSuite defines a comprehensive test suite for provider certification
type CertificationTestSuite struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	TestCases    []CertificationTest   `json:"test_cases"`
	PassingScore float64               `json:"passing_score"`    // Minimum score to pass
	Timeout      time.Duration         `json:"timeout"`
	Metadata     map[string]string     `json:"metadata"`
}

// CertificationTest defines a single test case for provider certification
type CertificationTest struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Category        TestCategory      `json:"category"`
	Input           TestInput         `json:"input"`
	ExpectedOutput  TestExpected      `json:"expected_output"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptance_criteria"`
	Weight          float64           `json:"weight"`           // Test weight in overall score
	Timeout         time.Duration     `json:"timeout"`
	Metadata        map[string]string `json:"metadata"`
}

// TestCategory defines categories of certification tests
type TestCategory string

const (
	CategoryAccuracy     TestCategory = "accuracy"      // Accuracy and correctness tests
	CategoryLatency      TestCategory = "latency"       // Performance and speed tests
	CategoryReliability  TestCategory = "reliability"   // Consistency and stability tests
	CategoryComplexity   TestCategory = "complexity"    // Complex reasoning tests
	CategorySafety       TestCategory = "safety"        // Safety and ethical tests
	CategorySpecialized  TestCategory = "specialized"   // Domain-specific tests
)

// TestInput defines the input for a certification test
type TestInput struct {
	Instruction      string            `json:"instruction"`
	Context          ai.TutorialContext `json:"context"`
	Parameters       map[string]interface{} `json:"parameters"`
	RequiredCapabilities []string      `json:"required_capabilities"`
}

// TestExpected defines the expected output for a certification test
type TestExpected struct {
	Type             ExpectedType      `json:"type"`
	Value            interface{}       `json:"value"`
	AccuracyThreshold float64          `json:"accuracy_threshold"`
	LatencyThreshold time.Duration     `json:"latency_threshold"`
	QualityThreshold float64           `json:"quality_threshold"`
}

// ExpectedType defines types of expected test outputs
type ExpectedType string

const (
	ExpectedExact      ExpectedType = "exact"       // Exact match required
	ExpectedSimilar    ExpectedType = "similar"     // Semantic similarity
	ExpectedPattern    ExpectedType = "pattern"     // Regex pattern match
	ExpectedStructured ExpectedType = "structured"  // Structured data validation
	ExpectedMetric     ExpectedType = "metric"      // Performance metric
)

// AcceptanceCriterion defines acceptance criteria for test results
type AcceptanceCriterion struct {
	Metric     string  `json:"metric"`
	Operator   string  `json:"operator"`  // >, <, >=, <=, ==, !=
	Threshold  float64 `json:"threshold"`
	Weight     float64 `json:"weight"`
	Required   bool    `json:"required"`  // Must pass for test to pass
}

// TestResult represents the result of a certification test
type TestResult struct {
	TestID           string        `json:"test_id"`
	ProviderName     string        `json:"provider_name"`
	Passed           bool          `json:"passed"`
	Score            float64       `json:"score"`
	ActualOutput     interface{}   `json:"actual_output"`
	ExecutionTime    time.Duration `json:"execution_time"`
	Cost             float64       `json:"cost"`
	ErrorMessage     string        `json:"error_message,omitempty"`
	Metrics          TestMetrics   `json:"metrics"`
	Timestamp        time.Time     `json:"timestamp"`
	Metadata         map[string]string `json:"metadata"`
}

// TestMetrics contains detailed metrics from test execution
type TestMetrics struct {
	Accuracy         float64 `json:"accuracy"`
	Latency          float64 `json:"latency"`
	Throughput       float64 `json:"throughput"`
	ResourceUsage    float64 `json:"resource_usage"`
	CostEfficiency   float64 `json:"cost_efficiency"`
	QualityScore     float64 `json:"quality_score"`
	ConsistencyScore float64 `json:"consistency_score"`
}

// BenchmarkResults contains comprehensive benchmark results for a provider
type BenchmarkResults struct {
	ProviderName         string                    `json:"provider_name"`
	TotalTests           int                       `json:"total_tests"`
	PassedTests          int                       `json:"passed_tests"`
	OverallScore         float64                   `json:"overall_score"`
	CategoryScores       map[TestCategory]float64  `json:"category_scores"`
	PerformanceMetrics   PerformanceMetrics        `json:"performance_metrics"`
	CostAnalysis         CostAnalysis              `json:"cost_analysis"`
	ReliabilityMetrics   ReliabilityMetrics        `json:"reliability_metrics"`
	ComparisonRanking    int                       `json:"comparison_ranking"`
	BenchmarkDate        time.Time                 `json:"benchmark_date"`
	Recommendations      []string                  `json:"recommendations"`
}

// PerformanceMetrics contains detailed performance analysis
type PerformanceMetrics struct {
	AverageLatency    time.Duration `json:"average_latency"`
	P95Latency        time.Duration `json:"p95_latency"`
	P99Latency        time.Duration `json:"p99_latency"`
	Throughput        float64       `json:"throughput"`        // Requests per second
	ErrorRate         float64       `json:"error_rate"`
	TimeoutRate       float64       `json:"timeout_rate"`
	ConcurrencyLimit  int           `json:"concurrency_limit"`
}

// CostAnalysis contains cost-related metrics
type CostAnalysis struct {
	AverageCostPerRequest float64 `json:"average_cost_per_request"`
	CostPerToken          float64 `json:"cost_per_token"`
	CostEfficiencyRatio   float64 `json:"cost_efficiency_ratio"`
	MonthlyEstimate       float64 `json:"monthly_estimate"`
	CostOptimizationTips  []string `json:"cost_optimization_tips"`
}

// ReliabilityMetrics contains reliability and consistency metrics
type ReliabilityMetrics struct {
	ConsistencyScore      float64 `json:"consistency_score"`
	ReproducibilityRate   float64 `json:"reproducibility_rate"`
	FailureRecoveryTime   time.Duration `json:"failure_recovery_time"`
	ServiceAvailability   float64 `json:"service_availability"`
	DataIntegrityScore    float64 `json:"data_integrity_score"`
}

// NewQualityCertifier creates a new quality certification system
func NewQualityCertifier(config CertificationConfig) *QualityCertifier {
	if config.MinTestCases == 0 {
		config.MinTestCases = 50
	}
	if config.CertificationTimeout == 0 {
		config.CertificationTimeout = 30 * time.Minute
	}
	if config.RecertificationPeriod == 0 {
		config.RecertificationPeriod = 30 * 24 * time.Hour // 30 days
	}

	// Set default thresholds
	if config.AccuracyThresholds.Gold == 0 {
		config.AccuracyThresholds = Thresholds{Gold: 95.0, Silver: 90.0, Bronze: 80.0}
	}
	if config.LatencyThresholds.Gold == 0 {
		config.LatencyThresholds = Thresholds{Gold: 2.0, Silver: 5.0, Bronze: 10.0}
	}
	if config.ReliabilityThresholds.Gold == 0 {
		config.ReliabilityThresholds = Thresholds{Gold: 99.0, Silver: 95.0, Bronze: 90.0}
	}

	return &QualityCertifier{
		testSuites:       make(map[string]*CertificationTestSuite),
		certifications:   make(map[string]*ProviderCertification),
		benchmarkResults: make(map[string]*BenchmarkResults),
		config:           config,
		mutex:            sync.RWMutex{},
	}
}

// RegisterTestSuite registers a new test suite for provider certification
func (c *QualityCertifier) RegisterTestSuite(suite *CertificationTestSuite) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if suite.Name == "" {
		return fmt.Errorf("test suite name cannot be empty")
	}

	if len(suite.TestCases) < c.config.MinTestCases {
		return fmt.Errorf("test suite must contain at least %d test cases", c.config.MinTestCases)
	}

	c.testSuites[suite.Name] = suite
	return nil
}

// CertifyProvider performs comprehensive certification of an AI provider
func (c *QualityCertifier) CertifyProvider(ctx context.Context, provider ai.Provider, providerName string) (*ProviderCertification, error) {
	start := time.Now()
	
	// Create certification context with timeout
	certCtx, cancel := context.WithTimeout(ctx, c.config.CertificationTimeout)
	defer cancel()

	certification := &ProviderCertification{
		ProviderName: providerName,
		IssuedAt:     start,
		ExpiresAt:    start.Add(c.config.RecertificationPeriod),
		TestResults:  make([]TestResult, 0),
		Capabilities: make([]string, 0),
		Limitations:  make([]string, 0),
		Metadata:     make(map[string]string),
	}

	// Run all registered test suites
	allTestResults := make([]TestResult, 0)
	categoryScores := make(map[TestCategory]float64)
	
	for suiteName, suite := range c.testSuites {
		fmt.Printf("Running test suite: %s\n", suiteName)
		
		suiteResults, err := c.runTestSuite(certCtx, provider, providerName, suite)
		if err != nil {
			return nil, fmt.Errorf("failed to run test suite %s: %w", suiteName, err)
		}
		
		allTestResults = append(allTestResults, suiteResults...)
		
		// Calculate category scores
		for _, result := range suiteResults {
			if test := c.findTestByID(result.TestID); test != nil {
				if _, exists := categoryScores[test.Category]; !exists {
					categoryScores[test.Category] = 0
				}
				categoryScores[test.Category] += result.Score * test.Weight
			}
		}
	}

	certification.TestResults = allTestResults

	// Calculate overall scores and metrics
	certification.Score = c.calculateOverallScore(allTestResults)
	certification.AccuracyScore = c.calculateCategoryScore(allTestResults, CategoryAccuracy)
	certification.LatencyScore = c.calculateLatencyScore(allTestResults)
	certification.ReliabilityScore = c.calculateCategoryScore(allTestResults, CategoryReliability)
	certification.CostEfficiency = c.calculateCostEfficiency(allTestResults)

	// Determine certification level
	certification.Level = c.determineCertificationLevel(certification)

	// Identify capabilities and limitations
	certification.Capabilities = c.identifyCapabilities(allTestResults)
	certification.Limitations = c.identifyLimitations(allTestResults)

	// Calculate SLA compliance metrics
	certification.SLACompliance = c.calculateSLAMetrics(allTestResults)

	// Store certification
	c.mutex.Lock()
	c.certifications[providerName] = certification
	c.mutex.Unlock()

	return certification, nil
}

// GetCertification retrieves the current certification for a provider
func (c *QualityCertifier) GetCertification(providerName string) (*ProviderCertification, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	cert, exists := c.certifications[providerName]
	if !exists {
		return nil, false
	}
	
	// Check if certification has expired
	if time.Now().After(cert.ExpiresAt) {
		return nil, false
	}
	
	return cert, true
}

// GetAllCertifications returns all current provider certifications
func (c *QualityCertifier) GetAllCertifications() map[string]*ProviderCertification {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	result := make(map[string]*ProviderCertification)
	now := time.Now()
	
	for name, cert := range c.certifications {
		if now.Before(cert.ExpiresAt) {
			result[name] = cert
		}
	}
	
	return result
}

// Helper methods for certification calculation

func (c *QualityCertifier) runTestSuite(ctx context.Context, provider ai.Provider, providerName string, suite *CertificationTestSuite) ([]TestResult, error) {
	results := make([]TestResult, 0, len(suite.TestCases))
	
	for _, test := range suite.TestCases {
		result, err := c.runSingleTest(ctx, provider, providerName, &test)
		if err != nil {
			// Log error but continue with other tests
			result = TestResult{
				TestID:       test.ID,
				ProviderName: providerName,
				Passed:       false,
				Score:        0,
				ErrorMessage: err.Error(),
				Timestamp:    time.Now(),
				Metadata:     make(map[string]string),
			}
		}
		results = append(results, result)
	}
	
	return results, nil
}

func (c *QualityCertifier) runSingleTest(ctx context.Context, provider ai.Provider, providerName string, test *CertificationTest) (TestResult, error) {
	start := time.Now()
	
	// Create test context with timeout
	testCtx, cancel := context.WithTimeout(ctx, test.Timeout)
	defer cancel()

	result := TestResult{
		TestID:       test.ID,
		ProviderName: providerName,
		Timestamp:    start,
		Metadata:     make(map[string]string),
	}

	// Execute the test based on category
	var output interface{}
	var err error
	
	switch test.Category {
	case CategoryAccuracy:
		output, err = c.testAccuracy(testCtx, provider, test)
	case CategoryLatency:
		output, err = c.testLatency(testCtx, provider, test)
	case CategoryReliability:
		output, err = c.testReliability(testCtx, provider, test)
	case CategoryComplexity:
		output, err = c.testComplexity(testCtx, provider, test)
	case CategorySafety:
		output, err = c.testSafety(testCtx, provider, test)
	default:
		err = fmt.Errorf("unsupported test category: %s", test.Category)
	}

	result.ExecutionTime = time.Since(start)
	result.ActualOutput = output

	if err != nil {
		result.ErrorMessage = err.Error()
		result.Passed = false
		result.Score = 0
		return result, nil
	}

	// Evaluate test results
	result.Passed, result.Score = c.evaluateTestResult(test, output, result.ExecutionTime)
	result.Metrics = c.calculateTestMetrics(test, output, result.ExecutionTime)

	return result, nil
}

func (c *QualityCertifier) testAccuracy(ctx context.Context, provider ai.Provider, test *CertificationTest) (interface{}, error) {
	// Test instruction parsing accuracy
	parsed, err := provider.ParseInstruction(ctx, test.Input.Instruction, test.Input.Context)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (c *QualityCertifier) testLatency(ctx context.Context, provider ai.Provider, test *CertificationTest) (interface{}, error) {
	start := time.Now()
	parsed, err := provider.ParseInstruction(ctx, test.Input.Instruction, test.Input.Context)
	latency := time.Since(start)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"output":  parsed,
		"latency": latency,
	}, nil
}

func (c *QualityCertifier) testReliability(ctx context.Context, provider ai.Provider, test *CertificationTest) (interface{}, error) {
	// Run the same test multiple times to test consistency
	const iterations = 5
	results := make([]interface{}, iterations)
	
	for i := 0; i < iterations; i++ {
		parsed, err := provider.ParseInstruction(ctx, test.Input.Instruction, test.Input.Context)
		if err != nil {
			return nil, fmt.Errorf("iteration %d failed: %w", i, err)
		}
		results[i] = parsed
	}
	
	return map[string]interface{}{
		"results":     results,
		"consistency": c.calculateConsistency(results),
	}, nil
}

func (c *QualityCertifier) testComplexity(ctx context.Context, provider ai.Provider, test *CertificationTest) (interface{}, error) {
	// Test complex reasoning capabilities
	parsed, err := provider.ParseInstruction(ctx, test.Input.Instruction, test.Input.Context)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (c *QualityCertifier) testSafety(ctx context.Context, provider ai.Provider, test *CertificationTest) (interface{}, error) {
	// Test safety and ethical guidelines
	parsed, err := provider.ParseInstruction(ctx, test.Input.Instruction, test.Input.Context)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (c *QualityCertifier) evaluateTestResult(test *CertificationTest, output interface{}, executionTime time.Duration) (bool, float64) {
	// Simple evaluation logic - would be more sophisticated in production
	
	// Check latency requirement
	if test.ExpectedOutput.LatencyThreshold > 0 && executionTime > test.ExpectedOutput.LatencyThreshold {
		return false, 0.0
	}
	
	// Check if output exists (basic test)
	if output == nil {
		return false, 0.0
	}
	
	// Calculate score based on acceptance criteria
	score := 0.0
	totalWeight := 0.0
	allRequiredPassed := true
	
	for _, criterion := range test.AcceptanceCriteria {
		passed := c.evaluateCriterion(criterion, output, executionTime)
		if criterion.Required && !passed {
			allRequiredPassed = false
		}
		if passed {
			score += criterion.Weight
		}
		totalWeight += criterion.Weight
	}
	
	if !allRequiredPassed {
		return false, 0.0
	}
	
	if totalWeight > 0 {
		score = (score / totalWeight) * 100
	} else {
		score = 100.0 // Default pass if no criteria
	}
	
	return score >= test.ExpectedOutput.AccuracyThreshold, score
}

func (c *QualityCertifier) evaluateCriterion(criterion AcceptanceCriterion, output interface{}, executionTime time.Duration) bool {
	// Simplified criterion evaluation
	switch criterion.Metric {
	case "latency":
		latencyMs := float64(executionTime.Milliseconds())
		return c.compareMetric(latencyMs, criterion.Operator, criterion.Threshold)
	case "output_length":
		if str, ok := output.(string); ok {
			return c.compareMetric(float64(len(str)), criterion.Operator, criterion.Threshold)
		}
	}
	
	return true // Default pass for unknown metrics
}

func (c *QualityCertifier) compareMetric(value float64, operator string, threshold float64) bool {
	switch operator {
	case ">":
		return value > threshold
	case "<":
		return value < threshold
	case ">=":
		return value >= threshold
	case "<=":
		return value <= threshold
	case "==":
		return value == threshold
	case "!=":
		return value != threshold
	default:
		return false
	}
}

func (c *QualityCertifier) calculateTestMetrics(test *CertificationTest, output interface{}, executionTime time.Duration) TestMetrics {
	// Calculate various metrics for the test result
	return TestMetrics{
		Accuracy:         85.0, // Placeholder - would calculate based on actual comparison
		Latency:          float64(executionTime.Milliseconds()),
		Throughput:       1000.0 / float64(executionTime.Milliseconds()), // Requests per second
		ResourceUsage:    0.1,   // Placeholder
		CostEfficiency:   0.95,  // Placeholder
		QualityScore:     88.0,  // Placeholder
		ConsistencyScore: 92.0,  // Placeholder
	}
}

func (c *QualityCertifier) calculateOverallScore(results []TestResult) float64 {
	if len(results) == 0 {
		return 0.0
	}
	
	totalScore := 0.0
	for _, result := range results {
		totalScore += result.Score
	}
	
	return totalScore / float64(len(results))
}

func (c *QualityCertifier) calculateCategoryScore(results []TestResult, category TestCategory) float64 {
	categoryResults := make([]TestResult, 0)
	for _, result := range results {
		if test := c.findTestByID(result.TestID); test != nil && test.Category == category {
			categoryResults = append(categoryResults, result)
		}
	}
	
	return c.calculateOverallScore(categoryResults)
}

func (c *QualityCertifier) calculateLatencyScore(results []TestResult) float64 {
	if len(results) == 0 {
		return 0.0
	}
	
	totalLatency := 0.0
	count := 0
	
	for _, result := range results {
		if result.ExecutionTime > 0 {
			totalLatency += float64(result.ExecutionTime.Milliseconds())
			count++
		}
	}
	
	if count == 0 {
		return 0.0
	}
	
	avgLatency := totalLatency / float64(count)
	
	// Convert to score (lower latency = higher score)
	// Score = 100 - (latency_ms / 100) with minimum of 0
	score := 100.0 - (avgLatency / 100.0)
	if score < 0 {
		score = 0
	}
	
	return score
}

func (c *QualityCertifier) calculateCostEfficiency(results []TestResult) float64 {
	if len(results) == 0 {
		return 0.0
	}
	
	totalCost := 0.0
	totalScore := 0.0
	
	for _, result := range results {
		totalCost += result.Cost
		totalScore += result.Score
	}
	
	if totalCost == 0 {
		return 100.0
	}
	
	// Cost efficiency = quality per dollar
	return (totalScore / float64(len(results))) / (totalCost * 1000) // Scale for readability
}

func (c *QualityCertifier) determineCertificationLevel(cert *ProviderCertification) CertificationLevel {
	// Check Gold requirements
	if cert.AccuracyScore >= c.config.AccuracyThresholds.Gold &&
		cert.LatencyScore >= (100.0-c.config.LatencyThresholds.Gold) &&
		cert.ReliabilityScore >= c.config.ReliabilityThresholds.Gold {
		return CertificationGold
	}
	
	// Check Silver requirements
	if cert.AccuracyScore >= c.config.AccuracyThresholds.Silver &&
		cert.LatencyScore >= (100.0-c.config.LatencyThresholds.Silver) &&
		cert.ReliabilityScore >= c.config.ReliabilityThresholds.Silver {
		return CertificationSilver
	}
	
	// Check Bronze requirements
	if cert.AccuracyScore >= c.config.AccuracyThresholds.Bronze &&
		cert.LatencyScore >= (100.0-c.config.LatencyThresholds.Bronze) &&
		cert.ReliabilityScore >= c.config.ReliabilityThresholds.Bronze {
		return CertificationBronze
	}
	
	return CertificationUnverified
}

func (c *QualityCertifier) identifyCapabilities(results []TestResult) []string {
	capabilities := make([]string, 0)
	categoryPassed := make(map[TestCategory]bool)
	
	for _, result := range results {
		if result.Passed {
			if test := c.findTestByID(result.TestID); test != nil {
				categoryPassed[test.Category] = true
			}
		}
	}
	
	for category, passed := range categoryPassed {
		if passed {
			capabilities = append(capabilities, string(category))
		}
	}
	
	return capabilities
}

func (c *QualityCertifier) identifyLimitations(results []TestResult) []string {
	limitations := make([]string, 0)
	categoryFailed := make(map[TestCategory]bool)
	
	for _, result := range results {
		if !result.Passed {
			if test := c.findTestByID(result.TestID); test != nil {
				categoryFailed[test.Category] = true
			}
		}
	}
	
	for category, failed := range categoryFailed {
		if failed {
			limitations = append(limitations, fmt.Sprintf("Limited %s capability", category))
		}
	}
	
	return limitations
}

func (c *QualityCertifier) calculateSLAMetrics(results []TestResult) SLAMetrics {
	if len(results) == 0 {
		return SLAMetrics{}
	}
	
	passedCount := 0
	totalLatency := time.Duration(0)
	errorCount := 0
	
	for _, result := range results {
		if result.Passed {
			passedCount++
		} else if result.ErrorMessage != "" {
			errorCount++
		}
		totalLatency += result.ExecutionTime
	}
	
	return SLAMetrics{
		UptimePercentage:    float64(passedCount) / float64(len(results)) * 100,
		AverageResponseTime: totalLatency / time.Duration(len(results)),
		ErrorRate:           float64(errorCount) / float64(len(results)) * 100,
		RateLimitViolations: 0, // Would track from actual errors
		ServiceInterruptions: 0, // Would track from monitoring
	}
}

func (c *QualityCertifier) calculateConsistency(results []interface{}) float64 {
	// Simplified consistency calculation
	// In production, this would compare semantic similarity of outputs
	if len(results) <= 1 {
		return 100.0
	}
	
	// For now, return a placeholder consistency score
	return 85.0
}

func (c *QualityCertifier) findTestByID(testID string) *CertificationTest {
	for _, suite := range c.testSuites {
		for i := range suite.TestCases {
			if suite.TestCases[i].ID == testID {
				return &suite.TestCases[i]
			}
		}
	}
	return nil
}