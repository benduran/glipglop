package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glipglop",
	Short: "Managed all of your programming languages in a single place",
	Long: `Tired of using nvm, pynev, jabba, sdkman, and all those other language-specific tools for managing your various, installed tool versions?
Glipglop brings this all under a single roof and allows you to use them with ease.`,
	Run: func(cmd *cobra.Command, args []string) {
		// User didn't specify a command, so we'll show the help menu instead
		cmd.Help()
	},
}

func SetupCLI() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
