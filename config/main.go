package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

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
		_, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		_ = configFile.Name()[:strings.LastIndex(configFile.Name(), ".")]
		// FIXME: Execute the JS files to get the configs
	}
	return configs, nil
}
