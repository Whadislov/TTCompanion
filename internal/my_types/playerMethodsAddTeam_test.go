package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestPlayerMethodsAddTeam(t *testing.T) {
	t1 := Team{
		ID:   uuid.New(),
		Name: "t1",
	}

	t2 := Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
		TeamIDs:   map[uuid.UUID]string{t1.ID: t1.Name},
	}

	expectedLen := 2
	expectedError := "player p1 is already in team t2"

	t.Run("Add team to player", func(t *testing.T) {
		err := p1.AddTeam(&t1)
		err2 := p1.AddTeam(&t2)
		if err == nil {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
		if err2 != nil || len(p1.TeamIDs) != expectedLen {
			t.Errorf("Expected len of TeamIDs %v, got %v", expectedLen, len(p1.TeamIDs))
		}
	})
}
