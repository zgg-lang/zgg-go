// +build !windows

package repl

import (
	"fmt"
	"os"
)

func (ConsoleReplContext) write(msg string) {
	tc := os.Getenv("ZGG_TEXT_STYLE")
	if tc == "" {
		tc = "36"
	}
	fmt.Printf("\033[%sm%s\033[0m\n", tc, msg)
}
