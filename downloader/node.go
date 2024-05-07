package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
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
	nodeBinaryExt := ""

	switch machineInfo.OS {
	case "windows":
		ext = ".zip"
		nodeBinaryExt = ".exe"
	case "darwin":
		ext = ".tar.xz"
	default:
		ext = ".tar.gz"
	}

	arch := ""

	if machineInfo.Arch == "amd64" {
		arch = "x64"
	} else {
		arch = "arm64"
	}

	filenamePrefix := fmt.Sprintf("node-v%s", version)
	filename := fmt.Sprintf("%s-%s-%s%s", filenamePrefix, machineInfo.OS, arch, ext)

	urlToDownload := fmt.Sprintf("https://nodejs.org/dist/v%s/%s", version, filename)

	downloadPath, err := internal.DownloadFileFromURL(urlToDownload)

	if err != nil {
		return "", err
	}

	// now we need to extract the archive
	extractedPath, err := internal.ExtractArchive(downloadPath)

	if err != nil {
		return "", err
	}

	// find the node binary
	nodeGlob := filepath.Join(extractedPath, "bin", fmt.Sprintf("node%s", nodeBinaryExt))

	if machineInfo.OS == "windows" {
		nodeGlob = filepath.Join(extractedPath, "**", "node.exe")
	}

	logger.Info(fmt.Sprintf("scanning for node binary with the following glob path: %s", nodeGlob))
	matches, err := filepath.Glob(nodeGlob)

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no matching node binary was found in %s", extractedPath)
	}

	nodeBinary := matches[0]

	logger.Info(fmt.Sprintf("Found the node binary to be %s", nodeBinary))
	return internal.MoveBinaryToToolCache("node", version, nodeBinary)
}
