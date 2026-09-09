package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
)
// registerRequest is the expected JSON body for POST /auth/register.
type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleRegister creates a new user account.
// It hashes the password before saving anything - the plain password
// is never written to storage.
func (app *App) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid email format")
		return
	}
	req.Email = strings.ToLower(addr.Address)

	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}
	if len(req.Password) > 72 {
		writeError(w, http.StatusBadRequest, "password must be at most 72 characters")
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not process password")
		return
	}

	user, err := app.Store.CreateUser(req.Email, hash)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":    user.ID,
		"email": user.Email,
	})
}

// loginRequest is the expected JSON body for POST /auth/login.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleLogin checks the user's credentials and, if correct, returns a JWT
// that the client must send on every future protected request.
func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, ok := app.Store.GetUserByEmail(req.Email)
	if !ok || !CheckPassword(user.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
