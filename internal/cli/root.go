package cli

import (
	"log"
	"os"

	"github.com/lukasmetzner/clai/internal/cli/get"
	"github.com/lukasmetzner/clai/internal/cli/run"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clai-cli",
	Short: "CLI tool to interact with the Clai System",
}

func init() {
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(get.GetCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}
