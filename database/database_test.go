package database

import (
	"github.com/thirdknife/scoutingapp/database/player"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenOrCreate(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Create a database since none exists at the given path.
	db, err := OpenOrCreate(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.PlayerTable.Insert(&player.Player{ID: "123", Name: "abc"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	// Verify that the database was written to file.
	_, err = os.Stat(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	// Reopen the database.
	db, err = OpenOrCreate(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	p, err := db.PlayerTable.Get("123")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "abc" {
		t.Fatalf("got %s, wanted 'abc'", p.Name)
	}
}

func TestOpenInMemory(t *testing.T) {
	db, err := OpenOrCreate("")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.PlayerTable.Insert(&player.Player{ID: "123", Name: "abc"}); err != nil {
		t.Fatal(err)
	}
	p, err := db.PlayerTable.Get("123")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "abc" {
		t.Fatalf("got %s, wanted 'abc'", p.Name)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
