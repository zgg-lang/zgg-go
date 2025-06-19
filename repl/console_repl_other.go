//go:build !windows
// +build !windows

package repl

import (
	"fmt"
	"os"
)

func (c *ConsoleReplContext) write(msg string) {
	tc := os.Getenv("ZGG_TEXT_STYLE")
	if tc == "" {
		tc = "36"
	}
	fmt.Fprintf(c.readline.Stdout(), "\033[%sm%s\033[0m\n", tc, msg)
}
