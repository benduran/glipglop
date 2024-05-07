package internal

import (
	"os"
)

// gets the CWD for the process
func GetCWD() (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return cwd, nil
}
