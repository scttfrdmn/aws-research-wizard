package tenant

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Manager handles tenant operations and data isolation
type Manager struct {
	tenants    map[string]*TenantConfig
	users      map[string]*TenantUser
	deployments map[string]*TenantDeployment
	dataDir    string
	mutex      sync.RWMutex
}

// NewManager creates a new tenant manager
func NewManager(dataDir string) *Manager {
	return &Manager{
		tenants:     make(map[string]*TenantConfig),
		users:       make(map[string]*TenantUser),
		deployments: make(map[string]*TenantDeployment),
		dataDir:     dataDir,
	}
}

// CreateTenant creates a new tenant organization
func (m *Manager) CreateTenant(config *TenantConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if config.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}

	if _, exists := m.tenants[config.TenantID]; exists {
		return fmt.Errorf("tenant %s already exists", config.TenantID)
	}

	// Set timestamps
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	// Set default status
	if config.Status == "" {
		config.Status = TenantStatusActive
	}

	// Set default limits if not provided
	if config.UserLimits.MaxUsers == 0 {
		config.UserLimits.MaxUsers = 100
	}
	if config.UserLimits.MaxDeployments == 0 {
		config.UserLimits.MaxDeployments = 50
	}
	if config.UserLimits.MaxConcurrentDeploy == 0 {
		config.UserLimits.MaxConcurrentDeploy = 10
	}

	// Store in memory
	m.tenants[config.TenantID] = config

	// Persist to disk
	return m.saveTenant(config)
}

// GetTenant retrieves a tenant by ID
func (m *Manager) GetTenant(tenantID string) (*TenantConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tenant, exists := m.tenants[tenantID]
	if !exists {
		return nil, fmt.Errorf("tenant %s not found", tenantID)
	}

	return tenant, nil
}

// ListTenants returns all tenants
func (m *Manager) ListTenants() []*TenantConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tenants := make([]*TenantConfig, 0, len(m.tenants))
	for _, tenant := range m.tenants {
		tenants = append(tenants, tenant)
	}

	return tenants
}

// UpdateTenant updates an existing tenant
func (m *Manager) UpdateTenant(tenantID string, updates *TenantConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	tenant, exists := m.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}

	// Update fields
	if updates.OrgName != "" {
		tenant.OrgName = updates.OrgName
	}
	if updates.DisplayName != "" {
		tenant.DisplayName = updates.DisplayName
	}
	if len(updates.Domains) > 0 {
		tenant.Domains = updates.Domains
	}

	// Update timestamp
	tenant.UpdatedAt = time.Now()

	// Persist to disk
	return m.saveTenant(tenant)
}

// CreateUser creates a new user within a tenant
func (m *Manager) CreateUser(user *TenantUser) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Verify tenant exists
	tenant, exists := m.tenants[user.TenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", user.TenantID)
	}

	// Check user limits
	currentUsers := m.countUsersForTenant(user.TenantID)
	if currentUsers >= tenant.UserLimits.MaxUsers {
		return fmt.Errorf("tenant %s has reached maximum user limit of %d", user.TenantID, tenant.UserLimits.MaxUsers)
	}

	if user.UserID == "" {
		user.UserID = fmt.Sprintf("user-%s-%d", user.TenantID, time.Now().Unix())
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now

	// Set default status
	if user.Status == "" {
		user.Status = UserStatusActive
	}

	// Store user
	m.users[user.UserID] = user

	return m.saveUser(user)
}

// GetUser retrieves a user by ID
func (m *Manager) GetUser(userID string) (*TenantUser, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	user, exists := m.users[userID]
	if !exists {
		return nil, fmt.Errorf("user %s not found", userID)
	}

	return user, nil
}

// GetUsersForTenant returns all users for a specific tenant
func (m *Manager) GetUsersForTenant(tenantID string) []*TenantUser {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var users []*TenantUser
	for _, user := range m.users {
		if user.TenantID == tenantID {
			users = append(users, user)
		}
	}

	return users
}

// CreateDeployment creates a new deployment within a tenant
func (m *Manager) CreateDeployment(deployment *TenantDeployment) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Verify tenant exists
	tenant, exists := m.tenants[deployment.TenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", deployment.TenantID)
	}

	// Check deployment limits
	currentDeployments := m.countDeploymentsForTenant(deployment.TenantID)
	if currentDeployments >= tenant.UserLimits.MaxDeployments {
		return fmt.Errorf("tenant %s has reached maximum deployment limit of %d", deployment.TenantID, tenant.UserLimits.MaxDeployments)
	}

	// Set timestamp
	deployment.CreatedAt = time.Now()

	// Store deployment
	m.deployments[deployment.DeploymentID] = deployment

	return m.saveDeployment(deployment)
}

// GetDeploymentsForTenant returns all deployments for a specific tenant
func (m *Manager) GetDeploymentsForTenant(tenantID string) []*TenantDeployment {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var deployments []*TenantDeployment
	for _, deployment := range m.deployments {
		if deployment.TenantID == tenantID {
			deployments = append(deployments, deployment)
		}
	}

	return deployments
}

// GetTenantStats returns usage statistics for a tenant
func (m *Manager) GetTenantStats(tenantID string) (*TenantStats, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if _, exists := m.tenants[tenantID]; !exists {
		return nil, fmt.Errorf("tenant %s not found", tenantID)
	}

	stats := &TenantStats{
		TenantID: tenantID,
	}

	// Count active users
	for _, user := range m.users {
		if user.TenantID == tenantID && user.Status == UserStatusActive {
			stats.ActiveUsers++
		}
	}

	// Count deployments and calculate costs
	var totalCost float64
	for _, deployment := range m.deployments {
		if deployment.TenantID == tenantID {
			stats.TotalDeployments++
			if deployment.Status == "running" {
				stats.ActiveDeployments++
			}
			totalCost += deployment.Cost.TotalCostUSD
		}
	}

	stats.CurrentMonthlyCost = totalCost
	stats.LastActivityAt = time.Now()

	return stats, nil
}

// ValidateAccess checks if a user has access to a resource within their tenant
func (m *Manager) ValidateAccess(userID, tenantID string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	user, exists := m.users[userID]
	if !exists {
		return fmt.Errorf("user %s not found", userID)
	}

	if user.TenantID != tenantID {
		return fmt.Errorf("user %s does not belong to tenant %s", userID, tenantID)
	}

	if user.Status != UserStatusActive {
		return fmt.Errorf("user %s is not active", userID)
	}

	return nil
}

// countUsersForTenant counts active users for a tenant (internal method)
func (m *Manager) countUsersForTenant(tenantID string) int {
	count := 0
	for _, user := range m.users {
		if user.TenantID == tenantID && user.Status == UserStatusActive {
			count++
		}
	}
	return count
}

// countDeploymentsForTenant counts deployments for a tenant (internal method)
func (m *Manager) countDeploymentsForTenant(tenantID string) int {
	count := 0
	for _, deployment := range m.deployments {
		if deployment.TenantID == tenantID {
			count++
		}
	}
	return count
}

// saveTenant persists tenant data to disk
func (m *Manager) saveTenant(tenant *TenantConfig) error {
	if m.dataDir == "" {
		return nil // No persistence if no data directory
	}

	// Create tenant directory
	tenantDir := filepath.Join(m.dataDir, "tenants")
	if err := os.MkdirAll(tenantDir, 0755); err != nil {
		return err
	}

	// Save tenant config
	tenantFile := filepath.Join(tenantDir, tenant.TenantID+".json")
	data, err := json.MarshalIndent(tenant, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tenantFile, data, 0644)
}

// saveUser persists user data to disk
func (m *Manager) saveUser(user *TenantUser) error {
	if m.dataDir == "" {
		return nil // No persistence if no data directory
	}

	// Create users directory
	usersDir := filepath.Join(m.dataDir, "users")
	if err := os.MkdirAll(usersDir, 0755); err != nil {
		return err
	}

	// Save user data
	userFile := filepath.Join(usersDir, user.UserID+".json")
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(userFile, data, 0644)
}

// saveDeployment persists deployment data to disk
func (m *Manager) saveDeployment(deployment *TenantDeployment) error {
	if m.dataDir == "" {
		return nil // No persistence if no data directory
	}

	// Create deployments directory
	deploymentsDir := filepath.Join(m.dataDir, "deployments")
	if err := os.MkdirAll(deploymentsDir, 0755); err != nil {
		return err
	}

	// Save deployment data
	deploymentFile := filepath.Join(deploymentsDir, deployment.DeploymentID+".json")
	data, err := json.MarshalIndent(deployment, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(deploymentFile, data, 0644)
}

// LoadFromDisk loads tenant data from disk
func (m *Manager) LoadFromDisk() error {
	if m.dataDir == "" {
		return nil // No persistence if no data directory
	}

	// Load tenants
	tenantsDir := filepath.Join(m.dataDir, "tenants")
	if entries, err := os.ReadDir(tenantsDir); err == nil {
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) == ".json" {
				tenantFile := filepath.Join(tenantsDir, entry.Name())
				data, err := os.ReadFile(tenantFile)
				if err != nil {
					continue
				}

				var tenant TenantConfig
				if err := json.Unmarshal(data, &tenant); err != nil {
					continue
				}

				m.tenants[tenant.TenantID] = &tenant
			}
		}
	}

	// Load users
	usersDir := filepath.Join(m.dataDir, "users")
	if entries, err := os.ReadDir(usersDir); err == nil {
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) == ".json" {
				userFile := filepath.Join(usersDir, entry.Name())
				data, err := os.ReadFile(userFile)
				if err != nil {
					continue
				}

				var user TenantUser
				if err := json.Unmarshal(data, &user); err != nil {
					continue
				}

				m.users[user.UserID] = &user
			}
		}
	}

	// Load deployments
	deploymentsDir := filepath.Join(m.dataDir, "deployments")
	if entries, err := os.ReadDir(deploymentsDir); err == nil {
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) == ".json" {
				deploymentFile := filepath.Join(deploymentsDir, entry.Name())
				data, err := os.ReadFile(deploymentFile)
				if err != nil {
					continue
				}

				var deployment TenantDeployment
				if err := json.Unmarshal(data, &deployment); err != nil {
					continue
				}

				m.deployments[deployment.DeploymentID] = &deployment
			}
		}
	}

	return nil
}