package clubmethods_test

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestAddTeam(t *testing.T) {
	t1ID := uuid.New()
	t1Name := "t1"

	c1 := mt.Club{
		ID:      uuid.New(),
		Name:    "c1",
		TeamIDs: map[uuid.UUID]string{t1ID: t1Name},
	}

	t1 := mt.Team{
		ID:     t1ID,
		Name:   t1Name,
		ClubID: map[uuid.UUID]string{},
	}

	t2 := mt.Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	expectedError := "team t1 is already in club c1"
	expectedLen := 2

	t.Run("Add team to club", func(t *testing.T) {
		err := c1.AddTeam(&t1)
		err2 := c1.AddTeam(&t2)
		if err == nil {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
		if err2 != nil || len(c1.TeamIDs) != expectedLen {
			t.Errorf("Expected len of ClubID %v, got %v", expectedLen, len(c1.TeamIDs))
		}
	})
}
