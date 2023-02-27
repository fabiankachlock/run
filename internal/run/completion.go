package run

import (
	"os"
	"strings"
)

func GenerateCompletionScript() {
	os.Stdout.Write([]byte(`
_run_completions()
{
  OLD=$COMP_WORDBREAKS
  COMP_WORDBREAKS=" "
  COMPREPLY=($(compgen -W "$(run --generate-completion-list)" -- "${COMP_WORDS[1]}"))
  COMP_WORDBREAKS=$OLD
}

complete -F _run_completions run`))
}

func GenerateCompletionSuggestions() {
	cwd, err := os.Getwd()
	if err == nil {
		scripts := ListScriptNames(cwd)
		os.Stdout.Write([]byte(strings.Join(scripts, " ")))
	}
}
