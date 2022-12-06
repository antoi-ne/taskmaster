package shell

import (
	"os"
	"strings"

	"golang.org/x/term"
)

type Shell struct {
	prompt string
	tty    *term.Terminal
	state  *term.State
	sigs   chan os.Signal
	exit   bool
}

func New(prompt string) *Shell {
	s := new(Shell)

	s.prompt = strings.TrimSpace(prompt) + " "

	s.sigs = make(chan os.Signal)

	return s
}

func (s *Shell) Run() error {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	s.state = state

	go s.signalHandler()

	s.tty = term.NewTerminal(os.Stdin, s.prompt)

	for !s.exit {
		if err := s.readLine(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Shell) readLine() error {
	_, err := s.tty.ReadLine()
	if err != nil {
		return err
	}

	return nil
}

func (s *Shell) SignalChan() chan os.Signal {
	return s.sigs
}

func (s *Shell) signalHandler() {
	<-s.sigs
	s.exit = true
}

func (s *Shell) Restore() {
	if s.state != nil {
		term.Restore(int(os.Stdin.Fd()), s.state)
	}
}
