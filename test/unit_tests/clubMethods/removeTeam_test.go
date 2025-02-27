package clubmethods_test

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestRemoveTeam(t *testing.T) {
	t1Id := uuid.New()
	t1Name := "t1"

	c1 := mt.Club{
		ID:      uuid.New(),
		Name:    "c1",
		TeamIDs: map[uuid.UUID]string{t1Id: t1Name},
	}

	t1 := mt.Team{
		ID:   t1Id,
		Name: t1Name,
	}

	t2 := mt.Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	expectedLen1 := 0
	expectedError2 := "team t1 is not in club c2"

	t.Run("Remove team from player", func(t *testing.T) {
		err := c1.RemoveTeam(&t1)
		err2 := c1.RemoveTeam(&t2)
		if err != nil || len(c1.TeamIDs) != expectedLen1 {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError2, err2)
		}
	})
}
