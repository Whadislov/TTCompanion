package mydb

import (
	"fmt"
	"log"

	mt "github.com/Whadislov/TTCompanion/internal/my_types"
	"github.com/google/uuid"
)

// LoadUsers loads users from the database into the user map.
func (db *Database) LoadUser() (map[uuid.UUID]*mt.User, error) {
	rows, err := db.Conn.Query("SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1", userIDOfSession)
	if err != nil {
		return nil, fmt.Errorf("failed to load users: %w", err)
	}
	defer rows.Close()

	var users = make(map[uuid.UUID]*mt.User)
	for rows.Next() {
		var user mt.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users[user.ID] = &user
	}

	return users, rows.Err()
}

// LoadPlayers loads players from the database into the player map.
func (db *Database) LoadPlayers() (map[uuid.UUID]*mt.Player, error) {
	rows, err := db.Conn.Query("SELECT id, firstname, lastname, age, ranking, forehand, backhand, blade FROM players WHERE user_id = $1", userIDOfSession)
	if err != nil {
		return nil, fmt.Errorf("failed to load players: %w", err)
	}
	defer rows.Close()

	var players = make(map[uuid.UUID]*mt.Player)
	for rows.Next() {
		var player mt.Player
		player.Material = []string{"", "", ""}
		player.TeamIDs = make(map[uuid.UUID]string)
		player.ClubIDs = make(map[uuid.UUID]string)

		err := rows.Scan(&player.ID, &player.Firstname, &player.Lastname, &player.Age, &player.Ranking, &player.Material[0], &player.Material[1], &player.Material[2])
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}
		players[player.ID] = &player
	}

	return players, rows.Err()
}

// LoadTeams loads teams from the database into the team map.
func (db *Database) LoadTeams() (map[uuid.UUID]*mt.Team, error) {
	rows, err := db.Conn.Query("SELECT id, name FROM teams WHERE user_id = $1", userIDOfSession)
	if err != nil {
		return nil, fmt.Errorf("failed to load teams: %w", err)
	}
	defer rows.Close()

	var teams = make(map[uuid.UUID]*mt.Team)
	for rows.Next() {
		var team mt.Team
		team.PlayerIDs = make(map[uuid.UUID]string)
		team.ClubID = make(map[uuid.UUID]string)
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		teams[team.ID] = &team
	}
	return teams, rows.Err()
}

// LoadClubs loads clubs from the database into the club map.
func (db *Database) LoadClubs() (map[uuid.UUID]*mt.Club, error) {
	rows, err := db.Conn.Query("SELECT id, name FROM clubs WHERE user_id = $1", userIDOfSession)
	if err != nil {
		return nil, fmt.Errorf("failed to load clubs: %w", err)
	}
	defer rows.Close()

	var clubs = make(map[uuid.UUID]*mt.Club)
	for rows.Next() {
		var club mt.Club
		club.PlayerIDs = make(map[uuid.UUID]string)
		club.TeamIDs = make(map[uuid.UUID]string)
		err := rows.Scan(&club.ID, &club.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan club: %w", err)
		}
		clubs[club.ID] = &club
	}

	return clubs, rows.Err()
}

// LoadPlayerClubs loads the player-club relationships from the database.
func (db *Database) LoadPlayerClubs(players map[uuid.UUID]*mt.Player, clubs map[uuid.UUID]*mt.Club) error {
	rows, err := db.Conn.Query("SELECT player_id, club_id FROM player_club WHERE user_id = $1", userIDOfSession)
	if err != nil {
		return fmt.Errorf("failed to load player_club relationships: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playerID, clubID uuid.UUID
		err := rows.Scan(&playerID, &clubID)
		if err != nil {
			return fmt.Errorf("failed to scan player_club relationship: %w", err)
		}
		if player, ok := players[playerID]; ok {
			player.ClubIDs[clubID] = clubs[clubID].Name
		}
		if club, ok := clubs[clubID]; ok {
			club.PlayerIDs[playerID] = fmt.Sprintf("%v %v", players[playerID].Firstname, players[playerID].Lastname)
		}
	}
	return rows.Err()
}

// LoadPlayerTeams loads the player-team relationships from the database.
func (db *Database) LoadPlayerTeams(players map[uuid.UUID]*mt.Player, teams map[uuid.UUID]*mt.Team) error {
	rows, err := db.Conn.Query("SELECT player_id, team_id FROM player_team WHERE user_id = $1", userIDOfSession)
	if err != nil {
		return fmt.Errorf("failed to load player_team relationships: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playerID, teamID uuid.UUID
		err := rows.Scan(&playerID, &teamID)
		if err != nil {
			return fmt.Errorf("failed to scan player_team relationship: %w", err)
		}
		if player, ok := players[playerID]; ok {
			player.TeamIDs[teamID] = teams[teamID].Name
		}
		if team, ok := teams[teamID]; ok {
			team.PlayerIDs[playerID] = fmt.Sprintf("%v %v", players[playerID].Firstname, players[playerID].Lastname)
		}
	}

	return rows.Err()
}

// LoadTeamClubs loads the team-club relationships from the database.
func (db *Database) LoadTeamClubs(teams map[uuid.UUID]*mt.Team, clubs map[uuid.UUID]*mt.Club) error {
	rows, err := db.Conn.Query("SELECT team_id, club_id FROM team_club WHERE user_id = $1", userIDOfSession)
	if err != nil {
		return fmt.Errorf("failed to load team_club relationships: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var teamID, clubID uuid.UUID
		err := rows.Scan(&teamID, &clubID)
		if err != nil {
			return fmt.Errorf("failed to scan team_club relationship: %w", err)
		}
		if team, ok := teams[teamID]; ok {
			team.ClubID[clubID] = clubs[clubID].Name
		}
		if club, ok := clubs[clubID]; ok {
			club.TeamIDs[teamID] = teams[teamID].Name
		}
	}
	return rows.Err()
}

// LoadDB loads the database.
func LoadDB() (*mt.Database, error) {
	db, err := ConnectToDB()
	if err != nil {
		fmt.Println("Error loading postgresql database:", err)
		return nil, err
	}
	defer db.Close()

	log.Println("Loading user")
	user, err := db.LoadUser()
	if err != nil {
		return nil, err
	}

	log.Println("Loading players")
	players, err := db.LoadPlayers()
	if err != nil {
		return nil, err
	}
	log.Println("Loading teams")
	teams, err := db.LoadTeams()
	if err != nil {
		return nil, err
	}
	log.Println("Loading clubs")
	clubs, err := db.LoadClubs()
	if err != nil {
		return nil, err
	}
	log.Println("Loading player team relationships")
	err = db.LoadPlayerTeams(players, teams)
	if err != nil {
		return nil, err
	}
	log.Println("Loading player club relationships")
	err = db.LoadPlayerClubs(players, clubs)
	if err != nil {
		return nil, err
	}
	log.Println("Loading team club relationships")
	err = db.LoadTeamClubs(teams, clubs)
	if err != nil {
		return nil, err
	}
	golangDB := &mt.Database{
		Users:           user,
		Players:         players,
		Teams:           teams,
		Clubs:           clubs,
		DeletedElements: map[string][]uuid.UUID{},
	}
	log.Println("Database loaded successfully")
	return golangDB, nil
}

// LoadUsers loads users from the database into the user map.
func (db *Database) LoadAllUsers() (map[uuid.UUID]*mt.User, error) {
	log.Println("Loading all users")
	rows, err := db.Conn.Query("SELECT id, username, email, password_hash, created_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to load users: %w", err)
	}
	defer rows.Close()

	var users = make(map[uuid.UUID]*mt.User)
	for rows.Next() {
		var user mt.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users[user.ID] = &user
	}

	return users, rows.Err()
}

// LoadUsers loads users from the database into the user map.
func LoadUsersOnly() (*mt.Database, error) {
	db, err := ConnectToDB()
	if err != nil {
		fmt.Println("Error loading postgresql database:", err)
		return nil, err
	}
	defer db.Close()

	users, err := db.LoadAllUsers()
	if err != nil {
		return nil, err
	}

	golangDB := &mt.Database{
		Users: users,
	}
	return golangDB, nil
}
