package repos

import (
	"github.com/goccy/go-yaml"
)

type PluginMetadata struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type InvalidPluginMetadataError struct{ FileName string }

func (e InvalidPluginMetadataError) Error() string {
	return "invalid metadata for plugin " + e.FileName + ": missing name or version"
}

// ParsePluginMetadata retrieves the plugin name and version from its metadata file.
func ParsePluginMetadata(filename string, metadataFile []byte) (PluginMetadata, error) {
	var pluginMetadata PluginMetadata
	err := yaml.Unmarshal(metadataFile, &pluginMetadata)
	if err != nil {
		return PluginMetadata{}, err
	} else if pluginMetadata.Name == "" || pluginMetadata.Version == "" {
		return PluginMetadata{}, InvalidPluginMetadataError{FileName: filename}
	}
	return pluginMetadata, nil
}
