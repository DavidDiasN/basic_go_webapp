package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"strings"
	//	"log"
	//	s "main/server"
	//	"net/http"
	"os"
)

/*
Make PostgresPlayerStore
I am thinking we will need to store a connection?
I will need to read up on how to do this since working with databases
in this way will be new to me.
How many connections do I want to be able to support?
I have a lot of questions to ask myself and there is a lot to implement here.
GO BACK TO THE VIDEO PGX Top to Bottom
There is more to know from this video, your current method of collecting row data is wrong.


*/

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func main() {

	//	myServer := &s.PlayerServer{Store: NewInMemoryPlayerStore()}
	//	log.Fatal(http.ListenAndServe(":5000", myServer))

	rdr := bufio.NewReader(os.Stdin)
	userInput, _ := rdr.ReadString('\n')
	userInput = strings.TrimSuffix(userInput, "\n")

	queryTemplate := fmt.Sprintf(`SELECT gamesWon FROM scores WHERE name = '%s';`, userInput)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("WEB_APP_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	//queryResult, err := conn.Exec(ctx, `INSERT INTO scores (name, gamesWon)
	//VALUES ('Jerry', 5);`)

	//fmt.Println(queryResult.String())

	rows, err := conn.Query(ctx, queryTemplate)
	numbers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (int32, error) {
		var n int32
		err := row.Scan(&n)
		return n, err
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(numbers)
}
