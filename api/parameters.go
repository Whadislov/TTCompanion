package api

import (
	"github.com/gorilla/sessions"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	ServerPort    string `json:"server_port"`
}

var jwtSecret []byte

// sign cookies with the secret key
var cookieStore *sessions.CookieStore

// set the secret key for JWT generation and sign cookies with the key
func SetJWTSecretKey(jwtSecretString string) {
	jwtSecret = []byte(jwtSecretString)
	cookieStore = sessions.NewCookieStore(jwtSecret)
}
