package npm

import (
	"encoding/json"
	"os"
	"path"

	"github.com/fabiankachlock/exec/internal/loader"
)

type npmLoader struct{}

type packageJson struct {
	Scripts map[string]string `json:"scripts"`
}

func NewLoader() loader.Loader {
	return &npmLoader{}
}

func (n *npmLoader) GetScope() string {
	return "npm"
}

func (n *npmLoader) LoadConfig(dir string) map[string]string {
	file, err := os.ReadFile(path.Join(dir, "package.json"))
	if err != nil {
		return map[string]string{}
	}

	var parsedPackageJson packageJson
	err = json.Unmarshal(file, &parsedPackageJson)
	if err != nil {
		return map[string]string{}
	}
	return parsedPackageJson.Scripts
}