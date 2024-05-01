package cmd

import (
	"fmt"
	"os"

	"github.com/benduran/glipglop/schema"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Short: "Downloads all of the tools you have specified for your project",
	Use:   "download",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Downloading all of your tools now...")
		_, err := schema.ReadUserSchema()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
