package notifier

import (
	"fmt"
    "scalable-notification/internal/models"
)

func Send(notif models.Notification) error {

	switch notif.Type {
    case "otp":
        // Here call SMS provider API (fake for now)
        fmt.Printf("Sending OTP to %s: %s\n", notif.To, notif.Message)
    case "promo":
        // Here call Email provider API
        fmt.Printf("Sending Promo Email to %s: %s\n", notif.To, notif.Message)
    default:
        fmt.Printf("Unknown notification type %s to %s: %s\n", notif.Type, notif.To, notif.Message)
    }
    return nil

}