package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestRemovePlayerFromClub(t *testing.T) {
	p1ID := uuid.New()
	p2ID := uuid.New()
	p3ID := uuid.New()
	p4ID := uuid.New()
	c1ID := uuid.New()
	c2ID := uuid.New()
	p1Name := "p1"
	p2Name := "p2"
	p3Name := "p3"
	p4Name := "p4"
	c1Name := "c1"
	c2Name := "c2"

	c1 := mt.Club{
		ID:        c1ID,
		Name:      c1Name,
		PlayerIDs: map[uuid.UUID]string{p1ID: p1Name, p3ID: p3Name},
	}

	p1 := mt.Player{
		ID:        p1ID,
		Firstname: p1Name,
		ClubIDs:   map[uuid.UUID]string{c1ID: c1Name},
	}

	p2 := mt.Player{
		ID:        p2ID,
		Firstname: p2Name,
	}

	p3 := mt.Player{
		ID:        p3ID,
		Firstname: p3Name,
	}

	c2 := mt.Club{
		ID:   c2ID,
		Name: c2Name,
	}

	p4 := mt.Player{
		ID:        p4ID,
		Firstname: p4Name,
	}

	expectedLen1 := 1
	expectedError2 := "p2 has not been successfully removed from c1. Reason : p2 is not in club c1"
	expectedError3 := "p3 has not been successfully removed from c1. Reason : p3 is not in club c1"
	expectedError4 := "p4 has not been successfully removed from c2. Reason : p4 is not in club c2"

	t.Run("Remove Player from club", func(t *testing.T) {
		err := RemovePlayerFromClub(&p1, &c1)
		err2 := RemovePlayerFromClub(&p2, &c1)
		if err != nil || len(c1.PlayerIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}

		err3 := RemovePlayerFromClub(&p3, &c1)
		if err3 == nil {
			t.Errorf("Expected error %v, got %v", expectedError3, err3)
		}
		err4 := RemovePlayerFromClub(&p4, &c2)
		if err4 == nil {
			t.Errorf("Expected error %v, got %v", expectedError4, err4)
		}
	})
}
