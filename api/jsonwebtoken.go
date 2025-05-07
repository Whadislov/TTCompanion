package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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

/*
// Middleware to verify JSON Web token (JWT)
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the Credential token after "Bearer "
		credTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		credToken, err := jwt.Parse(credTokenString, func(credToken *jwt.Token) (interface{}, error) {
			// Checks the algorithm to prevent attacks
			if _, ok := credToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !credToken.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract User ID from the credToken
		claims, ok := credToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid credToken", http.StatusUnauthorized)
			return
		}

		// Check credToken expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				http.Error(w, "credToken expired", http.StatusUnauthorized)
				return
			}
		}

		// Now get user ID (string)
		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
}
*/

// verifyJWT verifies the JWT and updates the session if the token is valid
// returns userID or error
func verifyJWT(w http.ResponseWriter, r *http.Request, credTokenString string, session *sessions.Session) (string, error) {
	credToken, err := jwt.Parse(credTokenString, func(credToken *jwt.Token) (interface{}, error) {
		if _, ok := credToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !credToken.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := credToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", fmt.Errorf("token expired")
		}
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID in token")
	}

	// Update session
	session.Values["authenticated"] = true
	session.Values["user_id"] = userID
	session.Values["jwt"] = credTokenString
	if err := session.Save(r, w); err != nil {
		return "", fmt.Errorf("could not save session: %w", err)
	}

	return userID, nil
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
