package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/models/payment"
	"github.com/fractalpal/eventflow-example/payment-query/app"
	"github.com/fractalpal/eventflow"
)

const (
	beneficiaryParty = "beneficiary_party"
	debtorParty      = "debtor_party"
)

// Aggregator is a simple in memory holder
// normally should be used more persistent storage (redis?)
// and be able to handle async communication (ex. Kafka)
//
// Storing raw events on query part not implemented
// that could be helpful for rebuild/reapply aggregator again from local store
//
// This aggregator applies incoming events as they come
// building current materialized view of payments in the system
type Aggregator struct {
	fields     logrus.Fields
	repository app.Repository
	Ctx        context.Context
	cache      PaymentsCache
}

func NewMemory(ctx context.Context, repository app.Repository) *Aggregator {
	fields := logrus.Fields{}
	fields["aggregator"] = []string{"NewMemory"}
	return &Aggregator{
		fields:     fields,
		repository: repository,
		Ctx:        ctx,
		cache:      NewInMemoryCache(),
	}
}

func (a *Aggregator) Apply(e eventflow.Event) (err error) {
	defer log.AddFieldsForErr(a.Ctx, a.fields, err)
	paymentID, ok := e.Columns[app.AggregatePayments].(string)
	if !ok {
		err = fmt.Errorf("expected column '%s' is not a string", app.AggregatePayments)
		return
	}
	fields := logrus.Fields{}
	fields["type"] = e.Type
	fields["id"] = paymentID
	log.AddFields(a.Ctx, fields)

	l := log.FromContext(a.Ctx)
	l = l.WithFields(a.fields)

	// default party path
	partyPath := debtorParty
	switch e.Type {
	case app.PaymentCreated:
		var p app.Payment
		if err = json.Unmarshal(e.Data, &p); err != nil {
			return errors.Wrap(err, "couldn't unmarshal data")
		}
		p.LastUpdateTimestamp = e.Timestamp

		if err := a.repository.Insert(a.Ctx, p); err != nil {
			return errors.Wrap(err, "couldn't create in repository")
		}
		a.cache.Set(paymentID, p)

	case app.PaymentDeleted:
		if err := a.repository.Delete(a.Ctx, paymentID); err != nil {
			return errors.Wrap(err, "couldn't create in repository")
		}
		a.cache.Remove(paymentID)
		break
	case app.BeneficiaryUpdated:
		// change party path and fallthrough
		partyPath = beneficiaryParty
		fallthrough
	case app.DebtorUpdated:
		var party app.ThirdParty
		if err = json.Unmarshal(e.Data, &party); err != nil {
			return errors.Wrap(err, "couldn't unmarshal data")
		}
		party.PaymentID = paymentID
		party.Timestamp = e.Timestamp

		if err = a.updateThirdParty(a.Ctx, party, partyPath); err != nil {
			return
		}
		curr := a.cache.Get(paymentID)
		if curr == nil {
			return errors.New("couldn't get from cache")
		}
		curr.Attributes.BeneficiaryParty = payment.ThirdParty{
			PaymentID:     party.PaymentID,
			AccountName:   party.AccountName,
			AccountNumber: party.AccountNumber}

		a.cache.Set(paymentID, *curr)
	default:
		return app.ErrPartyNotSupported
	}

	l.Debug("event applied")
	return nil
}

func (a *Aggregator) updateThirdParty(ctx context.Context, party app.ThirdParty, partyKey string) (err error) {
	fields := logrus.Fields{}
	fields["party"] = partyKey
	defer log.AddFieldsForErr(ctx, fields, err)

	switch partyKey {
	case beneficiaryParty, debtorParty:
		if err := a.repository.UpdateThirdParty(ctx, party, partyKey); err != nil {
			return errors.Wrap(err, "couldn't update in repository")
		}
		break
	default:
		err = app.ErrPartyNotSupported
		return
	}
	return
}
