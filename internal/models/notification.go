package models

import "time"

type Notification struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`     // otp, promo, alert, etc.
    To        string    `json:"to"`
    Message   string    `json:"message"`
    Priority  int       `json:"priority"` // 10 = High, 1 = Low
    CreatedAt time.Time `json:"created_at"`
}

func Now() time.Time {
    return time.Now().UTC()
}
