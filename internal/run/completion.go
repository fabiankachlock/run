package run

import (
	"os"
	"strings"
)

func GenerateCompletionScript() {
	os.Stdout.Write([]byte(`
_run_completions()
{
  COMPREPLY=($(compgen -W "$(run --generate-completion-list)" -- "${COMP_WORDS[1]}"))
}

complete -F _run_completions run`))
}

func GenerateCompletionSuggestions() {
	cwd, err := os.Getwd()
	if err == nil {
		scripts := ListScripts(cwd)
		os.Stdout.Write([]byte(strings.Join(scripts, " ")))
	}
}
