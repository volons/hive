package autopilot

import (
	"sync"
	"time"

	"github.com/volons/hive/libs/store"
)

type sessionKey string

var sessions sync.Map

// Session represents a vehicle controlling user's session
type Session struct {
	vehicleID string
	userID    string
	start     time.Time
	duration  time.Duration
	done      chan bool
	stop      chan bool
}

// StartSession starts a user session
func StartSession(vehicleID, userID string) *Session {
	session := newSession(vehicleID, userID, time.Minute*3)
	val, loaded := sessions.LoadOrStore(sessionKey(userID), session)
	if loaded {
		// if session already exists restore previous session
		session := val.(*Session)
		return session
	}

	go session.Start()
	return session
}

// newSession creates a new session
func newSession(vehicleID, userID string, duration time.Duration) *Session {
	return &Session{
		vehicleID: vehicleID,
		userID:    userID,
		start:     time.Now(),
		duration:  duration,
		done:      make(chan bool),
		stop:      make(chan bool),
	}
}

// Start runns the session
func (s *Session) Start() {
	select {
	case <-time.After(s.duration):
	case <-s.stop:
	}
	sessions.Delete(sessionKey(s.userID))
	user := store.Users.Get(s.userID)
	store.Users.Delete(user)
	close(s.done)
}

// Stop stops the session
// returns true if it was running
func (s *Session) Stop() bool {
	select {
	case s.stop <- true:
		return true
	default:
		return false
	}
}

// TimeLeft returns the time left on this session
func (s *Session) TimeLeft() time.Duration {
	elapsed := time.Since(s.start)
	left := s.duration - elapsed
	if left < 0 {
		left = 0
	}

	return left
}
