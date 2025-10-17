package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestPlayerMethodsRemoveClub(t *testing.T) {
	c1 := Club{
		ID:   uuid.New(),
		Name: "c1",
	}

	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
		ClubIDs:   map[uuid.UUID]string{c1.ID: c1.Name},
	}

	c2 := Club{
		ID:   uuid.New(),
		Name: "c2",
	}

	expectedLen1 := 0
	expectedError2 := "player p1 is not in club c2"

	t.Run("Remove club from player", func(t *testing.T) {
		err := p1.RemoveClub(&c1)
		err2 := p1.RemoveClub(&c2)
		if err != nil || len(p1.ClubIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
	})
}
