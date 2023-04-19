package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jayfaust3/auth.api/pkg/models/application/user"
	"github.com/jayfaust3/auth.api/pkg/models/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetUserFromEmail(email string) (res user.User, err error) {
	rabbitUserName := os.Getenv("RABBITMQ_USERNAME")
	rabbitPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUserName, rabbitPassword, rabbitHost, rabbitPort)

	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	rabbitExchange := "GetUserByEmail"
	rabbitQueue := "GetUserByEmail"

	q, err := ch.QueueDeclare(
		rabbitQueue, // name
		false,       // durable
		false,       // delete when unused
		true,        // exclusive
		false,       // noWait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response := ""
	err = ch.PublishWithContext(ctx,
		rabbitExchange, // exchange
		rabbitQueue,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: "",
			ReplyTo:       q.Name,
			Body:          []byte(response),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		messageDataBytes := d.Body
		// messageDataJSON := string(messageDataBytes)
		var messageData messaging.Message[user.User]
		err := json.Unmarshal(messageDataBytes, &messageData)

		failOnError(err, "Failed to extract user from message")
		res = messageData.Data
		break
	}

	return
}
