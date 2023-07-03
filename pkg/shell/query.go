package shell

import (
	"fmt"
	"strings"
)

// Query represents a single command line sent by the user.
type Query struct {
	line  string
	shell *Shell
}

// Line returns the line sent by the user
func (q *Query) Line() string {
	return q.line
}

// Argv returns an array of whitespace separated arguments from the user's command.
func (q *Query) Argv() []string {
	return strings.Fields(q.line)
}

// Print writes the string to the terminal
func (q Query) Print(a ...any) (int, error) {
	return fmt.Fprint(q.shell.tty, a...)
}

// Println writes the string to the terminal then adds a newline.
func (q *Query) Println(a ...any) (int, error) {
	return fmt.Fprintln(q.shell.tty, a...)
}

func (q *Query) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(q.shell.tty, format, a...)
}

// Exit stops the interactive shell after the handler returns.
func (q *Query) Exit() {
	q.shell.exit = true
}
