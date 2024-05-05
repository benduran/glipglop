package cmd

import (
	"fmt"

	logger "github.com/benduran/glipglop/log"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	DisableFlagParsing: true,
	Use:                "exec",
	Short:              "Runs a command against a specific language or tool",
	Long: `Want to run a Node.js or Python script, start a Java application, or something similar
against your project's specific language or tool requirement? Use exec`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		argsForTool := args[1:]
		logger.Info(fmt.Sprintf("Executing a command %s with args %s", tool, argsForTool))
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
