package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestPlayerMethodsHasClub(t *testing.T) {
	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
		ClubIDs:   map[uuid.UUID]string{uuid.New(): "c1"},
	}

	p2 := Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	expectedBool1 := true
	expectedBool2 := false

	t.Run("Has player a club", func(t *testing.T) {
		bool1 := p1.HasClub()
		bool2 := p2.HasClub()
		if bool1 != expectedBool1 {
			t.Errorf("Expected bool %v, got %v", expectedBool1, bool1)
		}
		if bool2 != expectedBool2 {
			t.Errorf("Expected bool %v, got %v", expectedBool2, bool2)
		}
	})
}
