package mq

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Channel *amqp.Channel
var Queue amqp.Queue

func InitMQ() {

	mqUser := os.Getenv("MQ_USER")
	mqPassword := os.Getenv("MQ_PASSWORD")
	mqHost := os.Getenv("MQ_HOST")
	mqPort := os.Getenv("MQ_PORT")

	mqcs := fmt.Sprintf("amqp://%s:%s@%s:%s/", mqUser, mqPassword, mqHost, mqPort)

	conn, err := amqp.Dial(mqcs)

	if err != nil {
		log.Fatalf("%s", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s", err)
	}

	q, err := ch.QueueDeclare("clai", true, false, false, false, nil)

	if err != nil {
		log.Fatalf("%s", err)
	}

	Channel = ch
	Queue = q
}
