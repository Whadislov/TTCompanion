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

func SetJWTSecretKey(jwtSecretString string) {
	jwtSecret = []byte(jwtSecretString)
	cookieStore = sessions.NewCookieStore(jwtSecret)
}
