package poker_test

import (
	"bytes"
	"fmt"
	"io"
	"poker"
	"strings"
	"testing"
	"time"
)

var (
	dummyBlindAlerter = &SpyBlindAlerter{}
	dummyPlayerStore  = &poker.StubPlayerStore{}
	dummyStdIn        = &bytes.Buffer{}
	dummyStdOut       = &bytes.Buffer{}
)

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{at, amount})
}

func TestCLI(t *testing.T) {
	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("Start game with 8 players and record 'Cleo' as the winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")

	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("It prints error when given the wrong winner string", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}
		in := userSends("9", "Lloyd is a killer")

		cli := poker.NewCLI(in, stdout, game)
		err := cli.PlayPoker()
		if err == nil {
			t.Error("Expected to run into an error")
		}
	})

}

func userSends(vars ...string) io.Reader {
	var resString string
	for i := range vars {
		resString += (vars[i] + "\n")
	}
	return strings.NewReader(resString)
}

func assertGameStartedWith(t *testing.T, game *GameSpy, want int) {
	t.Helper()

	got := game.StartedWith
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, want string) {
	t.Helper()

	got := game.FinishedWith
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertGameNotStarted(t *testing.T, game *GameSpy) {
	t.Helper()

	if game.StartCalled {
		t.Errorf("Game was called when it shouldn't have been.")
	}
}

func assertScheduledAlert(t *testing.T, got, want ScheduledAlert) {
	t.Helper()

	if got.Amount != want.Amount {
		t.Errorf("got amount %d, want %d", got.Amount, want.Amount)
	}

	if got.At != want.At {
		t.Errorf("got scheduled time of %v, want %v", got.At, want.At)
	}
}

func checkSchedulingCases(cases []ScheduledAlert, t *testing.T, gots []ScheduledAlert) {
	t.Helper()

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(gots) <= i {
				t.Fatalf("alert %d want not scheduled %v", i, gots)
			}

			got := gots[i]
			assertScheduledAlert(t, got, want)
		})
	}
}
