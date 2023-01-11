package exec

import (
	"embed"
	"os"
	"path"
)

//go:embed template/*
var templates embed.FS

func Init(cwd string) {
	path := path.Join(cwd, CONFIG_FILE)
	file, err := os.Create(path)
	handleError(err)

	content, err := templates.ReadFile("template/empty-config.json")
	handleError(err)

	_, err = file.Write(content)
	handleError(err)
}
