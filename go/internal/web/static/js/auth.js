// Enhanced GUI Phase 4: Enterprise Authentication & User Management
// SSO integration, user sessions, and role-based access control

// Authentication Manager Component
function AuthManager({ onAuthChange }) {
    const [authState, setAuthState] = React.useState({
        isAuthenticated: false,
        user: null,
        permissions: [],
        sessionExpiry: null
    });
    const [authConfig, setAuthConfig] = React.useState({
        enableSSO: true,
        providers: ['okta', 'azure', 'google', 'globus'],
        sessionTimeout: 8, // hours
        requireMFA: false
    });

    React.useEffect(() => {
        // Check for existing session on component mount
        checkExistingSession();

        // Set up session monitoring
        const sessionCheck = setInterval(checkSessionValidity, 60000); // Check every minute

        return () => clearInterval(sessionCheck);
    }, []);

    const checkExistingSession = async () => {
        try {
            const response = await fetch('/api/auth/session', {
                credentials: 'include'
            });

            if (response.ok) {
                const sessionData = await response.json();
                setAuthState({
                    isAuthenticated: true,
                    user: sessionData.user,
                    permissions: sessionData.permissions,
                    sessionExpiry: new Date(sessionData.expires)
                });
                onAuthChange(true, sessionData.user);
            }
        } catch (error) {
            console.log('No existing session found');
        }
    };

    const checkSessionValidity = () => {
        if (authState.sessionExpiry && new Date() > authState.sessionExpiry) {
            handleLogout();
        }
    };

    const handleSSOLogin = async (provider) => {
        try {
            // Redirect to SSO provider
            window.location.href = `/api/auth/sso/${provider}`;
        } catch (error) {
            console.error('SSO login failed:', error);
        }
    };

    const handleLocalLogin = async (credentials) => {
        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(credentials)
            });

            if (response.ok) {
                const authData = await response.json();
                setAuthState({
                    isAuthenticated: true,
                    user: authData.user,
                    permissions: authData.permissions,
                    sessionExpiry: new Date(authData.expires)
                });
                onAuthChange(true, authData.user);
            } else {
                throw new Error('Login failed');
            }
        } catch (error) {
            console.error('Local login failed:', error);
            throw error;
        }
    };

    const handleLogout = async () => {
        try {
            await fetch('/api/auth/logout', {
                method: 'POST',
                credentials: 'include'
            });
        } catch (error) {
            console.error('Logout request failed:', error);
        } finally {
            setAuthState({
                isAuthenticated: false,
                user: null,
                permissions: [],
                sessionExpiry: null
            });
            onAuthChange(false, null);
        }
    };

    return {
        authState,
        authConfig,
        handleSSOLogin,
        handleLocalLogin,
        handleLogout,
        setAuthConfig
    };
}

// Login Component
function LoginForm({ authManager }) {
    const [loginMode, setLoginMode] = React.useState('sso');
    const [credentials, setCredentials] = React.useState({
        username: '',
        password: '',
        mfaCode: ''
    });
    const [isLogging, setIsLogging] = React.useState(false);
    const [error, setError] = React.useState('');

    const handleLocalLogin = async (e) => {
        e.preventDefault();
        setIsLogging(true);
        setError('');

        try {
            await authManager.handleLocalLogin(credentials);
        } catch (error) {
            setError('Invalid credentials. Please try again.');
        } finally {
            setIsLogging(false);
        }
    };

    return (
        <div className="login-container">
            <div className="login-form">
                <h2>🔐 AWS Research Wizard Login</h2>

                <div className="login-mode-selector">
                    <button
                        className={`mode-btn ${loginMode === 'sso' ? 'active' : ''}`}
                        onClick={() => setLoginMode('sso')}
                    >
                        SSO Login
                    </button>
                    <button
                        className={`mode-btn ${loginMode === 'local' ? 'active' : ''}`}
                        onClick={() => setLoginMode('local')}
                    >
                        Local Login
                    </button>
                </div>

                {loginMode === 'sso' && (
                    <div className="sso-providers">
                        <h3>Single Sign-On</h3>
                        <div className="provider-buttons">
                            {authManager.authConfig.providers.map(provider => (
                                <button
                                    key={provider}
                                    className={`sso-btn ${provider}`}
                                    onClick={() => authManager.handleSSOLogin(provider)}
                                >
                                    Sign in with {provider.charAt(0).toUpperCase() + provider.slice(1)}
                                </button>
                            ))}
                        </div>
                    </div>
                )}

                {loginMode === 'local' && (
                    <form onSubmit={handleLocalLogin} className="local-login-form">
                        <div className="form-group">
                            <label>Username:</label>
                            <input
                                type="text"
                                value={credentials.username}
                                onChange={(e) => setCredentials(prev => ({
                                    ...prev,
                                    username: e.target.value
                                }))}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label>Password:</label>
                            <input
                                type="password"
                                value={credentials.password}
                                onChange={(e) => setCredentials(prev => ({
                                    ...prev,
                                    password: e.target.value
                                }))}
                                required
                            />
                        </div>

                        {authManager.authConfig.requireMFA && (
                            <div className="form-group">
                                <label>MFA Code:</label>
                                <input
                                    type="text"
                                    value={credentials.mfaCode}
                                    onChange={(e) => setCredentials(prev => ({
                                        ...prev,
                                        mfaCode: e.target.value
                                    }))}
                                    placeholder="6-digit code"
                                />
                            </div>
                        )}

                        {error && <div className="error-message">{error}</div>}

                        <button
                            type="submit"
                            className="login-btn"
                            disabled={isLogging}
                        >
                            {isLogging ? '🔄 Signing In...' : '🔑 Sign In'}
                        </button>
                    </form>
                )}
            </div>
        </div>
    );
}

// User Profile Component
function UserProfile({ user, authManager }) {
    const [showProfile, setShowProfile] = React.useState(false);

    const formatPermissions = (permissions) => {
        return permissions.map(perm => ({
            resource: perm.split(':')[0],
            actions: perm.split(':')[1]?.split(',') || ['read']
        }));
    };

    return (
        <div className="user-profile">
            <button
                className="profile-trigger"
                onClick={() => setShowProfile(!showProfile)}
            >
                👤 {user.name || user.username}
            </button>

            {showProfile && (
                <div className="profile-dropdown">
                    <div className="profile-info">
                        <h3>User Information</h3>
                        <div className="user-details">
                            <div className="detail-item">
                                <span className="label">Name:</span>
                                <span className="value">{user.name || 'Not provided'}</span>
                            </div>
                            <div className="detail-item">
                                <span className="label">Email:</span>
                                <span className="value">{user.email || 'Not provided'}</span>
                            </div>
                            <div className="detail-item">
                                <span className="label">Role:</span>
                                <span className="value">{user.role || 'User'}</span>
                            </div>
                            <div className="detail-item">
                                <span className="label">Session Expires:</span>
                                <span className="value">
                                    {authManager.authState.sessionExpiry?.toLocaleString() || 'Unknown'}
                                </span>
                            </div>
                        </div>
                    </div>

                    <div className="permissions-info">
                        <h3>Permissions</h3>
                        <div className="permissions-list">
                            {formatPermissions(authManager.authState.permissions).map((perm, index) => (
                                <div key={index} className="permission-item">
                                    <span className="resource">{perm.resource}</span>
                                    <div className="actions">
                                        {perm.actions.map(action => (
                                            <span key={action} className="action-tag">{action}</span>
                                        ))}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="profile-actions">
                        <button
                            className="logout-btn"
                            onClick={authManager.handleLogout}
                        >
                            🚪 Sign Out
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}

// Permission Guard Component
function PermissionGuard({ permissions, requiredPermission, children, fallback }) {
    const hasPermission = (required) => {
        if (!required) return true;

        return permissions.some(perm => {
            const [resource, actions] = perm.split(':');
            const [reqResource, reqAction] = required.split(':');

            if (resource === '*' || resource === reqResource) {
                if (!reqAction) return true;
                const permActions = actions?.split(',') || [];
                return permActions.includes('*') || permActions.includes(reqAction);
            }

            return false;
        });
    };

    if (hasPermission(requiredPermission)) {
        return children;
    }

    return fallback || (
        <div className="permission-denied">
            <h3>🚫 Access Denied</h3>
            <p>You don't have permission to access this feature.</p>
            <p>Required permission: <code>{requiredPermission}</code></p>
        </div>
    );
}

// Session Status Component
function SessionStatus({ authState }) {
    const [timeRemaining, setTimeRemaining] = React.useState('');

    React.useEffect(() => {
        if (!authState.sessionExpiry) return;

        const updateTimer = () => {
            const now = new Date();
            const expiry = authState.sessionExpiry;
            const diff = expiry - now;

            if (diff <= 0) {
                setTimeRemaining('Expired');
                return;
            }

            const hours = Math.floor(diff / (1000 * 60 * 60));
            const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
            setTimeRemaining(`${hours}h ${minutes}m`);
        };

        updateTimer();
        const interval = setInterval(updateTimer, 60000);

        return () => clearInterval(interval);
    }, [authState.sessionExpiry]);

    if (!authState.isAuthenticated) return null;

    return (
        <div className="session-status">
            <span className="session-indicator">
                🕐 Session: {timeRemaining}
            </span>
        </div>
    );
}

// Role-Based Access Control Component
function RoleBasedAccess({ userRole, allowedRoles, children, fallback }) {
    const hasAccess = allowedRoles.includes(userRole) || allowedRoles.includes('*');

    if (hasAccess) {
        return children;
    }

    return fallback || (
        <div className="role-denied">
            <h3>🔒 Role Access Required</h3>
            <p>This feature requires one of the following roles:</p>
            <ul>
                {allowedRoles.map(role => (
                    <li key={role}><code>{role}</code></li>
                ))}
            </ul>
            <p>Your current role: <code>{userRole}</code></p>
        </div>
    );
}
