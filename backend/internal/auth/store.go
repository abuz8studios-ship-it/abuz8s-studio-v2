package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Store manages user authentication storage
type Store struct {
	mu       sync.RWMutex
	dataDir  string
	users    map[string]*User
	sessions map[string]string // sessionID -> userID
}

// NewStore creates a new auth store
func NewStore(dataDir string) *Store {
	return &Store{
		dataDir:  dataDir,
		users:    make(map[string]*User),
		sessions: make(map[string]string),
	}
}

// SaveUser saves a user to storage
func (s *Store) SaveUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.ID] = user

	// Save to file
	data, err := json.MarshalIndent(s.users, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(s.dataDir, "users.json")
	return os.WriteFile(path, data, 0600)
}

// GetUser retrieves a user by ID
func (s *Store) GetUser(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Store) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// CreateSession creates a new session for a user
func (s *Store) CreateSession(userID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := generateSessionID()
	s.sessions[sessionID] = userID

	return sessionID, nil
}

// GetSession retrieves the user ID for a session
func (s *Store) GetSession(sessionID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, ok := s.sessions[sessionID]
	if !ok {
		return "", ErrSessionNotFound
	}
	return userID, nil
}

// DeleteSession removes a session
func (s *Store) DeleteSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
	return nil
}

// Load loads users from storage
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dataDir, "users.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.users)
}

// Errors
var (
	ErrUserNotFound    = &AuthError{Message: "user not found"}
	ErrSessionNotFound = &AuthError{Message: "session not found"}
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
