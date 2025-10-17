package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestAddPlayerToClub(t *testing.T) {
	c1 := mt.Club{
		ID:        uuid.New(),
		Name:      "c1",
		PlayerIDs: map[uuid.UUID]string{},
	}

	p1 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p1",
		ClubIDs:   map[uuid.UUID]string{},
	}

	p2 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p2",
		ClubIDs:   map[uuid.UUID]string{c1.ID: c1.Name},
	}

	p3 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p3",
		ClubIDs:   map[uuid.UUID]string{},
	}

	c2 := mt.Club{
		ID:        uuid.New(),
		Name:      "c2",
		PlayerIDs: map[uuid.UUID]string{p3.ID: p3.Firstname},
	}

	expectedLen1 := 1
	expectedError2 := "error when adding player p2 to club c1: player p2 is already in club c1"
	expectedError3 := "error when adding player p3 to club c2: player p3 is already in club c2"

	t.Run("Add Player to club", func(t *testing.T) {
		err := AddPlayerToClub(&p1, &c1)
		err2 := AddPlayerToClub(&p2, &c1)
		if err != nil || len(c1.PlayerIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}

		err3 := AddPlayerToClub(&p3, &c2)
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
	})
}
