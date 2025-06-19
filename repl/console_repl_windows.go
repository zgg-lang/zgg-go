package repl

import (
	"fmt"
)

func (c *ConsoleReplContext) write(msg string) {
	fmt.Fprintln(c.readline.Stdout(), msg)
}
