package repos

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mythicmc/peridot/utils"
)

type Repositories map[string]Repository

type Repository struct {
	Name     string
	Software map[string]Software
	Plugins  map[string]Plugin
}

type Software struct {
	Type      string
	Path      string
	UpdatedAt string
	Checksum  string
}

type Plugin struct {
	Name     string
	Path     string
	Version  string
	Checksum string
}

type RepositoryLoadError struct{ Name string }

func (e RepositoryLoadError) Error() string { return "failed to load repository " + e.Name }

func LoadRepositories() (Repositories, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	repositories := make(Repositories)
	repoFolder := filepath.Join(wd, "repos")
	repoFolders, err := os.ReadDir(repoFolder)
	if err != nil && os.IsNotExist(err) {
		log.Println("Warning: repos/ folder does not exist! No repositories loaded.")
		return repositories, nil
	} else if err != nil {
		return nil, err
	}
	for _, folder := range repoFolders {
		if !folder.IsDir() {
			continue
		}
		repoPath := filepath.Join(repoFolder, folder.Name())
		repositories[folder.Name()], err = LoadRepository(repoPath, folder.Name())
		if err != nil {
			return nil, errors.Join(RepositoryLoadError{Name: folder.Name()}, err)
		}
	}
	return repositories, nil
}

func LoadRepository(path, name string) (Repository, error) {
	repo := Repository{
		Name:     name,
		Plugins:  make(map[string]Plugin),
		Software: make(map[string]Software),
	}
	jars, err := os.ReadDir(path)
	if err != nil {
		return repo, err
	}
	for _, file := range jars {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jar") {
			continue
		}
		jarPath := filepath.Join(path, file.Name())
		jarData, err := os.ReadFile(jarPath)
		if err != nil {
			return repo, err
		}
		jarType, metadataFile, err := utils.DetermineJarType(jarData)
		if err != nil {
			if errors.Is(err, utils.ErrUnknownJarType) {
				log.Printf("Warning in repo %s: %s is not a recognized JAR type, skipping...\n", name, jarPath)
				continue
			}
			return repo, err
		}
		hash := utils.HashData(jarData)
		if jarType == "vanilla" || jarType == "paper" || jarType == "velocity" {
			stat, err := os.Stat(jarPath)
			if err != nil {
				return Repository{}, err
			}
			updatedAt := stat.ModTime().String()
			if existingSoftware, exists := repo.Software[jarType]; exists {
				if strings.Compare(existingSoftware.UpdatedAt, updatedAt) < 0 {
					log.Printf("Warning in repo %s: Replacing software %s with timestamp %s with newer timestamp %s\n",
						name, jarType, existingSoftware.UpdatedAt, updatedAt)
				} else {
					log.Printf("Warning in repo %s: Skipping software %s with timestamp %s (already have timestamp %s!)\n",
						name, jarType, updatedAt, existingSoftware.UpdatedAt)
					continue
				}
			}
			repo.Software[jarType] = Software{
				Type:      jarType,
				Path:      jarPath,
				UpdatedAt: updatedAt,
				Checksum:  hash,
			}
		} else {
			pluginMetadata, err := utils.ParsePluginMetadata(file.Name(), metadataFile)
			if err != nil {
				log.Printf("Warning in repo %s: Failed to load plugin metadata from %s, skipping: %v\n",
					name, jarPath, err)
				continue
			}
			if existingPlugin, exists := repo.Plugins[pluginMetadata.Name]; exists {
				if strings.Compare(existingPlugin.Version, pluginMetadata.Version) < 0 {
					log.Printf("Warning in repo %s: Replacing plugin %s with version %s with newer version %s\n",
						name, pluginMetadata.Name, existingPlugin.Version, pluginMetadata.Version)
				} else {
					log.Printf("Warning in repo %s: Skipping plugin %s with version %s (already have version %s!)\n",
						name, pluginMetadata.Name, pluginMetadata.Version, existingPlugin.Version)
					continue
				}
			}
			repo.Plugins[pluginMetadata.Name] = Plugin{
				Name:     pluginMetadata.Name,
				Path:     jarPath,
				Version:  pluginMetadata.Version,
				Checksum: hash,
			}
		}
	}
	return repo, nil
}

var ErrPluginNotInRepos = errors.New("plugin not found in repositories")

func (r Repositories) GetPlugin(name string, repos []string) (Plugin, error) {
	var plugin Plugin
	found := false
	for _, repoName := range repos {
		if repo, exists := r[repoName]; exists {
			if p, exists := repo.Plugins[name]; exists {
				plugin = p
				found = true
			}
		}
	}
	if found {
		return plugin, nil
	}
	return Plugin{}, ErrPluginNotInRepos
}
