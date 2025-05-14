package myfrontend

import (
	"encoding/json"
	"fmt"
	"log"
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
		var errorResponse struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return false, nil, uuid.UUID{}, fmt.Errorf("error decoding server response: %w", err)
		}

		if errorResponse.Code == "PERSISTENCE_CHECK_ERROR" {
			return false, nil, uuid.UUID{}, fmt.Errorf("server response for persistence check :%s", errorResponse.Error)
		}
	}

	type response struct {
		Authenticated bool         `json:"authenticated"`
		Database      *mt.Database `json:"database"`
		UserID        uuid.UUID    `json:"user_id"`
	}

	var res response

	errDecode := json.NewDecoder(resp.Body).Decode(&res)
	if errDecode != nil {
		return false, nil, uuid.UUID{}, fmt.Errorf("error decoding JSON: %w", err)
	}

	log.Println(fmt.Sprintf(`Finished checking for persistence. Results  :
	Authenticated : %v
	Database == nil : %v
	UserID == uuid.UUID{} : %v`, res.Authenticated, res.Database == nil, res.UserID == uuid.UUID{}))
	return res.Authenticated, res.Database, res.UserID, nil
}
