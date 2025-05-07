package myfrontend

import (
	"fmt"
	"net/http"
)

// Logout requests a logout to the API. The user's session and JWT cookie will be deleted
func Logout(credToken string) error {

	resp, err := http.Get(apiURL + "logout")
	if err != nil {
		return fmt.Errorf("could not logout. Reason: %w", err)
	}
	defer resp.Body.Close()
	return nil

}
