package downloader

import (
	"fmt"

	logger "github.com/benduran/glipglop/log"
)

// Downloads a specific tool
func DownloadTool(tool string, version string) (string, error) {
	if tool == "node" {
		return DownloadNode(version)
	}

	err := fmt.Errorf("%s not downloading because it is not currently supported", tool)
	logger.Error(err)
	return "", err
}
