package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GlipGlopManifest struct {
	Tools map[string]string `json:"tools"`
}

func ReadUserSchema(cwd string) (*GlipGlopManifest, error) {
	// walk up to the nearest schema and use that.
	// if it doesn't exist, return an error

	splitPath := strings.Split(cwd, string(os.PathSeparator))

	scannedPaths := []string{}

	for i := len(splitPath); i >= 0; i-- {
		segment := strings.Join(splitPath[0:i], string(os.PathSeparator))
		segmentWithGlipglop := filepath.Join(segment, ".glipglop.json")

		if string(segmentWithGlipglop[0]) != string(os.PathSeparator) {
			segmentWithGlipglop = fmt.Sprintf("/%s", segmentWithGlipglop)
		}

		scannedPaths = append(scannedPaths, segmentWithGlipglop)

		stat, err := os.Stat(segmentWithGlipglop)
		if os.IsNotExist(err) || stat.IsDir() {
			// the file doesn't exist at this point,
			// or the user messes up and created a folder
			// with the same name as the file
			continue
		}
		// if we got here, the file exists, so we should return it
		var schema GlipGlopManifest

		fileContents, err := os.ReadFile(segmentWithGlipglop)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(fileContents, &schema)

		if err != nil {
			return nil, err
		}

		return &schema, nil
	}

	// TODO: this EOL character is not portable or friendly for Windows,
	// so we'll get back to this later
	return nil, fmt.Errorf("unable to read your glipglop schema because no glipglop.json file was found when scanning the following paths:\n%s", strings.Join(scannedPaths, "\n"))
}
