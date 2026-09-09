package main

import "errors"

// These are the "known" error cases the service can run into.
// Keeping them as named errors makes handler code easy to read,
// e.g. `if err == ErrEmailTaken`.
var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTicketNotFound     = errors.New("ticket not found")
	ErrInvalidTransition  = errors.New("invalid status transition")
)
