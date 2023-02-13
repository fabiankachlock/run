package run

import (
	"errors"
	"fmt"
)

var (
	ErrCantFindScript = errors.New("cant find script")
)

func handleError(err error, msg string) {
	if err != nil {
		fmt.Printf("$run: an error happened:\n-> %s\n-> %s\n", msg, err)
		panic(err)
	}
}
