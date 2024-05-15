package cache

import (
	"fmt"
	"os"
	"path/filepath"

	logger "github.com/benduran/glipglop/log"
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

// shortcut for getting a formatted folder that holds the tool binary
func GetToolCacheLocationForTool(toolCacheLocation, toolName, toolVersion string) string {
	filenamePrefix := fmt.Sprintf("%s-v%s", toolName, toolVersion)

	return filepath.Join(toolCacheLocation, filenamePrefix, toolName)
}

// checks if a specific tool version already exists in the tools cache.
// if it does, that path is returned.
// otherwise, an empty string is returned
func CheckBinaryInToolCache(toolName, toolVersion string) string {
	toolCacheLocation, err := GetToolCacheLocation()

	if err != nil {
		return ""
	}

	toolLocation := GetToolCacheLocationForTool(toolCacheLocation, toolName, toolVersion)

	// need to have a bunch of special cases for the various different versions of tools
	switch toolName {
	case "go":
		toolLocation = filepath.Join(toolLocation, "bin", "go")
	}

	stat, err := os.Stat(toolLocation)
	if err != nil || stat.IsDir() {
		return ""
	}

	logger.Debug(fmt.Sprintf("found %s location to be %s", toolName, toolLocation))

	return toolLocation
}
