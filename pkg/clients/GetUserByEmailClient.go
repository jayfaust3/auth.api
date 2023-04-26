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

type getUserByEmailRequest struct {
	Email string `json:"email"`
}

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

	msgs, err := ch.Consume(
		rabbitQueue, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// var request getUserByEmailRequest
	// request.Email = email
	// var requestMessage rabitMessage[getUserByEmailRequest]
	// requestMessage.Data = request
	var request getUserByEmailRequest
	request.Email = email
	var requestMessage messaging.Message[getUserByEmailRequest]
	requestMessage.Data = request

	encodedMessage, err := json.Marshal(requestMessage)

	if err == nil {
		log.Printf("publishing message: %s", string(encodedMessage))

		err = ch.PublishWithContext(ctx,
			rabbitExchange, // exchange
			rabbitQueue,    // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: "",
				ReplyTo:       rabbitQueue,
				Body:          []byte(string(encodedMessage)),
			})

		failOnError(err, "Failed to publish a message")

		for d := range msgs {
			log.Printf("processing message")
			messageDataBytes := d.Body
			log.Printf("message data: %s", string(messageDataBytes))
			var messageData messaging.Message[user.User]
			err := json.Unmarshal(messageDataBytes, &messageData)

			failOnError(err, "Failed to extract user from message")
			res = messageData.Data
			break
		}
	} else {
		failOnError(err, "Failed to encode message")
	}

	return
}
