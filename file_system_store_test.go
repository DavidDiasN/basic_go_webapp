package poker

import (
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
      {"Name": "Chris", "Wins":33},
      {"Name": "Cleo", "Wins": 10}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

		got := store.GetLeague()
		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)

		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("Get Player Score", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
      {"Name": "Chris", "Wins":33},
      {"Name": "Cleo", "Wins": 10}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 33
		AssertScoreEquals(t, got, want)
	})

	t.Run("Store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
      {"Name": "Chris", "Wins":33},
      {"Name": "Cleo", "Wins": 10}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		AssertScoreEquals(t, got, want)

	})

	t.Run("Store wins for new players", func(t *testing.T) {

		database, cleanDatabase := CreateTempFile(t, `[
      {"Name": "Chris", "Wins":33},
      {"Name": "Cleo", "Wins": 10}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

		store.RecordWin("Pepper")

		got := store.GetPLayerScore("Pepper")
		want := 1
		AssertScoreEquals(t, got, want)
	})

	t.Run("Testing empty file case", func(t *testing.T) {

		database, cleanDatabase := CreateTempFile(t, ``)
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

	})
}
