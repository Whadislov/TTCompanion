package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestTeamMethodsHasClub(t *testing.T) {
	t1 := Team{
		ID:     uuid.New(),
		Name:   "t1",
		ClubID: map[uuid.UUID]string{uuid.New(): "c1"},
	}

	t2 := Team{
		ID:     uuid.New(),
		Name:   "t2",
		ClubID: map[uuid.UUID]string{},
	}

	expectedBool1 := true
	expectedBool2 := false

	t.Run("Has team a club", func(t *testing.T) {
		bool1 := t1.HasClub()
		bool2 := t2.HasClub()
		if bool1 != expectedBool1 {
			t.Errorf("Expected bool %v, got %v", expectedBool1, bool1)
		}
		if bool2 != expectedBool2 {
			t.Errorf("Expected bool %v, got %v", expectedBool2, bool2)
		}
	})
}
