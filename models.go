package main

import "time"

// User represents a registered person who can create and manage their own tickets.
// The password is never stored in plain text - only its bcrypt hash is kept.
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // "-" means this field is never sent back in JSON responses
}

// Ticket represents a single support ticket created by a user.
// Only the user who created a ticket is allowed to view or update it.
type Ticket struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // one of: open, in_progress, closed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Allowed ticket statuses.
const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusClosed     = "closed"
)

// allowedTransitions defines which status changes are permitted.
// The assignment requires: open -> in_progress -> closed, and a closed
// ticket can never be reopened or moved back. We enforce that strictly:
// each ticket must move through the statuses one step at a time.
var allowedTransitions = map[string]string{
	StatusOpen:       StatusInProgress,
	StatusInProgress: StatusClosed,
	// StatusClosed has no entry, meaning no further transition is allowed.
}
