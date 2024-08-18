package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lukasmetzner/clai/pkg/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var scriptPath, reqPath string

var RunCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a Python job",
	Run: func(cmd *cobra.Command, args []string) {
		if scriptPath == "" {
			fmt.Printf("Provide a script path via --script or -s")
			os.Exit(1)
		}

		scriptFile, err := os.ReadFile(scriptPath)
		if err != nil {
			fmt.Printf("%s error with the file %s", err, scriptPath)
			os.Exit(1)
		}

		reqStr := ""

		if reqPath != "" {
			reqFile, err := os.ReadFile(reqPath)
			if err != nil {
				fmt.Printf("%s error with the file %s", err, reqPath)
				os.Exit(1)
			}
			reqStr = string(reqFile)
		} else {
			fmt.Println("You did not provide a requirements.txt path!")
		}

		job := models.JobRequest{
			Script:       string(scriptFile),
			Requirements: reqStr,
			Type:         models.ScriptJob,
		}

		body, _ := json.Marshal(job)
		reader := bytes.NewReader(body)
		res, err := http.Post("http://localhost:8080/api/jobs/", "application/json", reader)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		if res.StatusCode != 201 {
			fmt.Printf("Failed! Response code is %s", res.Status)
			os.Exit(1)
		}

		buffer := bytes.Buffer{}
		buffer.ReadFrom(res.Body)

		var resJob models.Job
		if err := json.Unmarshal(buffer.Bytes(), &resJob); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		jobIDStr := fmt.Sprintf("%d", resJob.ID)

		pterm.DefaultBasicText.Println("Your job has been scheduled with the ID: " + pterm.LightMagenta(jobIDStr))
	},
}

func init() {
	RunCmd.Flags().StringVarP(&scriptPath, "script", "s", "", "Path to the script file")
	RunCmd.Flags().StringVarP(&reqPath, "requirements", "r", "", "Path to the requirements file")
}
