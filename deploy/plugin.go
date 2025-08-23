package deploy

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/repos"
	"github.com/mythicmc/peridot/utils"
)

type PluginUpdateOperation struct {
	PluginName  string
	CurrentPath string
	UpdatePath  string
	PrevVersion string
	NewVersion  string
}

func PrepareAllPluginUpdates(
	repos repos.Repositories, configs config.Configs,
) (map[string]map[string]PluginUpdateOperation, error) {
	operations := make(map[string]map[string]PluginUpdateOperation)
	for name, config := range configs {
		operation, err := PreparePluginUpdates(repos, name, config)
		if err != nil {
			return nil, errors.Join(PrepareUpdateError{Name: name, Type: "plugin"}, err)
		} else if len(operation) > 0 {
			operations[name] = operation
		}
	}
	return operations, nil
}

func PreparePluginUpdates(
	repositories repos.Repositories, server string, config config.Config,
) (map[string]PluginUpdateOperation, error) {
	// Get all plugins in server
	files, err := os.ReadDir(filepath.Join(config.Location, "plugins"))
	if err != nil {
		return nil, err
	}
	updates := make(map[string]PluginUpdateOperation)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jar") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(config.Location, "plugins", file.Name()))
		if err != nil {
			return nil, err
		}

		// Check the plugin version... if we can't determine it, skip this suspicious file
		_, metadataFile, err := utils.DetermineJarType(data)
		if err != nil {
			if errors.Is(err, utils.ErrUnknownJarType) {
				log.Printf("Warning: %s is not a recognized JAR type, skipping...\n", file.Name())
				continue
			}
			return nil, err
		}
		metadata, err := utils.ParsePluginMetadata(file.Name(), metadataFile)
		if err != nil {
			log.Printf("Warning: Failed to load plugin metadata from %s, skipping: %v\n", file.Name(), err)
			continue
		}
		hash := utils.HashData(data)

		// FIXME: If the plugin isn't in the config, remove it

		// If the plugin is out of date, create an update op
		plugin, err := repositories.GetPlugin(metadata.Name, config.Repos)
		if err != nil {
			// FIXME: This is unexpected
		} else if plugin.Version != metadata.Version || plugin.Checksum != hash {
			updates[plugin.Name] = PluginUpdateOperation{
				// FIXME
			}
		}
	}

	// FIXME: If any plugin is missing, create an update op
	return updates, nil
}
