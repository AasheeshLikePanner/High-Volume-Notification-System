package main

import (
	"encoding/json"
    "log"
    "net/http"

    "scalable-notification/internal/queue"  
    "scalable-notification/internal/models"
)


func main(){
	err := queue.ConnectServer();

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queue.CloseServer();

	http.HandleFunc("/send-notification", handleSendNotification);

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSendNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var notif models.Notification;
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	notif.CreatedAt = models.Now();

	if err := queue.PublishNotification(notif); err != nil {
		http.Error(w, "Failed to publish notification", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Notification sent successfully"))
	log.Printf("Notification sent: %+v", notif)
}
