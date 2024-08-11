package claiagent

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/lukasmetzner/clai/pkg/k8s"
	"github.com/lukasmetzner/clai/pkg/models"
	"github.com/lukasmetzner/clai/pkg/mq"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Start() {
	log.Println("Starting clai-agent...")

	// Init Kubernetes
	k8s.Init()

	// Init RabbitMQ
	mq.InitMQ()

	msgs, err := mq.Channel.Consume(
		mq.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("%s", err)
	}

	for msg := range msgs {
		log.Printf("Received a new job %s", msg.Body)
		job := models.Job{}
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("%s", err)
			continue
		}

		go startJob(job)
	}
}

func startJob(job models.Job) {
	idStr := strconv.FormatUint(uint64(job.ID), 10)
	k8sJob := k8s.GetJob(idStr)

	result, err := k8s.JobClient.Create(context.TODO(), &k8sJob, metav1.CreateOptions{})
	if err != nil {
		log.Printf("%s", err)
	}

	log.Printf("Created job: %s\n", result.GetObjectMeta().GetName())
}
