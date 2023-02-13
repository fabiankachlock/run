package yarn

import (
	"encoding/json"
	"os"
	"path"

	"github.com/fabiankachlock/run/internal/loader"
)

type yarnLoader struct{}

type packageJson struct {
	Scripts map[string]string `json:"scripts"`
}

func NewLoader() loader.Loader {
	return &yarnLoader{}
}

func (n *yarnLoader) GetScope() string {
	return "yarn"
}

func (n *yarnLoader) LoadConfig(dir string) map[string]string {
	file, err := os.ReadFile(path.Join(dir, "package.json"))
	if err != nil {
		return map[string]string{}
	}

	var parsedPackageJson packageJson
	err = json.Unmarshal(file, &parsedPackageJson)
	if err != nil {
		return map[string]string{}
	}

	remappedScripts := map[string]string{}
	for alias := range parsedPackageJson.Scripts {
		remappedScripts[alias] = "yarn run " + alias
	}
	return remappedScripts
}
