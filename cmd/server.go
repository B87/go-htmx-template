package cmd

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

// runServerCmd represents the runServer command
var runServerCmd = &cobra.Command{
	Use:   "runServer",
	Short: "Run HTTP server",
	Long:  `Run HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		httpServer.Serve()
	},
}

var uploadStaticCmd = &cobra.Command{
	Use:   "uploadStatic",
	Short: "Upload static files",
	Long:  `Upload static files`,
	Run: func(cmd *cobra.Command, args []string) {
		httpServer.UploadStaticFiles()
	},
}

func init() {
	RootCmd.AddCommand(uploadStaticCmd)
	RootCmd.AddCommand(runServerCmd)
}
