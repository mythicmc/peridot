package repos

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
		Name:    name,
		Plugins: make(map[string]Plugin),
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
		jarType, metadataFile, err := DetermineJarType(jarData)
		if err != nil {
			if errors.Is(err, ErrUnknownJarType) {
				log.Printf("Warning: %s is not a recognized JAR type, skipping...\n", jarPath)
				continue
			}
			return repo, err
		}
		hash := sha256.Sum256(jarData)
		if jarType == "vanilla" || jarType == "paper" || jarType == "velocity" {
			stat, err := os.Stat(jarPath)
			if err != nil {
				return Repository{}, err
			}
			updatedAt := stat.ModTime().String()
			if existingSoftware, exists := repo.Software[jarType]; exists {
				if strings.Compare(existingSoftware.UpdatedAt, updatedAt) < 0 {
					log.Printf("Warning: Replacing software %s with timestamp %s with newer timestamp %s\n",
						jarType, existingSoftware.UpdatedAt, updatedAt)
				} else {
					log.Printf("Warning: Skipping software %s with timestamp %s (already have timestamp %s!)\n",
						jarType, updatedAt, existingSoftware.UpdatedAt)
					continue
				}
			}
			repo.Software[jarType] = Software{
				Type:      jarType,
				Path:      jarPath,
				UpdatedAt: updatedAt,
				Checksum:  strings.ToLower(hex.EncodeToString(hash[:])),
			}
		} else {
			pluginMetadata, err := ParsePluginMetadata(file.Name(), metadataFile)
			if err != nil {
				log.Printf("Warning: Failed to load plugin metadata from %s, skipping: %v\n", jarPath, err)
				continue
			}
			if existingPlugin, exists := repo.Plugins[pluginMetadata.Name]; exists {
				if strings.Compare(existingPlugin.Version, pluginMetadata.Version) < 0 {
					log.Printf("Warning: Replacing plugin %s with version %s with newer version %s\n",
						pluginMetadata.Name, existingPlugin.Version, pluginMetadata.Version)
				} else {
					log.Printf("Warning: Skipping plugin %s with version %s (already have version %s!)\n",
						pluginMetadata.Name, pluginMetadata.Version, existingPlugin.Version)
					continue
				}
			}
			repo.Plugins[pluginMetadata.Name] = Plugin{
				Name:     pluginMetadata.Name,
				Path:     jarPath,
				Version:  pluginMetadata.Version,
				Checksum: strings.ToLower(hex.EncodeToString(hash[:])),
			}
		}
	}
	return repo, nil
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
