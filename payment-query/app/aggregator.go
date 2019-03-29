package app

import "github.com/fractalpal/eventflow"

const (
	AggregatePayments = "payments"
)

const (
	PaymentCreated     = eventflow.EventType("PaymentCreated")
	PaymentDeleted     = eventflow.EventType("PaymentDeleted")
	BeneficiaryUpdated = eventflow.EventType("BeneficiaryUpdated")
	DebtorUpdated      = eventflow.EventType("DebtorUpdated")
)
