package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestNewTeam(t *testing.T) {
	d := mt.Database{
		Teams: map[uuid.UUID]*mt.Team{},
	}

	expectedLen := 1
	expectedError1 := "Team name cannot be empty"
	expectedError2 := "name must be valid (letters, figures and one space are allowed)"
	expectedError3 := "string is too long"

	t.Run("Delete Team", func(t *testing.T) {
		_, err1 := NewTeam("", &d)
		_, err2 := NewTeam("$$$zefz$&é", &d)
		_, err3 := NewTeam("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", &d)
		_, err4 := NewTeam("t2", &d)

		if err1 == nil {
			t.Errorf("Expected error %v, got %v", expectedError1, err1)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
		if err4 != nil || len(d.Teams) != expectedLen {
			t.Errorf("Expected error %v, got %v", nil, err4)
		}
	})
}
