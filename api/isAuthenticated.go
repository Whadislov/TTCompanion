package api

import (
	"fmt"
	"net/http"
	"strings"
)

// isAuthenticated verify that the user is auth via session or cookie
// returns (isAuthenticated, userID, error).
func isAuthenticated(w http.ResponseWriter, r *http.Request) (bool, string, error) {
	// session verification
	session, err := cookieStore.Get(r, "session-name")
	if err != nil {
		return false, "", fmt.Errorf("invalid session: %w", err)
	}

	// is the user already auth via session
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		userID, ok := session.Values["user_id"].(string)
		if !ok {
			return false, "", fmt.Errorf("invalid session data")
		}
		return true, userID, nil
	}

	// cookie verification
	cookie, err := r.Cookie("jwt_token")
	if err == nil {
		if userID, err := verifyJWT(w, r, cookie.Value, session); err == nil {
			return true, userID, nil
		}
	}

	// authorization header verification for compatibility
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		credTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if credTokenString != "" {
			if userID, err := verifyJWT(w, r, credTokenString, session); err == nil {
				return true, userID, nil
			}
		}
	}

	return false, "", nil
}
