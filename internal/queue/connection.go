package queue

import (
    amqp "github.com/rabbitmq/amqp091-go"
)

var serverConn *amqp.Connection
var serverCh *amqp.Channel
var workerConn *amqp.Connection
var workerCh *amqp.Channel

func ConnectServer() error {
    var err error
    serverConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return err
    }

    serverCh, err = serverConn.Channel()
    if err != nil {
        return err
    }
    
	_,err = serverCh.QueueDeclare(
		"notifications_exchange", // <- same name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
        return err
    }

	return nil
}

func ConnectWorker() error {
    var err error
    workerConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return err
    }

    workerCh, err = workerConn.Channel()
    if err != nil {
        return err
    }

	_, err = workerCh.QueueDeclare(
        "notifications_exchange_retry",
        true,  // durable
        false, // auto-delete
        false, // exclusive
        false, // no-wait
        amqp.Table{
            "x-message-ttl": 30000, // 30 seconds retry delay
            "x-dead-letter-exchange": "",
            "x-dead-letter-routing-key": "notifications_exchange", // Route back to main queue
        },
    )
    if err != nil {
        return err
    }
    
    // Declare dead letter queue
    _, err = workerCh.QueueDeclare(
        "notifications_exchange_dead",
        true,  // durable
        false, // auto-delete
        false, // exclusive
        false, // no-wait
        nil,   // no special arguments
    )
	if err != nil {
		return err
	}

    return nil
}

func CloseServer() {
    if serverCh != nil {
        serverCh.Close()
    }
    if serverConn != nil {
        serverConn.Close()
    }
}

func CloseWorker() {
    if workerCh != nil {
        workerCh.Close()
    }
    if workerConn != nil {
        workerConn.Close()
    }
}
