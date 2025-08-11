package deploy

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type PrepareUpdateError struct {
	Name string
	Type string
}

func (e PrepareUpdateError) Error() string {
	return "failed to prepare " + e.Type + " update for " + e.Name
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil))[:8], nil
}
