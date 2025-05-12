package api

import (
	"fmt"
	"net/http"
)

// middleware that allows only get and post requests
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

// middleware that allows only authenticated sessions (valid session and cookie or valid header)
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth, userID, err := isAuthenticated(w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusUnauthorized)
			return
		}
		if isAuth {
			// adding userID for next headers
			r.Header.Set("User-ID", userID)
		}

		next(w, r)
	}
}
