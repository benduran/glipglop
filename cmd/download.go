package cmd

import (
	"fmt"
	"os"

	"github.com/benduran/glipglop/downloader"
	"github.com/benduran/glipglop/schema"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Short: "Downloads all of the tools you have specified for your project",
	Use:   "download",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := rootCmd.Flags().GetString("cwd")
		fmt.Println("Downloading all of your tools now...")
		schema, err := schema.ReadUserSchema(cwd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for key, val := range schema.Tools {
			downloader.DownloadTool(key, val)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
