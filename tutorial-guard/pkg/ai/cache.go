/*
 * Tutorial Guard: AI-Powered Documentation Validation
 * Copyright © 2025 Scott Friedman. All rights reserved.
 *
 * This software is proprietary and confidential. Unauthorized copying,
 * distribution, or use is strictly prohibited.
 */

package ai

import (
	"sync"
	"time"
)

// ResponseCache implements an in-memory cache for AI responses
type ResponseCache struct {
	cache map[string]*CacheEntry
	mutex sync.RWMutex
}

// NewResponseCache creates a new response cache
func NewResponseCache() *ResponseCache {
	cache := &ResponseCache{
		cache: make(map[string]*CacheEntry),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Get retrieves a cached response
func (c *ResponseCache) Get(key string) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		// Remove expired entry (will be cleaned up by background goroutine)
		return nil
	}

	// Increment hit count
	entry.HitCount++

	return entry.Response
}

// Set stores a response in the cache
func (c *ResponseCache) Set(key string, response interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry := &CacheEntry{
		Key:       key,
		Response:  response,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		HitCount:  0,
		Metadata:  make(map[string]string),
	}

	c.cache[key] = entry
}

// Delete removes an entry from the cache
func (c *ResponseCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}

// Clear removes all entries from the cache
func (c *ResponseCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*CacheEntry)
}

// GetStats returns cache statistics
func (c *ResponseCache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	stats := CacheStats{
		TotalEntries: len(c.cache),
		TotalHits:    0,
	}

	for _, entry := range c.cache {
		stats.TotalHits += entry.HitCount
	}

	return stats
}

// cleanupExpired removes expired entries from the cache
func (c *ResponseCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		c.mutex.Lock()
		for key, entry := range c.cache {
			if now.After(entry.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

// CacheStats represents cache performance statistics
type CacheStats struct {
	TotalEntries int `json:"total_entries"`
	TotalHits    int `json:"total_hits"`
}

// ProviderMetrics tracks performance metrics for AI providers
type ProviderMetrics struct {
	requestCount int
	errorCount   int
	totalTokens  int
	totalCost    float64
	totalLatency time.Duration
	mutex        sync.RWMutex
}

// NewProviderMetrics creates a new metrics tracker
func NewProviderMetrics() *ProviderMetrics {
	return &ProviderMetrics{}
}

// RecordRequest records a successful request
func (m *ProviderMetrics) RecordRequest(inputTokens, outputTokens int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.requestCount++
	m.totalTokens += inputTokens + outputTokens
}

// RecordError records a failed request
func (m *ProviderMetrics) RecordError() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.errorCount++
}

// RecordLatency records request latency
func (m *ProviderMetrics) RecordLatency(latency time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.totalLatency += latency
}

// GetMetrics returns current performance metrics
func (m *ProviderMetrics) GetMetrics() PerformanceMetrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var avgLatency time.Duration
	var errorRate float64

	if m.requestCount > 0 {
		avgLatency = m.totalLatency / time.Duration(m.requestCount)
		errorRate = float64(m.errorCount) / float64(m.requestCount+m.errorCount)
	}

	return PerformanceMetrics{
		AverageLatency:      avgLatency,
		AverageCost:         m.totalCost / float64(m.requestCount),
		AccuracyRate:        1.0 - errorRate, // Simplified calculation
		ErrorRate:           errorRate,
		Availability:        1.0 - errorRate, // Simplified calculation
		CostEfficiencyScore: 0.8,             // Would need more complex calculation
	}
}
