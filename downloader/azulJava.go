package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/benduran/glipglop/cache"
	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
)

type JavaVersionInfo struct {
	AvailabilityType   string  `json:"availability_type"`
	DistroVersion      []int64 `json:"distro_version"`
	DownloadURL        string  `json:"download_url"`
	JavaVersion        []int64 `json:"java_version"`
	Latest             bool    `json:"latest"`
	Name               string  `json:"name"`
	OpenjdkBuildNumber int64   `json:"openjdk_build_number"`
	PackageUUID        string  `json:"package_uuid"`
	Product            string  `json:"product"`
}

// fetches java package information from the official Azul
// metadata API
func fetchJavaVersionInfo(version string, machineInfo *internal.MachineInfo) ([]JavaVersionInfo, error) {
	arch := ""

	fmt.Println(machineInfo)

	if machineInfo.Arch == "x64" {
		arch = "amd64"
	} else {
		arch = "aarch64"
	}

	os := ""

	if machineInfo.OS == "darwin" {
		os = "macos"
	} else {
		os = machineInfo.OS
	}

	javaUrl := fmt.Sprintf("https://api.azul.com/metadata/v1/zulu/packages?java_version=%s&os=%s&arch=%s&archive_type=zip&java_package_type=jdk&javafx_bundled=false&latest=true&release_status=ga&availability_types=CA&certifications=tck", version, os, arch)
	logger.Info(fmt.Sprintf("fetching java from %s", javaUrl))
	response, err := http.Get(javaUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var info []JavaVersionInfo
	err = json.Unmarshal(body, &info)

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%v", info))

	return info, nil
}

// Downloads an Azul JDK
func DownloadAzulJava(version string) (string, error) {
	// given the java version, we need to query the Azul metadata API
	// to get the latest distro for this version number, then download that
	// via the provided URL in the response
	// https://cdn.azul.com/zulu/bin/zulu22.30.13-ca-jdk22.0.1-macosx_aarch64.zip
	// https://cdn.azul.com/zulu/bin/zulu17.50.19-ca-jdk17.0.11-linux_x64.zip

	// if the tool is already in the cache, just return the path to that immediately
	existingPathToTool := cache.CheckBinaryInToolCache("azul-java", version)

	if len(existingPathToTool) > 0 {
		return existingPathToTool, nil
	}

	machineInfo := internal.GetMachineInfo()

	if machineInfo.Unsupported {
		return "", fmt.Errorf("unable to download azul-java because you are running an unsupported os + arch")
	}

	infos, err := fetchJavaVersionInfo(version, machineInfo)

	if err != nil {
		return "", err
	}

	if len(infos) == 0 {
		return "", fmt.Errorf("unable to download azul-java because no version matches your os + arch combo of %s:%s", machineInfo.OS, machineInfo.Arch)
	}

	// the first entry is the only one we want
	info := infos[0]

	downloadPath, err := internal.DownloadFileFromURL(info.DownloadURL)

	if err != nil {
		return "", err
	}

	// now we need to extract the archive
	extractedPath, err := internal.ExtractArchive("azul-java", version, downloadPath)

	if err != nil {
		return "", err
	}

	// the bin folder needs to be moved and placed into the glipglop folder.
	// we will determine the right folder by trying to find the javac binary
	var javacBinary string
	err = filepath.Walk(extractedPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "javac" {
			javacBinary = path
			return filepath.SkipDir // Stop walking the directory tree
		}
		return nil
	})

	logger.Info("FART")
	logger.Info(javacBinary)
	logger.Info("FART")

	if err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("Found the azul-javac binary to be %s", javacBinary))

	allBinaries, err := filepath.Glob(filepath.Join(filepath.Dir(javacBinary), "*"))

	if err != nil {
		return "", err
	}

	fuck := filepath.Join(extractedPath, "azul-java")

	if _, err := os.Stat(fuck); err != nil {
		return "", err
	}

	return internal.MoveFolderToToolCache("azul-java", version, fuck, allBinaries)
}
