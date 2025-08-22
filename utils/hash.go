package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

func HashFilePath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return strings.ToLower(hex.EncodeToString(hash.Sum(nil))), nil
}

func HashData(data []byte) string {
	hash := sha256.Sum256(data)
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
