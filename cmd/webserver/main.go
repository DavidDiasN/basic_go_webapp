package main

import (
	"log"
	"net/http"
	"os"
	"github.com/DavidDiasN/learn-with-tests-poker"
)

const dbFileName = "game.db.json"

func main() {
  db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

  if err != nil {
    log.Fatalf("problem opening %s %v", dbFileName, err)
  }
  
	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatal("Problem creating file system player store, %v", err)
	}
  
  game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)

	server, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatal("Problem creating player server %v", err)
	}
  
  log.Fatal(http.ListenAndServe(":5009", server))
}
