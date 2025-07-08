package config

import (
	"encoding/json"
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
		configs[configName] = config
	}
	return configs, nil
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
