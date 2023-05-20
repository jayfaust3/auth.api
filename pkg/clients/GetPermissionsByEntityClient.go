package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jayfaust3/auth.api/pkg/models/application/permission"
	"github.com/jayfaust3/auth.api/pkg/models/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type getPermissionsByEntityRequest struct {
	ActorType int    `json:"actorType"`
	EntityId  string `json:"entityId"`
}

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

func GetPermissionsByEntity(entityId string, actorType int) (res []permission.Scope, err error) {
	if actorType != 1 {
		actorType = 0
	}

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

	rabbitExchange := "exchange:rpc"
	rabbitQueue := "queue:get-permissions-by-entity"

	replyToQueue, err := ch.QueueDeclare(
		fmt.Sprintf("%s-reply-to", rabbitQueue), // name
		false,                                   // durable
		false,                                   // delete when unused
		true,                                    // exclusive
		false,                                   // noWait
		nil,                                     // arguments
	)
	replyToQueueName := replyToQueue.Name

	ch.QueueBind(replyToQueueName, replyToQueueName, replyToQueueName, false, nil)

	msgs, err := ch.Consume(
		replyToQueue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	requestMessage := messaging.Message[getPermissionsByEntityRequest]{
		Data: getPermissionsByEntityRequest{
			ActorType: actorType,
			EntityId:  entityId,
		},
	}

	encodedMessage, err := json.Marshal(requestMessage)

	if err == nil {
		log.Printf("publishing message: %s", string(encodedMessage))

		err = ch.PublishWithContext(ctx,
			rabbitExchange, // exchange
			rabbitQueue,    // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: corrId,
				DeliveryMode:  amqp.Transient,
				MessageId:     uuid.New().String(),
				Timestamp:     time.Now(),
				ReplyTo:       replyToQueueName,
				Body:          []byte(string(encodedMessage)),
			})

		failOnError(err, "Failed to publish a message")

		for msg := range msgs {
			log.Printf("processing message")

			if msg.CorrelationId == corrId {
				log.Printf("correlation id matches")

				messageDataBytes := msg.Body

				var messageData messaging.Message[[]permission.Scope]
				err := json.Unmarshal(messageDataBytes, &messageData)

				failOnError(err, "Failed to extract permissions from message")
				res = messageData.Data
				break
			}
		}
	} else {
		failOnError(err, "Failed to encode message")
	}

	return
}
