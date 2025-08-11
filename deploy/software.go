package deploy

import (
	"errors"
	"path/filepath"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/repos"
)

type SoftwareUpdateOperation struct {
	SoftwareType string
	CurrentPath  string
	UpdatePath   string
	PrevVersion  string
	NewVersion   string
}

var ErrSoftwareNotInRepos = errors.New("software not found in repositories")

func PrepareAllSoftwareUpdate(
	repos repos.Repositories, configs config.Configs,
) (map[string]SoftwareUpdateOperation, error) {
	operations := make(map[string]SoftwareUpdateOperation)
	for name, config := range configs {
		operation, err := PreparePluginUpdates(repos, name, config)
		if err != nil {
			return nil, errors.Join(PrepareUpdateError{Name: name, Type: "software"}, err)
		} else if operation.CurrentPath != "" {
			operations[name] = operation
		}
	}
	return operations, nil
}

func PrepareSoftwareUpdate(
	repositories repos.Repositories, server string, config config.Config,
) (SoftwareUpdateOperation, error) {
	var software repos.Software
	for _, repoName := range config.Repos {
		repo := repositories[repoName]
		software = repo.Software[config.Software]
	}
	if software.Path == "" {
		return SoftwareUpdateOperation{}, ErrSoftwareNotInRepos
	}

	operation := SoftwareUpdateOperation{
		SoftwareType: config.Software,
		CurrentPath:  filepath.Join(config.Location, config.Software+".jar"),
		UpdatePath:   software.Path,
		NewVersion:   software.Checksum[:8],
	}
	prevHash, err := hashFile(operation.CurrentPath)
	if err != nil {
		return SoftwareUpdateOperation{}, err
	}
	operation.PrevVersion = prevHash[:8]

	if prevHash == software.Checksum {
		return SoftwareUpdateOperation{}, nil
	}
	return operation, nil
}
