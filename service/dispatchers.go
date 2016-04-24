package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	"github.com/streadway/amqp"
)

// AmqpDispatcher is used as anchor for dispatch messsage method for real
// AMQP channels
type AmqpDispatcher struct {
	channel       queuePublishableChannel
	queueName     string
	mandatorySend bool
}

// NewAMQPDispatcher returns a new AMQP dispatcher wrapped around a single
// publishing channel.
func NewAMQPDispatcher(publishChannel queuePublishableChannel, name string, mandatory bool) *AmqpDispatcher {
	return &AmqpDispatcher{channel: publishChannel, queueName: name, mandatorySend: mandatory}
}

type queuePublishableChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// DispatchMessage implementation of dispatch message interface method
func (q *AmqpDispatcher) DispatchMessage(message interface{}) (err error) {
	fmt.Printf("Dispatching message to queue %s\n", q.queueName)
	body, err := json.Marshal(message)
	if err == nil {
		err = q.channel.Publish(
			"",              // exchange
			q.queueName,     // routing key
			q.mandatorySend, // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			fmt.Printf("Failed to dispatch message: %s\n", err)
		}
	} else {
		fmt.Printf("Failed to marshal message %v (%s)\n", message, err)
	}
	return
}

// BuildDispatcher builds a new QueueDispatcher
func BuildDispatcher(queueName string, appEnv *cfenv.App) QueueDispatcher {
	url := resolveAMQPURL(appEnv)
	return createAMQPDispatcher(queueName, url)
}

func createAMQPDispatcher(queueName string, url string) QueueDispatcher {
	fmt.Printf("\nUsing URL (%s) for Rabbit.\n", url)

	conn, err := amqp.Dial(url)
	log.Fatalf("%s: %s", err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	log.Fatalf("%s: %s", err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	log.Fatalf("%s: %s", err, "Failed to declare a queue")
	dispatcher := NewAMQPDispatcher(ch, q.Name, false)
	return dispatcher
}

func resolveAMQPURL(appEnv *cfenv.App) string {
	url, err := cftools.GetVCAPServiceProperty("rabbit", "uri", appEnv)
	if err != nil {
		fmt.Println("Failed to detect bound service for rabbit. Falling back to in-memory dispatcher (fake)")
		return "fake://foo"
	}
	if len(url) < 10 {
		fmt.Printf("URL detected for bound rabbit service not valid, was '%s'. Falling back to in-memory fake.\n", url)
		return "fake://foo"
	}
	return url
}

// Fake
type fakeQueueDispatcher struct {
	Messages []interface{}
}

func newFakeQueueDispatcher() (dispatcher *fakeQueueDispatcher) {
	dispatcher = &fakeQueueDispatcher{}
	dispatcher.Messages = make([]interface{}, 0)
	return
}

// DispatchMessage implementation of dispatch message interface method
func (q *fakeQueueDispatcher) DispatchMessage(message interface{}) (err error) {
	q.Messages = append(q.Messages, message)
	return
}
