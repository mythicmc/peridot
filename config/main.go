package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/mythicmc/peridot/repos"
)

type Config struct {
	Location string   `json:"location"`
	Repos    []string `json:"repos"`
	Software string   `json:"software"` // Supported: "vanilla", "paper", "velocity"
	Plugins  []string `json:"plugins"`
}

type Configs map[string]Config

func LoadConfigs(repositories repos.Repositories) (Configs, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configs := make(Configs)
	configFolder := filepath.Join(wd, "configs")
	configFiles, err := os.ReadDir(configFolder)
	if err != nil && os.IsNotExist(err) {
		log.Println("Warning: configs/ folder does not exist!")
		return configs, nil
	} else if err != nil {
		return nil, err
	}
	for _, configFile := range configFiles {
		if configFile.IsDir() || !strings.HasSuffix(configFile.Name(), ".js") {
			continue
		}
		filePath := filepath.Join(configFolder, configFile.Name())
		configJs, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		configName := configFile.Name()[:strings.LastIndex(configFile.Name(), ".")]
		config, err := ExecuteConfig(configFolder, string(configJs))
		if err != nil {
			return nil, err
		}
		err = ValidateConfig(config, repositories)
		if err != nil {
			return nil, err
		}
		configs[configName] = config
	}
	return configs, nil
}

var ErrInvalidLocation = errors.New("invalid location: must be an absolute path")

var ErrLocationNotExists = errors.New("invalid location: its parent directory does not exist")

var ErrInvalidSoftware = errors.New("invalid software: must be one of 'vanilla', 'paper', or 'velocity'")

var ErrInvalidRepos = errors.New("invalid repositories: at least one repository must be specified")

var ErrUnknownRepo = errors.New("unknown repository: must be one of the configured repositories")

var ErrInvalidPlugin = errors.New("invalid plugin: must be a non-empty string")

var ErrUnknownPluginSoftware = errors.New("unknown plugin or software: not found in the configured repositories")

func ValidateConfig(config Config, repositories repos.Repositories) error {
	if config.Location == "" || !filepath.IsAbs(config.Location) {
		return ErrInvalidLocation
	} else if stat, err := os.Stat(filepath.Dir(config.Location)); os.IsNotExist(err) {
		return ErrLocationNotExists
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return ErrLocationNotExists
	}
	if len(config.Repos) == 0 {
		return ErrInvalidRepos
	}
	repos := make([]repos.Repository, len(repositories))
	for _, repo := range config.Repos {
		if _, ok := repositories[repo]; !ok {
			return ErrUnknownRepo
		} else {
			repos = append(repos, repositories[repo])
		}
	}
	if config.Software != "vanilla" && config.Software != "paper" && config.Software != "velocity" {
		return ErrInvalidSoftware
	}
	for _, plugin := range config.Plugins {
		if plugin == "" {
			return ErrInvalidPlugin
		}
	}

	// Validate plugins and software against repositories
	for _, plugin := range config.Plugins {
		pluginFound := false
		for _, repo := range repos {
			if _, ok := repo.Plugins[plugin]; ok {
				pluginFound = true
				break
			}
		}
		if !pluginFound {
			return ErrUnknownPluginSoftware
		}
	}
	softwareFound := false
	for _, repo := range repos {
		if _, ok := repo.Software[config.Software]; ok {
			softwareFound = true
			break
		}
	}
	if !softwareFound {
		return ErrUnknownPluginSoftware
	}
	return nil
}

func ExecuteConfig(configFolder, configFile string) (Config, error) {
	vm := goja.New()

	vm.GlobalObject().Set("module", vm.NewObject())
	vm.GlobalObject().Set("require", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) != 1 {
			return vm.NewTypeError("require() expects a single argument")
		}
		modulePath := call.Argument(0).String()
		// Resolve module path
		if !filepath.IsAbs(modulePath) {
			modulePath = filepath.Join(configFolder, modulePath)
		}
		// Read the module file
		moduleData, err := os.ReadFile(modulePath)
		if err != nil {
			return vm.NewGoError(err)
		}
		// Execute the module code
		_, err = vm.RunString(string(moduleData))
		if err != nil {
			return vm.NewGoError(err)
		}
		// Return the exports of the module
		value := vm.GlobalObject().Get("module").ToObject(vm).Get("exports")
		return value
	})
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	_, err := vm.RunString(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	exportsJson, err := vm.RunString("JSON.stringify(module.exports)")
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal([]byte(exportsJson.String()), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
