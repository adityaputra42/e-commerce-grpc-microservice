package services

import (
	"context"
	pb "e-commerce-microservice/payment/internal/pb"
	"e-commerce-microservice/payment/internal/token"
	"e-commerce-microservice/payment/internal/worker"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nats-io/nats.go"
)

type PaymentService struct {
	tokenMaker token.Maker
	natsConn   *nats.Conn
}

func NewPaymentService(tokenMaker token.Maker, natsURL string) (*PaymentService, error) {
	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	service := &PaymentService{
		tokenMaker: tokenMaker,
		natsConn:   natsConn,
	}

	// Subscribe to order.created events
	_, err = natsConn.Subscribe(worker.OrderCreatedSubject, service.handleOrderCreated)
	if err != nil {
		natsConn.Close()
		return nil, fmt.Errorf("failed to subscribe to order events: %v", err)
	}

	return service, nil
}

func (s *PaymentService) handleOrderCreated(msg *nats.Msg) {
	var orderPayload worker.OrderCreatedPayload
	if err := json.Unmarshal(msg.Data, &orderPayload); err != nil {
		fmt.Printf("Failed to unmarshal order created payload: %v\n", err)
		return
	}

	// Create payment request
	req := &pb.CreatePaymentRequest{
		OrderId:  orderPayload.OrderID,
		Username: orderPayload.UserID,
		Amount:   strconv.FormatFloat(orderPayload.TotalAmount, 'f', 2, 64),
	}

	// Process payment
	err := s.ProcessPayment(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed to process payment for order %s: %v\n", orderPayload.OrderID, err)
		return
	}

	fmt.Printf("Successfully processed payment for order %s\n", orderPayload.OrderID)
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.CreatePaymentRequest) error {
	// Add your payment processing logic here
	// This could involve:
	// 1. Validating the payment request
	// 2. Calling a payment gateway
	// 3. Recording the payment in the database
	// 4. Publishing a payment.completed event

	// For now, we'll just simulate a successful payment

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %v", err)
	}
	paymentCompletedPayload := worker.PaymentCompletedPayload{
		OrderID:   req.OrderId,
		UserID:    req.Username,
		Amount:    amount,
		Status:    "completed",
		PaymentID: "payment_" + req.OrderId, // You would generate this in a real implementation
	}

	data, err := json.Marshal(paymentCompletedPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payment completed payload: %v", err)
	}

	err = s.natsConn.Publish(worker.PaymentCompletedSubject, data)
	if err != nil {
		return fmt.Errorf("failed to publish payment completed event: %v", err)
	}

	return nil
}

func (s *PaymentService) Close() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
}
