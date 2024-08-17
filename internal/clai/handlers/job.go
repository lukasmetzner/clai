package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukasmetzner/clai/pkg/database"
	"github.com/lukasmetzner/clai/pkg/models"
	"github.com/lukasmetzner/clai/pkg/mq"
	"github.com/rabbitmq/amqp091-go"
)

func CreateJob(c *gin.Context) {
	var job models.Job
	if err := json.NewDecoder(c.Request.Body).Decode(&job); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&job).Error; err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		log.Printf("Database error: %s", err)
		return
	}

	bytes, err := json.Marshal(job)
	if err != nil {
		log.Printf("%s", err)
	}

	if err := mq.Channel.Publish(
		"",
		mq.Queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusCreated, "application/json", bytes)
}

func GetJobs(c *gin.Context) {
	var jobs []models.Job
	if err := database.DB.Find(&jobs).Error; err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func GetJob(c *gin.Context) {
	var job models.Job
	jobId, exists := c.Params.Get("id")

	if !exists {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := database.DB.First(&job, jobId).Error; err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, job)
}

func UpdateJob(c *gin.Context) {
	var job models.Job
	jobId, exists := c.Params.Get("id")

	if !exists {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := database.DB.First(&job, jobId).Error; err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&job); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Save(&job).Error; err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, job)
}

func DeleteJob(c *gin.Context) {
	jobId, exists := c.Params.Get("id")

	if !exists {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := database.DB.Delete(&models.Job{}, jobId).Error; err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
