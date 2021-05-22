package auth

import (
	"crypto/subtle"
)

type Auth struct {
	Username string
	Password string
}

// Authenticate does a constant time comparison of a users password and username.
// This is critical to prevent timing attacks whereby a user might be able guess the
// credentials based on timing how long it takes to receive an error
func (a *Auth) Authenticate(username, password string) bool {
	if subtle.ConstantTimeCompare([]byte(username), []byte(a.Username)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(a.Password)) == 1 {
		return true
	}
	return false
}
