package exec

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/fabiankachlock/exec/internal/loader"
	"github.com/fabiankachlock/exec/internal/loader/npm"
)

type ScopeOptions struct {
	Alias  string `json:"alias"`
	Ignore bool   `json:"ignore"`
}

type Config struct {
	Scripts  map[string]string       `json:"scripts"`
	Extends  []string                `json:"extends"`
	Scopes   map[string]ScopeOptions `json:"scopes"`
	Location string                  `json:"-"`
}

type Script struct {
	Command string
	Wd      string
}

var Loaders []loader.Loader = []loader.Loader{
	npm.NewLoader(),
}

func FindConfig(cwd string) (Config, error) {
	dir := cwd
	for {
		filePath := path.Join(dir, CONFIG_FILE)
		if _, err := os.Stat(filePath); err == nil {
			return readConfig(filePath)
		} else if errors.Is(err, os.ErrNotExist) {
			dir = path.Join(dir, "..")
		}
	}
}

func readConfig(filePath string) (Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	config.Location = path.Dir(filePath)
	return config, nil
}

func (c Config) FindScript(name string) Script {
	return c.findScriptRecursive(name)
}

func (c Config) listScopedCommands() map[string]string {
	scripts := map[string]string{}

	for key, script := range c.Scripts {
		scripts[key] = script
	}

	for _, loader := range Loaders {
		prefix := loader.GetScope()
		for key, script := range loader.LoadConfig(c.Location) {
			scopedKey := prefix + ":" + key
			if _, found := scripts[scopedKey]; !found {
				scripts[scopedKey] = script
			}
		}
	}

	return scripts
}

func (c Config) findScriptRecursive(name string) Script {
	scripts := c.listScopedCommands()
	command, ok := scripts[name]
	if !ok {
		return Script{}
	}
	return Script{
		Command: command,
		Wd:      c.Location,
	}
}
