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

type Script struct {
	Command string
	Wd      string
}

var Loaders []loader.Loader = []loader.Loader{
	npm.NewLoader(),
	yarn.NewLoader(),
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

	config.Location = filepath.Dir(filePath)
	return config, nil
}

func FindScript(cwd string, targetScript string) (*Script, error) {
	var loadedConfigs *map[string]bool = &map[string]bool{}

	dir := cwd
	for {
		var script *Script
		filePath := filepath.Join(dir, CONFIG_FILE)
		log.Printf("[info] try reading config %s", filePath)

		// check if config file exists in folder
		_, err := os.Stat(filePath)
		// if exists search that config
		if err == nil {
			script, err = findScriptInConfig(filePath, targetScript, loadedConfigs)
		} else {
			log.Printf("[error] while reading config %s: %s", filePath, err)
			err = nil
		}

		// if found a scrip
		if err == nil && script != nil {
			return script, nil
		} else if err != nil {
			log.Printf("[error] while searching script in config %s: %s", filePath, err)
		}

		// try go up one folder
		newDir := filepath.Dir(dir)
		// dir doesn't change when root
		if newDir == dir {
			log.Printf("[info] reached root - stopping search")
			return nil, ErrCantFindScript
		} else {
			dir = newDir
		}
	}
}

func findScriptInConfig(filePath string, targetScript string, alreadyLoaded *map[string]bool) (*Script, error) {
	log.Printf("[info] reading config %s", filePath)
	dir := filepath.Dir(filePath)
	config, err := readConfig(filePath)
	if err != nil {
		return nil, err
	}
	(*alreadyLoaded)[filePath] = true

	config = computeConfigScopes(config, filePath)
	script, ok := config.Scripts[targetScript]
	if ok {
		return &Script{
			Command: script,
			Wd:      config.Location,
		}, nil
	} else {
		log.Printf("[info] [%s] config script don't include target script", filePath)
	}

	log.Printf("[info] [%s] loading vendor scripts", filePath)
	// search all vendors
	for _, vendor := range Loaders {
		log.Printf("[info] [%s] [%s] loading vendor script", filePath, vendor.GetScope())
		for alias, script := range vendor.LoadConfig(dir) {
			// targetScript should match {vendorScope}:{vendorScript} (scoped version of vendor script)
			if vendor.GetScope()+":"+alias == targetScript {
				return &Script{
					Command: script,
					Wd:      config.Location,
				}, nil
			}
		}
		log.Printf("[info] [%s] [%s] vendor script don't include target script", filePath, vendor.GetScope())
	}

	log.Printf("[info] [%s] loading reference scripts", filePath)
	// search all reference
	for _, ref := range config.Extends {
		// only load config if not already loaded (against cyclic refs)
		referencePath := filepath.Join(dir, ref)
		if _, ok := (*alreadyLoaded)[referencePath]; !ok {
			log.Printf("[info] [%s] [%s] loading reference at %s", filePath, ref, referencePath)
			foundScript, err := findScriptInConfig(referencePath, targetScript, alreadyLoaded)
			if foundScript != nil && err == nil {
				return foundScript, nil
			} else if err != nil {
				log.Printf("[error] [%s] [%s] while loading reference: %s", filePath, ref, err)
			}
		} else {
			log.Printf("[info] [%s] [%s] not loading - already loaded", filePath, ref)
		}
	}

	return nil, ErrCantFindScript
}

func computeConfigScopes(config Config, filePath string) Config {
	log.Printf("[info] [%s] computing scopes", filePath)

	if scope, ok := config.Scopes[SELF_SCOPE]; ok {
		var alias string
		switch v := scope.(type) {
		case bool:
			if v {
				alias = filepath.Base(config.Location)
			} else {
				return config
			}
		case string:
			alias = v
		default:
			return config
		}
		log.Printf("[info] [%s] rescoped self as '%s'", filePath, alias)
		for key, script := range config.Scripts {
			delete(config.Scripts, key)
			config.Scripts[alias+":"+key] = script
		}
		log.Println(config)
	}
	return config
}
