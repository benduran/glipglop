package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
)

// downloads one of the official go archives
// from the go website, and does all the magic
// to get it setup on disk
func DownloadGo(version string) (string, error) {
	// linux format:   https://go.dev/dl/go1.22.3.linux-386.tar.gz
	// mac format:     https://go.dev/dl/go1.22.3.darwin-arm64.tar.gz
	// windows format: https://go.dev/dl/go1.22.3.windows-amd64.zip

	// if the tool is already in the cache, just return the path to that immediately
	existingPathToTool := cache.CheckBinaryInToolCache("node", version)

	if len(existingPathToTool) > 0 {
		return existingPathToTool, nil
	}

	machineInfo := internal.GetMachineInfo()

	if machineInfo.Unsupported {
		return "", fmt.Errorf("unable to download go because you are running an unsupported os + arch")
	}

	arch := machineInfo.Arch
	ext := ""
	goBinaryExt := ""

	switch machineInfo.OS {
	case "windows":
		ext = ".zip"
		goBinaryExt = ".exe"
	case "linux":
		arch = "386"
	default:
		ext = ".tar.gz"
	}

	filename := fmt.Sprintf("go%s.%s-%s%s", version, machineInfo.OS, arch, ext)

	urlToDownload := fmt.Sprintf("https://go.dev/dl/%s", filename)

	downloadPath, err := internal.DownloadFileFromURL(urlToDownload)

	if err != nil {
		return "", err
	}

	// now we need to extract the archive
	extractedPath, err := internal.ExtractArchive("go", version, downloadPath)

	if err != nil {
		return "", err
	}

	// find the node binary
	goGlob := filepath.Join(extractedPath, "bin", fmt.Sprintf("go%s", goBinaryExt))

	logger.Info(fmt.Sprintf("scanning for go binary with the following glob path: %s", goGlob))
	matches, err := filepath.Glob(goGlob)

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no matching go binary was found in %s", extractedPath)
	}

	nodeBinary := matches[0]

	logger.Info(fmt.Sprintf("Found the go binary to be %s", nodeBinary))
	return internal.MoveBinaryToToolCache("go", version, nodeBinary)
}
