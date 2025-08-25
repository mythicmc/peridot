package deploy

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/repos"
	"github.com/mythicmc/peridot/utils"
)

type SoftwareUpdateOperation struct {
	SoftwareType string
	CurrentPath  string
	UpdatePath   string
	PrevHash     string
	NewHash      string
}

var ErrSoftwareNotInRepos = errors.New("software not found in repositories")

func PrepareAllSoftwareUpdate(
	repos repos.Repositories, configs config.Configs,
) (map[string]SoftwareUpdateOperation, error) {
	operations := make(map[string]SoftwareUpdateOperation)
	for name, config := range configs {
		operation, err := PrepareSoftwareUpdate(repos, name, config)
		if err != nil {
			return nil, errors.Join(PrepareUpdateError{Name: name, Type: "software"}, err)
		} else if operation != (SoftwareUpdateOperation{}) {
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
		NewHash:      software.Checksum[:8],
	}
	prevHash, err := utils.HashFilePath(operation.CurrentPath)
	if err != nil && os.IsNotExist(err) {
		operation.CurrentPath = ""
		operation.PrevHash = ""
	} else if err != nil {
		return SoftwareUpdateOperation{}, err
	} else {
		operation.PrevHash = prevHash[:8]
	}

	if prevHash == software.Checksum {
		return SoftwareUpdateOperation{}, nil
	}
	return operation, nil
}
