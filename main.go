package main

import (
	"log"
	"os"

	"github.com/fabiankachlock/run/internal/run"
)

type voidWriter struct{}

func (v *voidWriter) Write(bytes []byte) (n int, err error) {
	return len(bytes), nil
}

// main TODO: docs
func main() {
	args := os.Args
	if len(args) <= 1 {
		run.Help()
		return
	}

	log.SetFlags(log.Ltime)
	if !run.HasDebugFlag(os.Args) {
		log.SetOutput(&voidWriter{})
	}

	cleanArgs := run.GetCleanArgs(args)
	if cleanArgs[1] == "--init" {
		run.Init()
	} else if cleanArgs[1] == "--help" || cleanArgs[1] == "-h" {
		run.Help()
	} else if cleanArgs[1] == "--completion" {
		run.GenerateCompletionScript()
	} else if cleanArgs[1] == "--generate-completion-list" {
		run.GenerateCompletionSuggestions()
	} else {
		run.Execute(cleanArgs[1])
	}
}

// TODO
// - allow vendors to return errors for better debugging
