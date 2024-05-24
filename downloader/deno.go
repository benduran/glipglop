package downloader

import (
	"fmt"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
)

// Downloads the Deno runtime version
// the user specified in their config
// and extracts it
func DownloadDeno(version string) (string, error) {
	// linux format:      https://github.com/denoland/deno/releases/download/v1.43.5/deno-x86_64-unknown-linux-gnu.zip
	// mac ARM64 format:  https://github.com/denoland/deno/releases/download/v1.43.5/deno-aarch64-apple-darwin.zip
	// mac x64 format:    https://github.com/denoland/deno/releases/download/v1.43.5/deno-x86_64-apple-darwin.zip
	// windows format:    https://github.com/denoland/deno/releases/download/v1.43.5/denor-x86_64-pc-windows-msvc.zip

	// if the tool is already in the cache, just return the path to that immediately
	existingPathToTool := cache.CheckBinaryInToolCache("deno", version)

	if len(existingPathToTool) > 0 {
		return existingPathToTool, nil
	}

	machineInfo := internal.GetMachineInfo()

	if machineInfo.Unsupported {
		return "", fmt.Errorf("unable to download deno because you are running an unsupported os + arch")
	}

	ext := ".zip"
	denoBinaryExt := ""

	if machineInfo.OS == "windows" {
		denoBinaryExt = ".exe"
	}

	arch := ""

	if machineInfo.Arch == "amd64" {
		arch = "x86_64"
	} else {
		arch = "aarch64"
	}
	filename := ""

	if machineInfo.OS == "darwin" {
		filename = fmt.Sprintf("deno-%s-apple-darwin%s", arch, ext)
	} else if machineInfo.OS == "linux" {
		filename = fmt.Sprintf("deno-x86_64-unknown-linux-gnu%s", ext)
	} else {
		// windows
		filename = fmt.Sprintf("deno-x86_64-pc-windows-msvc%s", ext)
	}

	urlToDownload := fmt.Sprintf("https://github.com/denoland/deno/releases/download/v%s/%s", version, filename)

	downloadPath, err := internal.DownloadFileFromURL(urlToDownload)

	if err != nil {
		return "", err
	}

	// now we need to extract the archive
	extractedPath, err := internal.ExtractArchive("deno", version, downloadPath)

	if err != nil {
		return "", err
	}

	// find deno
	denoGlob := filepath.Join(extractedPath, fmt.Sprintf("deno%s", denoBinaryExt))

	logger.Info(fmt.Sprintf("scanning for deno binary with the following glob path: %s", denoGlob))
	matches, err := filepath.Glob(denoGlob)

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no matching deno binary was found in %s", extractedPath)
	}

	denoBinary := matches[0]

	logger.Info(fmt.Sprintf("Found the deno binary to be %s", denoBinary))
	return internal.MoveBinaryToToolCache("deno", version, denoBinary)
}
