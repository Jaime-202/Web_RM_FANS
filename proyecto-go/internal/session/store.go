package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	UserID    string
	Role      string
	ExpiresAt time.Time
}

type Store struct {
	sessions map[string]Session
	mu       sync.RWMutex
}

var GlobalStore = &Store{
	sessions: make(map[string]Session),
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Store) CreateSession(userID, role string) (string, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := generateToken()
	expires := time.Now().Add(24 * time.Hour)
	s.sessions[token] = Session{
		UserID:    userID,
		Role:      role,
		ExpiresAt: expires,
	}
	return token, expires
}

func (s *Store) GetSession(token string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return Session{}, false
	}
	if time.Now().After(session.ExpiresAt) {
		// Session expired
		// Optionally we can delete it here, but read lock prevents it.
		// A background job could clean expired sessions.
		return Session{}, false
	}
	return session, true
}

func (s *Store) DeleteSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}
