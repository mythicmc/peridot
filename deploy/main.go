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

type PrepareUpdateError struct {
	Name string
	Type string
}

func (e PrepareUpdateError) Error() string {
	return "failed to prepare " + e.Type + " update for " + e.Name
}

func PrepareAllSoftwareUpdate(
	repos repos.Repositories, configs config.Configs,
) (map[string]SoftwareUpdateOperation, error) {
	operations := make(map[string]SoftwareUpdateOperation)
	for name, config := range configs {
		operation, err := PrepareSoftwareUpdate(repos, name, config)
		if err != nil {
			return nil, errors.Join(PrepareUpdateError{Name: name, Type: "software"}, err)
		} else if operation.CurrentPath != "" {
			operations[name] = operation
		}
	}
	return operations, nil
}

func PrepareSoftwareUpdate(
	repos repos.Repositories, server string, config config.Config,
) (SoftwareUpdateOperation, error) {
	operation := SoftwareUpdateOperation{
		SoftwareType: config.Software,
		CurrentPath:  filepath.Join(config.Location, config.Software+".jar"),
	}
	prevHash, err := hashFile(operation.CurrentPath)
	if err != nil {
		return SoftwareUpdateOperation{}, err
	}
	operation.PrevVersion = prevHash[:8]

	for _, repoName := range config.Repos {
		repo := repos[repoName]
		operation.UpdatePath = repo.Software[config.Software].Path
	}
	updateHash, err := hashFile(operation.UpdatePath)
	if err != nil {
		return SoftwareUpdateOperation{}, err
	}
	operation.NewVersion = updateHash[:8]

	if prevHash == updateHash {
		return SoftwareUpdateOperation{}, nil
	}
	return operation, nil
}
