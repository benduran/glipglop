package internal

import (
	"runtime"
)

type MachineInfo struct {
	Arch        string
	OS          string
	Unsupported bool
}

// Determines information about the
// machine where glipglop is currently
// executing
func GetMachineInfo() *MachineInfo {
	detectedOS := runtime.GOOS
	detectedArch := runtime.GOARCH

	unsupported := false

	if detectedOS != "darwin" && detectedOS != "linux" && detectedOS != "windows" {
		unsupported = true
	}

	if detectedArch != "arm64" && detectedArch != "amd64" {
		unsupported = true
	}

	return &MachineInfo{
		Arch:        detectedArch,
		OS:          detectedOS,
		Unsupported: unsupported,
	}
}
