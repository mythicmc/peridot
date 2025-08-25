package deploy

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mythicmc/peridot/config"
)

type ServerPropertiesUpdateOperation struct {
	Property string
	OldValue string
	NewValue string
}

func PrepareAllerverPropertiesUpdate(
	configs config.Configs,
) (map[string][]ServerPropertiesUpdateOperation, error) {
	operations := make(map[string][]ServerPropertiesUpdateOperation)
	for name, config := range configs {
		operation, err := PrepareServerPropertiesUpdate(name, config)
		if err != nil {
			return nil, errors.Join(PrepareUpdateError{Name: name, Type: "server_properties"}, err)
		} else if len(operation) > 0 {
			operations[name] = operation
		}
	}
	return operations, nil
}

func PrepareServerPropertiesUpdate(
	server string, config config.Config,
) ([]ServerPropertiesUpdateOperation, error) {
	if config.Software != "vanilla" && config.Software != "paper" {
		return nil, nil
	}

	// Read server's server.properties file
	serverPropertiesFile, err := os.ReadFile(filepath.Join(config.Location, "server.properties"))
	if err != nil {
		return nil, err
	}
	serverProperties := make(map[string]string)
	for _, line := range strings.Split(string(serverPropertiesFile), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		serverProperties[key] = value
	}

	operations := make([]ServerPropertiesUpdateOperation, 0)
	for name, value := range config.ServerProperties {
		if serverProperties[name] != value {
			operations = append(operations, ServerPropertiesUpdateOperation{
				Property: name,
				OldValue: serverProperties[name],
				NewValue: value,
			})
		}
	}
	return operations, nil
}
