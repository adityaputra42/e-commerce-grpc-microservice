package utils

type OrderStatus int

const (
	Pending   = "pending"
	Confirmed = "confirmed"
	Waiting   = "waiting_payment"
	Paid      = "paid"
	Finish    = "finish"
	Candeled  = "canceled"
)
