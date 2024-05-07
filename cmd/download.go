package cmd

import (
	"github.com/benduran/glipglop/downloader"
	logger "github.com/benduran/glipglop/log"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Short: "Downloads all of the tools you have specified for your project",
	Use:   "download",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := rootCmd.Flags().GetString("cwd")
		if err := downloader.DownloadAllTools(cwd); err != nil {
			logger.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
