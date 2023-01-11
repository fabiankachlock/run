package exec

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

func FindConfig(path string) (Config, error) {

}

func (c Config) Merge(with Config) Config {

}

func (c Config) FindScript(name string) Script {
	return c.findScriptRecursive(name)
}

func (c Config) listScopedCommands() map[string]string {

}

func (c Config) findScriptRecursive(name string) Script {

}
