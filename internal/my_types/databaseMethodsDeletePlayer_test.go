package my_types

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsDeletePlayer(t *testing.T) {
	pID := uuid.New()
	p2ID := uuid.New()

	d := Database{
		Clubs: map[uuid.UUID]*Club{},
		Teams: map[uuid.UUID]*Team{},
		Players: map[uuid.UUID]*Player{pID: {
			ID:        pID,
			Firstname: "c",
		}},
	}

	expectedLen := 0
	expectedError := fmt.Sprintf("playerID %v does not exist", p2ID)

	t.Run("Delete player from database", func(t *testing.T) {
		err := d.DeletePlayer(pID)
		err2 := d.DeletePlayer(p2ID)
		lenToVerify := len(d.Players)
		if lenToVerify != expectedLen {
			t.Errorf("Expected len of Players %v, got %v", expectedLen, lenToVerify)
		}
		if err != nil {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if err2 == nil {
			t.Errorf("Expected error %v, got %v", expectedError, err2)
		}
	})
}
