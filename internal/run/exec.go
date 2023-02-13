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
	fmt.Println("run <script name> - runutes a script defined in run.json or referenced config files")
	fmt.Println("run --init        - creates a new config file")
	fmt.Println("run --help        - prints command usage")
}

func Execute(name string) {
	cwd, err := os.Getwd()
	handleError(err, "cant get cwd")

	handleError(err, "cant read config")

	script, err := FindScript(cwd, name)
	if script == nil || err != nil {
		fmt.Println("$run: an error happened: cant find a declaration for a script named '" + name + "'")
		log.Printf("[error]: loading script %s: %s", name, err)
		return
	}

	args := []string{"-c", script.Command}
	args = append(args, GetCleanArgs(os.Args[2:])...)

	fmt.Printf("$run: runuting: \"%s\" with args: %v\n", script.Command, os.Args[2:])
	start := time.Now()

	cmd := exec.Command("sh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = script.Wd

	err = cmd.Run()
	handleError(err, "cant runute command")

	elapsed := time.Since(start)
	fmt.Printf("$run: done in %s\n", elapsed)
}
