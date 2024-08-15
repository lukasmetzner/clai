package cli

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lukasmetzner/clai/pkg/models"
	"github.com/spf13/cobra"
)

var scriptPath, reqPath string

var rootCmd = &cobra.Command{
	Use:   "clai-cli",
	Short: "CLI tool to interact with the Clai System",
	Run: func(cmd *cobra.Command, args []string) {
		scriptFile, err := os.ReadFile(scriptPath)
		if err != nil {
			log.Fatalf("%s error with the file %s", err, scriptPath)
		}

		reqFile, err := os.ReadFile(reqPath)
		if err != nil {
			log.Fatalf("%s error with the file %s", err, reqPath)
		}

		job := models.JobRequest{
			Script:       string(scriptFile),
			Requirements: string(reqFile),
			Type:         models.ScriptJob,
		}

		body, _ := json.Marshal(job)
		reader := bytes.NewReader(body)
		res, err := http.Post("http://localhost:8080/api/jobs/", "application/json", reader)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if res.StatusCode != 201 {
			log.Fatalf("Failed! Response code is %s", res.Status)
		}

		log.Printf("Script queued for execution")
	},
}

func init() {
	rootCmd.Flags().StringVarP(&scriptPath, "script", "s", "", "Path to the script file")
	rootCmd.Flags().StringVarP(&reqPath, "requirements", "r", "", "Path to the requirements file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}
