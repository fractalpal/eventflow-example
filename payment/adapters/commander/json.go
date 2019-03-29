package commander

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fractalpal/eventflow-example/payment/app"
	"github.com/fractalpal/eventflow"
)

type jsonCommander struct {
}

func NewJson() *jsonCommander {
	return &jsonCommander{}
}

func (c *jsonCommander) Command(data interface{}) (eventflow.Event, error) {
	e := eventflow.Event{
		Columns: make(map[string]interface{}),
	}
	var payload interface{}
	switch data.(type) {
	case app.CreatePaymentCommand:
		payment := data.(app.CreatePaymentCommand).Payment
		e.Columns[app.AggregatePayments] = payment.ID
		e.Type = app.PaymentCreated
		payload = payment
	case app.DeletePaymentCommand:
		e.Columns[app.AggregatePayments] = data.(app.DeletePaymentCommand).ID
		e.Type = app.PaymentDeleted
	case app.UpdateBeneficiaryCommand:
		party := data.(app.UpdateBeneficiaryCommand).ThirdParty
		payload = party
		e.Columns[app.AggregatePayments] = party.PaymentID
		e.Type = app.BeneficiaryUpdated
	case app.UpdateDebtorCommand:
		party := data.(app.UpdateDebtorCommand).ThirdParty
		payload = party
		e.Columns[app.AggregatePayments] = party.PaymentID
		e.Type = app.DebtorUpdated
	default:
		return e, fmt.Errorf("unsupported command: '%v'", data)

	}

	e.Timestamp = time.Now().UnixNano()
	bytes, err := json.Marshal(payload)
	if err != nil {
		return e, err
	}
	e.Data = bytes
	e.Mapper = "json"
	return e, nil
}
