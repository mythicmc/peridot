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
		config, err := ExecuteConfig(string(configJs))
		if err != nil {
			return nil, err
		}
		configs[configName] = config
	}
	return configs, nil
}

func ExecuteConfig(configFile string) (Config, error) {
	vm := goja.New()
	// FIXME: Provide `require` CommonJS implementation...
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.GlobalObject().Set("module", vm.NewObject())
	_, err := vm.RunString(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	exportsJsonValue, err := vm.RunString("JSON.stringify(module.exports)")
	var exportsJson string
	err = vm.ExportTo(exportsJsonValue, &exportsJson)
	if err != nil {
		return Config{}, nil
	}
	err = json.Unmarshal([]byte(exportsJson), &config)
	if err != nil {
		return Config{}, nil
	}

	return config, nil
}
