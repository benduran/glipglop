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
