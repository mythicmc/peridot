package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"slices"

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

var ErrUnknownJarType = errors.New("unknown JAR type")

// DetermineJarType checks what type of software the given JAR is.
// Supported types are:
// - "vanilla"
// - "paper"
// - "velocity"
// - "plugin"
func DetermineJarType(jar []byte) (string, []byte, error) {
	r, err := zip.NewReader(bytes.NewReader(jar), int64(len(jar)))
	if err != nil {
		return "", nil, err
	}
	isVelocity := false
	isPaper := false
	isVanilla := false
	isPlugin := slices.IndexFunc(r.File, func(file *zip.File) bool {
		switch file.Name {
		case "com/velocitypowered/proxy/Velocity.class":
			isVelocity = true
		case "io/papermc/paperclip/Paperclip.class":
			isPaper = true
		case "net/minecraft/server/MinecraftServer.class":
			isVanilla = true
		case "net/minecraft/bundler/Main.class":
			isVanilla = true
		}
		return file.Name == "plugin.yml" ||
			file.Name == "bungee.yml" ||
			file.Name == "velocity-plugin.json"
	})
	if isPlugin != -1 {
		metadataFile, err := r.File[isPlugin].Open()
		if err != nil {
			return "", nil, err
		}
		defer metadataFile.Close()
		metadata, err := io.ReadAll(metadataFile)
		if err != nil {
			return "", nil, err
		}
		return "plugin", metadata, nil
	} else if isVelocity {
		return "velocity", nil, nil
	} else if isPaper {
		return "paper", nil, nil
	} else if isVanilla {
		return "vanilla", nil, nil
	}
	return "", nil, ErrUnknownJarType
}
