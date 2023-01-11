package exec

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

//go:embed template/*
var templates embed.FS

func Init() {
	cwd, err := os.Getwd()
	handleError(err, "cant get cwd")

	path := path.Join(cwd, CONFIG_FILE)
	file, err := os.Create(path)
	handleError(err, "cant create config file")

	content, err := templates.ReadFile("template/empty-config.json")
	handleError(err, "cant read default config")

	_, err = file.Write(content)
	handleError(err, "cant write to config file")
}

func Help() {
	fmt.Println("")
	fmt.Println("exec - usage:")
	fmt.Println("")
	fmt.Println("exec <script name> - executes a script defined in exec.json or referenced config files")
	fmt.Println("exec --init        - creates a new config file")
	fmt.Println("exec --help        - prints command usage")
}

func Execute(name string) {
	cwd, err := os.Getwd()
	handleError(err, "cant get cwd")

	config, err := FindConfig(cwd)
	handleError(err, "cant read config")

	script := config.FindScript(name)
	if script.Command == "" {
		fmt.Println("exec: an error happened: cant find a declaration for a script named '" + name + "'")
		return
	}

	script.Command = strings.ReplaceAll(script.Command, "'", "")
	script.Command = strings.ReplaceAll(script.Command, "\"", "")
	commandParts := strings.Split(script.Command, " ")
	if len(commandParts) == 0 {
		fmt.Println("exec: an error happened: invalid command '" + script.Command + "'")
	}

	args := []string{""}
	args = append(args, commandParts[1:]...)
	args = append(args, os.Args[2:]...)

	cmd := exec.Command(commandParts[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = script.Wd
	cmd.Args = args

	err = cmd.Run()
	handleError(err, "cant execute command")
}
