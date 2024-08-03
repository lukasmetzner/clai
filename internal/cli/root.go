package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clai-cli",
	Short: "CLI tool to interact with the Clai System",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Hello from Clai")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}
