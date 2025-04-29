package queue

import (
    "time"
    "os"
	"strconv"
)

type Config struct {
    URL             string
    ReconnectDelay  time.Duration
    MaxRetries      int
    HeartbeatInterval time.Duration
    QueueNames      struct {
        Main  string
        Retry string
        Dead  string
    }
}

func NewConfig() *Config {
    return &Config{
        URL:             getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"),
        ReconnectDelay:  time.Duration(getIntEnv("RABBITMQ_RECONNECT_DELAY_SECONDS", 5)) * time.Second,
        MaxRetries:      getIntEnv("RABBITMQ_MAX_RETRIES", 5),
        HeartbeatInterval: time.Duration(getIntEnv("RABBITMQ_HEARTBEAT_SECONDS", 10)) * time.Second,
        QueueNames: struct {
            Main  string
            Retry string
            Dead  string
        }{
            Main:  getEnv("RABBITMQ_QUEUE_MAIN", "notifications"),
            Retry: getEnv("RABBITMQ_QUEUE_RETRY", "notifications_retry"),
            Dead:  getEnv("RABBITMQ_QUEUE_DEAD", "notifications_dead"),
        },
    }
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}

func getIntEnv(key string, fallback int) int {
    if value, exists := os.LookupEnv(key); exists {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return fallback
}