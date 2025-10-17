package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsGetPlayer(t *testing.T) {
	pID := uuid.New()
	p2ID := uuid.New()

	d := Database{
		Players: map[uuid.UUID]*Player{pID: {
			ID:        pID,
			Firstname: "p",
		},
		},
	}
	expectedPlayer := Player{
		ID:        pID,
		Firstname: "p",
	}

	expectedError := "playerID 1 does not exist"

	t.Run("Get player from player ID", func(t *testing.T) {
		p, err := d.GetPlayer(pID)
		_, err2 := d.GetPlayer(p2ID)
		if p == nil {
			t.Errorf("Expected player %v, got %v", expectedPlayer, p)
		}
		if err != nil {
			t.Errorf("Expected err %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected err %v, got %v", expectedError, err2)
		}

	})

}
