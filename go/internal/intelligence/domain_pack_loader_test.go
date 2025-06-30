package intelligence

import (
	"testing"
)

func TestDomainPackLoader_LoadDomainPack(t *testing.T) {
	loader := NewDomainPackLoader()

	tests := []struct {
		name           string
		domain         string
		expectError    bool
		expectNil      bool
	}{
		{
			name:        "genomics domain",
			domain:      "genomics",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "machine_learning domain",
			domain:      "machine_learning",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "climate domain",
			domain:      "climate",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "unknown domain",
			domain:      "unknown_domain",
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "empty domain",
			domain:      "",
			expectError: false,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := loader.LoadDomainPack(tt.domain)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectNil && result != nil {
				t.Errorf("expected nil result but got: %v", result)
			}
			if !tt.expectNil && result == nil {
				t.Errorf("expected non-nil result but got nil")
			}

			// Validate structure if we got a result
			if result != nil {
				if result.Name == "" {
					t.Errorf("expected non-empty name")
				}
				if result.Description == "" {
					t.Errorf("expected non-empty description")
				}
			}
		})
	}
}

func TestDomainPackLoader_LoadAllDomainPacks(t *testing.T) {
	loader := NewDomainPackLoader()

	packs, err := loader.LoadAllDomainPacks()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(packs) == 0 {
		t.Errorf("expected at least one domain pack")
	}

	// Verify each pack has required fields
	for domain, pack := range packs {
		if domain == "" {
			t.Errorf("found empty domain key")
		}
		if pack == nil {
			t.Errorf("found nil pack for domain: %s", domain)
			continue
		}
		if pack.Name == "" {
			t.Errorf("pack for domain %s has empty name", domain)
		}
		if pack.Description == "" {
			t.Errorf("pack for domain %s has empty description", domain)
		}
	}
}

func TestDomainPackLoader_GetAvailableDomains(t *testing.T) {
	loader := NewDomainPackLoader()

	domains, err := loader.GetAvailableDomains()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(domains) == 0 {
		t.Errorf("expected at least one available domain")
	}

	expectedDomains := []string{"genomics", "machine_learning", "climate"}
	for _, expected := range expectedDomains {
		found := false
		for _, domain := range domains {
			if domain == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected domain %s not found in available domains", expected)
		}
	}
}

func TestDomainPackLoader_ValidateDomainPack(t *testing.T) {
	loader := NewDomainPackLoader()

	tests := []struct {
		name        string
		domainName  string
		expectValid bool
	}{
		{
			name:        "valid genomics domain",
			domainName:  "genomics",
			expectValid: true,
		},
		{
			name:        "valid machine_learning domain",
			domainName:  "machine_learning",
			expectValid: true,
		},
		{
			name:        "valid climate domain",
			domainName:  "climate",
			expectValid: true,
		},
		{
			name:        "invalid domain",
			domainName:  "nonexistent",
			expectValid: false,
		},
		{
			name:        "empty domain",
			domainName:  "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loader.ValidateDomainPack(tt.domainName)
			hasError := err != nil
			if hasError == tt.expectValid {
				t.Errorf("expected validation result %v, got error: %v", tt.expectValid, err)
			}
		})
	}
}

func TestDomainPackLoader_ClearCache(t *testing.T) {
	loader := NewDomainPackLoader()

	// Load something into cache first
	_, err := loader.LoadDomainPack("genomics")
	if err != nil {
		t.Fatalf("failed to load domain pack: %v", err)
	}

	// Clear cache should not return error
	loader.ClearCache()

	// Verify cache is cleared by loading again
	// This is mainly to ensure the function runs without error
	_, err = loader.LoadDomainPack("genomics")
	if err != nil {
		t.Errorf("failed to load after cache clear: %v", err)
	}
}

func TestDomainPackLoader_LoadErrorCases(t *testing.T) {
	loader := NewDomainPackLoader()

	// Test loading non-existent domain
	_, err := loader.LoadDomainPack("nonexistent-domain")
	if err == nil {
		t.Error("expected error for non-existent domain")
	}

	// Test validating non-existent domain
	err = loader.ValidateDomainPack("nonexistent-domain")
	if err == nil {
		t.Error("expected error for non-existent domain validation")
	}
}

func TestDomainPackLoader_EdgeCases(t *testing.T) {
	loader := NewDomainPackLoader()

	t.Run("load domain pack with special characters", func(t *testing.T) {
		_, err := loader.LoadDomainPack("domain-with-dashes_and_underscores")
		// Should not panic or cause issues
		if err != nil {
			// Error is acceptable for non-existent domains
		}
	})

	t.Run("load domain pack with very long name", func(t *testing.T) {
		longName := string(make([]byte, 1000))
		_, err := loader.LoadDomainPack(longName)
		// Should not panic
		if err != nil {
			// Error is acceptable for invalid domains
		}
	})

	t.Run("multiple cache operations", func(t *testing.T) {
		// Multiple clears should be safe
		loader.ClearCache()
		loader.ClearCache()
		loader.ClearCache()

		// Multiple loads should be consistent
		pack1, _ := loader.LoadDomainPack("genomics")
		pack2, _ := loader.LoadDomainPack("genomics")
		
		if pack1 != nil && pack2 != nil {
			if pack1.Name != pack2.Name {
				t.Errorf("inconsistent results from cache")
			}
		}
	})
}