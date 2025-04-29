package main;

import (
	"time"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	rate := 100000000 // messages per second
	duration := 10 * time.Second // test duration

	log.Println("Starting load test...")
	loadTestRabbitMQ(rate, duration)
}

func loadTestRabbitMQ(rate int, duration time.Duration) {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
    ch, _ := conn.Channel()
    defer conn.Close()
    defer ch.Close()

    ticker := time.NewTicker(time.Second / time.Duration(rate))
    timer := time.NewTimer(duration)

    for {
        select {
        case <-ticker.C:
            body := []byte(`{"type":"alert","to":"test@example.com","message":"Test","priority":5}`)
            err := ch.Publish("", "notifications", false, false, amqp.Publishing{
                ContentType: "application/json",
                Body:        body,
            })
            if err != nil {
                log.Println("Failed to publish:", err)
            }
        case <-timer.C:
            log.Println("Load test finished")
            return
        }
    }
}
