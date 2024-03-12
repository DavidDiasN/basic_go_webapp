package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		winCalls: nil,
	}

	myServer := NewPlayerServer(&store)

	t.Run("returns Pepper's score", func(t *testing.T) {

		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		myServer.ServeHTTP(response, request)
		assertResponseBody(t, response.Body.String(), "20")
		assertResponseHeader(t, response.Code, http.StatusOK)

	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		myServer.ServeHTTP(response, request)
		assertResponseBody(t, response.Body.String(), "10")
		assertResponseHeader(t, response.Code, http.StatusOK)

	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		myServer.ServeHTTP(response, request)

		assertResponseHeader(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
		league:   nil,
	}

	myServer := NewPlayerServer(&store)

	t.Run("It records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		myServer.ServeHTTP(response, request)

		assertResponseHeader(t, response.Code, http.StatusAccepted)
    assertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertResponseHeader(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, JsonContentType)

	})
}
