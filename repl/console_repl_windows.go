package repl

import (
	"fmt"
)

func (ConsoleReplContext) write(msg string) {
	fmt.Println(msg)
}
