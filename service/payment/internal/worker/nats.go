package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// PaymentService adalah interface yang harus diimplementasikan oleh service pembayaran
type PaymentService interface {
	CreatePayment(ctx context.Context, userID string, orderID string, amount float64) error
}

// PaymentWorker bertanggung jawab mendengarkan event order dan memproses pembayaran
type PaymentWorker struct {
	natsConn   *nats.Conn
	paymentSvc PaymentService
}

// NewPaymentWorker menginisialisasi koneksi ke NATS dan membuat worker
func NewPaymentWorker(natsURL string, paymentSvc PaymentService) (*PaymentWorker, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	return &PaymentWorker{
		natsConn:   conn,
		paymentSvc: paymentSvc,
	}, nil
}

// Start memulai worker untuk mendengarkan event "order.created"
func (w *PaymentWorker) Start(ctx context.Context) error {
	sub, err := w.natsConn.QueueSubscribe("order.created", "payment_workers", func(msg *nats.Msg) {
		var payload OrderCreatedPayload

		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal order.created payload: %v", err)
			return
		}

		// Buat pembayaran
		err := w.paymentSvc.CreatePayment(ctx, payload.UserID, payload.OrderID, payload.TotalAmount)
		if err != nil {
			log.Printf("Failed to create payment for order %s: %v", payload.OrderID, err)
			return
		}

		log.Println("Successfully created payment for order", payload.OrderID)

		// (Opsional) Kirim event payment.created
		paymentEvent := struct {
			OrderID string `json:"order_id"`
			Status  string `json:"status"`
		}{
			OrderID: payload.OrderID,
			Status:  "paid",
		}

		data, err := json.Marshal(paymentEvent)
		if err != nil {
			log.Printf("Failed to marshal payment.created event: %v", err)
			return
		}

		err = w.natsConn.Publish("payment.created", data)
		if err != nil {
			log.Printf("Failed to publish payment.created event: %v", err)
			return
		}

		log.Println("Published payment.created event for order", payload.OrderID)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to order.created events: %v", err)
	}

	// Pastikan koneksi dan subscription berhasil
	w.natsConn.Flush()
	if err := w.natsConn.LastError(); err != nil {
		return fmt.Errorf("NATS connection error: %v", err)
	}

	log.Println("PaymentWorker started and subscribed to 'order.created'")

	// Unsubscribe saat context dibatalkan
	defer sub.Unsubscribe()

	<-ctx.Done()

	log.Println("PaymentWorker shutting down")
	return ctx.Err()
}

// Close menutup koneksi NATS
func (w *PaymentWorker) Close() {
	if w.natsConn != nil && !w.natsConn.IsClosed() {
		w.natsConn.Close()
		log.Println("NATS connection closed")
	}
}
