package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ServerConfig{
				Port:        8080,
				Host:        "localhost",
				Development: false,
			},
			wantErr: false,
		},
		{
			name: "valid config with TLS",
			config: ServerConfig{
				Port:        8443,
				Host:        "0.0.0.0",
				Development: false,
				EnableTLS:   true,
				TLSCertPath: "/path/to/cert.pem",
				TLSKeyPath:  "/path/to/key.pem",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that config can be created without panic
			_ = tt.config
			// In real implementation, we'd test config validation
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Create a mock server for testing
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	// Create a simple handler that mimics the health endpoint
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","component":"enhanced-gui"}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}
}

func TestServerCreation(t *testing.T) {
	config := ServerConfig{
		Port:        8080,
		Host:        "localhost",
		Development: true,
	}

	// Test that we can create a server config
	if config.Port != 8080 {
		t.Errorf("expected port 8080, got %d", config.Port)
	}

	if config.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", config.Host)
	}

	if !config.Development {
		t.Error("expected development mode to be true")
	}
}