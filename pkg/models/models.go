package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type JobType string

const (
	ScriptJob JobType = "script"
	RepoJob   JobType = "repo"
)

type Job struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
	Priority        int       `json:"priority"`
	ResourceRequest string    `json:"resourceRequest"`
	Type            JobType   `json:"type"`
	Script          string    `json:"script" gorm:"type:text"`
	Requirements    string    `json:"requirements" gorm:"type:text"`
	RepoURL         string    `json:"repoUrl"`
	JobOutputID     *uint     `json:"job_output_id"`
	JobOutput       JobOutput `json:"job_output"`
}

type JobOutput struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type JobRequest struct {
	Script       string  `json:"script"`
	Requirements string  `json:"requirements"`
	Type         JobType `json:"type"`
}

func (job *Job) BeforeSave(tx *gorm.DB) (err error) {
	if job.Type == ScriptJob && job.Script == "" {
		return fmt.Errorf("script field cannot be empty for script job")
	}
	if job.Type == RepoJob && job.RepoURL == "" {
		return fmt.Errorf("repo_url field cannot be empty for repo job")
	}
	return nil
}
