package myfrontend

import (
	"os"
	"strings"
)

var apiURL string = "https://ttcompanion-prod-912172190800.europe-west9.run.app/api"

//var apiURL string = "https://ttcompanion.onrender.com/api/"

//var apiURL string = "http://localhost:8000/api/"

func init() {
	if url := os.Getenv("API_URL"); url != "" {
		apiURL = url
	}
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
}
