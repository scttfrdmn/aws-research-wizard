// Enhanced GUI Phase 4: Advanced Deployment Templates & Automation
// Pre-configured templates, automation workflows, and deployment optimization

// Template Manager Component
function TemplateManager() {
    const [templates, setTemplates] = React.useState([]);
    const [customTemplates, setCustomTemplates] = React.useState([]);

    React.useEffect(() => {
        // Load built-in templates
        setTemplates([
            {
                id: 'ml-gpu-large',
                name: 'Machine Learning - GPU Optimized',
                description: 'High-performance GPU instances for deep learning workloads',
                category: 'machine-learning',
                config: {
                    instanceSize: 'xlarge',
                    region: 'us-east-1',
                    useSpotInstances: true,
                    autoShutdown: true,
                    shutdownTimeout: 30,
                    enableBackup: true,
                    customConfiguration: {
                        instanceType: 'p3.2xlarge',
                        storage: '500GB',
                        gpuCount: 1,
                        frameworks: ['tensorflow', 'pytorch', 'keras']
                    }
                },
                estimatedCost: '$1.20/hour',
                popularity: 95,
                lastUsed: '2025-07-02'
            },
            {
                id: 'genomics-hpc',
                name: 'Genomics HPC Cluster',
                description: 'High-memory compute cluster for genomic analysis',
                category: 'genomics',
                config: {
                    instanceSize: 'large',
                    region: 'us-west-2',
                    useSpotInstances: false,
                    autoShutdown: true,
                    shutdownTimeout: 120,
                    enableBackup: true,
                    customConfiguration: {
                        instanceType: 'r5.4xlarge',
                        storage: '2TB',
                        memory: '128GB',
                        tools: ['samtools', 'bcftools', 'gatk', 'nextflow']
                    }
                },
                estimatedCost: '$1.85/hour',
                popularity: 88,
                lastUsed: '2025-07-01'
            },
            {
                id: 'climate-modeling',
                name: 'Climate Modeling Suite',
                description: 'Optimized for climate simulation and weather modeling',
                category: 'climate',
                config: {
                    instanceSize: 'xlarge',
                    region: 'eu-west-1',
                    useSpotInstances: true,
                    autoShutdown: true,
                    shutdownTimeout: 60,
                    enableBackup: true,
                    customConfiguration: {
                        instanceType: 'c5.9xlarge',
                        storage: '1TB',
                        cpu: '36 vCPU',
                        software: ['netcdf', 'grib', 'wrf', 'cesm']
                    }
                },
                estimatedCost: '$0.95/hour',
                popularity: 72,
                lastUsed: '2025-06-30'
            },
            {
                id: 'data-science-standard',
                name: 'Data Science Standard',
                description: 'Balanced configuration for general data science work',
                category: 'data-science',
                config: {
                    instanceSize: 'medium',
                    region: 'us-east-1',
                    useSpotInstances: true,
                    autoShutdown: true,
                    shutdownTimeout: 45,
                    enableBackup: false,
                    customConfiguration: {
                        instanceType: 'm5.2xlarge',
                        storage: '250GB',
                        memory: '32GB',
                        tools: ['jupyter', 'pandas', 'numpy', 'scikit-learn']
                    }
                },
                estimatedCost: '$0.35/hour',
                popularity: 91,
                lastUsed: '2025-07-03'
            }
        ]);

        // Load custom templates from local storage
        const savedTemplates = localStorage.getItem('customTemplates');
        if (savedTemplates) {
            setCustomTemplates(JSON.parse(savedTemplates));
        }
    }, []);

    const saveCustomTemplate = (template) => {
        const newTemplate = {
            ...template,
            id: `custom-${Date.now()}`,
            isCustom: true,
            createdAt: new Date().toISOString()
        };

        const updated = [...customTemplates, newTemplate];
        setCustomTemplates(updated);
        localStorage.setItem('customTemplates', JSON.stringify(updated));

        return newTemplate;
    };

    const deleteCustomTemplate = (templateId) => {
        const updated = customTemplates.filter(t => t.id !== templateId);
        setCustomTemplates(updated);
        localStorage.setItem('customTemplates', JSON.stringify(updated));
    };

    const getAllTemplates = () => [...templates, ...customTemplates];

    return {
        templates,
        customTemplates,
        getAllTemplates,
        saveCustomTemplate,
        deleteCustomTemplate
    };
}

// Template Selector Component
function TemplateSelector({ templateManager, onTemplateSelect, selectedTemplate }) {
    const [searchTerm, setSearchTerm] = React.useState('');
    const [categoryFilter, setCategoryFilter] = React.useState('all');
    const [sortBy, setSortBy] = React.useState('popularity');

    const allTemplates = templateManager.getAllTemplates();

    const filteredTemplates = allTemplates
        .filter(template => {
            const matchesSearch = template.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                                template.description.toLowerCase().includes(searchTerm.toLowerCase());
            const matchesCategory = categoryFilter === 'all' || template.category === categoryFilter;
            return matchesSearch && matchesCategory;
        })
        .sort((a, b) => {
            switch (sortBy) {
                case 'popularity':
                    return (b.popularity || 0) - (a.popularity || 0);
                case 'name':
                    return a.name.localeCompare(b.name);
                case 'cost':
                    const aCost = parseFloat(a.estimatedCost?.replace(/[$\/hour]/g, '') || '0');
                    const bCost = parseFloat(b.estimatedCost?.replace(/[$\/hour]/g, '') || '0');
                    return aCost - bCost;
                case 'recent':
                    return new Date(b.lastUsed || 0) - new Date(a.lastUsed || 0);
                default:
                    return 0;
            }
        });

    const categories = ['all', ...new Set(allTemplates.map(t => t.category))];

    return (
        <div className="template-selector">
            <h2>🚀 Deployment Templates</h2>

            <div className="template-controls">
                <div className="search-control">
                    <input
                        type="text"
                        placeholder="Search templates..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="template-search"
                    />
                </div>

                <div className="filter-control">
                    <select
                        value={categoryFilter}
                        onChange={(e) => setCategoryFilter(e.target.value)}
                        className="category-filter"
                    >
                        {categories.map(cat => (
                            <option key={cat} value={cat}>
                                {cat === 'all' ? 'All Categories' : cat.replace('-', ' ').toUpperCase()}
                            </option>
                        ))}
                    </select>
                </div>

                <div className="sort-control">
                    <select
                        value={sortBy}
                        onChange={(e) => setSortBy(e.target.value)}
                        className="sort-select"
                    >
                        <option value="popularity">Most Popular</option>
                        <option value="name">Name</option>
                        <option value="cost">Cost (Low to High)</option>
                        <option value="recent">Recently Used</option>
                    </select>
                </div>
            </div>

            <div className="template-grid">
                {filteredTemplates.map(template => (
                    <div
                        key={template.id}
                        className={`template-card ${selectedTemplate?.id === template.id ? 'selected' : ''}`}
                        onClick={() => onTemplateSelect(template)}
                    >
                        <div className="template-header">
                            <h3>{template.name}</h3>
                            {template.isCustom && <span className="custom-badge">Custom</span>}
                        </div>

                        <p className="template-description">{template.description}</p>

                        <div className="template-meta">
                            <div className="template-category">
                                📂 {template.category?.replace('-', ' ')}
                            </div>
                            <div className="template-cost">
                                💰 {template.estimatedCost}
                            </div>
                        </div>

                        <div className="template-specs">
                            <div className="spec-item">
                                <span>Instance:</span>
                                <span>{template.config.customConfiguration?.instanceType || template.config.instanceSize}</span>
                            </div>
                            <div className="spec-item">
                                <span>Storage:</span>
                                <span>{template.config.customConfiguration?.storage || 'Default'}</span>
                            </div>
                            <div className="spec-item">
                                <span>Region:</span>
                                <span>{template.config.region}</span>
                            </div>
                        </div>

                        {template.popularity && (
                            <div className="popularity-bar">
                                <div
                                    className="popularity-fill"
                                    style={{ width: `${template.popularity}%` }}
                                ></div>
                                <span className="popularity-text">{template.popularity}% popular</span>
                            </div>
                        )}

                        {template.isCustom && (
                            <div className="template-actions">
                                <button
                                    className="delete-template-btn"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        templateManager.deleteCustomTemplate(template.id);
                                    }}
                                >
                                    🗑️ Delete
                                </button>
                            </div>
                        )}
                    </div>
                ))}
            </div>

            {filteredTemplates.length === 0 && (
                <div className="no-templates">
                    <h3>No templates found</h3>
                    <p>Try adjusting your search or filter criteria.</p>
                </div>
            )}
        </div>
    );
}

// Template Creator Component
function TemplateCreator({ templateManager, onTemplateCreated }) {
    const [template, setTemplate] = React.useState({
        name: '',
        description: '',
        category: '',
        config: {
            instanceSize: 'medium',
            region: 'us-east-1',
            useSpotInstances: true,
            autoShutdown: true,
            shutdownTimeout: 30,
            enableBackup: false,
            customConfiguration: {}
        }
    });
    const [isCreating, setIsCreating] = React.useState(false);

    const handleConfigChange = (key, value) => {
        setTemplate(prev => ({
            ...prev,
            config: {
                ...prev.config,
                [key]: value
            }
        }));
    };

    const handleCustomConfigChange = (key, value) => {
        setTemplate(prev => ({
            ...prev,
            config: {
                ...prev.config,
                customConfiguration: {
                    ...prev.config.customConfiguration,
                    [key]: value
                }
            }
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsCreating(true);

        try {
            // Calculate estimated cost based on configuration
            const estimatedCost = calculateEstimatedCost(template.config);

            const newTemplate = {
                ...template,
                estimatedCost,
                popularity: 0,
                lastUsed: new Date().toISOString()
            };

            const created = templateManager.saveCustomTemplate(newTemplate);
            onTemplateCreated(created);

            // Reset form
            setTemplate({
                name: '',
                description: '',
                category: '',
                config: {
                    instanceSize: 'medium',
                    region: 'us-east-1',
                    useSpotInstances: true,
                    autoShutdown: true,
                    shutdownTimeout: 30,
                    enableBackup: false,
                    customConfiguration: {}
                }
            });
        } catch (error) {
            console.error('Failed to create template:', error);
        } finally {
            setIsCreating(false);
        }
    };

    const calculateEstimatedCost = (config) => {
        // Simplified cost calculation
        const baseCosts = {
            small: 0.15,
            medium: 0.35,
            large: 0.85,
            xlarge: 1.20
        };

        let cost = baseCosts[config.instanceSize] || 0.35;
        if (config.useSpotInstances) {
            cost *= 0.3; // 70% discount for spot instances
        }

        return `$${cost.toFixed(2)}/hour`;
    };

    return (
        <div className="template-creator">
            <h2>🛠️ Create Custom Template</h2>

            <form onSubmit={handleSubmit} className="template-form">
                <div className="form-section">
                    <h3>Basic Information</h3>
                    <div className="form-grid">
                        <div className="form-group">
                            <label>Template Name:</label>
                            <input
                                type="text"
                                value={template.name}
                                onChange={(e) => setTemplate(prev => ({ ...prev, name: e.target.value }))}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label>Category:</label>
                            <input
                                type="text"
                                value={template.category}
                                onChange={(e) => setTemplate(prev => ({ ...prev, category: e.target.value }))}
                                placeholder="e.g., machine-learning, genomics"
                                required
                            />
                        </div>

                        <div className="form-group full-width">
                            <label>Description:</label>
                            <textarea
                                value={template.description}
                                onChange={(e) => setTemplate(prev => ({ ...prev, description: e.target.value }))}
                                required
                            />
                        </div>
                    </div>
                </div>

                <div className="form-section">
                    <h3>Configuration</h3>
                    <div className="form-grid">
                        <div className="form-group">
                            <label>Instance Size:</label>
                            <select
                                value={template.config.instanceSize}
                                onChange={(e) => handleConfigChange('instanceSize', e.target.value)}
                            >
                                <option value="small">Small</option>
                                <option value="medium">Medium</option>
                                <option value="large">Large</option>
                                <option value="xlarge">XLarge</option>
                            </select>
                        </div>

                        <div className="form-group">
                            <label>Region:</label>
                            <select
                                value={template.config.region}
                                onChange={(e) => handleConfigChange('region', e.target.value)}
                            >
                                <option value="us-east-1">US East (N. Virginia)</option>
                                <option value="us-west-2">US West (Oregon)</option>
                                <option value="eu-west-1">Europe (Ireland)</option>
                                <option value="ap-southeast-1">Asia Pacific (Singapore)</option>
                            </select>
                        </div>

                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    checked={template.config.useSpotInstances}
                                    onChange={(e) => handleConfigChange('useSpotInstances', e.target.checked)}
                                />
                                Use Spot Instances (70% savings)
                            </label>
                        </div>

                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    checked={template.config.autoShutdown}
                                    onChange={(e) => handleConfigChange('autoShutdown', e.target.checked)}
                                />
                                Auto-shutdown
                            </label>
                        </div>

                        <div className="form-group">
                            <label>Shutdown Timeout (minutes):</label>
                            <input
                                type="number"
                                value={template.config.shutdownTimeout}
                                onChange={(e) => handleConfigChange('shutdownTimeout', parseInt(e.target.value))}
                                min="5"
                                max="120"
                            />
                        </div>

                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    checked={template.config.enableBackup}
                                    onChange={(e) => handleConfigChange('enableBackup', e.target.checked)}
                                />
                                Enable Backup
                            </label>
                        </div>
                    </div>
                </div>

                <div className="form-section">
                    <h3>Advanced Configuration (Optional)</h3>
                    <div className="form-grid">
                        <div className="form-group">
                            <label>Instance Type:</label>
                            <input
                                type="text"
                                value={template.config.customConfiguration.instanceType || ''}
                                onChange={(e) => handleCustomConfigChange('instanceType', e.target.value)}
                                placeholder="e.g., c5.4xlarge"
                            />
                        </div>

                        <div className="form-group">
                            <label>Storage:</label>
                            <input
                                type="text"
                                value={template.config.customConfiguration.storage || ''}
                                onChange={(e) => handleCustomConfigChange('storage', e.target.value)}
                                placeholder="e.g., 500GB"
                            />
                        </div>

                        <div className="form-group full-width">
                            <label>Software/Tools:</label>
                            <input
                                type="text"
                                value={template.config.customConfiguration.tools || ''}
                                onChange={(e) => handleCustomConfigChange('tools', e.target.value)}
                                placeholder="e.g., tensorflow, jupyter, pandas (comma-separated)"
                            />
                        </div>
                    </div>
                </div>

                <div className="form-actions">
                    <div className="estimated-cost">
                        Estimated Cost: {calculateEstimatedCost(template.config)}
                    </div>
                    <button
                        type="submit"
                        className="create-template-btn"
                        disabled={isCreating || !template.name || !template.category}
                    >
                        {isCreating ? '🔄 Creating...' : '🚀 Create Template'}
                    </button>
                </div>
            </form>
        </div>
    );
}

// Template Quick Deploy Component
function TemplateQuickDeploy({ template, onDeploy }) {
    const [isDeploying, setIsDeploying] = React.useState(false);

    const handleQuickDeploy = async () => {
        setIsDeploying(true);
        try {
            await onDeploy(template);
        } catch (error) {
            console.error('Quick deploy failed:', error);
        } finally {
            setIsDeploying(false);
        }
    };

    if (!template) {
        return (
            <div className="quick-deploy-empty">
                <p>Select a template to enable quick deploy</p>
            </div>
        );
    }

    return (
        <div className="template-quick-deploy">
            <h3>🚀 Quick Deploy: {template.name}</h3>
            <div className="deploy-summary">
                <div className="summary-item">
                    <span>Instance:</span>
                    <span>{template.config.customConfiguration?.instanceType || template.config.instanceSize}</span>
                </div>
                <div className="summary-item">
                    <span>Region:</span>
                    <span>{template.config.region}</span>
                </div>
                <div className="summary-item">
                    <span>Cost:</span>
                    <span>{template.estimatedCost}</span>
                </div>
                <div className="summary-item">
                    <span>Spot Instance:</span>
                    <span>{template.config.useSpotInstances ? 'Yes' : 'No'}</span>
                </div>
            </div>

            <button
                className="quick-deploy-btn"
                onClick={handleQuickDeploy}
                disabled={isDeploying}
            >
                {isDeploying ? '🔄 Deploying...' : '🚀 Deploy Now'}
            </button>
        </div>
    );
}
