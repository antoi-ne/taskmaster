package shell

import (
	"os"
	"strings"

	"golang.org/x/term"
)

type Shell struct {
	prompt  string
	handler func(*Query) error
	tty     *term.Terminal
	sigs    chan os.Signal
	exit    bool
}

func New(prompt string, handler func(*Query) error) *Shell {
	s := new(Shell)

	s.prompt = strings.TrimSpace(prompt) + " "
	s.handler = handler

	s.sigs = make(chan os.Signal)

	return s
}

func (s *Shell) Run() error {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), state)

	s.tty = term.NewTerminal(os.Stdin, s.prompt)

	for !s.exit {
		if err := s.readLine(); err != nil {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (s *Shell) readLine() error {
	line, err := s.tty.ReadLine()
	if err != nil {
		return err
	}

	if len(strings.Fields(line)) == 0 {
		return nil
	}

	if s.handler == nil {
		return nil
	}

	return s.handler(&Query{
		shell: s,
		line:  line,
	})

}
