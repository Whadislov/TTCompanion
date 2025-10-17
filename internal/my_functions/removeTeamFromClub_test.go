package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestRemoveTeamFromClub(t *testing.T) {
	t1ID := uuid.New()
	t2ID := uuid.New()
	t3ID := uuid.New()
	t4ID := uuid.New()
	c1ID := uuid.New()
	c2ID := uuid.New()
	t1Name := "t1"
	t2Name := "t2"
	t3Name := "t3"
	t4Name := "t4"
	c1Name := "c1"
	c2Name := "c2"

	c1 := mt.Club{
		ID:      c1ID,
		Name:    "c1",
		TeamIDs: map[uuid.UUID]string{t1ID: t1Name, t3ID: t3Name},
	}

	t1 := mt.Team{
		ID:     t1ID,
		Name:   t1Name,
		ClubID: map[uuid.UUID]string{c1ID: c1Name},
	}

	t2 := mt.Team{
		ID:   t2ID,
		Name: t2Name,
	}

	t3 := mt.Team{
		ID:   t3ID,
		Name: t3Name,
	}

	c2 := mt.Club{
		ID:   c2ID,
		Name: c2Name,
	}

	t4 := mt.Team{
		ID:     t4ID,
		Name:   t4Name,
		ClubID: map[uuid.UUID]string{c2ID: c2Name},
	}

	expectedLen1 := 1
	expectedError2 := "t2 has not been successfully removed from c1. Reason : t2 is not in club c1"
	expectedError3 := "t3 has not been successfully removed from c1. Reason : t3 is not in club c1"
	expectedError4 := "t4 has not been successfully removed from c2. Reason : t4 is not in club c2"

	t.Run("Remove team from club", func(t *testing.T) {
		err := RemoveTeamFromClub(&t1, &c1)
		err2 := RemoveTeamFromClub(&t2, &c1)
		if err != nil || len(c1.TeamIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}

		err3 := RemoveTeamFromClub(&t3, &c1)
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
		err4 := RemoveTeamFromClub(&t4, &c2)
		if err4 == nil {
			t.Errorf("Expected error %v, got %v", expectedError4, err4)
		}
	})
}
