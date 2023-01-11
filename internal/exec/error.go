package exec

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
