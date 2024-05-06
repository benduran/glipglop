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

// Gets glipglop's download location
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

// Gets glipgops's extracted forlder location
func GetExtractedCacheLocation() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	extractedLocation := filepath.Join(homedir, ".glipglop", "extracted")

	// ensure this directory is created
	if err := os.MkdirAll(extractedLocation, os.ModePerm); err != nil {
		return "", err
	}

	return extractedLocation, nil
}
