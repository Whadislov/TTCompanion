package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestTeamMethodsRemovePlayer(t *testing.T) {
	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
	}

	t1 := Team{
		ID:        uuid.New(),
		Name:      "t1",
		PlayerIDs: map[uuid.UUID]string{p1.ID: p1.Firstname},
	}

	p2 := Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	expectedLen1 := 0
	expectedError2 := "player p1 is not in team t2"

	t.Run("Remove team from player", func(t *testing.T) {
		err := t1.RemovePlayer(&p1)
		err2 := t1.RemovePlayer(&p2)
		if err != nil || len(t1.PlayerIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
	})
}
