package main

import (
	"strconv"
	"sync"
	"time"
)

// Store keeps all data in memory. No external database.
// RWMutex guards concurrent access since each request runs on its own goroutine.
type Store struct {
	mu sync.RWMutex

	usersByEmail map[string]*User
	usersByID    map[string]*User
	tickets      map[string]*Ticket

	nextUserID   int
	nextTicketID int
}

// NewStore creates an empty, ready-to-use in-memory store.
func NewStore() *Store {
	return &Store{
		usersByEmail: make(map[string]*User),
		usersByID:    make(map[string]*User),
		tickets:      make(map[string]*Ticket),
	}
}

// ---------- User operations ----------

// CreateUser saves a new user. It returns an error if the email is already taken.
func (s *Store) CreateUser(email, passwordHash string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByEmail[email]; exists {
		return nil, ErrEmailTaken
	}

	s.nextUserID++
	user := &User{
		ID:           "u" + strconv.Itoa(s.nextUserID),
		Email:        email,
		PasswordHash: passwordHash,
	}
	s.usersByEmail[email] = user
	s.usersByID[user.ID] = user
	return user, nil
}

// GetUserByEmail looks up a user by their email address (used during login).
func (s *Store) GetUserByEmail(email string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.usersByEmail[email]
	return u, ok
}

// ---------- Ticket operations ----------

// CreateTicket saves a new ticket owned by userID, starting in "open" status.
func (s *Store) CreateTicket(userID, title, description string) *Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextTicketID++
	now := time.Now().UTC()
	ticket := &Ticket{
		ID:          "t" + strconv.Itoa(s.nextTicketID),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tickets[ticket.ID] = ticket
	return ticket
}

// ListTicketsByUser returns every ticket that belongs to the given user.
func (s *Store) ListTicketsByUser(userID string) []*Ticket {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := []*Ticket{}
	for _, t := range s.tickets {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	return result
}

// GetTicketByID fetches a ticket by ID regardless of owner.
// Callers must check ownership themselves.
func (s *Store) GetTicketByID(id string) (*Ticket, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tickets[id]
	return t, ok
}

// UpdateTicketStatus changes a ticket's status in place.
func (s *Store) UpdateTicketStatus(ticket *Ticket, newStatus string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ticket.Status = newStatus
	ticket.UpdatedAt = time.Now().UTC()
}
