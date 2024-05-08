package internal

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/benduran/glipglop/cache"
	logger "github.com/benduran/glipglop/log"
	"github.com/xi2/xz"
)

// extracts an arbitrary archive to a folder on disk
func ExtractArchive(archivePath string) (string, error) {
	logger.Info(fmt.Sprintf("Extracting %s", archivePath))
	// Check if the archive path exists
	_, err := os.Stat(archivePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("archive file %s does not exist", archivePath)
	}

	extractDir, err := cache.GetExtractedCacheLocation()
	if err != nil {
		return "", err
	}

	// Determine the archive type based on the file extension
	endsWithTar, _ := regexp.Compile(`\.(tar)(\.(gz|xz))$`)
	endsWithZip, _ := regexp.Compile(`\.zip$`)

	isTar := endsWithTar.MatchString(archivePath)
	isZip := endsWithZip.MatchString(archivePath)

	if isTar {
		if filepath.Ext(archivePath) == ".gz" {
			err = extractTarGz(archivePath, extractDir)
		} else if filepath.Ext(archivePath) == ".xz" {
			err = extractTarXz(archivePath, extractDir)
		} else {
			err = extractTar(archivePath, extractDir)
		}
	} else if isZip {
		err = extractZip(archivePath, extractDir)
	} else {
		return "", fmt.Errorf("unable to extract %s because it is neither a tar nor a zip archive", archivePath)
	}

	if err != nil {
		return "", err
	}

	sansExt := endsWithTar.ReplaceAllString(filepath.Base(archivePath), "")
	sansExt = endsWithZip.ReplaceAllString(sansExt, "")

	logger.Info(fmt.Sprintf("Successfully extracted %s to %s", archivePath, extractDir))

	return filepath.Join(extractDir, sansExt), nil
}

func extractTar(archivePath, targetPath string) error {
	logger.Info(fmt.Sprintf("detected %s to be a regular tar archive", archivePath))
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			err = os.MkdirAll(targetFilePath, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// Create file
			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Write file contents
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractTarXz(archivePath, targetPath string) error {
	logger.Info(fmt.Sprintf("detected %s to be a tar-xz archive", archivePath))

	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	r, err := xz.NewReader(file, 0)
	if err != nil {
		return err
	}

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			err = os.MkdirAll(targetFilePath, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// Create file
			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Write file contents
			_, err = io.Copy(outFile, tr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractTarGz(archivePath, targetPath string) error {
	logger.Info(fmt.Sprintf("detected %s to be a tar-gz archive", archivePath))
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			err = os.MkdirAll(targetFilePath, os.ModePerm)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// Create file
			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Write file contents
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func extractZip(archivePath, targetPath string) error {
	logger.Info(fmt.Sprintf("detected %s to be a zip archive", archivePath))
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		targetFilePath := filepath.Join(targetPath, file.Name)

		if file.FileInfo().IsDir() {
			// Create directory
			err = os.MkdirAll(targetFilePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// Create file
		outFile, err := os.Create(targetFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Open file in zip archive
		inFile, err := file.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		// Write file contents
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
	}

	return nil
}
