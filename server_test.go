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

	server, err := NewPlayerServer(&store, dummyGame)
	if err != nil {
		t.Errorf("Problem creating server %v", err)
	}

	t.Run("returns Pepper's score", func(t *testing.T) {

		request := NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		AssertResponseBody(t, response.Body.String(), "20")
		AssertResponseHeader(t, response, http.StatusOK)

	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		AssertResponseBody(t, response.Body.String(), "10")
		AssertResponseHeader(t, response, http.StatusOK)

	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseHeader(t, response, http.StatusNotFound)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server, err := NewPlayerServer(&StubPlayerStore{}, dummyGame)
		if err != nil {
			t.Errorf("Problem creating server %v", err)
		}
		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseHeader(t, response, http.StatusOK)
	})

}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
		league:   nil,
	}

	server, err := NewPlayerServer(&store, dummyGame)
	if err != nil {
		t.Errorf("Problem creatin server %v", err)
	}

	t.Run("It records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseHeader(t, response, http.StatusAccepted)
		AssertPlayerWin(t, &store, player)
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
		server, err := NewPlayerServer(&store, dummyGame)
		if err != nil {
			t.Errorf("Problem creating server %v", err)
		}

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetLeagueFromResponse(t, response.Body)
		AssertResponseHeader(t, response, http.StatusOK)
		AssertLeague(t, got, wantedLeague)
		AssertContentType(t, response, JsonContentType)

	})
}
