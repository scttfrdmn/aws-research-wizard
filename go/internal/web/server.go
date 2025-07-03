package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/scttfrdmn/aws-research-wizard/go/internal/intelligence"
	"github.com/scttfrdmn/aws-research-wizard/go/internal/monitoring"
	"github.com/scttfrdmn/aws-research-wizard/go/internal/tenant"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the web GUI server
type Server struct {
	port              int
	server            *http.Server
	domainLoader      intelligence.DomainPackLoaderInterface
	tenantManager     *tenant.Manager
	tenantMiddleware  *tenant.Middleware
	monitoringManager *monitoring.Manager
	staticFS          fs.FS
	developmentMode   bool
}

// ServerConfig holds configuration for the web server
type ServerConfig struct {
	Port        int    `json:"port"`
	Host        string `json:"host"`
	Development bool   `json:"development"`
	TLSCertPath string `json:"tls_cert_path,omitempty"`
	TLSKeyPath  string `json:"tls_key_path,omitempty"`
	EnableTLS   bool   `json:"enable_tls"`
	DataDir     string `json:"data_dir,omitempty"`
}

// NewServer creates a new web server instance
func NewServer(config ServerConfig) (*Server, error) {
	// Initialize domain loader
	domainLoader := intelligence.NewDomainPackLoader()

	// Initialize tenant manager
	tenantManager := tenant.NewManager(config.DataDir)
	if err := tenantManager.LoadFromDisk(); err != nil {
		return nil, fmt.Errorf("failed to load tenant data: %w", err)
	}

	// Initialize tenant middleware
	tenantMiddleware := tenant.NewMiddleware(tenantManager)

	// Initialize monitoring manager
	monitoringManager := monitoring.NewManager(config.DataDir)

	// Setup static file system
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to setup static file system: %w", err)
	}

	server := &Server{
		port:              config.Port,
		domainLoader:      domainLoader,
		tenantManager:     tenantManager,
		tenantMiddleware:  tenantMiddleware,
		monitoringManager: monitoringManager,
		staticFS:          staticFS,
		developmentMode:   config.Development,
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	server.setupRoutes(mux)

	server.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      server.corsMiddleware(server.loggingMiddleware(server.tenantMiddleware.TenantIsolation(mux))),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}

// Start starts the web server
func (s *Server) Start() error {
	fmt.Printf("🌐 Starting Enhanced GUI server on http://localhost:%d\n", s.port)
	fmt.Printf("🎯 Development mode: %v\n", s.developmentMode)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

// Stop gracefully stops the web server
func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("🛑 Shutting down Enhanced GUI server...")
	return s.server.Shutdown(ctx)
}

// GetPort returns the server port
func (s *Server) GetPort() int {
	return s.port
}
