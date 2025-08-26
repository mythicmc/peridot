package deploy

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mythicmc/peridot/config"
)

type ServerPropertiesUpdateOperation struct {
	Property string
	OldValue string
	NewValue string
}

func PrepareAllServerPropertiesUpdates(
	configs config.Configs,
) (map[string][]ServerPropertiesUpdateOperation, error) {
	operations := make(map[string][]ServerPropertiesUpdateOperation)
	for name, config := range configs {
		operation, err := PrepareServerPropertiesUpdates(name, config)
		if err != nil {
			return nil, errors.Join(PrepareUpdateError{Name: name, Type: "server_properties"}, err)
		} else if len(operation) > 0 {
			operations[name] = operation
		}
	}
	return operations, nil
}

func PrepareServerPropertiesUpdates(
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
		operation := ServerPropertiesUpdateOperation{Property: name, OldValue: serverProperties[name]}
		if str, ok := value.(string); ok && serverProperties[name] != str {
			operation.NewValue = str
			operations = append(operations, operation)
		} else if num, ok := value.(float64); ok {
			valAsNum, err := strconv.ParseFloat(serverProperties[name], 64)
			if err != nil || valAsNum != num {
				operation.NewValue = strconv.FormatFloat(num, 'f', -1, 64)
				operations = append(operations, operation)
			}
		} else if boolVal, ok := value.(bool); ok && serverProperties[name] != strconv.FormatBool(boolVal) {
			operation.NewValue = strconv.FormatBool(boolVal)
			operations = append(operations, operation)
		}
	}
	if len(operations) == 0 {
		return nil, nil
	}
	return operations, nil
}
