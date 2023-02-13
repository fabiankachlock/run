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
	} else {
		run.Execute(cleanArgs[1])
	}
}
