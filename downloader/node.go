package downloader

import (
	"fmt"

	"github.com/benduran/glipglop/internal"
)

// Downloads the Node.js version
// the user specified in their config
// and extracts it
func DownloadNode(version string) (string, error) {
	// linux format: https://nodejs.org/dist/v20.12.0/node-v20.12.0-linux-x64.tar.xz
	// mac format: https://nodejs.org/dist/v20.12.0/node-v20.12.0-darwin-arm64.tar.gz
	// windows format: https://nodejs.org/dist/v20.12.0/node-v20.12.0-win-x64.zip

	machineInfo := internal.GetMachineInfo()

	if machineInfo.Unsupported {
		return "", fmt.Errorf("unable to download node.js because you are running an unsupported os + arch")
	}

	ext := ""

	if machineInfo.OS == "windows" {
		ext = ".zip"
	} else if machineInfo.OS == "darwin" {
		ext = ".tar.xz"
	} else {
		ext = ".tar.gz"
	}

	filename := fmt.Sprintf("node-v%s-%s-%s%s", version, machineInfo.OS, machineInfo.Arch, ext)

	urlToDownload := fmt.Sprintf("https://nodejs.org/dist/v%s/%s", version, filename)

	return internal.DownloadFileFromURL(urlToDownload)
}
