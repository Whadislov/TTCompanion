package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestRemovePlayerFromTeam(t *testing.T) {
	p1ID := uuid.New()
	p2ID := uuid.New()
	p3ID := uuid.New()
	p4ID := uuid.New()
	t1ID := uuid.New()
	t2ID := uuid.New()
	p1Name := "p1"
	p2Name := "p2"
	p3Name := "p3"
	p4Name := "p4"
	t1Name := "t1"
	t2Name := "t2"

	t1 := mt.Team{
		ID:        t1ID,
		Name:      t1Name,
		PlayerIDs: map[uuid.UUID]string{p1ID: p1Name, p3ID: p3Name},
	}

	p1 := mt.Player{
		ID:        p1ID,
		Firstname: p1Name,
		TeamIDs:   map[uuid.UUID]string{t1ID: t1Name},
	}

	p2 := mt.Player{
		ID:        p2ID,
		Firstname: p2Name,
	}

	p3 := mt.Player{
		ID:        p3ID,
		Firstname: p3Name,
	}

	t2 := mt.Team{
		ID:   t2ID,
		Name: t2Name,
	}

	p4 := mt.Player{
		ID:        p4ID,
		Firstname: p4Name,
	}

	expectedLen1 := 1
	expectedError2 := "p2 has not been successfully removed from t1. Reason : p2 is not in Team t1"
	expectedError3 := "p3 has not been successfully removed from t1. Reason : p3 is not in Team t1"
	expectedError4 := "p4 has not been successfully removed from t2. Reason : p4 is not in Team t2"

	t.Run("Remove Player from Team", func(t *testing.T) {
		err := RemovePlayerFromTeam(&p1, &t1)
		err2 := RemovePlayerFromTeam(&p2, &t1)
		if err != nil || len(t1.PlayerIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}

		err3 := RemovePlayerFromTeam(&p3, &t1)
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
		err4 := RemovePlayerFromTeam(&p4, &t2)
		if err4 == nil {
			t.Errorf("Expected error %v, got %v", expectedError4, err4)
		}
	})
}
