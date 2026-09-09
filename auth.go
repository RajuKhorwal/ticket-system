package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret returns the key used to sign/verify tokens.
// Set JWT_SECRET in production. Falls back to a dev default otherwise.
func jwtSecret() []byte {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	return []byte("dev-only-secret-change-me")
}

// HashPassword hashes a plain-text password with bcrypt (salted, one-way).
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a plain-text password against a stored hash.
// Returns true if they match.
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// tokenExpiry controls how long a login session (JWT) stays valid.
const tokenExpiry = 24 * time.Hour

// GenerateToken creates a signed JWT carrying the user's ID as its subject.
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(tokenExpiry).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// ParseToken validates a JWT string and, if valid, returns the user ID
// that was embedded in it at login time.
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Make sure the token was signed with the algorithm we expect.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("invalid token subject")
	}
	return userID, nil
}
