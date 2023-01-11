package exec

import "fmt"

func handleError(err error, msg string) {
	if err != nil {
		fmt.Printf("exec: an error happened:\n-> %s\n-> %s\n", msg, err)
		panic(err)
	}
}
