package clients

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetUserFromEmail(email string) (res int, err error) {
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

	for d := range msgs {
		messageBytes := d.Body

		response := ""
		err = ch.PublishWithContext(ctx,
			rabbitExchange, // exchange
			d.RoutingKey,   // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				ReplyTo:       q.Name,
				Body:          []byte(response),
			})
		failOnError(err, "Failed to publish a message")
	}

	return
}
