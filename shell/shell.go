package shell

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type Shell struct {
	prompt string
	tty    *term.Terminal
	state  *term.State
	sigs   chan os.Signal
}

func New(prompt string) *Shell {
	s := new(Shell)

	s.prompt = strings.TrimSpace(prompt) + " "

	return s
}

func (s *Shell) Run() error {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	s.state = state

	s.tty = term.NewTerminal(os.Stdin, s.prompt)

	return nil
}

func (s *Shell) HandleSignals() {
	signal.Notify(s.sigs, syscall.SIGINT)
}
