package playermethods_test

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestHasTeam(t *testing.T) {
	p1 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p1",
		TeamIDs:   map[uuid.UUID]string{uuid.New(): "t1"},
	}

	p2 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	expectedBool1 := true
	expectedBool2 := false

	t.Run("Has player a team", func(t *testing.T) {
		bool1 := p1.HasTeam()
		bool2 := p2.HasTeam()
		if bool1 != expectedBool1 {
			t.Errorf("Expected bool %v, got %v", expectedBool1, bool1)
		}
		if bool2 != expectedBool2 {
			t.Errorf("Expected bool %v, got %v", expectedBool2, bool2)
		}
	})
}
