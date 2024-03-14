package main

import (
	"fmt"
	"github.com/DavidDiasN/learn-with-tests-poker"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, closePS, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer closePS()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
  
  game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
	poker.NewCLI(os.Stdin, os.Stdout, game)
  cli.PlayPoker()
}
