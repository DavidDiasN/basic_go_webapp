package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const JsonContentType = "application/json"

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	Store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.Store = store
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	p.Handler = router

	return p
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	tmpl, err := template.ParseFiles("./poker/game.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("Problem loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", JsonContentType)
	json.NewEncoder(w).Encode(p.Store.GetLeague())
}

func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.Store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
