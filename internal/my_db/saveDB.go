package mydb

import (
	"fmt"
	"log"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
)

// SaveUsers saves users in the database.
func (db *Database) SaveUsers(users map[uuid.UUID]*mt.User) error {
	for _, user := range users {
		if user.IsNew {
			// Let postgresql creates its own ID for a new user
			var postgresUserID uuid.UUID
			query := `
        INSERT INTO users (username, email, password_hash, created_at)
        VALUES ($1, $2, $3, $4)
		RETURNING id;
        `
			err := db.Conn.QueryRow(query, user.Name, user.Email, user.PasswordHash, user.CreatedAt).Scan(&postgresUserID)
			if err != nil {
				return fmt.Errorf("failed to save user: %w", err)
			}
			// Change the ID for the relationship tables
			user.ID = postgresUserID
			userIDOfSession = user.ID
		} else {
			query := `
        UPDATE users
		SET username = $1, email = $2, password_hash = $3, created_at = $4
		WHERE id = $5;
        `
			_, err := db.Conn.Exec(query, user.Name, user.Email, user.PasswordHash, user.CreatedAt, user.ID)
			if err != nil {
				return fmt.Errorf("failed to save user: %w", err)
			}
		}
	}
	return nil
}

// SavePlayers saves players in the database.
func (db *Database) SavePlayers(golangDB *mt.Database) error {
	for _, player := range golangDB.Players {
		if player.IsNew {
			// Let postgresql creates its own ID for a new player
			var postgresPlayerID uuid.UUID
			query := `
			INSERT INTO players (firstname, lastname, age, ranking, forehand, backhand, blade, user_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;
			`
			err := db.Conn.QueryRow(query, player.Firstname, player.Lastname, player.Age, player.Ranking, player.Material[0], player.Material[1], player.Material[2], userIDOfSession).Scan(&postgresPlayerID)
			if err != nil {
				return fmt.Errorf("failed to save the new player: %w", err)
			}
			// Change the ID for the relationship tables
			player.ID = postgresPlayerID
		} else {
			// Modify the player if it's not new
			query := `
			UPDATE players 
			SET firstname = $1, lastname = $2, age = $3, ranking = $4, forehand = $5, backhand = $6, blade = $7
			WHERE id = $8;
			`
			_, err := db.Conn.Exec(query, player.Firstname, player.Lastname, player.Age, player.Ranking, player.Material[0], player.Material[1], player.Material[2], player.ID)
			if err != nil {
				return fmt.Errorf("failed to save the edited player: %w", err)
			}
		}
	}

	for _, player := range golangDB.Players {
		for oldTeamID := range player.TeamIDs {
			if newTeamID, exists := idMapping[oldTeamID]; exists {
				// Replace old ID with the new one
				delete(player.TeamIDs, oldTeamID)
				player.TeamIDs[newTeamID] = golangDB.Teams[oldTeamID].Name
			}
		}
	}

	return nil
}

// SaveTeams saves teams in the database.
func (db *Database) SaveTeams(golangDB *mt.Database) error {
	for _, team := range golangDB.Teams {
		if team.IsNew {
			// Let postgresql creates its own ID for a new team
			var postgresTeamID uuid.UUID
			query := `
			INSERT INTO teams (name, user_id)
			VALUES ($1, $2)
			RETURNING id;
			`
			err := db.Conn.QueryRow(query, team.Name, userIDOfSession).Scan(&postgresTeamID)
			if err != nil {
				return fmt.Errorf("failed to save the new team: %w", err)
			}
			// Store the link between the old and the new ID
			idMapping[team.ID] = postgresTeamID
			// Change the ID for the relationship tables
			team.ID = postgresTeamID
		} else {
			// Modify the team if it's not new
			query := `
			UPDATE teams 
			SET name = $1
			WHERE id = $2;
			`
			_, err := db.Conn.Exec(query, team.Name, team.ID)
			if err != nil {
				return fmt.Errorf("failed to save the edited team: %w", err)
			}
		}
	}
	return nil
}

// SaveClubs saves clubs in the database.
func (db *Database) SaveClubs(golangDB *mt.Database) error {
	for _, club := range golangDB.Clubs {
		if club.IsNew {
			// Let postgresql creates its own ID for a new club
			var postgresClubID uuid.UUID
			query := `
			INSERT INTO clubs (name, user_id)
			VALUES ($1, $2)
			RETURNING id;
			`
			err := db.Conn.QueryRow(query, club.Name, userIDOfSession).Scan(&postgresClubID)
			if err != nil {
				return fmt.Errorf("failed to save the new club: %w", err)
			}
			// Store the link between the old and the new ID
			idMapping[club.ID] = postgresClubID
			// Change the ID for the relationship tables
			club.ID = postgresClubID
		} else {
			// Modify the club if it's not new
			query := `
			UPDATE clubs 
			SET name = $1
			WHERE id = $2;
			`
			_, err := db.Conn.Exec(query, club.Name, club.ID)
			if err != nil {
				return fmt.Errorf("failed to save the edited club: %w", err)
			}
		}
	}

	for _, player := range golangDB.Players {
		for oldClubID := range player.ClubIDs {
			if newClubID, exists := idMapping[oldClubID]; exists {
				// Replace old ID with the new one
				delete(player.ClubIDs, oldClubID)
				player.ClubIDs[newClubID] = golangDB.Clubs[oldClubID].Name
			}
		}
	}

	for _, team := range golangDB.Teams {
		for oldClubID := range team.ClubID {
			if newClubID, exists := idMapping[oldClubID]; exists {
				// Replace old ID with the new one
				delete(team.ClubID, oldClubID)
				team.ClubID[newClubID] = golangDB.Clubs[oldClubID].Name
			}
		}
	}

	return nil
}

// SavePlayerClubs saves the player-club relationships in the database.
func (db *Database) SavePlayerClubs(players map[uuid.UUID]*mt.Player) error {
	for _, player := range players {
		for clubID := range player.ClubIDs {
			query := `
				INSERT INTO player_club (player_id, club_id, user_id)
				VALUES ($1, $2, $3)
				ON CONFLICT (player_id, club_id) DO NOTHING;
				`
			_, err := db.Conn.Exec(query, player.ID, clubID, userIDOfSession)
			if err != nil {
				return fmt.Errorf("failed to save the new player_club relationship: %w", err)
			}
		}
	}
	return nil
}

// SavePlayerTeams saves the player-team relationships in the database.
func (db *Database) SavePlayerTeams(players map[uuid.UUID]*mt.Player) error {
	for _, player := range players {
		for teamID := range player.TeamIDs {
			query := `
			INSERT INTO player_team (player_id, team_id, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (player_id, team_id) DO NOTHING;
			`
			_, err := db.Conn.Exec(query, player.ID, teamID, userIDOfSession)
			if err != nil {
				return fmt.Errorf("failed to save the new player_team relationship: %w", err)
			}
		}
	}
	return nil
}

// SaveTeamClubs saves the team-club relationships in the database.
func (db *Database) SaveTeamClubs(teams map[uuid.UUID]*mt.Team) error {
	for _, team := range teams {
		for clubID := range team.ClubID {
			query := `
			INSERT INTO team_club (team_id, club_id, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (team_id, club_id) DO NOTHING;
			`
			_, err := db.Conn.Exec(query, team.ID, clubID, userIDOfSession)
			if err != nil {
				return fmt.Errorf("failed to save the new team_club relationship: %w", err)
			}
		}
	}
	return nil
}

// SaveDeletions saves the deletion that have been made by the user in the database.
func (db *Database) SaveDeletions(DElements map[string][]uuid.UUID) error {
	for table, ids := range DElements {
		if table != "users" && table != "players" && table != "teams" && table != "clubs" {
			return fmt.Errorf("invalid table name: %s", table)
		} else {
			for _, id := range ids {
				query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", table)
				_, err := db.Conn.Exec(query, id)
				if err != nil {
					return fmt.Errorf("failed to save the deletion: %w", err)
				}
			}
		}
	}
	return nil
}

func (db *Database) ResetTables() error {
	_, err := db.Conn.Exec(resetTablesQuery, userIDOfSession)
	if err != nil {
		return fmt.Errorf("failed to reset user data: %w", err)
	}
	return nil
}

// SaveDB saves the database (serialize).
func SaveDB(golangDB *mt.Database) error {
	var err error

	sqlDB, err = ConnectToDB()
	if err != nil {
		fmt.Println("Error while connecting to postgres database:", err)
	}
	defer sqlDB.Close()
	sqlDB.ResetTables()
	// Respect order, user then clubs, then team, then players
	log.Println("Saving user")
	err = sqlDB.SaveUsers(golangDB.Users)
	if err != nil {
		return err
	}
	log.Println("Saving clubs")
	err = sqlDB.SaveClubs(golangDB)
	if err != nil {
		return err
	}
	log.Println("Saving teams")
	err = sqlDB.SaveTeams(golangDB)
	if err != nil {
		return err
	}
	log.Println("Saving players")
	err = sqlDB.SavePlayers(golangDB)
	if err != nil {
		return err
	}
	log.Println("Saving player team relationships")
	err = sqlDB.SavePlayerTeams(golangDB.Players)
	if err != nil {
		return err
	}
	log.Println("Saving player club relationships")
	err = sqlDB.SavePlayerClubs(golangDB.Players)
	if err != nil {
		return err
	}
	log.Println("Saving team club relationships")
	err = sqlDB.SaveTeamClubs(golangDB.Teams)
	if err != nil {
		return err
	}
	log.Println("Saving deleted elements")
	err = sqlDB.SaveDeletions(golangDB.DeletedElements)
	if err != nil {
		return err
	}
	log.Println("Database saved successfully.")
	return nil
}
