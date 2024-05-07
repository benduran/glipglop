package cmd

import (
	"github.com/benduran/glipglop/downloader"
	"github.com/benduran/glipglop/internal"
	logger "github.com/benduran/glipglop/log"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Short: "Downloads all of the tools you have specified for your project",
	Use:   "download",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := internal.GetCWD()
		if err := downloader.DownloadAllTools(cwd); err != nil {
			logger.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
