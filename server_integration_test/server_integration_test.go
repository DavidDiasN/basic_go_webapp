package server_integration_test

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "testing"
  s "main/server"
)


type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
  store := InMemoryPlayerStore{}
  server := s.PlayerServer{Store: &store}
  player := "Pepper"

  server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
  server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
  server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

  response := httptest.NewRecorder()
  server.ServeHTTP(response, newGetScoreRequest(player))
  assertResponseHeader(t, response.Code, http.StatusOK)

  assertResponseBody(t, response.Body.String(), "3")
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertResponseHeader(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response header is wrong, got %d want %d", got, want)
	}
}
