package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/benduran/glipglop/cache"
	"github.com/benduran/glipglop/downloader"
	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
	"github.com/benduran/glipglop/schema"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	DisableFlagParsing: true,
	Use:                "exec",
	Short:              "Runs a command against a specific language or tool",
	Long: `Want to run a Node.js or Python script, start a Java application, or something similar
against your project's specific language or tool requirement? Use exec`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := strings.TrimSpace(args[0])
		argsForTool := args[1:]
		cwd, _ := internal.GetCWD()
		logger.Info(fmt.Sprintf("Executing a command %s with args %s in %s", tool, argsForTool, cwd))

		manifest, err := schema.ReadUserSchema(cwd)

		if err != nil {
			logger.Error(err)
			return
		}

		// the downloader will only download the missing bits
		if err := downloader.DownloadAllTools(cwd); err != nil {
			logger.Error(err)
			return
		}

		// now we need to get the binary paths to all the tools in the cache
		// and apply them as a path variable to a child spawned shell
		path := ""
		for toolName, toolVersion := range manifest.Tools {
			binaryPath := filepath.Dir(cache.CheckBinaryInToolCache(toolName, toolVersion))
			path += fmt.Sprintf(":%s", binaryPath)
		}
		// Remove leading colon if path is not empty
		if path != "" {
			path = path[1:]
		}

		existingPath := os.Getenv("PATH")

		childCmd := exec.Command(tool, argsForTool...)
		childCmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s:%s", path, existingPath))
		childCmd.Stdin = os.Stdin
		childCmd.Stdout = os.Stdout
		childCmd.Stderr = os.Stderr

		childEnv := childCmd.Environ()
		sort.Strings(childEnv)
		// uncomment this line to see the PATH set for the child process
		// logger.Info(fmt.Sprintf("applying the following as $PATH: %s", strings.Join(childEnv, "\n")))

		childErr := childCmd.Start()

		if childErr != nil {
			logger.Error(childErr)
			return
		}

		childErr = childCmd.Wait()
		if childErr != nil {
			logger.Error(childErr)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
