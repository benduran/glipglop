package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
	logger "github.com/benduran/glipglop/log"
)

// Given info about a specific tool that's been successfully downloaded
// and extracted, updates its chmod permissions and moves it to the correct
// cache location
func MoveBinaryToToolCache(toolName, toolVersion, binaryPath string) (string, error) {
	logger.Info(fmt.Sprintf("Updating %s execution permissions...", binaryPath))
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", err
	}

	// great, we didn't blow up, so now we need
	// to simplify the extraction folder
	toolCacheLocation, err := cache.GetToolCacheLocation()
	if err != nil {
		return "", err
	}

	extractionLocation := cache.GetToolCacheLocationForTool(toolCacheLocation, toolName, toolVersion)

	if err := os.MkdirAll(filepath.Dir(extractionLocation), os.ModePerm); err != nil {
		return "", err
	}

	if err := os.Rename(binaryPath, extractionLocation); err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("%s %s is cached locally at %s", toolName, toolVersion, extractionLocation))

	return extractionLocation, nil
}

// some tools aren't good at operating as a single, standalone binary.
// that's where this function can be used: to copy all of the extracted
// contents to the tool cache.
// it requires a few more pieces to define, namely, the folder that needs to be moved,
// as well as an array of binary files that all need to have their execution permissions
// updated to 0755
func MoveFolderToToolCache(toolName, toolVersion, toolFolderPath string, binaries []string) (string, error) {
	for _, binaryPath := range binaries {
		logger.Info(fmt.Sprintf("Updating %s execution permissions...", binaryPath))
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return "", err
		}
	}

	toolCacheLocation, err := cache.GetToolCacheLocation()
	if err != nil {
		return "", err
	}

	extractionLocation := cache.GetToolCacheLocationForTool(toolCacheLocation, toolName, toolVersion)

	if err := os.MkdirAll(filepath.Dir(extractionLocation), os.ModePerm); err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("Moving %s to %s", toolFolderPath, extractionLocation))
	if err := os.Rename(toolFolderPath, extractionLocation); err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("%s %s is cached locally at %s", toolName, toolVersion, extractionLocation))
	return extractionLocation, nil
}
