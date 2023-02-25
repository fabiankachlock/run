package run

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed template/*
var templates embed.FS

func Init() {
	cwd, err := os.Getwd()
	handleError(err, "cant get cwd")

	path := filepath.Join(cwd, CONFIG_FILE)
	file, err := os.Create(path)
	handleError(err, "cant create config file")

	content, err := templates.ReadFile("template/empty-config.json")
	handleError(err, "cant read default config")

	_, err = file.Write(content)
	handleError(err, "cant write to config file")
}

func Help() {
	fmt.Println("")
	fmt.Println("run - usage:")
	fmt.Println("")
	fmt.Println("run <script> [-v|--debug] - executes a script defined in run.json or referenced config files")
	fmt.Println("              -v|--debug  - print debug information")
	fmt.Println("run --init                - creates a new config file")
	fmt.Println("run -h|--help             - prints command usage")
}

func List() {
	cwd, err := os.Getwd()
	handleError(err, "cant get cwd")
	fmt.Printf("$run: scripts available at %s\n", cwd)
	scripts := ListScriptsRaw(cwd)
	isDebug := HasDebugFlag(os.Args)
	for _, script := range scripts {
		if isDebug {
			fmt.Printf("  %s → %s (in: %s)\n", script.Key, script.Command, script.Wd)
		} else {
			fmt.Printf("  %s → %s\n", script.Key, script.Command)
		}
	}
	fmt.Println("$run: done")
}

func Execute(name string) {
	log.Printf("[info] searching script '%s'", name)
	cwd, err := os.Getwd()
	handleError(err, "cant get cwd")

	script, err := FindScript(cwd, name)
	if script == nil || err != nil {
		fmt.Println("$run: an error happened: cant find a declaration for a script named '" + name + "'")
		log.Printf("[error]: loading script %s: %s", name, err)
		return
	}

	args := []string{"-c", script.Command}
	args = append(args, GetCleanArgs(os.Args[2:])...)

	fmt.Printf("$run: executing: \"%s\" with args: %v\n", script.Command, GetCleanArgs(os.Args[1:])[1:])
	start := time.Now()

	cmd := exec.Command("sh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = script.Wd

	err = cmd.Run()
	handleError(err, "cant execute command")

	elapsed := time.Since(start)
	fmt.Printf("$run: done in %s\n", elapsed)
}
