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
	//AmountUpdated      = eventflow.EventType("AmountUpdated")
	//OtherHappened      = eventflow.EventType("OtherHappened")
)

// CreatePaymentCommand
type CreatePaymentCommand struct {
	Payment
}

// DeletePaymentCommand
type DeletePaymentCommand struct {
	ID string
}

// UpdateBeneficiaryCommand
type UpdateBeneficiaryCommand struct {
	ThirdParty
}

// UpdateDebtorCommand
type UpdateDebtorCommand struct {
	ThirdParty
}
