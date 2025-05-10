package myfrontend

import (
	"encoding/json"
	"fmt"
	"net/http"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"

	"github.com/google/uuid"
)

// checks if there is persistence
// returns (bool, *mt.Database, int, error)
func CheckPersistence() (bool, *mt.Database, uuid.UUID, error) {
	resp, err := http.Get(apiURL + "check-persistence")
	if err != nil {
		return false, nil, uuid.UUID{}, fmt.Errorf("Error fetching persistence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, uuid.UUID{}, fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	type response struct {
		Authenticated bool         `json:"authenticated"`
		Database      *mt.Database `json:"database"`
		UserID        uuid.UUID    `json:"user_id"`
	}

	var res response

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return false, nil, uuid.UUID{}, fmt.Errorf("error decoding JSON: %w", err)
	}

	return res.Authenticated, res.Database, res.UserID, nil
}
