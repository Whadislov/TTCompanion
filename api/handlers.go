package api

import (
	"encoding/json"
	"log"
	"net/http"

	mdb "github.com/Whadislov/TTCompanion/internal/my_db"
	mf "github.com/Whadislov/TTCompanion/internal/my_functions"
	mt "github.com/Whadislov/TTCompanion/internal/my_types"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// Handler for loading the database
func loadUserDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to load user DB")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "User is unidentified", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Failed to parse user ID", http.StatusUnauthorized)
		return
	}

	mdb.SetUserIDOfSession(id)
	db, err := mdb.LoadDB()
	if err != nil {
		http.Error(w, "Failed to connect to database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db)
}

// Handler for saving the local changes to the database
func saveDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var db mt.Database
	err := json.NewDecoder(r.Body).Decode(&db)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	err = mdb.SaveDB(&db)
	if err != nil {
		http.Error(w, "Failed to save database", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler to check if the API is ready to take requests (not yet used)
func IsApiReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is ready !"))
}

// loginHandler process the request to analyse the credentials, returns a Credential token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendJSONError(w, "Invalid request", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	userID, err := checkUserCredentials(creds.Username, creds.Password)
	if err != nil {
		sendJSONError(w, "Invalid username or password", "INVALID_USERNAME_OR_PASSWORD", http.StatusUnauthorized)
		return
	}

	credToken, err := generateJWT(userID)
	if err != nil {
		sendJSONError(w, "Could not generate credToken", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	// Create a new session with gorilla
	session, _ := cookieStore.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["jwt"] = credToken
	session.Values["user_id"] = userID.String()
	session.Options = &sessions.Options{
		Path:     "/" + userID.String(),
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// Send JWT in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    credToken,
		Path:     "/" + userID.String(),
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"cred_token": credToken})
}

// signUpHandler process the request to create a new user, returns a Credential token
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&signUpData); err != nil {
		sendJSONError(w, "Invalid request", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	// Load the whole database to register the new user. Could be optimised to request directly postgres
	db, err := mdb.LoadUsersOnly()
	if err != nil {
		sendJSONError(w, "Could not load database to create the new user", "UNABLE_TO_LOAD_DATABASE", http.StatusInternalServerError)
		return
	}

	// Verify if the user already exists
	value, _ := checkUserExists(signUpData.Username, signUpData.Email, db)
	if value == 1 {
		sendJSONError(w, "Email already used", "EMAIL_USED", http.StatusConflict)
		return
	} else if value == 2 {
		sendJSONError(w, "Username already exists", "USERNAME_EXISTS", http.StatusConflict)
		return
	}

	// Save the user in the database, need to enter x2 password (UI design for the local app)
	newUser, err := mf.NewUser(signUpData.Username, signUpData.Email, signUpData.Password, signUpData.Password, db)
	if err != nil {
		sendJSONError(w, "Could not create user", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	errS := mdb.SaveDB(db)
	if errS != nil {
		sendJSONError(w, "Could not save the new user", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	// User Id is the number of current registered users (there is possibility yet to delete a user, so this should work for now)
	credToken, err := generateJWT(newUser.ID)
	if err != nil {
		sendJSONError(w, "Could not generate credToken", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"cred_token": credToken})
}
