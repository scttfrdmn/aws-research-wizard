// Enhanced GUI Phase 4: Advanced Notification System
// Email/Slack integration, real-time alerts, and notification management

// Notification Manager Component
function NotificationManager() {
    const [notifications, setNotifications] = React.useState([]);
    const [notificationSettings, setNotificationSettings] = React.useState({
        enableEmail: true,
        enableSlack: false,
        enableInApp: true,
        emailAddress: '',
        slackWebhook: '',
        alertThresholds: {
            costLimit: 100,
            deploymentFailure: true,
            performanceIssues: true,
            securityAlerts: true
        }
    });

    React.useEffect(() => {
        // Simulate notification subscription
        const notificationInterval = setInterval(() => {
            if (Math.random() > 0.7) { // 30% chance of notification
                generateRandomNotification();
            }
        }, 15000); // Check every 15 seconds

        return () => clearInterval(notificationInterval);
    }, []);

    const generateRandomNotification = () => {
        const notificationTypes = [
            {
                type: 'deployment',
                severity: 'info',
                title: 'Deployment Complete',
                message: 'Your genomics research environment has been successfully deployed.',
                icon: '🚀'
            },
            {
                type: 'cost',
                severity: 'warning',
                title: 'Cost Alert',
                message: 'Monthly costs approaching 80% of budget limit ($80 of $100).',
                icon: '💰'
            },
            {
                type: 'performance',
                severity: 'warning',
                title: 'Performance Alert',
                message: 'CPU usage has exceeded 85% for more than 10 minutes.',
                icon: '⚡'
            },
            {
                type: 'security',
                severity: 'critical',
                title: 'Security Alert',
                message: 'Unusual access pattern detected from new IP address.',
                icon: '🔒'
            }
        ];

        const notification = notificationTypes[Math.floor(Math.random() * notificationTypes.length)];
        addNotification({
            ...notification,
            id: `notif-${Date.now()}`,
            timestamp: new Date().toISOString(),
            read: false
        });
    };

    const addNotification = (notification) => {
        setNotifications(prev => [notification, ...prev.slice(0, 19)]); // Keep last 20

        // Send external notifications if enabled
        if (notificationSettings.enableEmail && notificationSettings.emailAddress) {
            sendEmailNotification(notification);
        }
        if (notificationSettings.enableSlack && notificationSettings.slackWebhook) {
            sendSlackNotification(notification);
        }
    };

    const sendEmailNotification = async (notification) => {
        try {
            await fetch('/api/notifications/email', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    to: notificationSettings.emailAddress,
                    subject: `AWS Research Wizard: ${notification.title}`,
                    message: notification.message,
                    severity: notification.severity
                })
            });
        } catch (error) {
            console.error('Failed to send email notification:', error);
        }
    };

    const sendSlackNotification = async (notification) => {
        try {
            await fetch('/api/notifications/slack', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    webhook: notificationSettings.slackWebhook,
                    text: `${notification.icon} *${notification.title}*\n${notification.message}`,
                    severity: notification.severity
                })
            });
        } catch (error) {
            console.error('Failed to send Slack notification:', error);
        }
    };

    const markAsRead = (notificationId) => {
        setNotifications(prev =>
            prev.map(notif =>
                notif.id === notificationId ? { ...notif, read: true } : notif
            )
        );
    };

    const markAllAsRead = () => {
        setNotifications(prev =>
            prev.map(notif => ({ ...notif, read: true }))
        );
    };

    const deleteNotification = (notificationId) => {
        setNotifications(prev =>
            prev.filter(notif => notif.id !== notificationId)
        );
    };

    const clearAllNotifications = () => {
        setNotifications([]);
    };

    return {
        notifications,
        notificationSettings,
        setNotificationSettings,
        addNotification,
        markAsRead,
        markAllAsRead,
        deleteNotification,
        clearAllNotifications
    };
}

// Notification Bell Component
function NotificationBell({ notificationManager }) {
    const [showDropdown, setShowDropdown] = React.useState(false);
    const unreadCount = notificationManager.notifications.filter(n => !n.read).length;

    return (
        <div className="notification-bell">
            <button
                className="bell-trigger"
                onClick={() => setShowDropdown(!showDropdown)}
            >
                🔔
                {unreadCount > 0 && (
                    <span className="notification-badge">{unreadCount > 99 ? '99+' : unreadCount}</span>
                )}
            </button>

            {showDropdown && (
                <div className="notification-dropdown">
                    <div className="notification-header">
                        <h3>Notifications</h3>
                        <div className="notification-actions">
                            <button
                                className="mark-all-read"
                                onClick={notificationManager.markAllAsRead}
                                disabled={unreadCount === 0}
                            >
                                Mark All Read
                            </button>
                            <button
                                className="clear-all"
                                onClick={notificationManager.clearAllNotifications}
                                disabled={notificationManager.notifications.length === 0}
                            >
                                Clear All
                            </button>
                        </div>
                    </div>

                    <div className="notification-list">
                        {notificationManager.notifications.length === 0 && (
                            <div className="no-notifications">
                                <p>No notifications</p>
                            </div>
                        )}

                        {notificationManager.notifications.map(notification => (
                            <div
                                key={notification.id}
                                className={`notification-item ${notification.severity} ${notification.read ? 'read' : 'unread'}`}
                            >
                                <div className="notification-icon">
                                    {notification.icon}
                                </div>
                                <div className="notification-content">
                                    <div className="notification-title">
                                        {notification.title}
                                    </div>
                                    <div className="notification-message">
                                        {notification.message}
                                    </div>
                                    <div className="notification-time">
                                        {new Date(notification.timestamp).toLocaleString()}
                                    </div>
                                </div>
                                <div className="notification-controls">
                                    {!notification.read && (
                                        <button
                                            className="mark-read-btn"
                                            onClick={() => notificationManager.markAsRead(notification.id)}
                                        >
                                            ✓
                                        </button>
                                    )}
                                    <button
                                        className="delete-btn"
                                        onClick={() => notificationManager.deleteNotification(notification.id)}
                                    >
                                        ✕
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}

// Notification Settings Component
function NotificationSettings({ notificationManager }) {
    const { notificationSettings, setNotificationSettings } = notificationManager;

    const handleSettingChange = (setting, value) => {
        setNotificationSettings(prev => ({
            ...prev,
            [setting]: value
        }));
    };

    const handleThresholdChange = (threshold, value) => {
        setNotificationSettings(prev => ({
            ...prev,
            alertThresholds: {
                ...prev.alertThresholds,
                [threshold]: value
            }
        }));
    };

    const testNotification = () => {
        notificationManager.addNotification({
            id: `test-${Date.now()}`,
            type: 'test',
            severity: 'info',
            title: 'Test Notification',
            message: 'This is a test notification to verify your settings.',
            icon: '🧪',
            timestamp: new Date().toISOString(),
            read: false
        });
    };

    return (
        <div className="notification-settings">
            <h2>🔔 Notification Settings</h2>

            <div className="settings-section">
                <h3>Notification Channels</h3>
                <div className="settings-grid">
                    <div className="setting-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={notificationSettings.enableInApp}
                                onChange={(e) => handleSettingChange('enableInApp', e.target.checked)}
                            />
                            In-App Notifications
                        </label>
                    </div>

                    <div className="setting-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={notificationSettings.enableEmail}
                                onChange={(e) => handleSettingChange('enableEmail', e.target.checked)}
                            />
                            Email Notifications
                        </label>
                        {notificationSettings.enableEmail && (
                            <input
                                type="email"
                                placeholder="Enter email address"
                                value={notificationSettings.emailAddress}
                                onChange={(e) => handleSettingChange('emailAddress', e.target.value)}
                            />
                        )}
                    </div>

                    <div className="setting-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={notificationSettings.enableSlack}
                                onChange={(e) => handleSettingChange('enableSlack', e.target.checked)}
                            />
                            Slack Notifications
                        </label>
                        {notificationSettings.enableSlack && (
                            <input
                                type="url"
                                placeholder="Slack webhook URL"
                                value={notificationSettings.slackWebhook}
                                onChange={(e) => handleSettingChange('slackWebhook', e.target.value)}
                            />
                        )}
                    </div>
                </div>
            </div>

            <div className="settings-section">
                <h3>Alert Thresholds</h3>
                <div className="settings-grid">
                    <div className="setting-item">
                        <label>Monthly Cost Limit ($):</label>
                        <input
                            type="number"
                            value={notificationSettings.alertThresholds.costLimit}
                            onChange={(e) => handleThresholdChange('costLimit', parseInt(e.target.value))}
                            min="1"
                        />
                    </div>

                    <div className="setting-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={notificationSettings.alertThresholds.deploymentFailure}
                                onChange={(e) => handleThresholdChange('deploymentFailure', e.target.checked)}
                            />
                            Deployment Failures
                        </label>
                    </div>

                    <div className="setting-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={notificationSettings.alertThresholds.performanceIssues}
                                onChange={(e) => handleThresholdChange('performanceIssues', e.target.checked)}
                            />
                            Performance Issues
                        </label>
                    </div>

                    <div className="setting-item">
                        <label>
                            <input
                                type="checkbox"
                                checked={notificationSettings.alertThresholds.securityAlerts}
                                onChange={(e) => handleThresholdChange('securityAlerts', e.target.checked)}
                            />
                            Security Alerts
                        </label>
                    </div>
                </div>
            </div>

            <div className="settings-actions">
                <button
                    className="test-notification-btn"
                    onClick={testNotification}
                >
                    🧪 Send Test Notification
                </button>
            </div>
        </div>
    );
}

// Toast Notification Component
function ToastNotifications({ notificationManager }) {
    const [toasts, setToasts] = React.useState([]);

    React.useEffect(() => {
        // Show toast for new notifications
        const unreadNotifications = notificationManager.notifications.filter(n => !n.read);
        const newNotifications = unreadNotifications.slice(0, 3); // Show max 3 toasts

        setToasts(newNotifications.map(notif => ({
            ...notif,
            toastId: `toast-${notif.id}`
        })));

        // Auto-dismiss toasts after 5 seconds
        const timeouts = newNotifications.map(notif =>
            setTimeout(() => {
                setToasts(prev => prev.filter(toast => toast.id !== notif.id));
            }, 5000)
        );

        return () => timeouts.forEach(clearTimeout);
    }, [notificationManager.notifications]);

    const dismissToast = (notificationId) => {
        setToasts(prev => prev.filter(toast => toast.id !== notificationId));
    };

    return (
        <div className="toast-container">
            {toasts.map(toast => (
                <div
                    key={toast.toastId}
                    className={`toast-notification ${toast.severity}`}
                >
                    <div className="toast-icon">{toast.icon}</div>
                    <div className="toast-content">
                        <div className="toast-title">{toast.title}</div>
                        <div className="toast-message">{toast.message}</div>
                    </div>
                    <button
                        className="toast-dismiss"
                        onClick={() => dismissToast(toast.id)}
                    >
                        ✕
                    </button>
                </div>
            ))}
        </div>
    );
}
