package queue

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeNotification(queueName string, handler func([]byte) error) {
	// Set prefetch for fair dispatch
	err := workerCh.Qos(10, 0, false) // 10 unacked msgs per consumer max
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	for i := 0; i < 50; i++ { // 10 concurrent consumers
		go func(workerID int) {
			msgs, err := workerCh.Consume(
				queueName,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Fatalf("Failed to register consumer: %v", err)
			}

			for msg := range msgs {
				log.Printf("[Worker %d] Received a message: %s", workerID, msg.Body)

				headers := msg.Headers
				if headers == nil {
					headers = amqp.Table{}
				}

				retryCount := 0
				if count, exists := headers["x-retry-count"]; exists {
					switch v := count.(type) {
					case int:
						retryCount = v
					case int32:
						retryCount = int(v)
					case int64:
						retryCount = int(v)
					case float64:
						retryCount = int(v)
					}
				}

				err := handler(msg.Body)

				if err != nil {
					if retryCount < 3 {
						retryCount++
						retryDelay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
						nextRetryTime := time.Now().Add(retryDelay).UnixNano()

						headers["x-retry-count"] = retryCount
						headers["x-retry-timestamp"] = nextRetryTime

						err = workerCh.Publish(
							"",
							queueName+"_retry",
							false,
							false,
							amqp.Publishing{
								Headers: headers,
								Body:    msg.Body,
							},
						)
						if err != nil {
							log.Printf("[Worker %d] Failed to requeue message: %v", workerID, err)
						}
						msg.Reject(false)
					} else {
						err = workerCh.Publish(
							"",
							queueName+"_dead",
							false,
							false,
							amqp.Publishing{
								Headers: headers,
								Body:    msg.Body,
							},
						)
						if err != nil {
							log.Printf("[Worker %d] Failed to send to dead letter queue: %v", workerID, err)
						}
						msg.Reject(false)
						log.Printf("[Worker %d] Message sent to dead letter queue: %v", workerID, msg.Body)
					}
				} else {
					msg.Ack(false)
					log.Printf("[Worker %d] Message processed: %v", workerID, msg.Body)
				}
			}
		}(i)
	}
}
