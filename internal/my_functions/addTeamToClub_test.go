package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestAddTeamToClub(t *testing.T) {
	c1 := mt.Club{
		ID:      uuid.New(),
		Name:    "c1",
		TeamIDs: map[uuid.UUID]string{},
	}

	t1 := mt.Team{
		ID:     uuid.New(),
		Name:   "t1",
		ClubID: map[uuid.UUID]string{},
	}

	t2 := mt.Team{
		ID:     uuid.New(),
		Name:   "t2",
		ClubID: map[uuid.UUID]string{c1.ID: "c1"},
	}

	t3 := mt.Team{
		ID:     uuid.New(),
		Name:   "t3",
		ClubID: map[uuid.UUID]string{},
	}

	c2 := mt.Club{
		ID:      uuid.New(),
		Name:    "c2",
		TeamIDs: map[uuid.UUID]string{t3.ID: t3.Name},
	}

	expectedLen1 := 1
	expectedError2 := "error when adding Team t2 to club c1: Team t2 is already in club c1"
	expectedError3 := "error when adding Team t3 to club c2: Team t3 is already in club c2"

	t.Run("Add Team to club", func(t *testing.T) {
		err := AddTeamToClub(&t1, &c1)
		err2 := AddTeamToClub(&t2, &c1)
		if err != nil || len(c1.TeamIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}

		err3 := AddTeamToClub(&t3, &c2)
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
	})
}
