package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukasmetzner/clai/pkg/database"
	"github.com/lukasmetzner/clai/pkg/models"
	"github.com/lukasmetzner/clai/pkg/mq"
	"github.com/rabbitmq/amqp091-go"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&job).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(job)

	if err != nil {
		log.Fatalf("%s", err)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Job
	if err := database.DB.Find(&jobs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var job models.Job
	if err := database.DB.First(&job, params["id"]).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var job models.Job
	if err := database.DB.First(&job, params["id"]).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Save(&job).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := database.DB.Delete(&models.Job{}, params["id"]).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
