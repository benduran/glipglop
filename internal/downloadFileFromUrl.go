package internal

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
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

	res, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	f, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, res.Body)

	if err != nil {
		return "", err
	}

	return downloadPath, nil
}
