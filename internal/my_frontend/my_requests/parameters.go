package myfrontend

import (
	"os"
	"strings"
)

// Problem with fyne and the compilation of the backend : we can't grab the address from a config file, because it is unallowed. So manual switch is the solution for now

// Uncomment below for the main branch (production environment)
// var apiURL string = "https://ttcompanion-prod-912172190800.europe-west9.run.app/api"

// Uncomment below for the development branch (staging environment)
var apiURL string = "https://ttcompanion.onrender.com/api/"

// Uncomment below for the local testing
//var apiURL string = "http://localhost:8000/api/"

func init() {
	if url := os.Getenv("API_URL"); url != "" {
		apiURL = url
	}
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
}
