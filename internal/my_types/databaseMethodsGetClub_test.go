package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsGetClub(t *testing.T) {
	cID := uuid.New()
	c2ID := uuid.New()

	d := Database{
		Clubs: map[uuid.UUID]*Club{cID: {
			ID:   cID,
			Name: "c",
		},
		},
	}
	expectedClub := Club{
		ID:   cID,
		Name: "c",
	}

	expectedError := "clubID 1 does not exist"

	t.Run("Get club from club ID", func(t *testing.T) {
		c, err := d.GetClub(cID)
		_, err2 := d.GetClub(c2ID)
		if c == nil {
			t.Errorf("Expected club %v, got %v", expectedClub, c)
		}
		if err != nil {
			t.Errorf("Expected err %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected err %v, got %v", expectedError, err2)
		}

	})

}
