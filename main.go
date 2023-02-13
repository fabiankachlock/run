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

	if args[1] == "--init" {
		run.Init()
	} else if args[1] == "--help" {
		run.Help()
	} else {
		run.Execute(args[1])
	}
}
