package player

import (
	"database/sql"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kr/pretty"
	"strings"
	"testing"
)

func setupTable() (*Table, error) {
	db, err := sql.Open("sqlite3", "")
	if err != nil {
		return nil, err
	}
	table := &Table{db}
	if err := table.CreateDBTable(); err != nil {
		return nil, err
	}
	return table, nil
}

func TestGet(t *testing.T) {

	inPlayers := []*Player{
		{
			ID:   "TEST_ID",
			Name: "name",
		},
		{
			ID:   "TEST_ID",
			Name: "name",
		},
		{
			ID:   "TEST_ID",
			Name: "name",
		},
		{
			ID:   "TEST_ID",
			Name: "name",
		},
	}
	for i, input := range inPlayers {
		inp := input // local copy
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			table, err := setupTable()
			if err != nil {
				t.Fatal(err)
			}
			if err = table.Insert(inp); err != nil {
				t.Fatal(err)
			}

			outPlayer, err := table.Get("TEST_ID")
			if err != nil {
				t.Fatal(err)
			}
			// Use EquateEmpty to allow equality of empty vs nil slices.
			if !cmp.Equal(inp, outPlayer, cmpopts.EquateEmpty()) {
				t.Errorf("Unexpected difference in input and output player.")
				t.Errorf(strings.Join(pretty.Diff(inp, outPlayer), "\n"))
			}
		})
	}

}

func TestGetMap(t *testing.T) {
	table, err := setupTable()
	if err != nil {
		t.Fatal(err)
	}

	inPlayers := map[string]*Player{
		"ID_1": {
			ID:   "ID_1",
			Name: "name1",
		},
		"ID_2": {
			ID:   "ID_2",
			Name: "name2",
		},
		"ID_3": {
			ID:   "ID_3",
			Name: "name3",
		},
	}
	for _, tr := range inPlayers {
		err = table.Insert(tr)
		if err != nil {
			t.Fatal(err)
		}
	}

	outPlayers, err := table.GetMap("ID_1", "ID_2", "ID_3")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(inPlayers, outPlayers, cmpopts.EquateEmpty()) {
		t.Errorf("Unexpected difference in input and output players.")
		t.Errorf(strings.Join(pretty.Diff(inPlayers, outPlayers), "\n"))
	}
}

func TestGetAll(t *testing.T) {
	table, err := setupTable()
	if err != nil {
		t.Fatal(err)
	}

	inPlayers := map[string]*Player{
		"ID_1": {
			ID:   "ID_1",
			Name: "name1",
		},
		"ID_2": {
			ID:   "ID_2",
			Name: "name2",
		},
		"ID_3": {
			ID:   "ID_3",
			Name: "name3",
		},
	}
	for _, tr := range inPlayers {
		err = table.Insert(tr)
		if err != nil {
			t.Fatal(err)
		}
	}

	outPlayers, err := table.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(inPlayers, outPlayers, cmpopts.EquateEmpty()) {
		t.Errorf("Unexpected difference in input and output players.")
		t.Errorf(strings.Join(pretty.Diff(inPlayers, outPlayers), "\n"))
	}
}

func TestInsertCannotUpdate(t *testing.T) {
	table, err := setupTable()
	if err != nil {
		t.Fatal(err)
	}

	inPlayer := &Player{
		ID:   "TEST_ID",
		Name: "name",
	}
	err = table.Insert(inPlayer)
	if err != nil {
		t.Fatal(err)
	}
	err = table.Insert(inPlayer)
	if err == nil || err.Error() != "UNIQUE constraint failed: Players.id" {
		t.Errorf("Wanted error for insertion of player with the same ID, got: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	table, err := setupTable()
	if err != nil {
		t.Fatal(err)
	}
	originalEntry := &Player{
		ID:   "TEST_ID",
		Name: "name",
	}
	if err := table.Insert(originalEntry); err != nil {
		t.Fatal(err)
	}

	proposedUpdate := &Player{
		ID:   "TEST_ID",
		Name: "updated_name",
	}

	// The original input should be this replaced entry.
	replacedEntry, err := table.Update(proposedUpdate)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(originalEntry, replacedEntry) {
		t.Errorf("Unexpected difference in original input, and output from Update.")
		t.Errorf("%#v\n%#v", originalEntry, replacedEntry)
	}
	finalEntry, err := table.Get("TEST_ID")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(finalEntry, proposedUpdate) {
		t.Errorf("Unexpected diff between got / want:\n%#v\n%#v", proposedUpdate, finalEntry)
	}
}

func TestUpdateCannotInsert(t *testing.T) {
	table, err := setupTable()
	if err != nil {
		t.Fatal(err)
	}

	inPlayer := &Player{
		ID:   "TEST_ID",
		Name: "name",
	}
	if _, err = table.Update(inPlayer); err == nil || err.Error() != "cannot insert player using update" {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	table, err := setupTable()
	if err != nil {
		t.Fatal(err)
	}
	inPlayer := &Player{
		ID:   "TEST_ID",
		Name: "name",
	}
	err = table.Insert(inPlayer)
	if err != nil {
		t.Fatal(err)
	}
	_, err = table.Delete(inPlayer.ID)
	if err != nil {
		t.Fatal(err)
	}
	outPlayer, err := table.Get("TEST_ID")
	if err != nil {
		t.Fatal(err)
	}
	if outPlayer != nil {
		t.Errorf("Found player, but expected it to have been deleted.")
	}
}
