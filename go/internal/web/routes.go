package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// setupRoutes configures all HTTP routes for the web interface
func (s *Server) setupRoutes(mux *http.ServeMux) {
	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.staticFS))))

	// Main application route
	mux.HandleFunc("/", s.handleIndex)

	// API routes
	mux.HandleFunc("/api/domains", s.handleDomains)
	mux.HandleFunc("/api/domains/", s.handleDomainDetails)
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/version", s.handleVersion)

	// Enhanced GUI Phase 4: Authentication endpoints
	mux.HandleFunc("/api/auth/session", s.handleAuthSession)
	mux.HandleFunc("/api/auth/login", s.handleAuthLogin)
	mux.HandleFunc("/api/auth/logout", s.handleAuthLogout)
	mux.HandleFunc("/api/auth/sso/", s.handleSSOAuth)

	// Deployment and monitoring endpoints
	mux.HandleFunc("/api/deploy", s.handleDeploy)
	mux.HandleFunc("/api/monitor", s.handleMonitor)
	mux.HandleFunc("/api/costs", s.handleCosts)
}

// handleIndex serves the main React application
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Serve the React-based interface from static files
	indexFile, err := s.staticFS.Open("static/index.html")
	if err != nil {
		http.Error(w, "Failed to load application", http.StatusInternalServerError)
		return
	}
	defer indexFile.Close()

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Copy the file content to response
	_, err = io.Copy(w, indexFile)
	if err != nil {
		http.Error(w, "Failed to serve application", http.StatusInternalServerError)
		return
	}
}

// handleDomains returns list of all available research domains
func (s *Server) handleDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	domains, err := s.domainLoader.LoadAllDomainPacks()
	if err != nil {
		http.Error(w, "Failed to load domains", http.StatusInternalServerError)
		return
	}

	// Convert to API-friendly format
	apiDomains := make([]map[string]interface{}, 0, len(domains))
	for name, domain := range domains {
		apiDomain := map[string]interface{}{
			"name":        name,
			"displayName": domain.Name,
			"description": domain.Description,
			"version":     domain.Version,
			"categories":  domain.Categories,
		}
		apiDomains = append(apiDomains, apiDomain)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"domains": apiDomains,
		"total":   len(apiDomains),
	})
}

// handleDomainDetails returns detailed information about a specific domain
func (s *Server) handleDomainDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract domain name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/domains/")
	domainName := strings.Split(path, "/")[0]

	if domainName == "" {
		http.Error(w, "Domain name required", http.StatusBadRequest)
		return
	}

	domain, err := s.domainLoader.LoadDomainPack(domainName)
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "2.1.0-alpha",
		"component": "enhanced-gui",
		"phase":     "2-domain-interface",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// handleVersion returns version information
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	version := map[string]interface{}{
		"version":   "2.1.0-alpha",
		"component": "enhanced-gui",
		"phase":     "2-domain-interface",
		"buildTime": "2025-07-03",
		"goVersion": "1.24+",
		"features": []string{
			"domain-management",
			"api-endpoints",
			"static-serving",
			"cors-middleware",
			"react-frontend",
			"cost-calculator",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

// handleDeploy manages deployment operations
func (s *Server) handleDeploy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleDeploymentStart(w, r)
	case http.MethodGet:
		s.handleDeploymentStatus(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleDeploymentStart(w http.ResponseWriter, r *http.Request) {
	var deployRequest struct {
		Domain string `json:"domain"`
		Config struct {
			InstanceSize     string `json:"instanceSize"`
			Region           string `json:"region"`
			UseSpotInstances bool   `json:"useSpotInstances"`
			AutoShutdown     bool   `json:"autoShutdown"`
			ShutdownTimeout  int    `json:"shutdownTimeout"`
			EnableBackup     bool   `json:"enableBackup"`
		} `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deployRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate deployment ID
	deploymentId := fmt.Sprintf("deploy-%d", time.Now().Unix())

	// Estimate deployment time based on instance size
	var estimatedTime string
	switch deployRequest.Config.InstanceSize {
	case "small":
		estimatedTime = "3-5 minutes"
	case "medium":
		estimatedTime = "5-8 minutes"
	case "large":
		estimatedTime = "8-12 minutes"
	case "xlarge":
		estimatedTime = "12-15 minutes"
	default:
		estimatedTime = "5-10 minutes"
	}

	response := map[string]interface{}{
		"deploymentId":     deploymentId,
		"status":           "initiated",
		"estimatedTime":    estimatedTime,
		"domain":           deployRequest.Domain,
		"instanceSize":     deployRequest.Config.InstanceSize,
		"region":           deployRequest.Config.Region,
		"useSpotInstances": deployRequest.Config.UseSpotInstances,
		"timestamp":        time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDeploymentStatus(w http.ResponseWriter, r *http.Request) {
	deploymentId := r.URL.Query().Get("id")
	if deploymentId == "" {
		http.Error(w, "Deployment ID required", http.StatusBadRequest)
		return
	}

	// Simulate deployment status - in real implementation, this would check actual deployment status
	response := map[string]interface{}{
		"deploymentId": deploymentId,
		"status":       "running",
		"progress":     75,
		"currentStep":  "Installing software packages",
		"timestamp":    time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleMonitor(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Monitor API endpoint - Coming in Phase 3",
		"phase":   "3-deployment-monitoring",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCosts(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Cost analysis API endpoint - Coming in Phase 2",
		"phase":   "2-domain-interface",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Enhanced GUI Phase 4: Authentication Handlers

func (s *Server) handleAuthSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate session check - in production, this would validate JWT/session tokens
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       "demo-user-001",
			"username": "demo-researcher",
			"name":     "Demo Researcher",
			"email":    "demo@research.example.com",
			"role":     "researcher",
		},
		"permissions": []string{
			"domains:read",
			"costs:read",
			"deployments:create",
			"deployments:read",
			"deployments:update",
			"analytics:read",
			"templates:read",
			"templates:create",
			"settings:read",
			"settings:update",
		},
		"expires":   time.Now().Add(8 * time.Hour).Format(time.RFC3339),
		"sessionId": fmt.Sprintf("session-%d", time.Now().Unix()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		MfaCode  string `json:"mfaCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate authentication - in production, this would validate against identity provider
	if loginRequest.Username == "demo" && loginRequest.Password == "demo123" {
		response := map[string]interface{}{
			"user": map[string]interface{}{
				"id":       "demo-user-001",
				"username": loginRequest.Username,
				"name":     "Demo Researcher",
				"email":    "demo@research.example.com",
				"role":     "researcher",
			},
			"permissions": []string{
				"domains:read",
				"costs:read",
				"deployments:create",
				"deployments:read",
				"deployments:update",
				"analytics:read",
				"templates:read",
				"templates:create",
				"settings:read",
				"settings:update",
			},
			"expires":   time.Now().Add(8 * time.Hour).Format(time.RFC3339),
			"sessionId": fmt.Sprintf("session-%d", time.Now().Unix()),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func (s *Server) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate logout - in production, this would invalidate session/tokens
	response := map[string]interface{}{
		"message": "Logged out successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleSSOAuth(w http.ResponseWriter, r *http.Request) {
	// Extract SSO provider from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/auth/sso/")
	provider := strings.Split(path, "/")[0]

	if provider == "" {
		http.Error(w, "SSO provider required", http.StatusBadRequest)
		return
	}

	// Simulate SSO redirect - in production, this would redirect to actual SSO provider
	redirectURL := fmt.Sprintf("https://auth.%s.com/oauth/authorize?client_id=aws-research-wizard&redirect_uri=%s/api/auth/sso/%s/callback",
		provider, r.Host, provider)

	response := map[string]interface{}{
		"provider":    provider,
		"redirectURL": redirectURL,
		"message":     fmt.Sprintf("Redirect to %s SSO", provider),
		"demo":        true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
