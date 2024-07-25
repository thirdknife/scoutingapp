// This file contains LLM-generated test cases.
package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	err = db.AutoMigrate(
		&Player{},
		&PlayerAnalysis{},
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
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestPlayerCRUD(t *testing.T) {
	db := setupTestDB(t)

	// Create
	player := &Player{Name: "John Doe", Score: 85}
	if err := db.Create(player).Error; err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Read
	var readPlayer Player
	if err := db.First(&readPlayer, "id = ?", player.ID).Error; err != nil {
		t.Fatalf("Failed to read player: %v", err)
	}
	if readPlayer.Name != player.Name || readPlayer.Score != player.Score {
		t.Errorf("Read player does not match created player")
	}

	// Update
	player.Name = "Jane Doe"
	player.Score = 90
	if err := db.Save(player).Error; err != nil {
		t.Fatalf("Failed to update player: %v", err)
	}

	// Verify Update
	var updatedPlayer Player
	if err := db.First(&updatedPlayer, "id = ?", player.ID).Error; err != nil {
		t.Fatalf("Failed to read updated player: %v", err)
	}
	if updatedPlayer.Name != "Jane Doe" || updatedPlayer.Score != 90 {
		t.Errorf("Updated player does not have expected values")
	}

	// Delete
	if err := db.Delete(player).Error; err != nil {
		t.Fatalf("Failed to delete player: %v", err)
	}

	// Verify Delete
	var deletedPlayer Player
	err := db.First(&deletedPlayer, "id = ?", player.ID).Error
	if err == nil {
		t.Errorf("Expected player to be deleted, but it still exists")
	}
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestScoutOperations(t *testing.T) {
	db := setupTestDB(t)

	// Create a scout
	scout := &Scout{Username: "topscout", Email: "scout@example.com"}
	if err := db.Create(scout).Error; err != nil {
		t.Fatalf("Failed to create scout: %v", err)
	}

	// Read scout
	var readScout Scout
	if err := db.First(&readScout, "id = ?", scout.ID).Error; err != nil {
		t.Fatalf("Failed to read scout: %v", err)
	}
	if readScout.Username != scout.Username || readScout.Email != scout.Email {
		t.Errorf("Read scout does not match created scout")
	}

	// Update scout
	scout.Email = "newscout@example.com"
	if err := db.Save(scout).Error; err != nil {
		t.Fatalf("Failed to update scout: %v", err)
	}

	// Verify update
	var updatedScout Scout
	if err := db.First(&updatedScout, "id = ?", scout.ID).Error; err != nil {
		t.Fatalf("Failed to read updated scout: %v", err)
	}
	if updatedScout.Email != "newscout@example.com" {
		t.Errorf("Updated scout does not have expected email")
	}
}

func TestPlayerAnalysisRelationship(t *testing.T) {
	db := setupTestDB(t)

	// Create a player
	player := &Player{Name: "Bob", Score: 88}
	if err := db.Create(player).Error; err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Create a player analysis
	playerAnalysis := &PlayerAnalysis{
		PlayerID:    player.ID,
		Notes:       "Promising talent",
		Name:        "Bob",
		Birthdate:   "2000-01-01",
		Height:      180,
		Weight:      75000,
		Club:        "FC Test",
		Position:    "Midfielder",
		ManagerName: "Coach Smith",
		Telephone:   "+1234567890",
	}
	if err := db.Create(playerAnalysis).Error; err != nil {
		t.Fatalf("Failed to create player analysis: %v", err)
	}

	// Retrieve the player with their analysis
	var retrievedPlayer Player
	err := db.Preload("PlayerAnalysis").First(&retrievedPlayer, "id = ?", player.ID).Error
	if err != nil {
		t.Fatalf("Failed to retrieve player with analysis: %v", err)
	}

	// Verify the relationship
	if retrievedPlayer.ID != playerAnalysis.PlayerID {
		t.Errorf("Player is not correctly associated with their analysis")
	}
}

func TestMultipleAnalysesForPlayer(t *testing.T) {
	db := setupTestDB(t)

	// Create a player
	player := &Player{Name: "Charlie", Score: 92}
	if err := db.Create(player).Error; err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Create multiple analyses for the player
	for i := 0; i < 3; i++ {
		analysis := &Analysis{
			PlayerID:         player.ID,
			PlayTimeMinutes:  90,
			Date:             time.Now().AddDate(0, 0, -i).Format("2006-01-02 15:04"),
			WeatherCondition: "Cloudy",
			Venue:            "Away Stadium",
		}
		if err := db.Create(analysis).Error; err != nil {
			t.Fatalf("Failed to create analysis %d: %v", i+1, err)
		}
	}

	// Retrieve all analyses for the player
	var analyses []Analysis
	err := db.Where("player_id = ?", player.ID).Find(&analyses).Error
	if err != nil {
		t.Fatalf("Failed to retrieve analyses for player: %v", err)
	}

	// Verify multiple analyses
	if len(analyses) != 3 {
		t.Errorf("Expected 3 analyses for player, got %d", len(analyses))
	}
}

func TestSaveToFile(t *testing.T) {
	tests := []struct {
		name    string
		dbPath  string
		dbSetup func(dbPath string) (*gorm.DB, error)
		wantErr bool
	}{
		{
			name:   "successful save",
			dbPath: filepath.Join(t.TempDir(), "test.db"),
			dbSetup: func(dbPath string) (*gorm.DB, error) {
				return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			},
			wantErr: false,
		},
		{
			name:   "in-memory database",
			dbPath: ":memory:",
			dbSetup: func(dbPath string) (*gorm.DB, error) {
				return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			},
			wantErr: false,
		},
		{
			name:   "error opening database",
			dbPath: filepath.Join(t.TempDir(), "nonexistent_dir", "test.db"),
			dbSetup: func(dbPath string) (*gorm.DB, error) {
				return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.dbSetup(tt.dbPath)
			if err != nil {
				t.Fatalf("failed to set up test database: %v", err)
			}

			err = SaveToFile(db)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// For successful saves, check if the file exists (except for in-memory database)
			if !tt.wantErr && tt.name != "in-memory database" {
				if _, err := os.Stat(tt.dbPath); os.IsNotExist(err) {
					t.Errorf("database file was not created at %s", tt.dbPath)
				}
			}
		})
	}
}
