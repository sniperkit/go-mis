package services

import (
	"github.com/streadway/amqp"
	"fmt"
	"encoding/json"
)

type QueueServiceShape struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
}

var QueueService = QueueServiceShape{}

func InitQueueService(conn *amqp.Connection, channel *amqp.Channel) {
	QueueService.Connection = conn
	QueueService.Channel = channel
}

func (qs *QueueServiceShape) PublishQueue(queueName string, obj interface{}) error {
	queue, err := qs.Channel.QueueDeclare(queueName, true, true, false, false, nil)
	if err != nil {
		fmt.Printf("Error declaring queue%+v\n", err)
		return err
	}

	item, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error marshalling obj%+v\n", err)
		return err
	}

	err = qs.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(string(item)),
	})
	if err != nil {
		fmt.Printf("Error publishing queue%+v\n", err)
		return err
	}

	fmt.Printf("Queue published%+v\n", obj)

	return nil
}