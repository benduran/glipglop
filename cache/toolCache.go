package cache

import (
	"os"
	"path/filepath"
)

// Gets glipglop's tool cache location,
// which is shared by all glipglop instances
func GetToolCacheLocation() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheLocation := filepath.Join(homedir, ".glipglop", "cache")

	// ensure this directory is created
	if err := os.MkdirAll(cacheLocation, os.ModePerm); err != nil {
		return "", err
	}

	return cacheLocation, nil
}

func GetDownloadCacheLocation() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	downloadsLocation := filepath.Join(homedir, ".glipglop", "downloads")

	// ensure this directory is created
	if err := os.MkdirAll(downloadsLocation, os.ModePerm); err != nil {
		return "", err
	}

	return downloadsLocation, nil
}
