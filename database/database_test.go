package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"testing"
)

func createTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := Load("")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestAddPlayer(t *testing.T) {
	db := createTestDB(t)

	newPlayer := &Player{
		Name: "Name",
	}
	result := db.Create(newPlayer)
	if result.Error != nil {
		t.Fatalf("Failed to create Player: %v", result.Error)
	}
	if newPlayer.ID == uuid.Nil {
		t.Fatal("Expected user ID to be set after creation")
	}

	var retrievedPlayer *Player
	result = db.First(&retrievedPlayer, "id = ?", newPlayer.ID)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve player: %v", result.Error)
	}
	if retrievedPlayer == nil {
		t.Fatal("Failed to retrieve player: got nil")
	}
	if retrievedPlayer.Name != newPlayer.Name {
		t.Fatalf("Expected player to have name %q, got %q", newPlayer.Name, retrievedPlayer.Name)
	}
}

func TestUpdatePlayer(t *testing.T) {
	db := createTestDB(t)

	newPlayer := &Player{
		Name:  "Original Name",
		Score: 100,
	}
	db.Create(newPlayer)

	newPlayer.Name = "Updated Name"
	newPlayer.Score = 150
	result := db.Save(newPlayer)
	if result.Error != nil {
		t.Fatalf("Failed to update Player: %v", result.Error)
	}

	var updatedPlayer Player
	result = db.First(&updatedPlayer, newPlayer.ID)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve updated Player: %v", result.Error)
	}

	if updatedPlayer.Name != "Updated Name" {
		t.Errorf("Expected updated Player name to be 'Updated Name', got %q", updatedPlayer.Name)
	}
	if updatedPlayer.Score != 150 {
		t.Errorf("Expected updated Player score to be 150, got %d", updatedPlayer.Score)
	}
}

func TestAddSingleAnalysis(t *testing.T) {
	db := createTestDB(t)

	newPlayer := &Player{
		Name: "Name",
	}
	db.Create(newPlayer)

	newMatch := &Match{}
	db.Create(newMatch)

	newDefenderAnalysis := &DefenderAnalysis{
		BallControl:        1,
		HeadingDefensively: 2,
		DefendingGeneral:   3,
		Defending1v1:       4,
		Tackling:           5,
		LongPassing:        6,
		ShortPassing:       7,
		RightFoot:          8,
		LeftFoot:           9,
	}
	db.Create(newDefenderAnalysis)

	newAnalysis := &Analysis{
		PlayerID:           newPlayer.ID,
		MatchID:            newMatch.ID,
		DefenderAnalysisID: newDefenderAnalysis.ID,
		PlayTimeMinutes:    15,
		Date:               "2024-07-07",
		WeatherCondition:   "sunny and hot",
		Venue:              "Foo Stadium",
	}
	db.Create(newAnalysis)

	if db.Error != nil {
		t.Fatalf("Error during testdata setup: %v", db.Error)
	}

	gotAnalysis := &Analysis{}
	if result := db.First(gotAnalysis, "player_id = ?", newPlayer.ID); result.Error != nil {
		t.Fatalf("Failed to retrieve analysis from player: %v", result.Error)
	}
	if gotAnalysis.DefenderAnalysisID != newDefenderAnalysis.ID {
		t.Fatalf("Expected defender analysis to have the same ID as originally inserted, got: %v", gotAnalysis.DefenderAnalysisID)
	}

	gotDefenderAnalysis := &DefenderAnalysis{}
	if result := db.First(gotDefenderAnalysis, "id = ?", gotAnalysis.DefenderAnalysisID); result.Error != nil {
		t.Fatalf("Failed to retrieve defender analysis: %v", result.Error)
	}
	// Clear database specific fields in order to cleanly compare the values that matter.
	gotDefenderAnalysis.BaseModel = BaseModel{}
	newDefenderAnalysis.BaseModel = BaseModel{}
	if *gotDefenderAnalysis != *newDefenderAnalysis {
		t.Fatalf("unexpected difference in defender analysis result, got:\n%+v\nwanted:\n%+v\n", gotDefenderAnalysis, newDefenderAnalysis)
	}
}

func TestAddScout(t *testing.T) {
	db := createTestDB(t)

	newScout := &Scout{
		Username: "username",
		Email:    "username@example.com",
	}
	result := db.Create(newScout)
	if result.Error != nil {
		t.Fatalf("Failed to create Scout: %v", result.Error)
	}

	var retrievedScout Scout
	result = db.First(&retrievedScout, "username = ?", newScout.Username)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve Scout: %v", result.Error)
	}

	if retrievedScout.Email != newScout.Email {
		t.Errorf("Expected Scout email to be %q, got %q", newScout.Email, retrievedScout.Email)
	}
}
