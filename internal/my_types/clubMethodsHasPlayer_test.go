package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestClubMethodsHasPlayer(t *testing.T) {
	c1 := Club{
		ID:        uuid.New(),
		Name:      "c1",
		PlayerIDs: map[uuid.UUID]string{uuid.New(): "p1"},
	}

	c2 := Club{
		ID:   uuid.New(),
		Name: "c2",
	}

	expectedBool1 := true
	expectedBool2 := false

	t.Run("Has club a player", func(t *testing.T) {
		bool1 := c1.HasPlayer()
		bool2 := c2.HasPlayer()
		if bool1 != expectedBool1 {
			t.Errorf("Expected bool %v, got %v", expectedBool1, bool1)
		}
		if bool2 != expectedBool2 {
			t.Errorf("Expected bool %v, got %v", expectedBool2, bool2)
		}
	})
}
