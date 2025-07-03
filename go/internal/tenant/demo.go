package tenant

import (
	"fmt"
	"time"
)

// CreateDemoTenants creates sample tenant organizations for demonstration
func (m *Manager) CreateDemoTenants() error {
	// Create Research University tenant
	universityTenant := &TenantConfig{
		TenantID:    "research-university",
		OrgName:     "research-university",
		DisplayName: "Research University",
		Domains:     []string{"university.edu"},
		UserLimits: TenantLimits{
			MaxUsers:            500,
			MaxDeployments:      200,
			MaxConcurrentDeploy: 50,
			MaxStorageGB:        10000,
			MaxMonthlyCostUSD:   50000,
		},
		Billing: TenantBilling{
			BillingEmail:   "billing@university.edu",
			BillingContact: "Finance Department",
			PaymentMethod:  "University Account",
			BillingAddress: "123 University Ave, Academic City, AC 12345",
			TaxID:          "TAX-UNIV-001",
			CostCenter:     "Research Computing",
		},
		Permissions: TenantPermissions{
			AvailablePermissions: []string{
				"domains:read", "costs:read", "deployments:create",
				"deployments:read", "deployments:update", "deployments:delete",
				"analytics:read", "templates:read", "templates:create",
				"settings:read", "settings:update",
			},
			RestrictedFeatures:   []string{},
			AllowedRegions:       []string{"us-east-1", "us-west-2", "eu-west-1"},
			AllowedInstanceTypes: []string{"small", "medium", "large", "xlarge"},
		},
		Settings: TenantSettings{
			DefaultRegion:      "us-east-1",
			AllowSpotInstances: true,
			AutoShutdownHours:  8,
			NotificationConfig: NotificationConfig{
				EmailEnabled:    true,
				SlackEnabled:    true,
				SlackWebhookURL: "https://hooks.slack.com/services/EXAMPLE/WEBHOOK",
				EmailDomains:    []string{"university.edu"},
				AlertThresholds: AlertThresholds{
					CostThresholdUSD:       1000,
					CPUThresholdPercent:    80,
					MemoryThresholdPercent: 85,
				},
			},
			CustomBranding: CustomBranding{
				LogoURL:          "https://university.edu/logo.png",
				PrimaryColor:     "#003366",
				SecondaryColor:   "#0066CC",
				OrganizationName: "Research University",
			},
		},
		Status: TenantStatusActive,
	}

	if err := m.CreateTenant(universityTenant); err != nil {
		return fmt.Errorf("failed to create university tenant: %w", err)
	}

	// Create National Lab tenant
	nationalLabTenant := &TenantConfig{
		TenantID:    "national-lab",
		OrgName:     "national-lab",
		DisplayName: "National Research Laboratory",
		Domains:     []string{"nationallab.gov"},
		UserLimits: TenantLimits{
			MaxUsers:            1000,
			MaxDeployments:      500,
			MaxConcurrentDeploy: 100,
			MaxStorageGB:        50000,
			MaxMonthlyCostUSD:   100000,
		},
		Billing: TenantBilling{
			BillingEmail:   "finance@nationallab.gov",
			BillingContact: "Laboratory Finance Office",
			PaymentMethod:  "Government Account",
			BillingAddress: "456 Science Blvd, Research Park, RP 67890",
			TaxID:          "GOV-LAB-002",
			CostCenter:     "High Performance Computing",
		},
		Permissions: TenantPermissions{
			AvailablePermissions: []string{
				"domains:read", "costs:read", "deployments:create",
				"deployments:read", "deployments:update", "deployments:delete",
				"analytics:read", "analytics:create", "templates:read",
				"templates:create", "templates:update", "settings:read",
				"settings:update", "admin:read",
			},
			RestrictedFeatures:   []string{"spot-instances"}, // Government restrictions
			AllowedRegions:       []string{"us-gov-east-1", "us-gov-west-1"},
			AllowedInstanceTypes: []string{"large", "xlarge", "2xlarge"},
		},
		Settings: TenantSettings{
			DefaultRegion:      "us-gov-east-1",
			AllowSpotInstances: false, // Government compliance
			AutoShutdownHours:  12,
			NotificationConfig: NotificationConfig{
				EmailEnabled:    true,
				SlackEnabled:    false, // Security restrictions
				SlackWebhookURL: "",
				EmailDomains:    []string{"nationallab.gov"},
				AlertThresholds: AlertThresholds{
					CostThresholdUSD:       5000,
					CPUThresholdPercent:    75,
					MemoryThresholdPercent: 80,
				},
			},
			CustomBranding: CustomBranding{
				LogoURL:          "https://nationallab.gov/seal.png",
				PrimaryColor:     "#1F4E79",
				SecondaryColor:   "#5B9BD5",
				OrganizationName: "National Research Laboratory",
			},
		},
		Status: TenantStatusActive,
	}

	if err := m.CreateTenant(nationalLabTenant); err != nil {
		return fmt.Errorf("failed to create national lab tenant: %w", err)
	}

	// Create Enterprise tenant
	enterpriseTenant := &TenantConfig{
		TenantID:    "enterprise-research",
		OrgName:     "enterprise-research",
		DisplayName: "Enterprise Research Division",
		Domains:     []string{"research.enterprise.com"},
		UserLimits: TenantLimits{
			MaxUsers:            200,
			MaxDeployments:      100,
			MaxConcurrentDeploy: 25,
			MaxStorageGB:        5000,
			MaxMonthlyCostUSD:   25000,
		},
		Billing: TenantBilling{
			BillingEmail:   "research-billing@enterprise.com",
			BillingContact: "Research Finance Team",
			PaymentMethod:  "Corporate Credit Card",
			BillingAddress: "789 Innovation Drive, Tech Valley, TV 13579",
			TaxID:          "CORP-ENT-003",
			CostCenter:     "R&D Computing",
		},
		Permissions: TenantPermissions{
			AvailablePermissions: []string{
				"domains:read", "costs:read", "deployments:create",
				"deployments:read", "deployments:update",
				"analytics:read", "templates:read", "templates:create",
				"settings:read",
			},
			RestrictedFeatures:   []string{"admin:*"},
			AllowedRegions:       []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"},
			AllowedInstanceTypes: []string{"small", "medium", "large"},
		},
		Settings: TenantSettings{
			DefaultRegion:      "us-east-1",
			AllowSpotInstances: true,
			AutoShutdownHours:  6,
			NotificationConfig: NotificationConfig{
				EmailEnabled:    true,
				SlackEnabled:    true,
				SlackWebhookURL: "https://hooks.slack.com/services/ENTERPRISE/RESEARCH",
				EmailDomains:    []string{"enterprise.com"},
				AlertThresholds: AlertThresholds{
					CostThresholdUSD:       2000,
					CPUThresholdPercent:    85,
					MemoryThresholdPercent: 90,
				},
			},
			CustomBranding: CustomBranding{
				LogoURL:          "https://enterprise.com/research-logo.png",
				PrimaryColor:     "#8B0000",
				SecondaryColor:   "#CD5C5C",
				OrganizationName: "Enterprise Research Division",
			},
		},
		Status: TenantStatusActive,
	}

	if err := m.CreateTenant(enterpriseTenant); err != nil {
		return fmt.Errorf("failed to create enterprise tenant: %w", err)
	}

	return nil
}

// CreateDemoUsers creates sample users for each demo tenant
func (m *Manager) CreateDemoUsers() error {
	demoUsers := []TenantUser{
		// Research University users
		{
			TenantID:    "research-university",
			Username:    "prof.smith",
			Email:       "prof.smith@university.edu",
			Name:        "Professor John Smith",
			Role:        "faculty",
			Permissions: []string{"domains:read", "costs:read", "deployments:create", "deployments:read", "analytics:read"},
			Status:      UserStatusActive,
		},
		{
			TenantID:    "research-university",
			Username:    "grad.johnson",
			Email:       "grad.johnson@university.edu",
			Name:        "Sarah Johnson",
			Role:        "graduate-student",
			Permissions: []string{"domains:read", "costs:read", "deployments:create", "deployments:read"},
			Status:      UserStatusActive,
		},
		{
			TenantID:    "research-university",
			Username:    "admin.university",
			Email:       "admin@university.edu",
			Name:        "Research Computing Admin",
			Role:        "administrator",
			Permissions: []string{"*"}, // Full permissions
			Status:      UserStatusActive,
		},

		// National Lab users
		{
			TenantID:    "national-lab",
			Username:    "scientist.davis",
			Email:       "scientist.davis@nationallab.gov",
			Name:        "Dr. Michael Davis",
			Role:        "research-scientist",
			Permissions: []string{"domains:read", "costs:read", "deployments:create", "deployments:read", "analytics:read"},
			Status:      UserStatusActive,
		},
		{
			TenantID:    "national-lab",
			Username:    "hpc.admin",
			Email:       "hpc.admin@nationallab.gov",
			Name:        "HPC Administrator",
			Role:        "hpc-administrator",
			Permissions: []string{"*"}, // Full permissions
			Status:      UserStatusActive,
		},

		// Enterprise users
		{
			TenantID:    "enterprise-research",
			Username:    "researcher.wilson",
			Email:       "researcher.wilson@enterprise.com",
			Name:        "Dr. Emily Wilson",
			Role:        "senior-researcher",
			Permissions: []string{"domains:read", "costs:read", "deployments:create", "deployments:read", "analytics:read"},
			Status:      UserStatusActive,
		},
		{
			TenantID:    "enterprise-research",
			Username:    "team.lead",
			Email:       "team.lead@enterprise.com",
			Name:        "Research Team Lead",
			Role:        "team-lead",
			Permissions: []string{"domains:read", "costs:read", "deployments:create", "deployments:read", "deployments:update", "analytics:read", "templates:read"},
			Status:      UserStatusActive,
		},
	}

	for _, user := range demoUsers {
		if err := m.CreateUser(&user); err != nil {
			return fmt.Errorf("failed to create demo user %s: %w", user.Username, err)
		}
	}

	return nil
}

// CreateDemoDeployments creates sample deployments for demo tenants
func (m *Manager) CreateDemoDeployments() error {
	// Get user IDs for deployments
	users := m.users
	var universityUserID, labUserID, enterpriseUserID string

	for _, user := range users {
		switch user.TenantID {
		case "research-university":
			if user.Username == "prof.smith" {
				universityUserID = user.UserID
			}
		case "national-lab":
			if user.Username == "scientist.davis" {
				labUserID = user.UserID
			}
		case "enterprise-research":
			if user.Username == "researcher.wilson" {
				enterpriseUserID = user.UserID
			}
		}
	}

	demoDeployments := []TenantDeployment{
		// Research University deployments
		{
			DeploymentID: "deploy-university-001",
			TenantID:     "research-university",
			UserID:       universityUserID,
			Domain:       "genomics",
			Region:       "us-east-1",
			InstanceType: "large",
			Status:       "running",
			Cost: DeploymentCost{
				HourlyCostUSD: 2.50,
				TotalCostUSD:  60.00,
				LastUpdated:   time.Now(),
			},
		},
		{
			DeploymentID: "deploy-university-002",
			TenantID:     "research-university",
			UserID:       universityUserID,
			Domain:       "climate-modeling",
			Region:       "us-east-1",
			InstanceType: "xlarge",
			Status:       "stopped",
			Cost: DeploymentCost{
				HourlyCostUSD: 4.80,
				TotalCostUSD:  240.00,
				LastUpdated:   time.Now(),
			},
		},

		// National Lab deployments
		{
			DeploymentID: "deploy-lab-001",
			TenantID:     "national-lab",
			UserID:       labUserID,
			Domain:       "physics-computational",
			Region:       "us-gov-east-1",
			InstanceType: "2xlarge",
			Status:       "running",
			Cost: DeploymentCost{
				HourlyCostUSD: 8.90,
				TotalCostUSD:  890.00,
				LastUpdated:   time.Now(),
			},
		},

		// Enterprise deployments
		{
			DeploymentID: "deploy-enterprise-001",
			TenantID:     "enterprise-research",
			UserID:       enterpriseUserID,
			Domain:       "machine-learning",
			Region:       "us-east-1",
			InstanceType: "large",
			Status:       "running",
			Cost: DeploymentCost{
				HourlyCostUSD: 3.20,
				TotalCostUSD:  160.00,
				LastUpdated:   time.Now(),
			},
		},
	}

	for _, deployment := range demoDeployments {
		if err := m.CreateDeployment(&deployment); err != nil {
			return fmt.Errorf("failed to create demo deployment %s: %w", deployment.DeploymentID, err)
		}
	}

	return nil
}

// SetupDemoEnvironment creates a complete demo multi-tenant environment
func (m *Manager) SetupDemoEnvironment() error {
	fmt.Println("🏗️  Setting up multi-tenant demo environment...")

	// Create demo tenants
	if err := m.CreateDemoTenants(); err != nil {
		return fmt.Errorf("failed to create demo tenants: %w", err)
	}
	fmt.Println("✅ Created 3 demo tenant organizations")

	// Create demo users
	if err := m.CreateDemoUsers(); err != nil {
		return fmt.Errorf("failed to create demo users: %w", err)
	}
	fmt.Println("✅ Created 7 demo users across tenants")

	// Create demo deployments
	if err := m.CreateDemoDeployments(); err != nil {
		return fmt.Errorf("failed to create demo deployments: %w", err)
	}
	fmt.Println("✅ Created 4 demo deployments")

	fmt.Println("🎉 Multi-tenant demo environment setup complete!")
	fmt.Println("📊 Available tenants:")
	fmt.Println("   - research-university (Research University)")
	fmt.Println("   - national-lab (National Research Laboratory)")
	fmt.Println("   - enterprise-research (Enterprise Research Division)")

	return nil
}
