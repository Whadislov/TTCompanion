package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestAddPlayerToTeam(t *testing.T) {
	t1 := mt.Team{
		ID:        uuid.New(),
		Name:      "t1",
		PlayerIDs: map[uuid.UUID]string{},
	}

	p1 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p1",
		TeamIDs:   map[uuid.UUID]string{},
	}

	p2 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p2",
		TeamIDs:   map[uuid.UUID]string{t1.ID: t1.Name},
	}

	p3 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p3",
		TeamIDs:   map[uuid.UUID]string{},
	}

	t2 := mt.Team{
		ID:        uuid.New(),
		Name:      "y2",
		PlayerIDs: map[uuid.UUID]string{p3.ID: p3.Firstname},
	}

	expectedLen1 := 1
	expectedError2 := "error when adding Player p2 to Team t1: Player p2 is already in Team t1"
	expectedError3 := "error when adding Player p3 to Team y2: Player p3 is already in Team y2"

	t.Run("Add Player to Team", func(t *testing.T) {
		err := AddPlayerToTeam(&p1, &t1)
		err2 := AddPlayerToTeam(&p2, &t1)
		if err != nil || len(t1.PlayerIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}

		err3 := AddPlayerToTeam(&p3, &t2)
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
	})
}
