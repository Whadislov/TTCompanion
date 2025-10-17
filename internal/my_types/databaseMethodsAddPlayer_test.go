package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsAddPlayer(t *testing.T) {
	d := Database{
		Clubs:   map[uuid.UUID]*Club{},
		Teams:   map[uuid.UUID]*Team{},
		Players: map[uuid.UUID]*Player{},
	}

	p1 := Player{
		ID:        uuid.New(),
		Firstname: "p1",
	}

	p2 := Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	expectedLen := 2

	t.Run("Add player to database", func(t *testing.T) {
		d.AddPlayer(&p1)
		d.AddPlayer(&p2)
		lenToVerify := len(d.Players)
		if lenToVerify != expectedLen {
			t.Errorf("Expected len of Players %v, got %v", expectedLen, lenToVerify)
		}
	})
}
