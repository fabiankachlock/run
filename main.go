package main

import (
	"os"

	"github.com/fabiankachlock/exec/internal/exec"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		exec.Help()
		return
	}

	if args[1] == "--init" {
		exec.Init()
	} else if args[1] == "--help" {
		exec.Help()
	} else {
		exec.Execute(args[0])
	}
}
