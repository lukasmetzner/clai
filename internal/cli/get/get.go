package get

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

var jobID int32

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get information about a job and the output",
	Run: func(cmd *cobra.Command, args []string) {
		if jobID == -1 {
			fmt.Println("JobID not specified")
			os.Exit(1)
		}

		urlStr := fmt.Sprintf("%s/%d", "http://localhost:8080/api/jobs", jobID)

		res, err := http.Get(urlStr)

		if err != nil {
			fmt.Printf("Error fetching job information: %s", err)
			os.Exit(1)
		}

		if res.StatusCode != 200 {
			fmt.Printf("Status code is %s", res.Status)
			os.Exit(1)
		}

		buffer := bytes.Buffer{}
		buffer.ReadFrom(res.Body)

		var job models.Job

		if err := json.Unmarshal(buffer.Bytes(), &job); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		if job.JobOutput.Stdout == "" && job.JobOutput.Stderr == "" {
			pterm.DefaultBasicText.Println(pterm.Red("Stdout and Stderr are not yet available"))
			return
		}

		stdout := pterm.LightGreen("Stdout")
		stderr := pterm.LightRed("Stderr")
		pterm.DefaultBasicText.Println(stdout)
		pterm.DefaultBasicText.Print(job.JobOutput.Stdout)
		pterm.DefaultBasicText.Println(stderr)
		pterm.DefaultBasicText.Print(job.JobOutput.Stderr)
	},
}

func init() {
	GetCmd.Flags().Int32VarP(&jobID, "jobid", "i", -1, "Specify jobid")
}
