package main

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is a private type used to avoid collisions when storing
// values in a request's context.
type contextKey string

const userIDContextKey contextKey = "userID"

// RequireAuth wraps a handler so that it can only be reached with a valid
// JWT. It reads the "Authorization: Bearer <token>" header, verifies the
// token, and - if valid - passes the logged-in user's ID down to the
// actual handler via the request context.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			writeError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "Authorization header must be: Bearer <token>")
			return
		}

		userID, err := ParseToken(parts[1])
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Attach the user ID so handlers know who's calling.
		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// userIDFromContext reads the logged-in user's ID that RequireAuth stored earlier.
func userIDFromContext(r *http.Request) string {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	return userID
}
