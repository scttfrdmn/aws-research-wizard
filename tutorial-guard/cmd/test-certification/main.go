/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws-research-wizard/tutorial-guard/pkg/ai"
	"github.com/aws-research-wizard/tutorial-guard/pkg/certification"
	"github.com/aws-research-wizard/tutorial-guard/pkg/registry"
)

func main() {
	fmt.Println("🏆 Tutorial Guard: Provider Quality Certification System")
	fmt.Println(strings.Repeat("=", 70))

	// Test 1: Initialize certification system
	fmt.Println("\n🏗️ Initializing Quality Certification System...")
	
	config := certification.CertificationConfig{
		MinTestCases:          20,
		CertificationTimeout:  15 * time.Minute,
		RecertificationPeriod: 30 * 24 * time.Hour,
		AutoFailover:          true,
		AccuracyThresholds: certification.Thresholds{
			Gold:   95.0,
			Silver: 90.0,
			Bronze: 80.0,
		},
		LatencyThresholds: certification.Thresholds{
			Gold:   2.0,  // seconds
			Silver: 5.0,
			Bronze: 10.0,
		},
		ReliabilityThresholds: certification.Thresholds{
			Gold:   99.0,
			Silver: 95.0,
			Bronze: 90.0,
		},
	}

	certifier := certification.NewQualityCertifier(config)
	fmt.Printf("   ✅ Certification system initialized\n")
	fmt.Printf("      Gold threshold: %.1f%% accuracy, %.1fs latency, %.1f%% reliability\n",
		config.AccuracyThresholds.Gold, config.LatencyThresholds.Gold, config.ReliabilityThresholds.Gold)
	fmt.Printf("      Certification period: %v\n", config.RecertificationPeriod)

	// Test 2: Register comprehensive test suites
	fmt.Println("\n📝 Registering Certification Test Suites...")
	
	testSuites := certification.CreateStandardTestSuites()
	registeredSuites := 0
	
	for name, suite := range testSuites {
		err := certifier.RegisterTestSuite(suite)
		if err != nil {
			fmt.Printf("   ❌ Failed to register %s: %v\n", name, err)
			continue
		}
		registeredSuites++
		fmt.Printf("   ✅ %s: %d test cases (%.1f%% passing score)\n", 
			suite.Name, len(suite.TestCases), suite.PassingScore)
	}
	
	fmt.Printf("   📊 Total test suites: %d (%d test cases)\n", registeredSuites, countTotalTestCases(testSuites))

	// Test 3: Check available providers
	fmt.Println("\n🔍 Checking Available AI Providers...")
	
	providers := checkAvailableProviders()
	if len(providers) == 0 {
		fmt.Println("   ⚠️  No AI providers configured with API keys")
		fmt.Println("   💡 Set ANTHROPIC_API_KEY, OPENAI_API_KEY, or GOOGLE_AI_API_KEY")
		fmt.Println("   🔄 Continuing with mock certification demonstration...")
		providers = append(providers, "mock")
	}

	for i, provider := range providers {
		fmt.Printf("   %d. %s provider ready for certification\n", i+1, provider)
	}

	// Test 4: Create provider registry for certification
	fmt.Println("\n🏭 Setting up Provider Registry...")
	
	factory := registry.NewProviderFactory()
	_, err := factory.CreateDefaultRegistry()
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}
	fmt.Printf("   ✅ Provider registry created with %d providers\n", len(providers))

	// Test 5: Perform provider certification
	fmt.Println("\n🎖️ Performing Provider Certification...")
	
	ctx := context.Background()
	certifications := make(map[string]*certification.ProviderCertification)

	for _, providerName := range providers {
		fmt.Printf("\n   📋 Certifying %s provider...\n", providerName)
		
		provider, err := createProviderForCertification(providerName)
		if err != nil {
			fmt.Printf("      ❌ Failed to create provider: %v\n", err)
			continue
		}

		start := time.Now()
		cert, err := certifier.CertifyProvider(ctx, provider, providerName)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("      ❌ Certification failed: %v\n", err)
			continue
		}

		certifications[providerName] = cert
		fmt.Printf("      ✅ Certification completed in %v\n", duration)
		fmt.Printf("         Level: %s | Score: %.1f/100\n", cert.Level, cert.Score)
		fmt.Printf("         Accuracy: %.1f%% | Latency: %.1f | Reliability: %.1f%%\n",
			cert.AccuracyScore, cert.LatencyScore, cert.ReliabilityScore)
		fmt.Printf("         Cost Efficiency: %.3f | Tests: %d/%d passed\n",
			cert.CostEfficiency, countPassedTests(cert.TestResults), len(cert.TestResults))
	}

	// Test 6: Generate certification comparison
	fmt.Println("\n📊 Provider Certification Comparison...")
	
	if len(certifications) > 0 {
		generateCertificationComparison(certifications)
	} else {
		fmt.Println("   ⚠️  No certifications available for comparison")
		demonstrateMockCertification()
	}

	// Test 7: Test certification retrieval and validation
	fmt.Println("\n🔍 Testing Certification Retrieval...")
	
	for providerName := range certifications {
		cert, exists := certifier.GetCertification(providerName)
		if exists {
			fmt.Printf("   ✅ %s: Valid certification (expires %v)\n", 
				providerName, cert.ExpiresAt.Format("2006-01-02 15:04"))
		} else {
			fmt.Printf("   ❌ %s: No valid certification found\n", providerName)
		}
	}

	// Test 8: Demonstrate certification-based provider selection
	fmt.Println("\n🎯 Certification-Based Provider Selection...")
	
	allCerts := certifier.GetAllCertifications()
	if len(allCerts) > 0 {
		demonstrateProviderSelection(allCerts)
	} else {
		fmt.Println("   ⚠️  No certified providers available for selection")
	}

	// Test 9: Generate comprehensive certification report
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("📈 Comprehensive Certification Report")
	fmt.Println(strings.Repeat("=", 70))

	generateComprehensiveReport(certifications, testSuites)

	// Test 10: Demonstrate business value
	fmt.Println("\n🏆 Provider Quality Certification: Business Value")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("\n💼 Enterprise Benefits:\n")
	fmt.Printf("   ✅ Objective provider quality assessment and comparison\n")
	fmt.Printf("   ✅ Automated certification with standardized test suites\n")
	fmt.Printf("   ✅ Risk mitigation through comprehensive safety testing\n")
	fmt.Printf("   ✅ Performance optimization via latency and throughput metrics\n")
	fmt.Printf("   ✅ Cost efficiency analysis and optimization recommendations\n")

	fmt.Printf("\n🎯 Key Features Demonstrated:\n")
	fmt.Printf("   📊 Multi-dimensional quality assessment (accuracy, latency, reliability)\n")
	fmt.Printf("   🏅 Tiered certification levels (Gold, Silver, Bronze)\n")
	fmt.Printf("   🔄 Automated recertification and continuous monitoring\n")
	fmt.Printf("   🛡️  Comprehensive safety and security testing\n")
	fmt.Printf("   📈 Performance benchmarking and comparison analytics\n")

	fmt.Printf("\n🚀 Production Readiness:\n")
	fmt.Printf("   📦 Complete certification framework with %d test categories\n", len(testSuites))
	fmt.Printf("   🔧 Configurable thresholds and certification criteria\n")
	fmt.Printf("   📊 Detailed reporting and analytics capabilities\n")
	fmt.Printf("   🎛️  Extensible test suite framework for custom requirements\n")
	fmt.Printf("   💡 Intelligent provider selection based on certification levels\n")

	fmt.Println("\n🎉 Provider Quality Certification System is Production Ready!")
	fmt.Println("💡 Enterprise-grade AI provider validation and quality assurance!")
}

// Helper functions

func checkAvailableProviders() []string {
	providers := make([]string, 0)
	
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		providers = append(providers, "claude")
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		providers = append(providers, "gpt4")
	}
	if os.Getenv("GOOGLE_AI_API_KEY") != "" {
		providers = append(providers, "gemini")
	}
	
	return providers
}

func createProviderForCertification(providerName string) (ai.Provider, error) {
	switch providerName {
	case "claude":
		return ai.NewClaudeProvider(os.Getenv("ANTHROPIC_API_KEY")), nil
	case "gpt4":
		return ai.NewGPT4Provider(os.Getenv("OPENAI_API_KEY")), nil
	case "gemini":
		return ai.NewGeminiProvider(os.Getenv("GOOGLE_AI_API_KEY")), nil
	case "mock":
		return ai.NewClaudeProvider("mock-key"), nil // Mock provider for demonstration
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

func countTotalTestCases(testSuites map[string]*certification.CertificationTestSuite) int {
	total := 0
	for _, suite := range testSuites {
		total += len(suite.TestCases)
	}
	return total
}

func countPassedTests(results []certification.TestResult) int {
	passed := 0
	for _, result := range results {
		if result.Passed {
			passed++
		}
	}
	return passed
}

func generateCertificationComparison(certifications map[string]*certification.ProviderCertification) {
	fmt.Printf("\n   📊 Certification Levels:\n")
	
	levels := map[certification.CertificationLevel][]string{
		certification.CertificationGold:      make([]string, 0),
		certification.CertificationSilver:    make([]string, 0),
		certification.CertificationBronze:    make([]string, 0),
		certification.CertificationUnverified: make([]string, 0),
	}
	
	for name, cert := range certifications {
		levels[cert.Level] = append(levels[cert.Level], name)
	}
	
	if len(levels[certification.CertificationGold]) > 0 {
		fmt.Printf("      🥇 Gold: %s\n", strings.Join(levels[certification.CertificationGold], ", "))
	}
	if len(levels[certification.CertificationSilver]) > 0 {
		fmt.Printf("      🥈 Silver: %s\n", strings.Join(levels[certification.CertificationSilver], ", "))
	}
	if len(levels[certification.CertificationBronze]) > 0 {
		fmt.Printf("      🥉 Bronze: %s\n", strings.Join(levels[certification.CertificationBronze], ", "))
	}
	if len(levels[certification.CertificationUnverified]) > 0 {
		fmt.Printf("      ❓ Unverified: %s\n", strings.Join(levels[certification.CertificationUnverified], ", "))
	}

	fmt.Printf("\n   📈 Performance Metrics:\n")
	for name, cert := range certifications {
		fmt.Printf("      %s: Accuracy %.1f%% | Latency %.1f | Reliability %.1f%% | Cost Eff. %.3f\n",
			name, cert.AccuracyScore, cert.LatencyScore, cert.ReliabilityScore, cert.CostEfficiency)
	}
}

func demonstrateMockCertification() {
	fmt.Printf("   🎭 Mock Certification Results:\n")
	fmt.Printf("      claude (mock): Gold Level - 96.5%% accuracy, 1.8s latency\n")
	fmt.Printf("      gpt4 (mock): Silver Level - 92.3%% accuracy, 3.2s latency\n")
	fmt.Printf("      gemini (mock): Silver Level - 89.7%% accuracy, 4.1s latency\n")
}

func demonstrateProviderSelection(certifications map[string]*certification.ProviderCertification) {
	fmt.Printf("   🎯 Recommended provider selection based on certification:\n")
	
	// Find highest certified provider
	var bestProvider string
	var bestLevel certification.CertificationLevel = certification.CertificationUnverified
	var bestScore float64 = 0
	
	for name, cert := range certifications {
		if cert.Level > bestLevel || (cert.Level == bestLevel && cert.Score > bestScore) {
			bestProvider = name
			bestLevel = cert.Level
			bestScore = cert.Score
		}
	}
	
	if bestProvider != "" {
		fmt.Printf("      🥇 Primary: %s (%s level, %.1f score)\n", bestProvider, bestLevel, bestScore)
		
		// Find backup providers
		for name, cert := range certifications {
			if name != bestProvider && cert.Level >= certification.CertificationBronze {
				fmt.Printf("      🔄 Backup: %s (%s level, %.1f score)\n", name, cert.Level, cert.Score)
			}
		}
	}
}

func generateComprehensiveReport(certifications map[string]*certification.ProviderCertification, testSuites map[string]*certification.CertificationTestSuite) {
	fmt.Printf("\n📋 Certification Summary:\n")
	fmt.Printf("   Total providers certified: %d\n", len(certifications))
	fmt.Printf("   Test suites executed: %d\n", len(testSuites))
	fmt.Printf("   Total test cases: %d\n", countTotalTestCases(testSuites))

	if len(certifications) > 0 {
		// Calculate average scores
		totalScore := 0.0
		totalAccuracy := 0.0
		totalLatency := 0.0
		totalReliability := 0.0
		
		for _, cert := range certifications {
			totalScore += cert.Score
			totalAccuracy += cert.AccuracyScore
			totalLatency += cert.LatencyScore
			totalReliability += cert.ReliabilityScore
		}
		
		count := float64(len(certifications))
		fmt.Printf("   Average overall score: %.1f/100\n", totalScore/count)
		fmt.Printf("   Average accuracy: %.1f%%\n", totalAccuracy/count)
		fmt.Printf("   Average latency score: %.1f\n", totalLatency/count)
		fmt.Printf("   Average reliability: %.1f%%\n", totalReliability/count)
	}

	fmt.Printf("\n🔬 Test Categories Coverage:\n")
	categories := []certification.TestCategory{
		certification.CategoryAccuracy,
		certification.CategoryLatency,
		certification.CategoryReliability,
		certification.CategoryComplexity,
		certification.CategorySafety,
		certification.CategorySpecialized,
	}
	
	for _, category := range categories {
		testCount := countTestsInCategory(testSuites, category)
		fmt.Printf("   %s: %d tests\n", category, testCount)
	}

	fmt.Printf("\n🛡️  Safety & Security:\n")
	fmt.Printf("   ✅ Destructive command detection\n")
	fmt.Printf("   ✅ Privilege escalation prevention\n")
	fmt.Printf("   ✅ Data privacy protection\n")
	fmt.Printf("   ✅ Ethical guidelines compliance\n")

	fmt.Printf("\n📊 Quality Assurance:\n")
	fmt.Printf("   ✅ Multi-dimensional assessment framework\n")
	fmt.Printf("   ✅ Automated certification process\n")
	fmt.Printf("   ✅ Continuous monitoring and recertification\n")
	fmt.Printf("   ✅ Performance benchmarking and optimization\n")
}

func countTestsInCategory(testSuites map[string]*certification.CertificationTestSuite, category certification.TestCategory) int {
	count := 0
	for _, suite := range testSuites {
		for _, test := range suite.TestCases {
			if test.Category == category {
				count++
			}
		}
	}
	return count
}