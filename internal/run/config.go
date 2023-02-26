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
	Root     bool                   `json:"root"`
	Scripts  map[string]string      `json:"scripts"`
	Extends  []string               `json:"extends,omitempty"`
	Scopes   map[string]interface{} `json:"scopes,omitempty"`
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
	location := filepath.Dir(filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Config{Location: location}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{Location: location}, err
	}

	config.Location = filepath.Dir(filePath)
	return config, nil
}

func readGlobalConfig() (*GlobalConfig, error) {
	globalFile := getRunGlobalConfigFilePath()
	content, err := os.ReadFile(globalFile)
	if err != nil {
		return &GlobalConfig{}, err
	}

	var globalConfig GlobalConfig
	err = json.Unmarshal(content, &globalConfig)
	if err != nil {
		return &GlobalConfig{}, err
	}

	return &globalConfig, nil
}

func WalkConfigs(cwd string, handler func(script Script) bool) error {
	var loadedConfigs *map[string]bool = &map[string]bool{}

	globalConfig, err := readGlobalConfig()
	if err != nil {
		log.Printf("[error] cant read global config: %s", err)
	}

	dir := cwd
	for {
		filePath := filepath.Join(dir, CONFIG_FILE)
		log.Printf("[info] try reading config %s", filePath)

		var shouldStop bool
		shouldStop, err = walkConfigRecursive(filePath, globalConfig, loadedConfigs, handler)
		if shouldStop {
			log.Printf("[info] stopping config search")
			break
		}

		// log errors happened while reading the config
		if err != nil {
			log.Printf("[error] while reading config %s: %s", filePath, err)
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
	config, err := readConfigFile(filePath)

	(*alreadyLoaded)[filePath] = true
	allExtension := config.Extends

	// scan local config first
	if err != nil {
		log.Printf("[error] [%s] reading local config: %s", config.Location, err)
	} else {
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
		globalConfigDeceleration.Location = filepath.Dir(filePath)
		allExtension = append(allExtension, globalConfigDeceleration.Extends...)
		if globalConfigDeceleration.Root {
			config.Root = true
		}

		log.Printf("[info] [%s] scanning global config", globalConfigDeceleration.Location)

		globalConfigDeceleration = computeConfigScopes(globalConfigDeceleration)
		shouldStop := scanConfig(globalConfigDeceleration, handler)
		if shouldStop {
			return shouldStop, nil
		}
	}

	// search all reference
	log.Printf("[info] [%s] loading reference scripts (%d)", config.Location, len(allExtension))
	for _, ref := range allExtension {
		// only load config if not already loaded (against cyclic refs)
		if _, ok := (*alreadyLoaded)[ref]; !ok {
			log.Printf("[info] [%s] [%s] loading reference at %s", config.Location, ref, ref)
			shouldStop, err := walkConfigRecursive(ref, globalConfig, alreadyLoaded, handler)
			if shouldStop {
				return shouldStop, nil
			} else if err != nil {
				log.Printf("[error] [%s] [%s] while loading reference: %s", filePath, ref, err)
			}
		} else {
			log.Printf("[info] [%s] [%s] not loading - already loaded", filePath, ref)
		}
	}
	if config.Root {
		return stopConfigWalk, nil
	}
	return continueConfigWalk, nil
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
			if scope != "" {
				script.Key = scope + ":" + alias
			} else {
				script.Key = alias
			}
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
	if err != nil {
		return scriptToReturn, err
	}
	if scriptToReturn == nil {
		return scriptToReturn, ErrCantFindScript
	}
	return scriptToReturn, nil
}

func ListScriptNames(cwd string) []string {
	var allScripts []string = []string{}
	WalkConfigs(cwd, func(script Script) bool {
		allScripts = append(allScripts, script.Key)
		return continueConfigWalk
	})
	return allScripts
}

func ListScriptsRaw(cwd string) []Script {
	var allScripts []Script = []Script{}
	WalkConfigs(cwd, func(script Script) bool {
		allScripts = append(allScripts, script)
		return continueConfigWalk
	})
	return allScripts
}
