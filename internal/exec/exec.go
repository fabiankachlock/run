package exec

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
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

	handleError(err, "cant read config")

	script, err := FindScript(cwd, name)
	if script == nil || err != nil {
		fmt.Println("$exec: an error happened: cant find a declaration for a script named '" + name + "'")
		log.Printf("[error]: loading script %s: %s", name, err)
		return
	}

	args := []string{"-c", script.Command}
	args = append(args, os.Args[2:]...)

	fmt.Printf("$exec: executing: \"%s\" with args: %v\n", script.Command, os.Args[2:])
	start := time.Now()

	cmd := exec.Command("sh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = script.Wd

	err = cmd.Run()
	handleError(err, "cant execute command")

	elapsed := time.Since(start)
	fmt.Printf("$exec: done in %s\n", elapsed)
}
