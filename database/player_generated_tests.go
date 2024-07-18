// This file contains LLM-generated test cases.
package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestAllPlayers(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Auto-migrate the Player schema
	err = db.AutoMigrate(&Player{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Define test cases
	testCases := []struct {
		name          string
		seedData      []Player
		expectedCount int
		expectedError bool
	}{
		{
			name: "Successful retrieval",
			seedData: []Player{
				{Name: "Player 1", Score: 100},
				{Name: "Player 2", Score: 200},
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "Empty database",
			seedData:      []Player{},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the database before each test
			db.Exec("DELETE FROM players")

			// Seed the database with test data
			for _, player := range tc.seedData {
				if err := db.Create(&player).Error; err != nil {
					t.Fatalf("Failed to seed database: %v", err)
				}
			}

			// Call the function
			players, err := AllPlayers(db)

			// Check for errors
			if tc.expectedError && err == nil {
				t.Error("Expected an error, but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check the number of players retrieved
			if len(players) != tc.expectedCount {
				t.Errorf("Expected %d players, but got %d", tc.expectedCount, len(players))
			}

			// Verify player data
			for i, player := range players {
				if player.Name != tc.seedData[i].Name {
					t.Errorf("Expected player name %s, but got %s", tc.seedData[i].Name, player.Name)
				}
				if player.Score != tc.seedData[i].Score {
					t.Errorf("Expected player score %d, but got %d", tc.seedData[i].Score, player.Score)
				}
			}
		})
	}
}
