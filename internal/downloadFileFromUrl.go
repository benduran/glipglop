package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
	logger "github.com/benduran/glipglop/log"
	"github.com/schollz/progressbar/v3"
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
		logger.Debug(fmt.Sprintf("%s already exists. Skipping download!", downloadPath))
		return downloadPath, nil
	}

	logger.Debug(fmt.Sprintf("Downloading %s to %s", url, downloadPath))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// there was an error of some kind when downloading,
		// so 100% do not try to write a file
		return "", fmt.Errorf("%s returned a non-okay status code of %d", url, resp.StatusCode)
	}

	f, err := os.OpenFile(downloadPath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return "", err
	}

	defer f.Close()

	bar := progressbar.DefaultBytes(resp.ContentLength, fmt.Sprintf("downloading %s", url))

	io.Copy(io.MultiWriter(f, bar), resp.Body)

	logger.Debug(fmt.Sprintf("Downloaded %s to %s", url, downloadPath))

	return downloadPath, nil
}
