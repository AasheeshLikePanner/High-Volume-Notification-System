
# Scalable Notification System

This is a notification system built with Go and RabbitMQ. It supports:

- Priority-based delivery (e.g. OTPs first, promos later)
- Automatic retries with backoff (2s → 4s → 8s)
- Dead-letter queue for failed messages
- Horizontal scaling with multiple workers
- Sends Email, SMS, or Push (simulated)

## Test Metrics

- Handled: **1M+ messages**
- Throughput: **12.5K notifications/sec**
- 50 goroutines per worker


## Architecture
![Screenshot From 2025-04-30 02-39-01](https://github.com/user-attachments/assets/0db26007-ae82-406f-9b18-eac7d1e57810)

