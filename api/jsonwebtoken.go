package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// generateJWT generates a JWT Credential token
func generateJWT(userID uuid.UUID) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	credToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	credTokenString, err := credToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return credTokenString, nil
}

// sign cookies with the secret key
var cookieStore = sessions.NewCookieStore(jwtSecret)

// Middleware to check the session
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Check if the session is ok
		session, err := cookieStore.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		//Check if the user is already authenticated
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			userID, ok := session.Values["user_id"].(string)
			if !ok {
				http.Error(w, "Invalid session data", http.StatusUnauthorized)
				return
			}
			r.Header.Set("User-ID", userID)
			next(w, r)
			return
		}

		//If session is not ok, check the JWT cookie (user was connected and recently closed his session)
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			//If no cookie, check header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			credTokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if credTokenString == "" {
				http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
				return
			}
			checkJWT(w, r, credTokenString, session, next)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//If the session is ok, check the JWT cookie
		checkJWT(w, r, cookie.Value, session, next)
	}
}

// Check the JSON Web token (JWT) and calls next header
func checkJWT(w http.ResponseWriter, r *http.Request, credTokenString string, session *sessions.Session, next http.HandlerFunc) {
	credToken, err := jwt.Parse(credTokenString, func(credToken *jwt.Token) (interface{}, error) {
		if _, ok := credToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !credToken.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract User ID from the credToken
	claims, ok := credToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Check credToken expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}
	}

	// Now get user ID (string)
	userID, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	// Update session
	session.Values["authenticated"] = true
	session.Values["user_id"] = userID
	session.Values["jwt"] = credTokenString
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// Add user ID in the header to use it in API
	r.Header.Set("User-ID", userID)

	// Call next header
	next(w, r)
}

// sendJSONError send an error following JSON format
func sendJSONError(w http.ResponseWriter, message string, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code":  code,
	})
}
