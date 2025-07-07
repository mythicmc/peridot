package repos

import (
	"log"
	"os"
	"path/filepath"
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
		_, err = DetermineJarType(jarData)
		// FIXME: Garbage code
		plugin, err := LoadPlugin(jarPath)
		if err != nil {
			return repo, err
		}
		if _, ok := repo.Plugins[plugin.Name]; ok {
			// TODO: Handle conflicts
		}
		repo.Plugins[plugin.Name] = plugin
	}
	return repo, nil
}

// DetermineJarType checks what type of software the given JAR is.
// Supported types are:
// - "vanilla"
// - "paper"
// - "velocity"
// - "plugin"
func DetermineJarType(jar []byte) (string, error) {
	/* r, err := zip.NewReader(bytes.NewReader(jar), int64(len(jar)))
	if err != nil {
		return "", err
	}
	// FIXME: Do something meaningful
	slices.ContainsFunc(r.File, func(f *zip.File) bool {
		return f.Name == "META-INF/MANIFEST.MF"
	}) */
	return "", nil
}

func LoadPlugin(path string) (Plugin, error) {
	// TODO
	return Plugin{}, nil
}
