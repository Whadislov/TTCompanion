package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestClubMethodsRemovePlayer(t *testing.T) {
	p1ID := uuid.New()
	p1Name := "p1"

	c1 := Club{
		ID:        uuid.New(),
		Name:      "c1",
		PlayerIDs: map[uuid.UUID]string{p1ID: p1Name},
	}

	p1 := Player{
		ID:        p1ID,
		Firstname: p1Name,
	}

	p2 := Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	expectedLen1 := 0
	expectedError2 := "player p1 is not in club c2"

	t.Run("Remove team from player", func(t *testing.T) {
		err := c1.RemovePlayer(&p1)
		err2 := c1.RemovePlayer(&p2)
		if err != nil || len(c1.PlayerIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
	})
}
