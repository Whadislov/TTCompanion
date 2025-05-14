package api

import (
	"fmt"
	"log"
	"net/http"
)

// getUserIDFromCookie gets the user ID from the cookie
func getUserIDFromCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	// session verification, use same session name as in the login handler
	session, err := cookieStore.Get(r, "auth-session")
	if err != nil {
		log.Printf("getUserIDFromCookie: session err : %v", err)
		return "", fmt.Errorf("invalid session: %w", err)
	}

	// is the user already auth via session ?
	auth, ok := session.Values["authenticated"].(bool)
	if ok && auth {
		jwt, ok := session.Values["jwt"].(string)
		if ok {
			userID, err := verifyJWT(w, r, jwt, session)
			if err != nil {
				return "", fmt.Errorf("invalid session data")
			}
			return userID, nil
		}
	}
	return "", nil
}
