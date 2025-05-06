package mydb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Database represents the postgres database connection.
type Database struct {
	Conn *sql.DB
}

// ConnectToDB initializes a connection to the postgres database and creates it if the database does not exist.
func ConnectToDB() (*Database, error) {
	var conn *sql.DB
	var err error

	// Connect to postgres DB
	conn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	// Check if the database exists. pg_database = system view in PostgreSQL that contains info of all databases. datname = view that has all database names. $1 will be replaced by the first argument in the query
	query := `SELECT 1 FROM pg_database WHERE datname = $1`
	var exists int
	err = conn.QueryRow(query, dbName).Scan(&exists)
	if err == sql.ErrNoRows {
		// Database does not exist, create it
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		_, err = conn.Exec(createDBQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database %s created successfully.\n", dbName)
	} else if err != nil {
		return nil, fmt.Errorf("error checking database existence: %w", err)
	}
	// Verify the connection
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to verify connection to database: %w", err)
	}

	log.Printf("Connected to database successfully.")
	db := &Database{Conn: conn}

	// Create tables if they do not exist
	err = db.CreateTables()
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// CreateTables creates the necessary tables in the database for PostgreSQL.
func (db *Database) CreateTables() error {
	_, err := db.Conn.Exec(createAllTablesQuery)
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (db *Database) Close() error {
	err := db.Conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	log.Println("Database connection closed successfully.")
	return nil
}
