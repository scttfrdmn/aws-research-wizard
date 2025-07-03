# TODO: Globus Auth Integration for AWS Research Wizard

## Future Enhancement: Add Globus Auth SSO Support

### Background
Globus Auth is a critical authentication provider for the research community, especially for:
- National labs and research institutions
- High-performance computing facilities
- Data-intensive research workflows
- Multi-institutional collaborations

### Implementation Plan

#### 1. Globus Auth Provider Integration
- **Location**: `go/internal/web/routes.go` - Add Globus SSO handler
- **Frontend**: Already prepared with Globus provider in auth.js and CSS styling
- **Library**: Use existing Go-based Globus library for OAuth2 integration

#### 2. Required Changes

**Backend API Endpoint**:
```go
func (s *Server) handleGlobusAuth(w http.ResponseWriter, r *http.Request) {
    // Implement Globus OAuth2 flow using the Go Globus library
    // - Redirect to Globus Auth
    // - Handle callback and token exchange
    // - Extract user information and permissions
    // - Create local session with Globus identity
}
```

**Authentication Configuration**:
```go
type GlobusConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURI  string
    Scopes       []string // openid, profile, email, urn:globus:auth:scope:transfer.api.globus.org:all
}
```

#### 3. Research Community Benefits
- **Seamless SSO**: Researchers already have Globus identities
- **Federated Identity**: Cross-institutional collaboration support
- **Data Transfer Integration**: Natural fit with Globus Transfer for data movement
- **Trust Framework**: Leverages established research identity federation

#### 4. Frontend Integration Status
✅ **Already Implemented**:
- Globus provider added to SSO provider list
- CSS styling for Globus SSO button (brand color: #1f5582)
- UI components ready for Globus authentication flow

#### 5. Implementation Priority
- **High Priority** for research institutions
- **Essential** for national labs and HPC centers
- **Strategic** for research data management workflows

### Next Steps
1. Integrate Go-based Globus library
2. Implement OAuth2 flow in backend routes
3. Add Globus-specific user profile handling
4. Test with research institution Globus identities
5. Document Globus Auth configuration for administrators

### Notes
- Frontend UI already includes Globus as 4th SSO provider option
- Styling and branding consistent with Globus identity guidelines
- Integration will complete the enterprise authentication suite for research environments
