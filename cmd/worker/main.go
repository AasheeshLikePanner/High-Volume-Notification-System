package main;

import (
	"encoding/json"
    "log"

    "scalable-notification/internal/models"
    "scalable-notification/internal/queue"
    "scalable-notification/internal/notifier" 
)

func main() {
	err := queue.ConnectWorker();
	
	if err != nil {	
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queue.CloseWorker();
	log.Println("Worker started. Waiting for notifications...")

    queue.ConsumeNotification("notifications_exchange", handleNotification)
    queue.ConsumeNotification("notifications_exchange", handleNotification)

    // Block forever
    select {}
}

func handleNotification(msg []byte) error {
	var notif models.Notification;
	if err := json.Unmarshal(msg, &notif); err != nil {
		log.Printf("Failed to unmarshal notification: %v", err)
		return err;
	}
	if err := notifier.Send(notif); err != nil {
		log.Printf("Failed to send notification: %v", err)
		return err;
	}
	log.Printf("Notification sent: %+v", notif)
	return nil;
}