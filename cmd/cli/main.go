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
  fmt.Println("INSTRUCTIONS")
	fmt.Println("Type {Name} wins to record a win")
  fmt.Println("THAT'S IT!") 
  game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
  cli := poker.NewCLI(os.Stdin, os.Stdout, game)
  cli.PlayPoker()
}
