package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// Create a new session with gorilla. Session name = auth-session, the name needs to be static
func createSession(w http.ResponseWriter, r *http.Request, credToken string) {
	session, _ := cookieStore.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["jwt"] = credToken
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	session.Save(r, w)

	log.Println("API : session created.", fmt.Sprintf(`Values :
session.Values["authenticated"] : %v
session.Values["jwt"] : %v`, session.Values["authenticated"], session.Values["jwt"]))

}

// Delete the session of the user
func deleteSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := cookieStore.Get(r, "auth-session")
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		sendJSONError(w, "Could not clear session", "INTERNAL_ERROR", http.StatusInternalServerError)
		return fmt.Errorf("could not clear session")
	}
	return nil
}
