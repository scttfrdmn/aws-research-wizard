// Enhanced GUI Phase 5: Multi-Tenant Management Components

// TenantSelector component for organization switching
function TenantSelector() {
    const [tenants, setTenants] = React.useState([]);
    const [currentTenant, setCurrentTenant] = React.useState(null);
    const [loading, setLoading] = React.useState(true);
    const [showDropdown, setShowDropdown] = React.useState(false);

    React.useEffect(() => {
        loadTenants();
    }, []);

    const loadTenants = async () => {
        try {
            const response = await fetch('/api/tenants');
            const data = await response.json();
            setTenants(data.tenants || []);
            
            // Set first tenant as current if none selected
            if (data.tenants && data.tenants.length > 0 && !currentTenant) {
                setCurrentTenant(data.tenants[0]);
            }
        } catch (error) {
            console.error('Failed to load tenants:', error);
        } finally {
            setLoading(false);
        }
    };

    const switchTenant = async (tenant) => {
        try {
            const response = await fetch('/api/tenant/switch', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Tenant-ID': tenant.tenantId
                },
                body: JSON.stringify({ tenantId: tenant.tenantId })
            });

            if (response.ok) {
                setCurrentTenant(tenant);
                setShowDropdown(false);
                // Trigger global tenant change event
                window.dispatchEvent(new CustomEvent('tenantChanged', { detail: tenant }));
            }
        } catch (error) {
            console.error('Failed to switch tenant:', error);
        }
    };

    if (loading) {
        return React.createElement('div', { className: 'tenant-selector loading' },
            React.createElement('div', { className: 'spinner' })
        );
    }

    return React.createElement('div', { className: 'tenant-selector' },
        React.createElement('div', { 
            className: 'tenant-current',
            onClick: () => setShowDropdown(!showDropdown)
        },
            React.createElement('div', { className: 'tenant-info' },
                React.createElement('div', { className: 'tenant-name' }, 
                    currentTenant ? currentTenant.displayName : 'Select Organization'
                ),
                React.createElement('div', { className: 'tenant-id' }, 
                    currentTenant ? `(${currentTenant.tenantId})` : ''
                )
            ),
            React.createElement('div', { className: 'tenant-arrow' }, '▼')
        ),
        
        showDropdown && React.createElement('div', { className: 'tenant-dropdown' },
            tenants.map(tenant => 
                React.createElement('div', {
                    key: tenant.tenantId,
                    className: `tenant-option ${currentTenant?.tenantId === tenant.tenantId ? 'active' : ''}`,
                    onClick: () => switchTenant(tenant)
                },
                    React.createElement('div', { className: 'tenant-name' }, tenant.displayName),
                    React.createElement('div', { className: 'tenant-details' },
                        React.createElement('span', { className: 'tenant-id' }, tenant.tenantId),
                        React.createElement('span', { className: 'tenant-status' }, tenant.status)
                    )
                )
            )
        )
    );
}

// TenantDashboard component for organization overview
function TenantDashboard() {
    const [tenantStats, setTenantStats] = React.useState(null);
    const [tenantConfig, setTenantConfig] = React.useState(null);
    const [loading, setLoading] = React.useState(true);
    const [currentTenant, setCurrentTenant] = React.useState(null);

    React.useEffect(() => {
        // Listen for tenant changes
        const handleTenantChange = (event) => {
            setCurrentTenant(event.detail);
            loadTenantData(event.detail.tenantId);
        };

        window.addEventListener('tenantChanged', handleTenantChange);
        return () => window.removeEventListener('tenantChanged', handleTenantChange);
    }, []);

    const loadTenantData = async (tenantId) => {
        if (!tenantId) return;
        
        setLoading(true);
        try {
            // Load tenant configuration
            const configResponse = await fetch(`/api/tenants/${tenantId}`, {
                headers: { 'X-Tenant-ID': tenantId }
            });
            const configData = await configResponse.json();
            setTenantConfig(configData);

            // Load tenant statistics
            const statsResponse = await fetch('/api/tenant/stats', {
                headers: { 'X-Tenant-ID': tenantId }
            });
            const statsData = await statsResponse.json();
            setTenantStats(statsData);

        } catch (error) {
            console.error('Failed to load tenant data:', error);
        } finally {
            setLoading(false);
        }
    };

    if (loading || !currentTenant) {
        return React.createElement('div', { className: 'tenant-dashboard loading' },
            React.createElement('h3', null, 'Loading tenant dashboard...'),
            React.createElement('div', { className: 'spinner' })
        );
    }

    return React.createElement('div', { className: 'tenant-dashboard' },
        React.createElement('div', { className: 'dashboard-header' },
            React.createElement('h2', null, `${tenantConfig?.displayName} Dashboard`),
            React.createElement('div', { className: 'tenant-status' },
                React.createElement('span', { 
                    className: `status-badge ${tenantConfig?.status}` 
                }, tenantConfig?.status)
            )
        ),

        React.createElement('div', { className: 'dashboard-stats' },
            React.createElement('div', { className: 'stat-card' },
                React.createElement('div', { className: 'stat-value' }, tenantStats?.activeUsers || 0),
                React.createElement('div', { className: 'stat-label' }, 'Active Users'),
                React.createElement('div', { className: 'stat-limit' }, 
                    `of ${tenantConfig?.userLimits?.maxUsers || 0} max`
                )
            ),
            React.createElement('div', { className: 'stat-card' },
                React.createElement('div', { className: 'stat-value' }, tenantStats?.activeDeployments || 0),
                React.createElement('div', { className: 'stat-label' }, 'Active Deployments'),
                React.createElement('div', { className: 'stat-limit' }, 
                    `of ${tenantConfig?.userLimits?.maxDeployments || 0} max`
                )
            ),
            React.createElement('div', { className: 'stat-card' },
                React.createElement('div', { className: 'stat-value' }, 
                    `$${(tenantStats?.currentMonthlyCost || 0).toFixed(2)}`
                ),
                React.createElement('div', { className: 'stat-label' }, 'Monthly Cost'),
                React.createElement('div', { className: 'stat-limit' }, 
                    `of $${tenantConfig?.userLimits?.maxMonthlyCostUSD || 0} budget`
                )
            ),
            React.createElement('div', { className: 'stat-card' },
                React.createElement('div', { className: 'stat-value' }, 
                    `${(tenantStats?.storageUsedGB || 0).toFixed(1)} GB`
                ),
                React.createElement('div', { className: 'stat-label' }, 'Storage Used'),
                React.createElement('div', { className: 'stat-limit' }, 
                    `of ${tenantConfig?.userLimits?.maxStorageGB || 0} GB limit`
                )
            )
        ),

        React.createElement('div', { className: 'dashboard-details' },
            React.createElement('div', { className: 'detail-section' },
                React.createElement('h3', null, 'Organization Details'),
                React.createElement('div', { className: 'detail-grid' },
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Organization:'),
                        React.createElement('span', { className: 'detail-value' }, tenantConfig?.displayName)
                    ),
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Tenant ID:'),
                        React.createElement('span', { className: 'detail-value' }, tenantConfig?.tenantId)
                    ),
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Default Region:'),
                        React.createElement('span', { className: 'detail-value' }, tenantConfig?.settings?.defaultRegion)
                    ),
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Spot Instances:'),
                        React.createElement('span', { className: 'detail-value' }, 
                            tenantConfig?.settings?.allowSpotInstances ? 'Enabled' : 'Disabled'
                        )
                    )
                )
            ),

            React.createElement('div', { className: 'detail-section' },
                React.createElement('h3', null, 'Billing Information'),
                React.createElement('div', { className: 'detail-grid' },
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Billing Contact:'),
                        React.createElement('span', { className: 'detail-value' }, tenantConfig?.billing?.billingContact)
                    ),
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Cost Center:'),
                        React.createElement('span', { className: 'detail-value' }, tenantConfig?.billing?.costCenter)
                    ),
                    React.createElement('div', { className: 'detail-item' },
                        React.createElement('span', { className: 'detail-label' }, 'Payment Method:'),
                        React.createElement('span', { className: 'detail-value' }, tenantConfig?.billing?.paymentMethod)
                    )
                )
            )
        )
    );
}

// TenantUserManagement component for per-organization user management
function TenantUserManagement() {
    const [users, setUsers] = React.useState([]);
    const [loading, setLoading] = React.useState(true);
    const [currentTenant, setCurrentTenant] = React.useState(null);
    const [showAddUser, setShowAddUser] = React.useState(false);

    React.useEffect(() => {
        // Listen for tenant changes
        const handleTenantChange = (event) => {
            setCurrentTenant(event.detail);
            loadUsers(event.detail.tenantId);
        };

        window.addEventListener('tenantChanged', handleTenantChange);
        return () => window.removeEventListener('tenantChanged', handleTenantChange);
    }, []);

    const loadUsers = async (tenantId) => {
        if (!tenantId) return;
        
        setLoading(true);
        try {
            const response = await fetch('/api/tenant/users', {
                headers: { 'X-Tenant-ID': tenantId }
            });
            const data = await response.json();
            setUsers(data.users || []);
        } catch (error) {
            console.error('Failed to load users:', error);
        } finally {
            setLoading(false);
        }
    };

    if (loading || !currentTenant) {
        return React.createElement('div', { className: 'tenant-users loading' },
            React.createElement('h3', null, 'Loading users...'),
            React.createElement('div', { className: 'spinner' })
        );
    }

    return React.createElement('div', { className: 'tenant-users' },
        React.createElement('div', { className: 'users-header' },
            React.createElement('h3', null, `Users - ${currentTenant.displayName}`),
            React.createElement('button', {
                className: 'add-user-btn',
                onClick: () => setShowAddUser(true)
            }, '+ Add User')
        ),

        React.createElement('div', { className: 'users-table' },
            React.createElement('div', { className: 'table-header' },
                React.createElement('div', { className: 'header-cell' }, 'Name'),
                React.createElement('div', { className: 'header-cell' }, 'Username'),
                React.createElement('div', { className: 'header-cell' }, 'Email'),
                React.createElement('div', { className: 'header-cell' }, 'Role'),
                React.createElement('div', { className: 'header-cell' }, 'Status'),
                React.createElement('div', { className: 'header-cell' }, 'Last Login')
            ),
            
            users.map(user => 
                React.createElement('div', { 
                    key: user.userId, 
                    className: 'table-row' 
                },
                    React.createElement('div', { className: 'table-cell' }, user.name),
                    React.createElement('div', { className: 'table-cell' }, user.username),
                    React.createElement('div', { className: 'table-cell' }, user.email),
                    React.createElement('div', { className: 'table-cell' }, 
                        React.createElement('span', { className: `role-badge ${user.role}` }, user.role)
                    ),
                    React.createElement('div', { className: 'table-cell' }, 
                        React.createElement('span', { className: `status-badge ${user.status}` }, user.status)
                    ),
                    React.createElement('div', { className: 'table-cell' }, 
                        user.lastLogin ? new Date(user.lastLogin).toLocaleDateString() : 'Never'
                    )
                )
            )
        )
    );
}

// TenantBilling component for organization-specific cost tracking
function TenantBilling() {
    const [deployments, setDeployments] = React.useState([]);
    const [loading, setLoading] = React.useState(true);
    const [currentTenant, setCurrentTenant] = React.useState(null);
    const [totalCost, setTotalCost] = React.useState(0);

    React.useEffect(() => {
        // Listen for tenant changes
        const handleTenantChange = (event) => {
            setCurrentTenant(event.detail);
            loadBillingData(event.detail.tenantId);
        };

        window.addEventListener('tenantChanged', handleTenantChange);
        return () => window.removeEventListener('tenantChanged', handleTenantChange);
    }, []);

    const loadBillingData = async (tenantId) => {
        if (!tenantId) return;
        
        setLoading(true);
        try {
            const response = await fetch('/api/tenant/deployments', {
                headers: { 'X-Tenant-ID': tenantId }
            });
            const data = await response.json();
            setDeployments(data.deployments || []);
            
            // Calculate total cost
            const total = (data.deployments || []).reduce((sum, deployment) => 
                sum + (deployment.cost?.totalCostUSD || 0), 0
            );
            setTotalCost(total);
        } catch (error) {
            console.error('Failed to load billing data:', error);
        } finally {
            setLoading(false);
        }
    };

    if (loading || !currentTenant) {
        return React.createElement('div', { className: 'tenant-billing loading' },
            React.createElement('h3', null, 'Loading billing data...'),
            React.createElement('div', { className: 'spinner' })
        );
    }

    return React.createElement('div', { className: 'tenant-billing' },
        React.createElement('div', { className: 'billing-header' },
            React.createElement('h3', null, `Billing - ${currentTenant.displayName}`),
            React.createElement('div', { className: 'total-cost' },
                React.createElement('span', { className: 'cost-label' }, 'Total Monthly Cost:'),
                React.createElement('span', { className: 'cost-value' }, `$${totalCost.toFixed(2)}`)
            )
        ),

        React.createElement('div', { className: 'deployments-table' },
            React.createElement('div', { className: 'table-header' },
                React.createElement('div', { className: 'header-cell' }, 'Deployment'),
                React.createElement('div', { className: 'header-cell' }, 'Domain'),
                React.createElement('div', { className: 'header-cell' }, 'Instance'),
                React.createElement('div', { className: 'header-cell' }, 'Region'),
                React.createElement('div', { className: 'header-cell' }, 'Status'),
                React.createElement('div', { className: 'header-cell' }, 'Hourly Cost'),
                React.createElement('div', { className: 'header-cell' }, 'Total Cost')
            ),
            
            deployments.map(deployment => 
                React.createElement('div', { 
                    key: deployment.deploymentId, 
                    className: 'table-row' 
                },
                    React.createElement('div', { className: 'table-cell' }, deployment.deploymentId),
                    React.createElement('div', { className: 'table-cell' }, deployment.domain),
                    React.createElement('div', { className: 'table-cell' }, deployment.instanceType),
                    React.createElement('div', { className: 'table-cell' }, deployment.region),
                    React.createElement('div', { className: 'table-cell' }, 
                        React.createElement('span', { className: `status-badge ${deployment.status}` }, deployment.status)
                    ),
                    React.createElement('div', { className: 'table-cell' }, 
                        `$${(deployment.cost?.hourlyCostUSD || 0).toFixed(2)}/hr`
                    ),
                    React.createElement('div', { className: 'table-cell' }, 
                        `$${(deployment.cost?.totalCostUSD || 0).toFixed(2)}`
                    )
                )
            )
        )
    );
}

// Global tenant management state
window.TenantManager = {
    TenantSelector,
    TenantDashboard,
    TenantUserManagement,
    TenantBilling
};