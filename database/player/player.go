// Package player encapsulates everything about a player in the database.
package player

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

// Table wraps a table of Player objects in a database.
type Table struct {
	DB *sql.DB
}

// CreateDBTable creates a table named Players containing Player structs.
// It will return an error if the database already has a Players table.
func (tt *Table) CreateDBTable() error {
	create := createPlayersTableSQL()
	if _, err := tt.DB.Exec(create); err != nil {
		return fmt.Errorf("failed to create Players table: %v", err)
	}
	return nil
}

// createPlayersTableSQL returns a SQL command to create a table named Players containing Player structs.
func createPlayersTableSQL() string {
	return `CREATE TABLE Players (
	id TEXT NOT NULL PRIMARY KEY,
	name TEXT
);`
}

func allPlayerTableFields() string {
	return `id, name`
}

type Player struct {
	// ID is a unique ID. It will never change after initial creation.
	ID string

	// Name is the full name.
	Name string
}

func (p *Player) String() string {
	return fmt.Sprintf("[%q: %q", p.ID, p.Name)
}

// Get returns the player with the provided ID.
// If no player is found, nil is returned and it is not considered an error.
func (tt *Table) Get(id string) (*Player, error) {
	players, err := tt.GetMap(id)
	if err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, nil
	}
	return players[id], nil
}

// GetMap returns a mapping of player ID to Player.
// Note that the size of the return will be smaller than the length of the inputs if duplicate input IDs are provided,
// or if any of the IDs are not found.
func (tt *Table) GetMap(ids ...string) (map[string]*Player, error) {
	// TODO(omar): This risks SQL injection if it takes an ID provided by a user. Bind parameters instead.
	query := fmt.Sprintf("SELECT %s FROM Players WHERE id IN (%s)", allPlayerTableFields(), cleanIDsCSV(quotedCSV(ids)))
	rows, err := tt.DB.Query(query)
	if err != nil {
		return nil, err
	}
	return getMapImpl(rows)
}

// cleanIDsCSV cleans a CSV input.
// This is intended for inserting user-provided inputs into a SQL query without risking injection attacks.
func cleanIDsCSV(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ',' || r == '"' || r == '_' {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (tt *Table) GetAll() (map[string]*Player, error) {
	query := fmt.Sprintf("SELECT %s FROM Players;", allPlayerTableFields())
	rows, err := tt.DB.Query(query)
	if err != nil {
		return nil, err
	}
	return getMapImpl(rows)
}

func getMapImpl(rows *sql.Rows) (map[string]*Player, error) {
	players, err := scanRows(rows)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Player)
	for _, p := range players {
		m[p.ID] = p
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return m, rows.Close()
}

func scanRows(rows *sql.Rows) ([]*Player, error) {
	var out []*Player
	for rows.Next() {
		p, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func scanRow(rows *sql.Rows) (*Player, error) {
	var p Player
	err := rows.Scan(&p.ID, &p.Name)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Insert inserts a player. It returns an error if there is a pre-existing player with the same ID.
func (tt *Table) Insert(p *Player) error {
	if p == nil {
		return fmt.Errorf("cannot insert nil player")
	}
	insert := `INSERT INTO Players VALUES (
		?, -- id STRING NOT NULL PRIMARY KEY,
		? -- name STRING
);`
	_, err := tt.DB.Exec(insert, p.ID, p.Name)
	return err
}

// ExternalUpdate updates a Player, but with limitations to the fields that can be modified.
// It is meant to be exposed to the frontend while preventing changes that might break features.
func (tt *Table) ExternalUpdate(p *Player) (*Player, error) {
	if p == nil {
		return nil, fmt.Errorf("cannot update nil player")
	}
	old, err := tt.Get(p.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving existing player with ID %q: %v", p.ID, err)
	}
	if old == nil {
		return nil, fmt.Errorf("cannot insert player using update")
	}

	// As of 2024-06-26, there are no fields that need to be protected.

	// ID difference is impossible since old was retrieved using the ID... but check to be sure.
	if old.ID != p.ID {
		return nil, fmt.Errorf("cannot modify protected fields: id")
	}
	return tt.Update(p)
}

// Update updates a player. It returns the old player, and fails if there was not an existing player.
func (tt *Table) Update(p *Player) (*Player, error) {
	if p == nil {
		return nil, fmt.Errorf("cannot update nil player")
	}
	old, err := tt.Get(p.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving existing player ID %q: %v", p.ID, err)
	}
	if old == nil {
		return nil, fmt.Errorf("cannot insert player using update")
	}

	// TODO(omar): This risks SQL injection if it takes an ID provided by a user. Bind parameters instead.
	p.Name = sanitize(p.Name)

	update := fmt.Sprintf(`UPDATE Players
	SET name = %q
	WHERE id = %q;
`, p.Name, p.ID)
	_, err = tt.DB.Exec(update)
	if err != nil {
		return nil, fmt.Errorf("error updating player with ID %q: %v", p.ID, err)
	}
	return old, nil
}

func sanitize(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsNumber(r) || unicode.IsLetter(r) || unicode.IsSpace(r) || r == ',' || r == '\'' || r == '_' {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// Delete deletes a Player. It returns the deleted Player, and is a no-op if there was not an existing Player.
func (tt *Table) Delete(id string) (*Player, error) {
	old, err := tt.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving existing Player ID %q: %v", id, err)
	}

	del := fmt.Sprintf(`DELETE FROM Players WHERE id = %q;`, id)
	_, err = tt.DB.Exec(del)
	if err != nil {
		return nil, fmt.Errorf("error deleting Player with ID %q: %v", id, err)
	}
	return old, nil
}

// quotedCSV returns the inputs quoted and escaped, and then concatenated with commas.
func quotedCSV(arr []string) string {
	var quoted []string
	for _, s := range arr {
		quoted = append(quoted, fmt.Sprintf("%q", s))
	}
	return strings.Join(quoted, ",")
}
