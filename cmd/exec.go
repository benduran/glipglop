package cmd

import (
	"fmt"
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
		fmt.Println(args)
		tool := strings.TrimSpace(args[0])
		argsForTool := args[1:]
		cwd, _ := internal.GetCWD()
		logger.Info(fmt.Sprintf("Executing a command %s with args %s in %s", tool, argsForTool, cwd))

		// the downloader will only download the missing bits
		err := downloader.DownloadAllTools(cwd)

		if err != nil {
			logger.Error(err)
			return
		}

		manifest, err := schema.ReadUserSchema(cwd)

		if err != nil {
			logger.Error(err)
			return
		}

		toolManifestEntry := manifest.Tools[tool]
		fmt.Println(toolManifestEntry)

		toolBinaryPath := cache.CheckBinaryInToolCache(tool, toolManifestEntry)

		if len(toolBinaryPath) == 0 {
			logger.Error(fmt.Errorf("unable to use %s because you don't have it declared in your glipglop.json manifest", tool))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
