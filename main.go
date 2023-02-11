package main

import (
	"log"
	"os"

	"github.com/fabiankachlock/exec/internal/exec"
)

type voidWriter struct{}

func (v *voidWriter) Write(bytes []byte) (n int, err error) {
	return len(bytes), nil
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		exec.Help()
		return
	}

	log.SetFlags(log.Ltime)
	if !exec.HasDebugFlag(os.Args) {
		log.SetOutput(&voidWriter{})
	}

	if args[1] == "--init" {
		exec.Init()
	} else if args[1] == "--help" {
		exec.Help()
	} else {
		exec.Execute(args[1])
	}
}
