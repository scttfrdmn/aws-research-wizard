package tenant

import (
	"time"
)

// TenantConfig represents a tenant organization configuration
type TenantConfig struct {
	TenantID    string            `json:"tenantId"`
	OrgName     string            `json:"orgName"`
	DisplayName string            `json:"displayName"`
	Domains     []string          `json:"domains"`
	UserLimits  TenantLimits      `json:"userLimits"`
	Billing     TenantBilling     `json:"billing"`
	Permissions TenantPermissions `json:"permissions"`
	Settings    TenantSettings    `json:"settings"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Status      TenantStatus      `json:"status"`
}

// TenantLimits defines resource limits for a tenant
type TenantLimits struct {
	MaxUsers            int   `json:"maxUsers"`
	MaxDeployments      int   `json:"maxDeployments"`
	MaxConcurrentDeploy int   `json:"maxConcurrentDeploy"`
	MaxStorageGB        int64 `json:"maxStorageGB"`
	MaxMonthlyCostUSD   int64 `json:"maxMonthlyCostUSD"`
}

// TenantBilling represents billing configuration for a tenant
type TenantBilling struct {
	BillingEmail   string `json:"billingEmail"`
	BillingContact string `json:"billingContact"`
	PaymentMethod  string `json:"paymentMethod"`
	BillingAddress string `json:"billingAddress"`
	TaxID          string `json:"taxId"`
	CostCenter     string `json:"costCenter"`
}

// TenantPermissions defines permissions available to a tenant
type TenantPermissions struct {
	AvailablePermissions []string `json:"availablePermissions"`
	RestrictedFeatures   []string `json:"restrictedFeatures"`
	AllowedRegions       []string `json:"allowedRegions"`
	AllowedInstanceTypes []string `json:"allowedInstanceTypes"`
}

// TenantSettings contains tenant-specific configuration
type TenantSettings struct {
	DefaultRegion      string             `json:"defaultRegion"`
	AllowSpotInstances bool               `json:"allowSpotInstances"`
	AutoShutdownHours  int                `json:"autoShutdownHours"`
	NotificationConfig NotificationConfig `json:"notificationConfig"`
	CustomBranding     CustomBranding     `json:"customBranding"`
}

// NotificationConfig defines notification settings for a tenant
type NotificationConfig struct {
	EmailEnabled    bool            `json:"emailEnabled"`
	SlackEnabled    bool            `json:"slackEnabled"`
	SlackWebhookURL string          `json:"slackWebhookUrl"`
	EmailDomains    []string        `json:"emailDomains"`
	AlertThresholds AlertThresholds `json:"alertThresholds"`
}

// AlertThresholds defines when to send alerts
type AlertThresholds struct {
	CostThresholdUSD       int `json:"costThresholdUsd"`
	CPUThresholdPercent    int `json:"cpuThresholdPercent"`
	MemoryThresholdPercent int `json:"memoryThresholdPercent"`
}

// CustomBranding allows tenant-specific UI customization
type CustomBranding struct {
	LogoURL          string `json:"logoUrl"`
	PrimaryColor     string `json:"primaryColor"`
	SecondaryColor   string `json:"secondaryColor"`
	OrganizationName string `json:"organizationName"`
}

// TenantStatus represents the current status of a tenant
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusPending   TenantStatus = "pending"
	TenantStatusInactive  TenantStatus = "inactive"
)

// TenantUser represents a user within a tenant
type TenantUser struct {
	UserID      string     `json:"userId"`
	TenantID    string     `json:"tenantId"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	Permissions []string   `json:"permissions"`
	CreatedAt   time.Time  `json:"createdAt"`
	LastLogin   time.Time  `json:"lastLogin"`
	Status      UserStatus `json:"status"`
}

// UserStatus represents the status of a user within a tenant
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
)

// TenantDeployment represents a deployment within a tenant
type TenantDeployment struct {
	DeploymentID string         `json:"deploymentId"`
	TenantID     string         `json:"tenantId"`
	UserID       string         `json:"userId"`
	Domain       string         `json:"domain"`
	Region       string         `json:"region"`
	InstanceType string         `json:"instanceType"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"createdAt"`
	Cost         DeploymentCost `json:"cost"`
}

// DeploymentCost tracks deployment costs
type DeploymentCost struct {
	HourlyCostUSD float64   `json:"hourlyCostUsd"`
	TotalCostUSD  float64   `json:"totalCostUsd"`
	LastUpdated   time.Time `json:"lastUpdated"`
}

// TenantStats provides tenant usage statistics
type TenantStats struct {
	TenantID           string    `json:"tenantId"`
	ActiveUsers        int       `json:"activeUsers"`
	ActiveDeployments  int       `json:"activeDeployments"`
	TotalDeployments   int       `json:"totalDeployments"`
	CurrentMonthlyCost float64   `json:"currentMonthlyCost"`
	StorageUsedGB      int64     `json:"storageUsedGb"`
	LastActivityAt     time.Time `json:"lastActivityAt"`
}
