package database

import (
	"fmt"
	"gorm.io/gorm"
)

func AllPlayers(db *gorm.DB) ([]*Player, error) {
	var players []*Player
	result := db.Find(&players) // SELECT * FROM Players;
	if result.Error != nil {
		return nil, fmt.Errorf("retrieving all Players failed: %w", result.Error)
	}
	return players, nil
}
