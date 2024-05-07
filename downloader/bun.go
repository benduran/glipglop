package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
)

// Downloads Bun JS runtime version
// the user specified in their config
// and extracts it
func DownloadBun(version string) (string, error) {
	// linux format:      https://github.com/oven-sh/bun/releases/download/bun-v1.1.6/bun-linux-x64.zip
	// mac ARM64 format:  https://github.com/oven-sh/bun/releases/download/bun-v1.1.6/bun-darwin-aarch64.zip
	// mac ARMx64 format: https://github.com/oven-sh/bun/releases/download/bun-v1.1.6/bun-darwin-x64.zip
	// windows format:    https://github.com/oven-sh/bun/releases/download/bun-v1.1.6/bun-windows-x64.zip

	machineInfo := internal.GetMachineInfo()

	if machineInfo.Unsupported {
		return "", fmt.Errorf("unable to download bun because you are running an unsupported os + arch")
	}

	ext := ".zip"
	bunBinaryExt := ""

	if machineInfo.OS == "windows" {
		bunBinaryExt = ".exe"
	}

	arch := ""

	if machineInfo.Arch == "amd64" {
		arch = "x64"
	} else {
		arch = "aarch64"
	}

	filenamePrefix := fmt.Sprintf("bun-v%s", version)
	filename := fmt.Sprintf("bun-%s-%s%s", machineInfo.OS, arch, ext)

	urlToDownload := fmt.Sprintf("https://github.com/oven-sh/bun/releases/download/%s/%s", filenamePrefix, filename)

	downloadPath, err := internal.DownloadFileFromURL(urlToDownload)

	if err != nil {
		return "", err
	}

	// now we need to extract the archive
	extractedPath, err := internal.ExtractArchive(downloadPath)

	if err != nil {
		return "", err
	}

	// find the bun
	bunGlob := filepath.Join(extractedPath, fmt.Sprintf("bun%s", bunBinaryExt))

	logger.Info(fmt.Sprintf("scanning for bun binary with the following glob path: %s", bunGlob))
	matches, err := filepath.Glob(bunGlob)

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no matching bun binary was found in %s", extractedPath)
	}

	bunBinary := matches[0]

	logger.Info(fmt.Sprintf("Found the bun binary to be %s", bunBinary))
	return internal.MoveBinaryToToolCache("bun", version, bunBinary)
}
