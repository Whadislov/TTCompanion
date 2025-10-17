package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestDeletePlayerComplexCase(t *testing.T) {

	// Easy case

	p01 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p1",
	}

	p02 := mt.Player{
		ID:        uuid.New(),
		Firstname: "p2",
	}

	d0 := mt.Database{
		Players:         map[uuid.UUID]*mt.Player{p01.ID: &p01},
		DeletedElements: map[string][]uuid.UUID{},
	}

	expectedLen01 := 0
	expectedError02 := "error when deleting player p2: playerID 1 does not exist"

	// Complex case

	p1ID := uuid.New()
	t1Id := uuid.New()
	c1Id := uuid.New()
	p1Name := "p1"
	t1Name := "t1"
	c1Name := "c1"

	p1 := mt.Player{
		ID:        p1ID,
		Firstname: p1Name,
		TeamIDs:   map[uuid.UUID]string{t1Id: t1Name},
		ClubIDs:   map[uuid.UUID]string{c1Id: c1Name},
	}

	t1 := mt.Team{
		ID:        t1Id,
		Name:      t1Name,
		PlayerIDs: map[uuid.UUID]string{p1ID: p1Name},
		ClubID:    map[uuid.UUID]string{c1Id: c1Name},
	}

	c1 := mt.Club{
		ID:        c1Id,
		Name:      c1Name,
		PlayerIDs: map[uuid.UUID]string{p1.ID: p1.Firstname},
		TeamIDs:   map[uuid.UUID]string{t1.ID: t1.Name},
	}

	d := mt.Database{
		Clubs:           map[uuid.UUID]*mt.Club{c1Id: &c1},
		Teams:           map[uuid.UUID]*mt.Team{t1Id: &t1},
		Players:         map[uuid.UUID]*mt.Player{p1ID: &p1},
		DeletedElements: map[string][]uuid.UUID{},
	}

	expectedLenDPlayers := 0
	expectedLenClubPlayerIDs := 0
	expectedLenTeamPlayerIDs := 0

	t.Run("Delete Club", func(t *testing.T) {

		// Easy case

		err0 := DeletePlayer(&p01, &d0)
		err02 := DeletePlayer(&p02, &d0)
		if err0 != nil || len(d0.Players) != expectedLen01 {
			t.Errorf("Expected error %v, got %v", nil, err0)
		}
		if err02 == nil {
			t.Errorf("Expected error %v, got %v", expectedError02, err02)
		}

		// Complex case

		err := DeletePlayer(&p1, &d)
		if err != nil {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if len(d.Players) != expectedLenDPlayers {
			t.Errorf("Expected LenDClubs %v, got %v", expectedLenDPlayers, len(d.Players))
		}
		if len(c1.PlayerIDs) != expectedLenClubPlayerIDs {
			t.Errorf("Expected LenPlayerClubIDs %v, got %v", expectedLenClubPlayerIDs, len(c1.PlayerIDs))
		}
		if len(t1.PlayerIDs) != expectedLenTeamPlayerIDs {
			t.Errorf("Expected LenPlayerClubIDs %v, got %v", expectedLenTeamPlayerIDs, len(t1.PlayerIDs))
		}
	})
}
