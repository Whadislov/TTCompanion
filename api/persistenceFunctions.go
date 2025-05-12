package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// Create a new session with gorilla. Session name = auth-session, the name needs to be static
func createSession(r *http.Request, credToken string, userID uuid.UUID) {
	session, _ := cookieStore.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["jwt"] = credToken
	session.Values["user_id"] = userID.String()
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
}

// Delete the session of the user
func deleteSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := cookieStore.Get(r, "auth-name")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		sendJSONError(w, "Could not clear session", "INTERNAL_ERROR", http.StatusInternalServerError)
		return fmt.Errorf("could not clear session")
	}
	return nil
}

// Creates a cookie that has the credtoken in its value and lives 3600 seconds
func createCookie(w http.ResponseWriter, credToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    credToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

}

// Delete the cookie of the user
func deleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}
