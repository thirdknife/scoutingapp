package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Load opens the database at the given path. If `path == ""` then a new in-memory database is returned.
func Load(path string) (*gorm.DB, error) {
	// Default to an in-memory database.
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	if path == "" {
		path = "file::memory:?cache=shared"
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto Migrate the schemas
	err = db.AutoMigrate(
		&Player{},
		&Analysis{},
		&DefenderAnalysis{},
		&MidfielderAnalysis{},
		&ForwardAnalysis{},
		&TacticalAnalysis{},
		&AthleticAnalysis{},
		&CharacterAnalysis{},
		&Match{},
		&Scout{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database to current schema: %w", err)
	}

	return db, nil
}

// SaveToFile saves the database to the same file path used when opening it. In-memory databases cannot be saved to
// files, but will not return an error.
func SaveToFile(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to save database to file: %w", err)
	}
	return sqlDB.Close()
}
