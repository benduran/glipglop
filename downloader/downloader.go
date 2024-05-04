package downloader

import "fmt"

// Downloads a specific tool
func DownloadTool(tool string, version string) (string, error) {
	if tool == "node" {
		return DownloadNode(version)
	}

	return "", fmt.Errorf("%s not downloading because it is not currently supported", tool)
}
