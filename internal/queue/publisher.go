package queue

import (
    "encoding/json"
    "scalable-notification/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)


func PublishNotification(notif models.Notification) error {
	body, err := json.Marshal(notif)
	if err != nil{
		return err
	}
	
	priority := uint8(notif.Priority);

	err = serverCh.Publish(
		"notifications_exchange",
		notif.Type,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
            Body:        body,
            Priority:    priority,
		},
	)
	return err;
}

func CloseRabbitMQ() {
    if serverCh != nil {
        serverCh.Close()
    }
    if serverConn != nil {
        serverConn.Close()
    }
}