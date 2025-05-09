package myfrontend

import (
	"encoding/json"
	"fmt"
	"net/http"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"

	"github.com/google/uuid"
)

// checks if there is persistence
// returns (bool, *mt.Database, int, error), if int == -1, it means there is no userID given
func CheckPersistence() (bool, *mt.Database, uuid.UUID, error) {
	var golangDB *mt.Database
	var isPersisOn bool
	var id uuid.UUID

	resp, err := http.Get(apiURL + "check-persistence")
	if err != nil {
		return false, nil, id, fmt.Errorf("Error fetching persistence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, id, fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	response := map[string]any{
		"authenticated": isPersisOn,
		"database":      golangDB,
		"userID":        id,
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, nil, id, fmt.Errorf("error decoding JSON: %w", err)
	}
	return true, golangDB, id, nil

}
