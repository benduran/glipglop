package cmd

import (
	"os"
	"path/filepath"

	logger "github.com/benduran/glipglop/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glipglop",
	Short: "Managed all of your programming languages in a single place",
	Long: `Tired of using nvm, pynev, jabba, sdkman, and all those other language-specific tools for managing your various, installed tool versions?
Glipglop brings this all under a single roof and allows you to use them with ease.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// transform the CWD flag, regardless of what it is, to an absolute path
		cwd, _ := cmd.Flags().GetString("cwd")

		cwd, _ = filepath.Abs(cwd)

		cmd.Flags().Set("cwd", cwd)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// User didn't specify a command, so we'll show the help menu instead
		cmd.Help()
	},
}

func SetupCLI() {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error(err)
	}

	rootCmd.PersistentFlags().String("cwd", cwd, "--cwd <path/to/dir>")

	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
	}
}
