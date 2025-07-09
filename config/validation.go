package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mythicmc/peridot/repos"
)

func ValidateConfig(config Config, repositories repos.Repositories) error {
	if err := validateConfigLocation(config); err != nil {
		return err
	}

	if err := validateConfigRepositories(config, repositories); err != nil {
		return err
	}

	if err := validateConfigSoftware(config, repositories); err != nil {
		return err
	}

	if err := validateConfigPlugins(config, repositories); err != nil {
		return err
	}

	return nil
}

var ErrInvalidLocationPath = errors.New("invalid location: must be an absolute path")

var ErrLocationNotExists = errors.New("invalid location: parent of location is not a valid folder")

func validateConfigLocation(config Config) error {
	if config.Location == "" || !filepath.IsAbs(config.Location) {
		return ErrInvalidLocationPath
	} else if stat, err := os.Stat(filepath.Dir(config.Location)); os.IsNotExist(err) {
		return ErrLocationNotExists
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return ErrLocationNotExists
	}
	return nil
}

var ErrNoReposConfigured = errors.New("invalid repositories: at least one repository must be specified")

type UnknownRepoError struct{ Name string }

func (e UnknownRepoError) Error() string { return "unknown repository specified: " + e.Name }

func validateConfigRepositories(config Config, repositories repos.Repositories) error {
	if len(config.Repos) == 0 {
		return ErrNoReposConfigured
	}
	repos := make([]repos.Repository, len(repositories))
	for _, repo := range config.Repos {
		if _, ok := repositories[repo]; !ok {
			return UnknownRepoError{Name: repo}
		} else {
			repos = append(repos, repositories[repo])
		}
	}
	return nil
}

var ErrInvalidSoftware = errors.New("invalid software type: must be one of 'vanilla', 'paper', or 'velocity'")

type UnknownPluginSoftwareError struct{ Name string }

func (e UnknownPluginSoftwareError) Error() string {
	return "unknown plugin/software specified: " + e.Name + " not found in configured repositories"
}

func validateConfigSoftware(config Config, repositories repos.Repositories) error {
	if config.Software != "vanilla" && config.Software != "paper" && config.Software != "velocity" {
		return ErrInvalidSoftware
	}
	softwareFound := false
	for _, repository := range repositories {
		if _, ok := repository.Software[config.Software]; ok {
			softwareFound = true
			break
		}
	}
	if !softwareFound {
		return UnknownPluginSoftwareError{Name: config.Software}
	}
	return nil
}

var ErrInvalidPlugin = errors.New("invalid plugin: empty string found in plugins list")

func validateConfigPlugins(config Config, repositories repos.Repositories) error {
	for _, plugin := range config.Plugins {
		if plugin == "" {
			return ErrInvalidPlugin
		}
		pluginFound := false
		for _, repo := range repositories {
			if _, ok := repo.Plugins[plugin]; ok {
				pluginFound = true
				break
			}
		}
		if !pluginFound {
			return UnknownPluginSoftwareError{Name: plugin}
		}
	}
	return nil
}
