package main

import (
	"fmt"
	"github.com/DavidDiasN/learn-with-tests-poker"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("Problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("Problem creating file system player store, %v ", err)
	}

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}
