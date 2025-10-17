package my_functions

import (
	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
	"testing"
)

func TestDeleteTeamComplexCase(t *testing.T) {

	// Easy case

	t01 := mt.Team{
		ID:   uuid.New(),
		Name: "t1",
	}

	t02 := mt.Team{
		ID:   uuid.New(),
		Name: "t2",
	}

	d0 := mt.Database{
		Teams:           map[uuid.UUID]*mt.Team{t01.ID: &t01},
		DeletedElements: map[string][]uuid.UUID{},
	}

	expectedLen01 := 0
	expectedError02 := "error when deleting Team t2: TeamID 1 does not exist"

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

	expectedLenDTeams := 0
	expectedLenClubTeamIDs := 0
	expectedLenPlayerTeamIDs := 0

	t.Run("Delete Club", func(t *testing.T) {

		// Easy case

		err0 := DeleteTeam(&t01, &d0)
		err02 := DeleteTeam(&t02, &d0)
		if err0 != nil || len(d0.Teams) != expectedLen01 {
			t.Errorf("Expected error %v, got %v", nil, err0)
		}
		if err02 == nil {
			t.Errorf("Expected error %v, got %v", expectedError02, err02)
		}

		// Complex case

		err := DeleteTeam(&t1, &d)
		if err != nil {
			t.Errorf("Expected error %v, got %v", nil, err)
		}
		if len(d.Teams) != expectedLenDTeams {
			t.Errorf("Expected LenDTeams %v, got %v", expectedLenDTeams, len(d.Teams))
		}
		if len(c1.TeamIDs) != expectedLenClubTeamIDs {
			t.Errorf("Expected LenClubTeamIDs %v, got %v", expectedLenClubTeamIDs, len(c1.TeamIDs))
		}
		if len(p1.TeamIDs) != expectedLenPlayerTeamIDs {
			t.Errorf("Expected LenPlayerTeamIDs %v, got %v", expectedLenPlayerTeamIDs, len(p1.TeamIDs))
		}
	})
}
