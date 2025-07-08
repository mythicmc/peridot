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
	Type     string
	Path     string
	Checksum string
}

type Plugin struct {
	Name     string
	Path     string
	Version  string
	Checksum string
}

func LoadRepositories() (Repositories, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	repoFolder := filepath.Join(wd, "repos")
	repositories := make(Repositories)
	repoFolders, err := os.ReadDir(repoFolder)
	if err != nil && os.IsNotExist(err) {
		log.Println("Warning: repos/ folder does not exist!")
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
			return nil, err
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
		jarType, metadata, err := DetermineJarType(jarData)
		if err != nil {
			if errors.Is(err, ErrUnknownJarType) {
				log.Printf("Warning: %s is not a recognized JAR type, skipping...\n", jarPath)
				continue
			}
			return repo, err
		}
		// TODO: Handle conflicts
		hash := sha256.Sum256(jarData)
		if jarType == "vanilla" || jarType == "paper" || jarType == "velocity" {
			repo.Software[jarType] = Software{
				Type:     jarType,
				Path:     jarPath,
				Checksum: strings.ToLower(hex.EncodeToString(hash[:])),
			}
		} else {
			name, version, err := ParsePluginMetadata(metadata)
			if err != nil {
				log.Printf("Warning: Failed to load plugin metadata from %s, skipping: %v\n", jarPath, err)
				continue
			}
			repo.Plugins[name] = Plugin{
				Name:     name,
				Path:     jarPath,
				Version:  version,
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

// ParsePluginMetadata retrieves the plugin name and version from its metadata file.
func ParsePluginMetadata(metadata []byte) (string, string, error) {
	// TODO: Parse YAML
	return "", "", errors.New("not implemented")
}
