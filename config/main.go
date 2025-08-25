package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mythicmc/peridot/repos"
)

type Config struct {
	Location         string            `json:"location"`
	Repos            []string          `json:"repos"`
	Software         string            `json:"software"`          // Supported: "vanilla", "paper", "velocity"
	ServerProperties map[string]string `json:"server_properties"` // FIXME: support float64, bool
	Plugins          []string          `json:"plugins"`
}

type Configs map[string]Config

type ConfigLoadError struct {
	Name string
	Step string
}

func (e ConfigLoadError) Error() string {
	return "failed to load config " + e.Name + " while " + e.Step
}

func LoadConfigs(repositories repos.Repositories) (Configs, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configs := make(Configs)
	configFolder := filepath.Join(wd, "configs")
	configFiles, err := os.ReadDir(configFolder)
	if err != nil && os.IsNotExist(err) {
		log.Println("Warning: configs/ folder does not exist! No configurations loaded.")
		return configs, nil
	} else if err != nil {
		return nil, err
	}
	for _, configFile := range configFiles {
		if configFile.IsDir() || !strings.HasSuffix(configFile.Name(), ".js") {
			continue
		}
		filePath := filepath.Join(configFolder, configFile.Name())
		configName := configFile.Name()[:strings.LastIndex(configFile.Name(), ".")]
		configJs, err := os.ReadFile(filePath)
		if err != nil {
			return nil, errors.Join(ConfigLoadError{Name: configName, Step: "reading file"}, err)
		}
		config, err := ExecuteConfig(configFolder, string(configJs))
		if err != nil {
			return nil, errors.Join(ConfigLoadError{Name: configName, Step: "executing JS"}, err)
		}
		err = ValidateConfig(config, repositories)
		if err != nil {
			return nil, errors.Join(ConfigLoadError{Name: configName, Step: "validating config"}, err)
		}
		configs[configName] = config
	}
	return configs, nil
}
