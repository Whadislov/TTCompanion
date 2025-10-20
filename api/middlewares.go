package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware that allows only get and post requests
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		// Do not forget to add new Headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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

		// Add user ID in the header to use it in API
		r.Header.Set("User-ID", userID)

		// Call next header
		next(w, r)
	}
}

// Middleware that checks persistency (valid session and cookie) and adds a header "User-ID"
func persisMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromCookie(w, r)
		// no persistency
		if err == fmt.Errorf("http: named cookie not present") {
			return
		} else if err != nil {
			sendJSONError(w, fmt.Sprintf("Persistence check error: %v", err), "PERSISTENCE_CHECK_ERROR", http.StatusConflict)
			return
		}

		// adding userID on the header for the load database function. userID is empty if there is no valid session/cookie
		r.Header.Set("User-ID", userID)

		next(w, r)
	}
}
