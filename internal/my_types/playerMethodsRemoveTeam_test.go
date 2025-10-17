package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestPlayerMethodsRemoveTeam(t *testing.T) {
	t1 := Team{
		ID:   uuid.New(),
		Name: "t1",
	}

	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
		TeamIDs:   map[uuid.UUID]string{t1.ID: t1.Name},
	}

	t2 := Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	expectedLen1 := 0
	expectedError2 := "player p1 is not in team t2"

	t.Run("Remove team from player", func(t *testing.T) {
		err := p1.RemoveTeam(&t1)
		err2 := p1.RemoveTeam(&t2)
		if err != nil || len(p1.TeamIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
	})
}
