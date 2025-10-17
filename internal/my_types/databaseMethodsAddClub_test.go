package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsAddClub(t *testing.T) {
	d := Database{
		Clubs:   map[uuid.UUID]*Club{},
		Teams:   map[uuid.UUID]*Team{},
		Players: map[uuid.UUID]*Player{},
	}

	c1 := Club{
		ID:   uuid.New(),
		Name: "c1",
	}

	c2 := Club{
		ID:   uuid.New(),
		Name: "c2",
	}

	expectedLen := 2

	t.Run("Add club to database", func(t *testing.T) {
		d.AddClub(&c1)
		d.AddClub(&c2)
		lenToVerify := len(d.Clubs)
		if lenToVerify != expectedLen {
			t.Errorf("Expected len of Clubs %v, got %v", expectedLen, lenToVerify)
		}
	})
}
