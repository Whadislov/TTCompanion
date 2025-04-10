package myfrontend

import (
	"os"
	"strings"
)

var apiURL string = "http://127.0.0.1:8080/api/"

func init() {
	if url := os.Getenv("API_URL"); url != "" {
		apiURL = url
	}
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
}
