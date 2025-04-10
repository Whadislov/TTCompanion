package myfrontend

import (
	"os"
	"strings"
)

var apiURL string = "http://localhost:3000/"

func init() {
	if url := os.Getenv("API_URL"); url != "" {
		apiURL = url
	}
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
}
