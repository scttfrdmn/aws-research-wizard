// Enhanced GUI Phase 4: Advanced Analytics & Visualization Components
// Advanced dashboards, batch operations, and performance monitoring

// Analytics Dashboard Component
function AnalyticsDashboard({ deployments = [] }) {
    const [analyticsData, setAnalyticsData] = React.useState({});
    const [selectedTimeframe, setSelectedTimeframe] = React.useState('24h');
    const [selectedMetric, setSelectedMetric] = React.useState('cost');

    React.useEffect(() => {
        // Simulate analytics data generation
        const generateAnalyticsData = () => {
            const timeframes = {
                '1h': 12,   // 12 data points (5-minute intervals)
                '6h': 24,   // 24 data points (15-minute intervals)
                '24h': 48,  // 48 data points (30-minute intervals)
                '7d': 168,  // 168 data points (1-hour intervals)
                '30d': 720  // 720 data points (1-hour intervals)
            };

            const dataPoints = timeframes[selectedTimeframe] || 48;
            const costData = [];
            const performanceData = [];
            const utilizationData = [];

            for (let i = 0; i < dataPoints; i++) {
                const timestamp = Date.now() - (dataPoints - i) * (selectedTimeframe === '30d' ? 3600000 :
                                                                   selectedTimeframe === '7d' ? 3600000 :
                                                                   selectedTimeframe === '24h' ? 1800000 :
                                                                   selectedTimeframe === '6h' ? 900000 : 300000);

                costData.push({
                    timestamp,
                    value: Math.random() * 50 + 10,
                    cumulative: (i + 1) * (Math.random() * 2 + 1)
                });

                performanceData.push({
                    timestamp,
                    cpu: Math.random() * 80 + 10,
                    memory: Math.random() * 70 + 20,
                    network: Math.random() * 100 + 50
                });

                utilizationData.push({
                    timestamp,
                    efficiency: Math.random() * 30 + 70,
                    uptime: Math.min(100, 95 + Math.random() * 5)
                });
            }

            return {
                costData,
                performanceData,
                utilizationData,
                summary: {
                    totalCost: costData.reduce((sum, point) => sum + point.value, 0).toFixed(2),
                    avgCpuUsage: (performanceData.reduce((sum, point) => sum + point.cpu, 0) / dataPoints).toFixed(1),
                    avgMemoryUsage: (performanceData.reduce((sum, point) => sum + point.memory, 0) / dataPoints).toFixed(1),
                    totalUptime: (utilizationData.reduce((sum, point) => sum + point.uptime, 0) / dataPoints).toFixed(1),
                    efficiency: (utilizationData.reduce((sum, point) => sum + point.efficiency, 0) / dataPoints).toFixed(1)
                }
            };
        };

        setAnalyticsData(generateAnalyticsData());
    }, [selectedTimeframe, deployments]);

    const renderChart = (data, type) => {
        const maxValue = Math.max(...data.map(d => type === 'cost' ? d.value :
                                                   type === 'cpu' ? d.cpu :
                                                   type === 'memory' ? d.memory : d.network));

        return (
            <div className="chart-container">
                <svg width="100%" height="200" viewBox="0 0 400 200">
                    <defs>
                        <linearGradient id={`gradient-${type}`} x1="0%" y1="0%" x2="0%" y2="100%">
                            <stop offset="0%" stopColor="var(--primary-color)" stopOpacity="0.3"/>
                            <stop offset="100%" stopColor="var(--primary-color)" stopOpacity="0.1"/>
                        </linearGradient>
                    </defs>

                    {/* Chart Lines */}
                    <polyline
                        fill={`url(#gradient-${type})`}
                        stroke="var(--primary-color)"
                        strokeWidth="2"
                        points={data.map((point, index) => {
                            const x = (index / (data.length - 1)) * 380 + 10;
                            const value = type === 'cost' ? point.value :
                                         type === 'cpu' ? point.cpu :
                                         type === 'memory' ? point.memory : point.network;
                            const y = 180 - (value / maxValue) * 160;
                            return `${x},${y}`;
                        }).join(' ') + ' 390,180 10,180'}
                    />

                    {/* Data Points */}
                    {data.map((point, index) => {
                        const x = (index / (data.length - 1)) * 380 + 10;
                        const value = type === 'cost' ? point.value :
                                     type === 'cpu' ? point.cpu :
                                     type === 'memory' ? point.memory : point.network;
                        const y = 180 - (value / maxValue) * 160;
                        return (
                            <circle
                                key={index}
                                cx={x}
                                cy={y}
                                r="3"
                                fill="var(--primary-color)"
                                className="chart-point"
                            />
                        );
                    })}
                </svg>
            </div>
        );
    };

    return (
        <div className="analytics-dashboard">
            <h2>📊 Advanced Analytics Dashboard</h2>

            <div className="analytics-controls">
                <div className="control-group">
                    <label>Time Range:</label>
                    <select
                        value={selectedTimeframe}
                        onChange={(e) => setSelectedTimeframe(e.target.value)}
                    >
                        <option value="1h">Last Hour</option>
                        <option value="6h">Last 6 Hours</option>
                        <option value="24h">Last 24 Hours</option>
                        <option value="7d">Last 7 Days</option>
                        <option value="30d">Last 30 Days</option>
                    </select>
                </div>

                <div className="control-group">
                    <label>Primary Metric:</label>
                    <select
                        value={selectedMetric}
                        onChange={(e) => setSelectedMetric(e.target.value)}
                    >
                        <option value="cost">Cost Analysis</option>
                        <option value="performance">Performance Metrics</option>
                        <option value="utilization">Resource Utilization</option>
                    </select>
                </div>
            </div>

            {analyticsData.summary && (
                <div className="analytics-summary">
                    <div className="summary-card">
                        <h3>💰 Total Cost</h3>
                        <div className="summary-value">${analyticsData.summary.totalCost}</div>
                        <div className="summary-period">({selectedTimeframe})</div>
                    </div>
                    <div className="summary-card">
                        <h3>💻 Avg CPU Usage</h3>
                        <div className="summary-value">{analyticsData.summary.avgCpuUsage}%</div>
                        <div className="summary-period">({selectedTimeframe})</div>
                    </div>
                    <div className="summary-card">
                        <h3>🧠 Avg Memory Usage</h3>
                        <div className="summary-value">{analyticsData.summary.avgMemoryUsage}%</div>
                        <div className="summary-period">({selectedTimeframe})</div>
                    </div>
                    <div className="summary-card">
                        <h3>⏱️ Uptime</h3>
                        <div className="summary-value">{analyticsData.summary.totalUptime}%</div>
                        <div className="summary-period">({selectedTimeframe})</div>
                    </div>
                    <div className="summary-card">
                        <h3>⚡ Efficiency</h3>
                        <div className="summary-value">{analyticsData.summary.efficiency}%</div>
                        <div className="summary-period">({selectedTimeframe})</div>
                    </div>
                </div>
            )}

            <div className="analytics-charts">
                {selectedMetric === 'cost' && analyticsData.costData && (
                    <div className="chart-section">
                        <h3>💰 Cost Over Time</h3>
                        {renderChart(analyticsData.costData, 'cost')}
                    </div>
                )}

                {selectedMetric === 'performance' && analyticsData.performanceData && (
                    <div className="chart-section">
                        <h3>📈 Performance Metrics</h3>
                        <div className="chart-tabs">
                            <button className="chart-tab active">CPU Usage</button>
                            <button className="chart-tab">Memory Usage</button>
                            <button className="chart-tab">Network I/O</button>
                        </div>
                        {renderChart(analyticsData.performanceData, 'cpu')}
                    </div>
                )}

                {selectedMetric === 'utilization' && analyticsData.utilizationData && (
                    <div className="chart-section">
                        <h3>🎯 Resource Utilization</h3>
                        {renderChart(analyticsData.utilizationData, 'efficiency')}
                    </div>
                )}
            </div>

            <div className="analytics-insights">
                <h3>🧠 AI Insights</h3>
                <div className="insight-cards">
                    <div className="insight-card optimization">
                        <h4>💡 Cost Optimization</h4>
                        <p>Consider using spot instances during off-peak hours to reduce costs by up to 70%.</p>
                        <button className="insight-action">Apply Recommendation</button>
                    </div>
                    <div className="insight-card performance">
                        <h4>🚀 Performance Improvement</h4>
                        <p>CPU usage is consistently below 40%. Consider downsizing instances to save costs.</p>
                        <button className="insight-action">View Options</button>
                    </div>
                    <div className="insight-card efficiency">
                        <h4>⚡ Efficiency Alert</h4>
                        <p>Resource utilization efficiency is excellent at {analyticsData.summary?.efficiency}%.</p>
                        <button className="insight-action">Maintain Configuration</button>
                    </div>
                </div>
            </div>
        </div>
    );
}

// Batch Operations Component
function BatchOperations({ deployments = [], onBatchAction }) {
    const [selectedDeployments, setSelectedDeployments] = React.useState([]);
    const [batchAction, setBatchAction] = React.useState('');
    const [isProcessing, setIsProcessing] = React.useState(false);

    const handleSelectAll = () => {
        if (selectedDeployments.length === deployments.length) {
            setSelectedDeployments([]);
        } else {
            setSelectedDeployments(deployments.map(d => d.id));
        }
    };

    const handleSelectDeployment = (deploymentId) => {
        setSelectedDeployments(prev =>
            prev.includes(deploymentId)
                ? prev.filter(id => id !== deploymentId)
                : [...prev, deploymentId]
        );
    };

    const handleBatchAction = async () => {
        if (!batchAction || selectedDeployments.length === 0) return;

        setIsProcessing(true);

        try {
            // Simulate batch operation
            await new Promise(resolve => setTimeout(resolve, 2000));

            if (onBatchAction) {
                onBatchAction(batchAction, selectedDeployments);
            }

            setSelectedDeployments([]);
            setBatchAction('');
        } catch (error) {
            console.error('Batch operation failed:', error);
        } finally {
            setIsProcessing(false);
        }
    };

    return (
        <div className="batch-operations">
            <h2>🔄 Batch Operations</h2>

            {deployments.length === 0 && (
                <div className="no-deployments">
                    <p>No deployments available for batch operations.</p>
                    <p>Deploy some research environments first to use batch features.</p>
                </div>
            )}

            {deployments.length > 0 && (
                <>
                    <div className="batch-controls">
                        <div className="selection-controls">
                            <button
                                className="select-all-btn"
                                onClick={handleSelectAll}
                            >
                                {selectedDeployments.length === deployments.length ? 'Deselect All' : 'Select All'}
                            </button>
                            <span className="selection-count">
                                {selectedDeployments.length} of {deployments.length} selected
                            </span>
                        </div>

                        <div className="action-controls">
                            <select
                                value={batchAction}
                                onChange={(e) => setBatchAction(e.target.value)}
                                disabled={selectedDeployments.length === 0}
                            >
                                <option value="">Choose Action...</option>
                                <option value="start">Start All</option>
                                <option value="stop">Stop All</option>
                                <option value="restart">Restart All</option>
                                <option value="optimize">Optimize Costs</option>
                                <option value="backup">Create Backups</option>
                                <option value="terminate">Terminate All</option>
                            </select>

                            <button
                                className={`batch-action-btn ${batchAction}`}
                                onClick={handleBatchAction}
                                disabled={!batchAction || selectedDeployments.length === 0 || isProcessing}
                            >
                                {isProcessing ? '🔄 Processing...' : `Apply to ${selectedDeployments.length} deployment${selectedDeployments.length !== 1 ? 's' : ''}`}
                            </button>
                        </div>
                    </div>

                    <div className="deployment-selection">
                        <h3>Select Deployments</h3>
                        <div className="deployment-grid">
                            {deployments.map(deployment => (
                                <div
                                    key={deployment.id}
                                    className={`deployment-item ${selectedDeployments.includes(deployment.id) ? 'selected' : ''}`}
                                    onClick={() => handleSelectDeployment(deployment.id)}
                                >
                                    <div className="deployment-checkbox">
                                        <input
                                            type="checkbox"
                                            checked={selectedDeployments.includes(deployment.id)}
                                            readOnly
                                        />
                                    </div>
                                    <div className="deployment-info">
                                        <h4>{deployment.domain}</h4>
                                        <div className="deployment-meta">
                                            <span className={`status-badge ${deployment.status}`}>
                                                {deployment.status}
                                            </span>
                                            <span className="deployment-id">
                                                {deployment.id}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="batch-preview">
                        <h3>Action Preview</h3>
                        {batchAction && selectedDeployments.length > 0 ? (
                            <div className="preview-content">
                                <p>
                                    <strong>Action:</strong> {batchAction.charAt(0).toUpperCase() + batchAction.slice(1)}
                                </p>
                                <p>
                                    <strong>Deployments:</strong> {selectedDeployments.length} selected
                                </p>
                                <p>
                                    <strong>Estimated Time:</strong> {selectedDeployments.length * 30} seconds
                                </p>
                                {batchAction === 'terminate' && (
                                    <div className="warning-message">
                                        ⚠️ <strong>Warning:</strong> This action cannot be undone. All selected deployments will be permanently terminated.
                                    </div>
                                )}
                            </div>
                        ) : (
                            <p className="no-preview">Select deployments and an action to see preview.</p>
                        )}
                    </div>
                </>
            )}
        </div>
    );
}

// Performance Monitoring Component
function PerformanceMonitoring({ deployments = [] }) {
    const [alerts, setAlerts] = React.useState([]);
    const [monitoringSettings, setMonitoringSettings] = React.useState({
        cpuThreshold: 80,
        memoryThreshold: 85,
        diskThreshold: 90,
        networkThreshold: 1000,
        enableAlerts: true,
        alertEmail: ''
    });

    React.useEffect(() => {
        // Simulate performance monitoring and alert generation
        const checkPerformance = () => {
            const newAlerts = [];

            deployments.forEach(deployment => {
                const alerts = [];

                // Simulate performance issues
                if (Math.random() > 0.8) {
                    alerts.push({
                        id: `alert-${Date.now()}-${deployment.id}`,
                        deploymentId: deployment.id,
                        domain: deployment.domain,
                        type: 'cpu',
                        severity: 'warning',
                        message: `CPU usage above threshold (85% > ${monitoringSettings.cpuThreshold}%)`,
                        timestamp: new Date().toISOString(),
                        acknowledged: false
                    });
                }

                if (Math.random() > 0.9) {
                    alerts.push({
                        id: `alert-${Date.now()}-${deployment.id}-mem`,
                        deploymentId: deployment.id,
                        domain: deployment.domain,
                        type: 'memory',
                        severity: 'critical',
                        message: `Memory usage critical (92% > ${monitoringSettings.memoryThreshold}%)`,
                        timestamp: new Date().toISOString(),
                        acknowledged: false
                    });
                }

                newAlerts.push(...alerts);
            });

            setAlerts(prev => [...newAlerts, ...prev.slice(0, 10)]); // Keep last 10 alerts
        };

        if (monitoringSettings.enableAlerts && deployments.length > 0) {
            const interval = setInterval(checkPerformance, 10000); // Check every 10 seconds
            return () => clearInterval(interval);
        }
    }, [deployments, monitoringSettings]);

    const acknowledgeAlert = (alertId) => {
        setAlerts(prev => prev.map(alert =>
            alert.id === alertId ? { ...alert, acknowledged: true } : alert
        ));
    };

    const clearAllAlerts = () => {
        setAlerts([]);
    };

    return (
        <div className="performance-monitoring">
            <h2>⚡ Performance Monitoring</h2>

            <div className="monitoring-settings">
                <h3>Monitoring Configuration</h3>
                <div className="settings-grid">
                    <div className="setting-item">
                        <label>CPU Threshold (%):</label>
                        <input
                            type="number"
                            value={monitoringSettings.cpuThreshold}
                            onChange={(e) => setMonitoringSettings(prev => ({
                                ...prev,
                                cpuThreshold: parseInt(e.target.value)
                            }))}
                            min="1"
                            max="100"
                        />
                    </div>

                    <div className="setting-item">
                        <label>Memory Threshold (%):</label>
                        <input
                            type="number"
                            value={monitoringSettings.memoryThreshold}
                            onChange={(e) => setMonitoringSettings(prev => ({
                                ...prev,
                                memoryThreshold: parseInt(e.target.value)
                            }))}
                            min="1"
                            max="100"
                        />
                    </div>

                    <div className="setting-item">
                        <label>Disk Threshold (%):</label>
                        <input
                            type="number"
                            value={monitoringSettings.diskThreshold}
                            onChange={(e) => setMonitoringSettings(prev => ({
                                ...prev,
                                diskThreshold: parseInt(e.target.value)
                            }))}
                            min="1"
                            max="100"
                        />
                    </div>

                    <div className="setting-item checkbox-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={monitoringSettings.enableAlerts}
                                onChange={(e) => setMonitoringSettings(prev => ({
                                    ...prev,
                                    enableAlerts: e.target.checked
                                }))}
                            />
                            Enable Alerts
                        </label>
                    </div>
                </div>
            </div>

            <div className="alerts-section">
                <div className="alerts-header">
                    <h3>🚨 Active Alerts</h3>
                    <div className="alerts-actions">
                        <span className="alert-count">
                            {alerts.filter(a => !a.acknowledged).length} unacknowledged
                        </span>
                        <button
                            className="clear-alerts-btn"
                            onClick={clearAllAlerts}
                            disabled={alerts.length === 0}
                        >
                            Clear All
                        </button>
                    </div>
                </div>

                {alerts.length === 0 && (
                    <div className="no-alerts">
                        <p>✅ No active alerts. All systems performing within normal parameters.</p>
                    </div>
                )}

                {alerts.length > 0 && (
                    <div className="alerts-list">
                        {alerts.map(alert => (
                            <div
                                key={alert.id}
                                className={`alert-item ${alert.severity} ${alert.acknowledged ? 'acknowledged' : ''}`}
                            >
                                <div className="alert-icon">
                                    {alert.severity === 'critical' ? '🔴' : '⚠️'}
                                </div>
                                <div className="alert-content">
                                    <div className="alert-header">
                                        <span className="alert-domain">{alert.domain}</span>
                                        <span className="alert-type">{alert.type.toUpperCase()}</span>
                                        <span className="alert-time">
                                            {new Date(alert.timestamp).toLocaleTimeString()}
                                        </span>
                                    </div>
                                    <div className="alert-message">{alert.message}</div>
                                </div>
                                <div className="alert-actions">
                                    {!alert.acknowledged && (
                                        <button
                                            className="acknowledge-btn"
                                            onClick={() => acknowledgeAlert(alert.id)}
                                        >
                                            Acknowledge
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
