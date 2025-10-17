package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestClubMethodsAddPlayer(t *testing.T) {

	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
	}

	c1 := Club{
		ID:        uuid.New(),
		Name:      "c1",
		PlayerIDs: map[uuid.UUID]string{p1.ID: p1.Firstname},
	}

	p2 := Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	expectedLen := 2
	expectedError := "player p1 is already in club c1"

	t.Run("Add player to club", func(t *testing.T) {
		err := c1.AddPlayer(&p1)
		err2 := c1.AddPlayer(&p2)
		if err == nil {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
		if err2 != nil || len(c1.PlayerIDs) != expectedLen {
			t.Errorf("Expected len of PlayerIDs %v, got %v", expectedLen, len(c1.PlayerIDs))
		}
	})
}
