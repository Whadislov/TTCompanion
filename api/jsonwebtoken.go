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
	errSave := session.Save(r, w)
	if errSave != nil {
		return "", fmt.Errorf("could not save session: %w", errSave)
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
