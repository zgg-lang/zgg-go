package main

import (
	"fmt"
	"os"
)

func runAddDep(args []string) {
	f, err := os.OpenFile("zggdeps.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, arg := range args {
		fmt.Fprintln(f, arg)
	}
	runDeps([]string{})
}
