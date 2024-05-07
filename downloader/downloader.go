package downloader

import (
	"fmt"

	logger "github.com/benduran/glipglop/log"
	"github.com/benduran/glipglop/schema"
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

// downloads all of the tools in the user's glipglop manifest.
// if some of them have already been downloaded, skip those
func DownloadAllTools(cwd string) error {
	logger.Info("Downloading all of your tools now...")
	schema, err := schema.ReadUserSchema(cwd)

	if err != nil {
		return err
	}

	for toolName, toolVersion := range schema.Tools {
		_, err := DownloadTool(toolName, toolVersion)
		if err != nil {
			return err
		}
	}

	return nil
}
