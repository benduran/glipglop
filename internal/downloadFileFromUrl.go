package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
	logger "github.com/benduran/glipglop/log"
)

// Given a URL path to a file, downloads it
// to a cached location on disk
func DownloadFileFromURL(url string) (string, error) {
	downloadsLocation, err := cache.GetDownloadCacheLocation()

	if err != nil {
		return "", err
	}

	filename := filepath.Base(url)
	downloadPath := filepath.Join(downloadsLocation, filename)

	// check if the download already exits, for some reason,
	// and skip the download if so
	_, err = os.Stat(downloadPath)

	if err == nil {
		logger.Info(fmt.Sprintf("%s already exists. Skipping download!", downloadPath))
		return downloadPath, nil
	}

	logger.Info(fmt.Sprintf("Downloading %s to %s", url, downloadPath))

	res, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		// there was an error of some kind when downloading,
		// so 100% do not try to write a file
		return "", fmt.Errorf("%s returned a non-okay status code of %d", url, res.StatusCode)
	}

	f, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, res.Body)

	if err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("Downloaded %s", url))

	return downloadPath, nil
}
