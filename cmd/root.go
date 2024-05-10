package cmd

import (
	"os"

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

		logLevel, _ := cmd.PersistentFlags().GetString("log-level")

		if len(logLevel) == 0 {
			logLevel = "info"
		}

		os.Setenv("GLIPGLOP_LEVEL", logLevel)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// User didn't specify a command, so we'll show the help menu instead
		cmd.Help()
	},
}

func SetupCLI() {
	rootCmd.PersistentFlags().String("log-level", "Determines the log verbosity.", "info")
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
	}
}
