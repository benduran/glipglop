package downloader

import (
	"fmt"

	logger "github.com/benduran/glipglop/log"
)

// Downloads a specific tool
func DownloadTool(tool string, version string) (string, error) {
	switch tool {
	case "bun":
		return DownloadBun(version)
	case "node":
		return DownloadNode(version)
	default:
		err := fmt.Errorf("%s not downloading because it is not currently supported", tool)
		logger.Error(err)
		return "", err
	}
}
