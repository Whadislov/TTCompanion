package my_types

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsDeleteClub(t *testing.T) {
	cID := uuid.New()
	c2ID := uuid.New()

	d := Database{
		Clubs: map[uuid.UUID]*Club{cID: {
			ID:   cID,
			Name: "c",
		},
		},
		Teams:   map[uuid.UUID]*Team{},
		Players: map[uuid.UUID]*Player{},
	}

	expectedLen := 0
	expectedError := fmt.Sprintf("clubID %v does not exist", c2ID)

	t.Run("Delete club from database", func(t *testing.T) {
		err := d.DeleteClub(cID)
		err2 := d.DeleteClub(c2ID)
		lenToVerify := len(d.Clubs)
		if lenToVerify != expectedLen {
			t.Errorf("Expected len of Clubs %v, got %v", expectedLen, lenToVerify)
		}
		if err != nil {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError, err2)
		}
	})
}
