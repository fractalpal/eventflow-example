package app

import "github.com/fractalpal/eventflow"

const (
	PaymentCreated     = eventflow.EventType("PaymentCreated")
	PaymentDeleted     = eventflow.EventType("PaymentDeleted")
	BeneficiaryUpdated = eventflow.EventType("BeneficiaryUpdated")
	DebtorUpdated      = eventflow.EventType("DebtorUpdated")
	//AmountUpdated      = eventflow.EventType("AmountUpdated")
	//OtherHappened      = eventflow.EventType("OtherHappened")
)

type CreatePaymentCommand struct {
	Payment
}

type DeletePaymentCommand struct {
	ID string
}

type UpdateBeneficiaryCommand struct {
	ThirdParty
}

type UpdateDebtorCommand struct {
	ThirdParty
}
