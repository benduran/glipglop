package cmd

import (
	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glipglop",
	Short: "Managed all of your programming languages in a single place",
	Long: `Tired of using nvm, pynev, jabba, sdkman, and all those other language-specific tools for managing your various, installed tool versions?
Glipglop brings this all under a single roof and allows you to use them with ease.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// if we can't determine the CWD, glipglop is totally boned, so just abort early
		_, err := internal.GetCWD()

		if err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// User didn't specify a command, so we'll show the help menu instead
		cmd.Help()
	},
}

func SetupCLI() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
	}
}
