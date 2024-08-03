package claiagent

import (
	"log"

	"github.com/lukasmetzner/clai/pkg/mq"
)

func Start() {
	log.Println("Starting clai-agent...")

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

	go func() {
		for msg := range msgs {
			log.Printf("Received a new job %s", msg.Body)
		}
	}()

	var forever chan struct{}

	log.Printf("Starting to consume messages")

	<-forever
}
