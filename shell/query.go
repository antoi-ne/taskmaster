package shell

import "strings"

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
func (q Query) Print(s string) {
	q.shell.tty.Write([]byte(s))
}

// Println writes the string to the terminal then adds a newline.
func (q *Query) Println(s string) {
	q.Print(s + "\n")
}

// Exit stops the interactive shell after the handler returns.
func (q *Query) Exit() {
	q.shell.exit = true
}
