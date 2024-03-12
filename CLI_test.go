package poker_test

import (
	"github.com/DavidDiasN/learn-with-tests-poker"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("Record chris win from user input", func(t *testing.T) {

		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("Record Cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Cleo")
	})
}
