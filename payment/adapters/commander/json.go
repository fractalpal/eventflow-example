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
	e := eventflow.BaseEvent{}
	var payload interface{}
	switch data.(type) {
	case app.CreatePaymentCommand:
		payment := data.(app.CreatePaymentCommand).Payment
		e.AggregatorID = payment.ID
		e.Type = app.PaymentCreated
		payload = payment
	case app.DeletePaymentCommand:
		e.AggregatorID = data.(app.DeletePaymentCommand).ID
		e.Type = app.PaymentDeleted
	case app.UpdateBeneficiaryCommand:
		party := data.(app.UpdateBeneficiaryCommand).ThirdParty
		payload = party
		e.AggregatorID = party.PaymentID
		e.Type = app.BeneficiaryUpdated
	case app.UpdateDebtorCommand:
		party := data.(app.UpdateDebtorCommand).ThirdParty
		payload = party
		e.AggregatorID = party.PaymentID
		e.Type = app.DebtorUpdated
	default:
		return nil, fmt.Errorf("unsupported command: '%v'", data)

	}

	e.Time = time.Now()
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	e.Data = bytes
	e.Mapper = "json"
	return &e, nil
}
