package api

import (
	"os"

	"github.com/gorilla/sessions"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	ServerPort    string `json:"server_port"`
}

var jwtSecret []byte

func SetJWTSecretKey(jwtSecretString string) {
	jwtSecret = []byte(os.Getenv(jwtSecretString))
}

// sign cookies with the secret key
var cookieStore = sessions.NewCookieStore(jwtSecret)
