package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/scttfrdmn/aws-research-wizard/go/internal/monitoring"
	"github.com/scttfrdmn/aws-research-wizard/go/internal/tenant"
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

	// Enhanced GUI Phase 5: Multi-tenant management endpoints
	mux.HandleFunc("/api/tenants", s.handleTenants)
	mux.HandleFunc("/api/tenants/", s.handleTenantDetails)
	mux.HandleFunc("/api/tenant/users", s.handleTenantUsers)
	mux.HandleFunc("/api/tenant/deployments", s.handleTenantDeployments)
	mux.HandleFunc("/api/tenant/stats", s.handleTenantStats)
	mux.HandleFunc("/api/tenant/switch", s.handleTenantSwitch)

	// Enhanced GUI Phase 5: Advanced monitoring and SLA endpoints
	mux.HandleFunc("/api/monitoring/slas", s.handleSLAs)
	mux.HandleFunc("/api/monitoring/slas/", s.handleSLADetails)
	mux.HandleFunc("/api/monitoring/metrics", s.handleMetrics)
	mux.HandleFunc("/api/monitoring/alerts", s.handleAlerts)
	mux.HandleFunc("/api/monitoring/dashboards", s.handleDashboards)
	mux.HandleFunc("/api/monitoring/dashboards/", s.handleDashboardDetails)
	mux.HandleFunc("/api/monitoring/compliance", s.handleCompliance)
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

// Enhanced GUI Phase 5: Multi-Tenant Management Handlers

// handleTenants manages tenant operations
func (s *Server) handleTenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListTenants(w, r)
	case http.MethodPost:
		s.handleCreateTenant(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleListTenants(w http.ResponseWriter, r *http.Request) {
	tenants := s.tenantManager.ListTenants()

	response := map[string]interface{}{
		"tenants": tenants,
		"total":   len(tenants),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCreateTenant(w http.ResponseWriter, r *http.Request) {
	var tenantConfig tenant.TenantConfig

	if err := json.NewDecoder(r.Body).Decode(&tenantConfig); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.tenantManager.CreateTenant(&tenantConfig); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create tenant: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Tenant created successfully",
		"tenantId": tenantConfig.TenantID,
	})
}

// handleTenantDetails manages individual tenant operations
func (s *Server) handleTenantDetails(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/tenants/")
	tenantID := strings.Split(path, "/")[0]

	if tenantID == "" {
		http.Error(w, "Tenant ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGetTenant(w, r, tenantID)
	case http.MethodPut:
		s.handleUpdateTenant(w, r, tenantID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetTenant(w http.ResponseWriter, r *http.Request, tenantID string) {
	tenantConfig, err := s.tenantManager.GetTenant(tenantID)
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenantConfig)
}

func (s *Server) handleUpdateTenant(w http.ResponseWriter, r *http.Request, tenantID string) {
	var updates tenant.TenantConfig

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.tenantManager.UpdateTenant(tenantID, &updates); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update tenant: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tenant updated successfully",
	})
}

// handleTenantUsers manages users within a tenant
func (s *Server) handleTenantUsers(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant context required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		users := s.tenantManager.GetUsersForTenant(tenantID)
		response := map[string]interface{}{
			"users":    users,
			"total":    len(users),
			"tenantId": tenantID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var user tenant.TenantUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user.TenantID = tenantID
		if err := s.tenantManager.CreateUser(&user); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User created successfully",
			"userId":  user.UserID,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTenantDeployments manages deployments within a tenant
func (s *Server) handleTenantDeployments(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant context required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		deployments := s.tenantManager.GetDeploymentsForTenant(tenantID)
		response := map[string]interface{}{
			"deployments": deployments,
			"total":       len(deployments),
			"tenantId":    tenantID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var deployment tenant.TenantDeployment
		if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		deployment.TenantID = tenantID
		deployment.UserID = tenant.GetUserID(r.Context())
		if err := s.tenantManager.CreateDeployment(&deployment); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create deployment: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Deployment created successfully",
			"deploymentId": deployment.DeploymentID,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTenantStats returns tenant usage statistics
func (s *Server) handleTenantStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := tenant.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant context required", http.StatusBadRequest)
		return
	}

	stats, err := s.tenantManager.GetTenantStats(tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get tenant stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleTenantSwitch handles tenant switching for users
func (s *Server) handleTenantSwitch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var switchRequest struct {
		TenantID string `json:"tenantId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&switchRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify tenant exists
	_, err := s.tenantManager.GetTenant(switchRequest.TenantID)
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	// In a real implementation, this would update the user's session
	// For demo purposes, we'll just return success
	response := map[string]interface{}{
		"message":  "Tenant switched successfully",
		"tenantId": switchRequest.TenantID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Enhanced GUI Phase 5: Advanced Monitoring & SLA Management Handlers

// handleSLAs manages SLA operations
func (s *Server) handleSLAs(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.GetTenantID(r.Context())

	switch r.Method {
	case http.MethodGet:
		slas := s.monitoringManager.ListSLAs(tenantID)
		response := map[string]interface{}{
			"slas":     slas,
			"total":    len(slas),
			"tenantId": tenantID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var sla monitoring.SLADefinition
		if err := json.NewDecoder(r.Body).Decode(&sla); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		sla.TenantID = tenantID
		if err := s.monitoringManager.CreateSLA(&sla); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create SLA: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "SLA created successfully",
			"slaId":   sla.ID,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSLADetails manages individual SLA operations
func (s *Server) handleSLADetails(w http.ResponseWriter, r *http.Request) {
	// Extract SLA ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/monitoring/slas/")
	slaID := strings.Split(path, "/")[0]

	if slaID == "" {
		http.Error(w, "SLA ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		sla, err := s.monitoringManager.GetSLA(slaID)
		if err != nil {
			http.Error(w, "SLA not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sla)

	case http.MethodPut:
		var updates monitoring.SLADefinition
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := s.monitoringManager.UpdateSLA(slaID, &updates); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update SLA: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "SLA updated successfully",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleMetrics manages metrics operations
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse query parameters for metric filtering
		query := monitoring.MetricQuery{
			MetricName: r.URL.Query().Get("name"),
		}

		if startTime := r.URL.Query().Get("start_time"); startTime != "" {
			if t, err := time.Parse(time.RFC3339, startTime); err == nil {
				query.StartTime = t
			}
		}

		if endTime := r.URL.Query().Get("end_time"); endTime != "" {
			if t, err := time.Parse(time.RFC3339, endTime); err == nil {
				query.EndTime = t
			}
		}

		metrics, err := s.monitoringManager.GetMetrics(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get metrics: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"metrics": metrics,
			"total":   len(metrics),
			"query":   query,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var metric monitoring.Metric
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := s.monitoringManager.RecordMetric(metric); err != nil {
			http.Error(w, fmt.Sprintf("Failed to record metric: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Metric recorded successfully",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleAlerts manages alert operations
func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.GetTenantID(r.Context())

	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		var alertStatus monitoring.AlertStatus
		if status != "" {
			alertStatus = monitoring.AlertStatus(status)
		}

		alerts := s.monitoringManager.GetAlerts(tenantID, alertStatus)
		response := map[string]interface{}{
			"alerts":   alerts,
			"total":    len(alerts),
			"tenantId": tenantID,
			"status":   status,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var alert monitoring.Alert
		if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		alert.TenantID = tenantID
		if err := s.monitoringManager.CreateAlert(&alert); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create alert: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Alert created successfully",
			"alertId": alert.ID,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDashboards manages dashboard operations
func (s *Server) handleDashboards(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.GetTenantID(r.Context())

	switch r.Method {
	case http.MethodGet:
		dashboards := s.monitoringManager.ListDashboards(tenantID)
		response := map[string]interface{}{
			"dashboards": dashboards,
			"total":      len(dashboards),
			"tenantId":   tenantID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var dashboard monitoring.Dashboard
		if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		dashboard.TenantID = tenantID
		dashboard.CreatedBy = tenant.GetUserID(r.Context())
		if err := s.monitoringManager.CreateDashboard(&dashboard); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create dashboard: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "Dashboard created successfully",
			"dashboardId": dashboard.ID,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDashboardDetails manages individual dashboard operations
func (s *Server) handleDashboardDetails(w http.ResponseWriter, r *http.Request) {
	// Extract dashboard ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/monitoring/dashboards/")
	dashboardID := strings.Split(path, "/")[0]

	if dashboardID == "" {
		http.Error(w, "Dashboard ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		dashboard, err := s.monitoringManager.GetDashboard(dashboardID)
		if err != nil {
			http.Error(w, "Dashboard not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dashboard)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCompliance manages compliance reporting
func (s *Server) handleCompliance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := tenant.GetTenantID(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant context required", http.StatusBadRequest)
		return
	}

	// Parse report period from query parameters
	reportType := r.URL.Query().Get("period")
	if reportType == "" {
		reportType = "monthly"
	}

	// Calculate report period based on type
	var period monitoring.ReportPeriod
	now := time.Now()

	switch reportType {
	case "daily":
		period = monitoring.ReportPeriod{
			StartTime: now.AddDate(0, 0, -1),
			EndTime:   now,
			Type:      "daily",
		}
	case "weekly":
		period = monitoring.ReportPeriod{
			StartTime: now.AddDate(0, 0, -7),
			EndTime:   now,
			Type:      "weekly",
		}
	case "monthly":
		period = monitoring.ReportPeriod{
			StartTime: now.AddDate(0, -1, 0),
			EndTime:   now,
			Type:      "monthly",
		}
	case "quarterly":
		period = monitoring.ReportPeriod{
			StartTime: now.AddDate(0, -3, 0),
			EndTime:   now,
			Type:      "quarterly",
		}
	default:
		http.Error(w, "Invalid report period", http.StatusBadRequest)
		return
	}

	report, err := s.monitoringManager.GenerateComplianceReport(tenantID, period)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate compliance report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
