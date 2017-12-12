package auth

import (
	"strings"

	"config"
)

// Auth as the type of the auth-system
type Auth string

// Verify ...
func (au Auth) Verify(auTest Auth) bool {
	return au == auTest
}

func (au Auth) size() int {
	return len(au)
}

func (au Auth) String() string {

	if config.DebugMode() {
		return string(au)
	}

	return strings.Repeat("*", au.size())
}
