package npm

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fabiankachlock/run/internal/loader"
)

type npmLoader struct{}

type packageJson struct {
	Scripts map[string]string `json:"scripts"`
}

func NewLoader() loader.Loader {
	return &npmLoader{}
}

func (n *npmLoader) LoadConfig(dir string) map[string]string {
	file, err := os.ReadFile(filepath.Join(dir, "package.json"))
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
		remappedScripts[alias] = "npm run " + alias
	}
	return remappedScripts
}
