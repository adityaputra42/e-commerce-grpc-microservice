package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	OrderCreatedSubject   = "order.created"
	OrderCreatedQueueName = "order_workers"
)

type OrderCreatedPayload struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type NatsWorker struct {
	natsConn *nats.Conn
}

func NewNatsWorker(natsURL string) (*NatsWorker, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	return &NatsWorker{
		natsConn: conn,
	}, nil
}

func (w *NatsWorker) PublishOrderCreated(ctx context.Context, payload OrderCreatedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal order created payload: %v", err)
	}

	err = w.natsConn.Publish(OrderCreatedSubject, data)
	if err != nil {
		return fmt.Errorf("failed to publish order created event: %v", err)
	}

	// Pastikan pesan dikirim
	err = w.natsConn.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush NATS connection: %v", err)
	}

	if err := w.natsConn.LastError(); err != nil {
		return fmt.Errorf("NATS reported an error: %v", err)
	}

	log.Printf("Published order.created event for order ID: %s", payload.OrderID)
	return nil
}

func (w *NatsWorker) Close() {
	if w.natsConn != nil && !w.natsConn.IsClosed() {
		w.natsConn.Close()
		log.Println("NATS connection closed")
	}
}
