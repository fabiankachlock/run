package exec

func GetCleanArgs(args []string) []string {
	newArgs := []string{}
	for _, arg := range args {
		if !isDebugFlag(arg) {
			newArgs = append(newArgs, arg)
		}
	}
	return newArgs
}

func isDebugFlag(arg string) bool {
	return arg == "-v" || arg == "--debug"
}

func HasDebugFlag(args []string) bool {
	for _, arg := range args {
		if isDebugFlag(arg) {
			return true
		}
	}
	return false
}
