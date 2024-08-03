package cli

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clai-cli",
	Short: "CLI tool to interact with the Clai System",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Hello from Clai")
		command := args[0]
		command_args := args[1:]
		cmdex := exec.Command(command, command_args...)

		stdout, err := cmdex.Output()

		if err != nil {
			log.Fatalf("%s", err)
		}

		log.Printf("%s", string(stdout))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}
