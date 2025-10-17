package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestTeamMethodsHasPlayer(t *testing.T) {
	t1 := Team{
		ID:        uuid.New(),
		Name:      "t1",
		PlayerIDs: map[uuid.UUID]string{uuid.New(): "p1"},
	}

	t2 := Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	expectedBool1 := true
	expectedBool2 := false

	t.Run("Has team a player", func(t *testing.T) {
		bool1 := t1.HasPlayer()
		bool2 := t2.HasPlayer()
		if bool1 != expectedBool1 {
			t.Errorf("Expected bool %v, got %v", expectedBool1, bool1)
		}
		if bool2 != expectedBool2 {
			t.Errorf("Expected bool %v, got %v", expectedBool2, bool2)
		}
	})
}
