package worker

import "time"

const (
	OrderCreatedSubject     = "order.created"
	PaymentCompletedSubject = "payment.completed"
)

type OrderCreatedPayload struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaymentCompletedPayload struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	PaymentID string  `json:"payment_id"`
}
