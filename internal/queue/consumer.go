package queue

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
	"math"
)

func ConsumeNotification(queueName string, handler func([]byte) error) {
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

	go func() {
		for msg := range msgs {

			headers := msg.Headers
            if headers == nil {
                headers = amqp.Table{}
            }
            
            retryCount := 0
            if count, exists := headers["x-retry-count"]; exists {
                retryCount = count.(int)
            }

			err := handler(msg.Body)

			if err != nil {
				if retryCount < 3 {
					retryCount++

					retryDelay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
                    nextRetryTime := time.Now().Add(retryDelay).UnixNano()
                    
                    // Update headers for retry
                    headers["x-retry-count"] = retryCount + 1
                    headers["x-retry-timestamp"] = nextRetryTime

					headers["x-retry-count"] = retryCount
					err = workerCh.Publish(
						"",
						queueName + "_retry",
						false,
						false,
						amqp.Publishing{
							Headers: headers,
							Body:    msg.Body,
						},
					)
					if err != nil {
						log.Printf("Failed to requeue message: %v", err)
					}
					msg.Reject(false)
				} else {
					err = workerCh.Publish(
						"",
						queueName + "_dead",
						false,
						false,
						amqp.Publishing{
							Headers: headers,
							Body:    msg.Body,
						},
					)
					if err != nil {
						log.Printf("Failed to send message to dead letter queue: %v", err)
					}
					msg.Reject(false);
					log.Printf("Message sent to dead letter queue: %v", msg.Body)
				}
			}else {
				msg.Ack(false);
				log.Printf("Message processed successfully: %v", msg.Body)
			}	
		}
	}()
}
