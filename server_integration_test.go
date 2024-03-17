package poker

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	if err != nil {
		log.Fatal(err)
	}

	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))

	t.Run("League sorted", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
      {"Name": "Cleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)

		got = store.GetLeague()
		AssertLeague(t, got, want)

	})

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		AssertResponseHeader(t, response, http.StatusOK)

		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest())
		AssertResponseHeader(t, response, http.StatusOK)

		got := GetLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		AssertLeague(t, got, want)
	})

}
