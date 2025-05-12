package api

import (
	"fmt"
	"net/http"
)

// isAuthenticated verify that the user is auth via session or cookie
// returns (isAuthenticated, userID, error).
func getUserIDFromCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	// session verification, use same session name as in the login handler
	session, err := cookieStore.Get(r, "auth-session")
	if err != nil {
		return "", fmt.Errorf("invalid session: %w", err)
	}

	// is the user already auth via session
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		userID, ok := session.Values["user_id"].(string)
		if !ok {
			return "", fmt.Errorf("invalid session data")
		}
		return userID, nil
	}

	// cookie verification
	cookie, err := r.Cookie("jwt_token")
	if err == nil {
		if userID, err := verifyJWT(w, r, cookie.Value, session); err == nil {
			return userID, nil
		}
	}

	return "", nil
}
