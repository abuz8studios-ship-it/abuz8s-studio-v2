package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthProvider interface for OAuth providers
type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*Token, error)
	GetUserInfo(ctx context.Context, token *Token) (*User, error)
}

// Token represents OAuth tokens
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// User represents OAuth user info
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Provider  string `json:"provider"`
	CreatedAt int64  `json:"created_at"`
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

// Manager handles OAuth authentication
type Manager struct {
	providers map[string]OAuthProvider
	states    map[string]time.Time // state verification
	users     map[string]*User     // user storage (in-memory for now)
}

// NewManager creates a new OAuth manager
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]OAuthProvider),
		states:    make(map[string]time.Time),
		users:     make(map[string]*User),
	}
}

// RegisterProvider registers an OAuth provider
func (m *Manager) RegisterProvider(name string, provider OAuthProvider) {
	m.providers[name] = provider
}

// GetAuthURL returns the OAuth authorization URL
func (m *Manager) GetAuthURL(provider string) (string, string, error) {
	p, ok := m.providers[provider]
	if !ok {
		return "", "", fmt.Errorf("provider not found: %s", provider)
	}

	state, err := generateState()
	if err != nil {
		return "", "", err
	}

	m.states[state] = time.Now().Add(10 * time.Minute)

	return p.GetAuthURL(state), state, nil
}

// HandleCallback handles OAuth callback
func (m *Manager) HandleCallback(ctx context.Context, provider, code, state string) (*User, error) {
	// Verify state
	if _, ok := m.states[state]; !ok {
		return nil, fmt.Errorf("invalid or expired state")
	}
	delete(m.states, state)

	p, ok := m.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", provider)
	}

	// Exchange code for token
	token, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info
	user, err := p.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user.Provider = provider
	user.CreatedAt = time.Now().Unix()

	// Store user
	m.users[user.ID] = user

	return user, nil
}

// GetUser returns a user by ID
func (m *Manager) GetUser(id string) (*User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// generateState generates a random state parameter
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ==================== GOOGLE OAUTH ====================

type GoogleProvider struct {
	config OAuthConfig
	client *http.Client
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		config: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
		},
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":     {p.config.ClientID},
		"redirect_uri":  {p.config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(p.config.Scopes, " ")},
		"state":         {state},
	}
	return p.config.AuthURL + "?" + v.Encode()
}

func (p *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*Token, error) {
	data := url.Values{
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {p.config.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
		TokenType:    result.TokenType,
	}, nil
}

func (p *GoogleProvider) GetUserInfo(ctx context.Context, token *Token) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &User{
		ID:     result.ID,
		Email:  result.Email,
		Name:   result.Name,
		Avatar: result.Picture,
	}, nil
}

// ==================== GITHUB OAUTH ====================

type GitHubProvider struct {
	config OAuthConfig
	client *http.Client
}

// NewGitHubProvider creates a new GitHub OAuth provider
func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		config: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email", "read:user"},
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
		},
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":    {p.config.ClientID},
		"redirect_uri": {p.config.RedirectURL},
		"scope":        {strings.Join(p.config.Scopes, " ")},
		"state":        {state},
	}
	return p.config.AuthURL + "?" + v.Encode()
}

func (p *GitHubProvider) ExchangeCode(ctx context.Context, code string) (*Token, error) {
	data := url.Values{
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {p.config.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &Token{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
	}, nil
}

func (p *GitHubProvider) GetUserInfo(ctx context.Context, token *Token) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token.AccessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Get primary email if not provided
	email := result.Email
	if email == "" {
		email = p.getPrimaryEmail(ctx, token)
	}

	return &User{
		ID:     fmt.Sprintf("%d", result.ID),
		Email:  email,
		Name:   result.Name,
		Avatar: result.AvatarURL,
	}, nil
}

func (p *GitHubProvider) getPrimaryEmail(ctx context.Context, token *Token) string {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "token "+token.AccessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	return ""
}

// ==================== DISCORD OAUTH ====================

type DiscordProvider struct {
	config OAuthConfig
	client *http.Client
}

// NewDiscordProvider creates a new Discord OAuth provider
func NewDiscordProvider(clientID, clientSecret, redirectURL string) *DiscordProvider {
	return &DiscordProvider{
		config: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"identify", "email"},
			AuthURL:      "https://discord.com/api/oauth2/authorize",
			TokenURL:     "https://discord.com/api/oauth2/token",
			UserInfoURL:  "https://discord.com/api/users/@me",
		},
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *DiscordProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":     {p.config.ClientID},
		"redirect_uri":  {p.config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(p.config.Scopes, " ")},
		"state":         {state},
	}
	return p.config.AuthURL + "?" + v.Encode()
}

func (p *DiscordProvider) ExchangeCode(ctx context.Context, code string) (*Token, error) {
	data := url.Values{
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {p.config.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
		TokenType:    result.TokenType,
	}, nil
}

func (p *DiscordProvider) GetUserInfo(ctx context.Context, token *Token) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		Avatar        string `json:"avatar"`
		Discriminator string `json:"discriminator"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	avatarURL := ""
	if result.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", result.ID, result.Avatar)
	}

	name := result.Username
	if result.Discriminator != "0" {
		name = fmt.Sprintf("%s#%s", result.Username, result.Discriminator)
	}

	return &User{
		ID:     result.ID,
		Email:  result.Email,
		Name:   name,
		Avatar: avatarURL,
	}, nil
}

// PKCE utils for enhanced security
func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
