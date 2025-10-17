package my_types

import (
	"github.com/google/uuid"
	"testing"
)

func TestDatabaseMethodsAddTeam(t *testing.T) {
	d := Database{
		Clubs:   map[uuid.UUID]*Club{},
		Teams:   map[uuid.UUID]*Team{},
		Players: map[uuid.UUID]*Player{},
	}

	t1 := Team{
		ID:   uuid.New(),
		Name: "t1",
	}

	t2 := Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	expectedLen := 2

	t.Run("Add team to database", func(t *testing.T) {
		d.AddTeam(&t1)
		d.AddTeam(&t2)
		lenToVerify := len(d.Teams)
		if lenToVerify != expectedLen {
			t.Errorf("Expected len of Teams %v, got %v", expectedLen, lenToVerify)
		}
	})
}
