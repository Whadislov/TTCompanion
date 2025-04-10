package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Whadislov/TTCompanion/api"
)

func main() {

	serverAddress, serverPort, errConfig := loadConfig("config_app.json")
	if errConfig != nil {
		log.Fatalf("Cannot read config file: %v", errConfig)
	}

	// Create multiplexer to manage all routes
	mux := http.NewServeMux()

	// API
	api.RegisterRoutes(mux)

	log.Printf("Starting app server on %v:%v", serverAddress, serverPort)

	errServer := http.ListenAndServe(serverAddress+":"+serverPort, mux)
	if errServer != nil {
		log.Fatalf("App server error: %v", errServer)
	}

	go func() {
		// Verify that the API is ready
		waitForAPI(serverPort, 10, 500*time.Millisecond)
	}()

}

func loadConfig(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	type Config struct {
		ServerAddress string `json:"server_address"`
		ServerPort    string `json:"server_port"`
	}

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return "", "", err
	}

	// Remove http:// if present
	config.ServerAddress = cleanAddress(config.ServerAddress)

	return config.ServerAddress, config.ServerPort, nil
}

func waitForAPI(apiPort string, retries int, delay time.Duration) {
	apiURL := "http://127.0.0.1:" + apiPort + "/api/healthz"
	for range retries {
		resp, err := http.Get(apiURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("API is ready!")
			return
		}
		log.Println("Waiting for API to be ready...")
		time.Sleep(delay)
	}
	log.Fatal("API did not start in time")
}

func cleanAddress(address string) string {
	if strings.HasPrefix(address, "http://") {
		return address[7:]
	}
	return address
}
