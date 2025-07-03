package web

import (
	"encoding/json"
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

	// Future API endpoints (Phase 2+)
	mux.HandleFunc("/api/deploy", s.handleDeploy)
	mux.HandleFunc("/api/monitor", s.handleMonitor)
	mux.HandleFunc("/api/costs", s.handleCosts)
}

// handleIndex serves the main application
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// For now, serve a simple HTML page that will be enhanced in subsequent phases
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AWS Research Wizard - Enhanced GUI</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            padding: 30px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 40px;
        }
        .logo {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 10px;
        }
        .subtitle {
            color: #666;
            font-size: 1.2em;
        }
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 40px;
        }
        .feature {
            padding: 20px;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            background: #f9f9f9;
        }
        .feature h3 {
            color: #333;
            margin-top: 0;
        }
        .status {
            padding: 20px;
            background: #e8f5e8;
            border: 1px solid #c3e6c3;
            border-radius: 8px;
            margin: 20px 0;
        }
        .api-link {
            display: inline-block;
            margin: 5px 10px;
            padding: 8px 16px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 4px;
        }
        .api-link:hover {
            background: #5a6fd8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🔬 AWS Research Wizard</div>
            <div class="subtitle">Enhanced GUI - Phase 1 Foundation</div>
        </div>

        <div class="status">
            <h3>✅ Phase 1 Foundation Active</h3>
            <p>The Enhanced GUI foundation is now running! This is the beginning of the 17-week development plan.</p>
            <p><strong>Current Features:</strong></p>
            <ul>
                <li>Web server foundation with embedded static files</li>
                <li>RESTful API endpoints for domain management</li>
                <li>CORS and security middleware</li>
                <li>Development/production mode support</li>
            </ul>
        </div>

        <div class="features">
            <div class="feature">
                <h3>🏗️ Foundation Ready</h3>
                <p>Web server infrastructure with Go backend integration complete. Ready for React frontend development.</p>
            </div>

            <div class="feature">
                <h3>📦 22 Research Domains</h3>
                <p>All research domains accessible via API. Domain pack loader integrated and tested.</p>
                <div style="margin-top: 10px;">
                    <a href="/api/domains" class="api-link">View Domains API</a>
                </div>
            </div>

            <div class="feature">
                <h3>🚀 Production Ready Backend</h3>
                <p>Robust Go backend with 131 passing tests, intelligent data movement, and cost optimization.</p>
                <div style="margin-top: 10px;">
                    <a href="/api/health" class="api-link">Health Check</a>
                    <a href="/api/version" class="api-link">Version Info</a>
                </div>
            </div>
        </div>

        <div style="margin-top: 40px; text-align: center; color: #666;">
            <p>🔄 <strong>Next:</strong> Phase 1 continues with React frontend development and domain interface components</p>
            <p>📅 <strong>Timeline:</strong> 17 weeks total development (7 weeks Phase 1-2, 6 weeks Phase 3-4, 4 weeks Phase 5)</p>
        </div>
    </div>

    <script>
        // Simple health check and API connectivity test
        fetch('/api/health')
            .then(response => response.json())
            .then(data => {
                console.log('✅ API Health Check:', data);
            })
            .catch(error => {
                console.error('❌ API Health Check failed:', error);
            });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
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
		"phase":     "1-foundation",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// handleVersion returns version information
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	version := map[string]interface{}{
		"version":   "2.1.0-alpha",
		"component": "enhanced-gui",
		"phase":     "1-foundation",
		"buildTime": "2025-07-03",
		"goVersion": "1.24+",
		"features": []string{
			"domain-management",
			"api-endpoints",
			"static-serving",
			"cors-middleware",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

// Placeholder handlers for future development
func (s *Server) handleDeploy(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Deploy API endpoint - Coming in Phase 3",
		"phase":   "3-deployment-monitoring",
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
