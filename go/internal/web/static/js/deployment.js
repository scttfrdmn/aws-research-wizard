// Enhanced GUI Phase 3: Deployment & Monitoring Components
// Deployment workflow management and real-time monitoring

// Deployment Workflow Component
function DeploymentWorkflow({ selectedDomain, onDeploymentStart }) {
    const [deploymentConfig, setDeploymentConfig] = React.useState({
        instanceSize: 'medium',
        region: 'us-east-1',
        useSpotInstances: false,
        autoShutdown: true,
        shutdownTimeout: 30,
        enableBackup: true
    });
    
    const [deploymentStatus, setDeploymentStatus] = React.useState('ready'); // ready, deploying, success, error
    const [deploymentLog, setDeploymentLog] = React.useState([]);
    const [estimatedTime, setEstimatedTime] = React.useState(null);

    const handleDeploy = async () => {
        if (!selectedDomain) {
            alert('Please select a research domain first');
            return;
        }

        setDeploymentStatus('deploying');
        setDeploymentLog([]);
        
        try {
            // Add initial log entry
            addLogEntry('🚀 Starting deployment process...');
            addLogEntry(`📦 Domain: ${selectedDomain}`);
            addLogEntry(`🏗️ Instance: ${deploymentConfig.instanceSize} in ${deploymentConfig.region}`);
            
            // Simulate deployment API call
            const response = await fetch('/api/deploy', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    domain: selectedDomain,
                    config: deploymentConfig
                })
            });

            if (response.ok) {
                const result = await response.json();
                addLogEntry('✅ Deployment initiated successfully');
                setEstimatedTime(result.estimatedTime || '5-10 minutes');
                
                // Start monitoring deployment progress
                monitorDeployment(result.deploymentId);
                
                if (onDeploymentStart) {
                    onDeploymentStart(result.deploymentId);
                }
            } else {
                throw new Error('Deployment request failed');
            }
        } catch (error) {
            addLogEntry('❌ Deployment failed: ' + error.message);
            setDeploymentStatus('error');
        }
    };

    const addLogEntry = (message) => {
        const timestamp = new Date().toLocaleTimeString();
        setDeploymentLog(prev => [...prev, { timestamp, message }]);
    };

    const monitorDeployment = async (deploymentId) => {
        // Simulate deployment monitoring
        const steps = [
            '🔍 Validating configuration...',
            '📊 Calculating resource requirements...',
            '🏗️ Creating CloudFormation stack...',
            '💾 Provisioning storage volumes...',
            '🖥️ Launching compute instances...',
            '📦 Installing software packages...',
            '🔧 Configuring research environment...',
            '✅ Deployment completed successfully!'
        ];

        for (let i = 0; i < steps.length; i++) {
            await new Promise(resolve => setTimeout(resolve, 2000)); // 2 second delay
            addLogEntry(steps[i]);
            
            if (i === steps.length - 1) {
                setDeploymentStatus('success');
            }
        }
    };

    return (
        <div className="deployment-workflow">
            <h2>🚀 Deployment Workflow</h2>
            
            {!selectedDomain && (
                <div className="warning-message">
                    <p>⚠️ Please select a research domain from the Domains tab to continue.</p>
                </div>
            )}

            {selectedDomain && (
                <>
                    <div className="deployment-config">
                        <h3>Configuration</h3>
                        <div className="config-grid">
                            <div className="config-item">
                                <label>Instance Size:</label>
                                <select 
                                    value={deploymentConfig.instanceSize}
                                    onChange={(e) => setDeploymentConfig(prev => ({...prev, instanceSize: e.target.value}))}
                                    disabled={deploymentStatus === 'deploying'}
                                >
                                    <option value="small">Small (Development)</option>
                                    <option value="medium">Medium (Standard)</option>
                                    <option value="large">Large (Production)</option>
                                    <option value="xlarge">XLarge (High Performance)</option>
                                </select>
                            </div>

                            <div className="config-item">
                                <label>AWS Region:</label>
                                <select 
                                    value={deploymentConfig.region}
                                    onChange={(e) => setDeploymentConfig(prev => ({...prev, region: e.target.value}))}
                                    disabled={deploymentStatus === 'deploying'}
                                >
                                    <option value="us-east-1">US East (N. Virginia)</option>
                                    <option value="us-west-2">US West (Oregon)</option>
                                    <option value="eu-west-1">Europe (Ireland)</option>
                                    <option value="ap-southeast-1">Asia Pacific (Singapore)</option>
                                </select>
                            </div>

                            <div className="config-item checkbox-item">
                                <label>
                                    <input 
                                        type="checkbox"
                                        checked={deploymentConfig.useSpotInstances}
                                        onChange={(e) => setDeploymentConfig(prev => ({...prev, useSpotInstances: e.target.checked}))}
                                        disabled={deploymentStatus === 'deploying'}
                                    />
                                    Use Spot Instances (70% cost savings)
                                </label>
                            </div>

                            <div className="config-item checkbox-item">
                                <label>
                                    <input 
                                        type="checkbox"
                                        checked={deploymentConfig.autoShutdown}
                                        onChange={(e) => setDeploymentConfig(prev => ({...prev, autoShutdown: e.target.checked}))}
                                        disabled={deploymentStatus === 'deploying'}
                                    />
                                    Auto-shutdown when idle
                                </label>
                            </div>

                            {deploymentConfig.autoShutdown && (
                                <div className="config-item">
                                    <label>Shutdown Timeout (minutes):</label>
                                    <input 
                                        type="number"
                                        value={deploymentConfig.shutdownTimeout}
                                        onChange={(e) => setDeploymentConfig(prev => ({...prev, shutdownTimeout: parseInt(e.target.value)}))}
                                        disabled={deploymentStatus === 'deploying'}
                                        min="5"
                                        max="120"
                                    />
                                </div>
                            )}

                            <div className="config-item checkbox-item">
                                <label>
                                    <input 
                                        type="checkbox"
                                        checked={deploymentConfig.enableBackup}
                                        onChange={(e) => setDeploymentConfig(prev => ({...prev, enableBackup: e.target.checked}))}
                                        disabled={deploymentStatus === 'deploying'}
                                    />
                                    Enable automatic backups
                                </label>
                            </div>
                        </div>
                    </div>

                    <div className="deployment-actions">
                        <button 
                            className={`deploy-button ${deploymentStatus}`}
                            onClick={handleDeploy}
                            disabled={deploymentStatus === 'deploying'}
                        >
                            {deploymentStatus === 'deploying' ? '🔄 Deploying...' : '🚀 Deploy Environment'}
                        </button>

                        {estimatedTime && deploymentStatus === 'deploying' && (
                            <p className="estimated-time">⏱️ Estimated time: {estimatedTime}</p>
                        )}
                    </div>

                    {deploymentLog.length > 0 && (
                        <div className="deployment-log">
                            <h3>Deployment Log</h3>
                            <div className="log-entries">
                                {deploymentLog.map((entry, index) => (
                                    <div key={index} className="log-entry">
                                        <span className="timestamp">{entry.timestamp}</span>
                                        <span className="message">{entry.message}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </>
            )}
        </div>
    );
}

// Deployment Monitoring Component
function DeploymentMonitor({ deployments = [] }) {
    const [monitoringData, setMonitoringData] = React.useState({});
    const [selectedDeployment, setSelectedDeployment] = React.useState(null);

    React.useEffect(() => {
        // Simulate real-time monitoring data
        const interval = setInterval(() => {
            const now = Date.now();
            const newData = {};
            
            deployments.forEach(deployment => {
                newData[deployment.id] = {
                    cpuUsage: Math.floor(Math.random() * 80) + 10,
                    memoryUsage: Math.floor(Math.random() * 70) + 20,
                    diskUsage: Math.floor(Math.random() * 50) + 30,
                    networkIn: Math.floor(Math.random() * 1000) + 100,
                    networkOut: Math.floor(Math.random() * 800) + 50,
                    cost: (Math.random() * 10 + 5).toFixed(2),
                    uptime: Math.floor((now - deployment.startTime) / 1000),
                    status: deployment.status || 'running'
                };
            });
            
            setMonitoringData(newData);
        }, 3000); // Update every 3 seconds

        return () => clearInterval(interval);
    }, [deployments]);

    const formatUptime = (seconds) => {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        return `${hours}h ${minutes}m`;
    };

    return (
        <div className="deployment-monitor">
            <h2>📊 Deployment Monitoring</h2>
            
            {deployments.length === 0 && (
                <div className="no-deployments">
                    <p>No active deployments to monitor.</p>
                    <p>Deploy a research environment to see real-time monitoring data.</p>
                </div>
            )}

            {deployments.length > 0 && (
                <>
                    <div className="deployment-overview">
                        <h3>Active Deployments</h3>
                        <div className="deployment-cards">
                            {deployments.map(deployment => {
                                const data = monitoringData[deployment.id] || {};
                                return (
                                    <div 
                                        key={deployment.id}
                                        className={`deployment-card ${selectedDeployment === deployment.id ? 'selected' : ''}`}
                                        onClick={() => setSelectedDeployment(deployment.id)}
                                    >
                                        <h4>{deployment.domain}</h4>
                                        <div className="status-indicator">
                                            <span className={`status-dot ${data.status}`}></span>
                                            <span>{data.status || 'unknown'}</span>
                                        </div>
                                        <div className="quick-stats">
                                            <div className="stat">
                                                <span>CPU:</span>
                                                <span>{data.cpuUsage || 0}%</span>
                                            </div>
                                            <div className="stat">
                                                <span>Memory:</span>
                                                <span>{data.memoryUsage || 0}%</span>
                                            </div>
                                            <div className="stat">
                                                <span>Cost:</span>
                                                <span>${data.cost || '0.00'}/hr</span>
                                            </div>
                                        </div>
                                        <div className="uptime">
                                            ⏱️ {formatUptime(data.uptime || 0)}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </div>

                    {selectedDeployment && monitoringData[selectedDeployment] && (
                        <div className="detailed-monitoring">
                            <h3>Detailed Monitoring - {deployments.find(d => d.id === selectedDeployment)?.domain}</h3>
                            <div className="monitoring-grid">
                                <div className="metric-card">
                                    <h4>💻 CPU Usage</h4>
                                    <div className="metric-value">{monitoringData[selectedDeployment].cpuUsage}%</div>
                                    <div className="metric-bar">
                                        <div 
                                            className="metric-fill cpu"
                                            style={{width: `${monitoringData[selectedDeployment].cpuUsage}%`}}
                                        ></div>
                                    </div>
                                </div>

                                <div className="metric-card">
                                    <h4>🧠 Memory Usage</h4>
                                    <div className="metric-value">{monitoringData[selectedDeployment].memoryUsage}%</div>
                                    <div className="metric-bar">
                                        <div 
                                            className="metric-fill memory"
                                            style={{width: `${monitoringData[selectedDeployment].memoryUsage}%`}}
                                        ></div>
                                    </div>
                                </div>

                                <div className="metric-card">
                                    <h4>💾 Disk Usage</h4>
                                    <div className="metric-value">{monitoringData[selectedDeployment].diskUsage}%</div>
                                    <div className="metric-bar">
                                        <div 
                                            className="metric-fill disk"
                                            style={{width: `${monitoringData[selectedDeployment].diskUsage}%`}}
                                        ></div>
                                    </div>
                                </div>

                                <div className="metric-card">
                                    <h4>🌐 Network In</h4>
                                    <div className="metric-value">{monitoringData[selectedDeployment].networkIn} MB/s</div>
                                </div>

                                <div className="metric-card">
                                    <h4>📡 Network Out</h4>
                                    <div className="metric-value">{monitoringData[selectedDeployment].networkOut} MB/s</div>
                                </div>

                                <div className="metric-card">
                                    <h4>💰 Current Cost</h4>
                                    <div className="metric-value">${monitoringData[selectedDeployment].cost}/hr</div>
                                </div>
                            </div>

                            <div className="deployment-actions">
                                <button className="action-button warning">⏸️ Stop Instance</button>
                                <button className="action-button">🔄 Restart Instance</button>
                                <button className="action-button danger">🗑️ Terminate Deployment</button>
                            </div>
                        </div>
                    )}
                </>
            )}
        </div>
    );
}