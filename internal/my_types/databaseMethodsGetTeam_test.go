package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsGetTeam(t *testing.T) {
	tID := uuid.New()
	t2ID := uuid.New()

	d := Database{
		Teams: map[uuid.UUID]*Team{tID: {
			ID:   tID,
			Name: "t",
		},
		},
	}
	expectedTeam := Team{
		ID:   tID,
		Name: "t",
	}

	expectedError := "teamID 1 does not exist"

	t.Run("Get team from team ID", func(t *testing.T) {
		team, err := d.GetTeam(tID)
		_, err2 := d.GetTeam(t2ID)
		if team == nil {
			t.Errorf("Expected team %v, got %v", expectedTeam, team)
		}
		if err != nil {
			t.Errorf("Expected err %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected err %v, got %v", expectedError, err2)
		}

	})

}
