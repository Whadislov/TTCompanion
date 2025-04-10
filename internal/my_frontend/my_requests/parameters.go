package myfrontend

import (
	"os"
	"strings"
)

var apiURL string = "https://ttcompanion.onrender.com/api/"

func init() {
	if url := os.Getenv("API_URL"); url != "" {
		apiURL = url
	}
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
}
