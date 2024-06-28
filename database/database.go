package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thirdknife/scoutingapp/database/player"
	"os"
)

type Database struct {
	db *sql.DB

	PlayerTable *player.Table
}

// OpenOrCreate attempts to open the provided file, but if it doesn't exist it creates
// one at the provided location. If no location is provided it creates an in-memory database.
func OpenOrCreate(path string) (*Database, error) {
	if path == "" {
		return New("")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return New(path)
	}
	return fromFile(path)
}

// New creates a new database at the given path. It fails if there is already a file at that path.
func New(path string) (*Database, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// The file not existing is good. This method creates it.
	} else {
		return nil, fmt.Errorf("database file %q already exists", path)
	}

	db, err := fromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database file %q: %v", path, err)
	}
	if err := db.PlayerTable.CreateDBTable(); err != nil {
		return nil, fmt.Errorf("failed to create database table %q: %v", path, err)
	}
	return db, nil
}

// Opens a database file.
func fromFile(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:          db,
		PlayerTable: &player.Table{DB: db},
	}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}
