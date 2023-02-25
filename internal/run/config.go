package run

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/fabiankachlock/run/internal/loader"
	"github.com/fabiankachlock/run/internal/loader/npm"
	"github.com/fabiankachlock/run/internal/loader/yarn"
)

type Config struct {
	Scripts  map[string]string      `json:"scripts"`
	Extends  []string               `json:"extends"`
	Scopes   map[string]interface{} `json:"scopes"`
	Location string                 `json:"-"`
}

type GlobalConfig map[string]Config

type Script struct {
	Key     string
	Command string
	Wd      string
}

// bool indicates whether the searched script has been found
// when true the recursive search stops
type ScriptHandler func(script Script) bool

var AllLoaders map[string]loader.Loader = map[string]loader.Loader{
	"npm":  npm.NewLoader(),
	"yarn": yarn.NewLoader(),
}

func readConfigFile(filePath string) (Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	config.Location = filepath.Dir(filePath)
	return config, nil
}

func readGlobalConfig() (*GlobalConfig, error) {
	return new(GlobalConfig), nil
}

func WalkConfigs(cwd string, handler func(script Script) bool) error {
	var loadedConfigs *map[string]bool = &map[string]bool{}

	globalConfig, err := readGlobalConfig()
	if err != nil {
		log.Printf("[error] cant read global config: %s", err)
		return err
	}

	dir := cwd
	for {
		filePath := filepath.Join(dir, CONFIG_FILE)
		log.Printf("[info] try reading config %s", filePath)

		// check if config file exists in folder
		_, err := os.Stat(filePath)
		// if exists search that config
		if err == nil {
			var shouldStop bool
			shouldStop, err = walkConfigRecursive(filePath, globalConfig, loadedConfigs, handler)
			if shouldStop {
				log.Printf("[info] stopping config search")
				break
			}
		} else {
			log.Printf("[error] while reading config %s: %s", filePath, err)
			err = nil
		}

		// log errors happened while reading the config
		if err != nil {
			log.Printf("[error] while searching script in config %s: %s", filePath, err)
		}

		// try go up one folder
		newDir := filepath.Dir(dir)
		// dir doesn't change when root
		if newDir == dir {
			log.Printf("[info] reached root - stopping search")
			return ErrCantFindScript
		} else {
			dir = newDir
		}
	}
	return nil
}

func walkConfigRecursive(filePath string, globalConfig *GlobalConfig, alreadyLoaded *map[string]bool, handler func(script Script) bool) (bool, error) {
	log.Printf("[info] reading config %s", filePath)
	dir := filepath.Dir(filePath)
	config, err := readConfigFile(filePath)
	(*alreadyLoaded)[filePath] = true

	// scan local config first
	if err == nil {
		log.Printf("[info] [%s] scanning local config", config.Location)
		config = computeConfigScopes(config)
		shouldStop := scanConfig(config, handler)
		if shouldStop {
			return shouldStop, nil
		}
	}

	// try reading global config for this path
	globalConfigDeceleration, foundGlobalConfig := (*globalConfig)[filePath]
	if foundGlobalConfig {
		log.Printf("[info] [%s] scanning global config", config.Location)
		globalConfigDeceleration.Location = filepath.Dir(filePath)
		globalConfigDeceleration = computeConfigScopes(globalConfigDeceleration)
		shouldStop := scanConfig(globalConfigDeceleration, handler)
		if shouldStop {
			return shouldStop, nil
		}
	}

	// search all reference
	log.Printf("[info] [%s] loading reference scripts", config.Location)
	for _, ref := range config.Extends {
		// only load config if not already loaded (against cyclic refs)
		referencePath := filepath.Join(dir, ref)
		if _, ok := (*alreadyLoaded)[referencePath]; !ok {
			log.Printf("[info] [%s] [%s] loading reference at %s", config.Location, ref, referencePath)
			shouldStop, err := walkConfigRecursive(referencePath, globalConfig, alreadyLoaded, handler)
			if shouldStop {
				return shouldStop, nil
			} else if err != nil {
				log.Printf("[error] [%s] [%s] while loading reference: %s", filePath, ref, err)
			}
		} else {
			log.Printf("[info] [%s] [%s] not loading - already loaded", filePath, ref)
		}
	}
	return continueConfigWalk, ErrCantFindScript
}

func scanConfig(config Config, handler func(script Script) bool) bool {
	var script = &Script{Wd: config.Location}

	for key, command := range config.Scripts {
		script.Command = command
		script.Key = key
		shouldStop := handler(*script)
		if shouldStop {
			return shouldStop
		}
	}

	log.Printf("[info] [%s] loading vendor scripts", config.Location)
	// search all vendors
	for scope, vendor := range getEnabledLoaders(config) {
		log.Printf("[info] [%s] [%s] loading vendor script", config.Location, scope)
		scripts, err := vendor.LoadConfig(config.Location)
		if err != nil {
			log.Printf("[error] [%s] [%s] cant load vendor: %s", config.Location, scope, err)
			continue
		}
		for alias, command := range scripts {
			// targetScript should match {vendorScope}:{vendorScript} (scoped version of vendor script)
			script.Command = command
			script.Key = scope + ":" + alias
			shouldStop := handler(*script)
			if shouldStop {
				return shouldStop
			}
		}
	}
	return false
}

func computeConfigScopes(config Config) Config {
	log.Printf("[info] [%s] computing scopes", config.Location)

	if scope, ok := config.Scopes[SELF_SCOPE]; ok {
		var alias string
		switch v := scope.(type) {
		case bool:
			if v {
				alias = filepath.Base(config.Location) // true -> use default alias (dirname)
			} else {
				return config // false -> don't alias self
			}
		case string:
			alias = v // string -> use provided alias
		default:
			return config
		}
		log.Printf("[info] [%s] rescoped self as '%s'", config.Location, alias)
		for key, script := range config.Scripts {
			delete(config.Scripts, key)
			config.Scripts[alias+":"+key] = script
		}
		log.Println(config)
	}
	return config
}

func getEnabledLoaders(config Config) map[string]loader.Loader {
	loaders := map[string]loader.Loader{}
	for key, scope := range config.Scopes {
		loader, ok := AllLoaders[key]
		if !ok {
			continue // skip unknown loaders
		}
		switch v := scope.(type) {
		case bool:
			if v {
				loaders[key] = loader
			}
		case string:
			loaders[v] = loader
		}
	}
	return loaders
}

func FindScript(cwd string, targetScript string) (*Script, error) {
	var scriptToReturn *Script
	err := WalkConfigs(cwd, func(script Script) bool {
		if script.Key == targetScript {
			scriptToReturn = &script
			return stopConfigWalk
		}
		return continueConfigWalk
	})
	return scriptToReturn, err
}

func ListScripts(cwd string) []string {
	var allScripts []string = []string{}
	WalkConfigs(cwd, func(script Script) bool {
		allScripts = append(allScripts, script.Key)
		return continueConfigWalk
	})
	return allScripts
}
