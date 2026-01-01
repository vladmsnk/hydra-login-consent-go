package idp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LDAPConfig holds the configuration for LDAP connection
type LDAPConfig struct {
	// Server address (e.g., "ldap.example.com:389" or "ldaps://ldap.example.com:636")
	Server string
	// BaseDN is the base distinguished name for searches (e.g., "dc=example,dc=com")
	BaseDN string
	// BindDN is the DN used to bind for searching users (e.g., "cn=admin,dc=example,dc=com")
	// Leave empty for anonymous bind
	BindDN string
	// BindPassword is the password for the bind DN
	BindPassword string
	// UserSearchFilter is the LDAP filter to find users, use %s as placeholder for username
	// Example: "(uid=%s)" or "(sAMAccountName=%s)" for AD
	UserSearchFilter string
	// UserSearchAttr is the attribute that contains the username (e.g., "uid", "cn", "sAMAccountName")
	UserSearchAttr string
	// UseTLS enables TLS connection (ldaps://)
	UseTLS bool
	// InsecureSkipVerify skips TLS certificate verification (not recommended for production)
	InsecureSkipVerify bool
	// ConnectionTimeout is the timeout for LDAP connection
	ConnectionTimeout time.Duration
}

// LDAPProvider implements IdentityProvider for LDAP authentication
type LDAPProvider struct {
	config LDAPConfig
}

// NewLDAPProvider creates a new LDAP identity provider
func NewLDAPProvider(config LDAPConfig) *LDAPProvider {
	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = 10 * time.Second
	}
	if config.UserSearchFilter == "" {
		config.UserSearchFilter = "(uid=%s)"
	}
	if config.UserSearchAttr == "" {
		config.UserSearchAttr = "uid"
	}
	return &LDAPProvider{config: config}
}

// Ping tests the LDAP connection and service account bind
func (p *LDAPProvider) Ping() error {
	conn, err := p.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to LDAP server %s: %w", p.config.Server, err)
	}
	defer conn.Close()

	// Test service account bind if configured
	if p.config.BindDN != "" {
		err = conn.Bind(p.config.BindDN, p.config.BindPassword)
		if err != nil {
			return fmt.Errorf("failed to bind with service account %s: %w", p.config.BindDN, err)
		}
	}

	return nil
}

// Authenticate verifies user credentials against LDAP
func (p *LDAPProvider) Authenticate(ctx context.Context, username, password string) (bool, error) {
	if username == "" || password == "" {
		return false, ErrInvalidCredentials
	}

	conn, err := p.connect()
	if err != nil {
		return false, fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()

	// First, bind with service account to search for user
	if p.config.BindDN != "" {
		err = conn.Bind(p.config.BindDN, p.config.BindPassword)
		if err != nil {
			return false, fmt.Errorf("ldap service bind: %w", err)
		}
	}

	// Search for the user
	userDN, err := p.findUserDN(conn, username)
	if err != nil {
		return false, err
	}

	// Attempt to bind with the user's credentials
	err = conn.Bind(userDN, password)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return false, ErrInvalidCredentials
		}
		return false, fmt.Errorf("ldap user bind: %w", err)
	}

	return true, nil
}

// connect establishes a connection to the LDAP server
func (p *LDAPProvider) connect() (*ldap.Conn, error) {
	var conn *ldap.Conn
	var err error

	dialer := &net.Dialer{
		Timeout: p.config.ConnectionTimeout,
	}

	if p.config.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: p.config.InsecureSkipVerify,
		}
		conn, err = ldap.DialURL(p.config.Server,
			ldap.DialWithTLSConfig(tlsConfig),
			ldap.DialWithDialer(dialer),
		)
	} else {
		conn, err = ldap.DialURL(p.config.Server,
			ldap.DialWithDialer(dialer),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("dial ldap server: %w", err)
	}

	return conn, nil
}

// findUserDN searches for a user and returns their DN
func (p *LDAPProvider) findUserDN(conn *ldap.Conn, username string) (string, error) {
	searchFilter := fmt.Sprintf(p.config.UserSearchFilter, ldap.EscapeFilter(username))

	searchRequest := ldap.NewSearchRequest(
		p.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,  // Size limit - we only need one result
		10, // Time limit in seconds
		false,
		searchFilter,
		[]string{"dn", p.config.UserSearchAttr},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return "", fmt.Errorf("ldap search: %w", err)
	}

	if len(result.Entries) == 0 {
		return "", ErrUserNotFound
	}

	if len(result.Entries) > 1 {
		return "", fmt.Errorf("multiple users found for username: %s", username)
	}

	return result.Entries[0].DN, nil
}
