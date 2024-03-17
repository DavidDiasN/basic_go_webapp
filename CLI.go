package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	PlayerPrompt         = "Please enter the number of players: "
	BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
)

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() error {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))
	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return err
	}

	cli.game.Start(numberOfPlayers, cli.out)

	winnerInput := cli.readLine()
	winner, err := extractWinner(winnerInput)
	if err != nil {
		return err
	}

	cli.game.Finish(winner)
	return nil
}

func extractWinner(userInput string) (string, error) {
	if strings.Contains(userInput, " wins") {
		return strings.Replace(userInput, " wins", "", 1), nil
	} else {
		return "", fmt.Errorf("The proper input format is: %v\n you entered: %v", PlayerPrompt, userInput)
	}

}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
