package aggregator

import (
	"context"
	"encoding/json"
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
	fields := logrus.Fields{}
	fields["type"] = e.EventType()
	fields["id"] = e.EventAggregatorID()
	log.AddFields(a.Ctx, fields)

	l := log.FromContext(a.Ctx)
	l = l.WithFields(a.fields)

	// default party path
	partyPath := debtorParty
	switch e.EventType() {
	case app.PaymentCreated:
		var p app.Payment
		if err = json.Unmarshal(e.EventData(), &p); err != nil {
			return errors.Wrap(err, "couldn't unmarshal data")
		}
		p.LastUpdate = e.EventTime()

		if err := a.repository.Insert(a.Ctx, p); err != nil {
			return errors.Wrap(err, "couldn't create in repository")
		}
		a.cache.Set(e.EventAggregatorID(), p)

	case app.PaymentDeleted:
		if err := a.repository.Delete(a.Ctx, e.EventAggregatorID()); err != nil {
			return errors.Wrap(err, "couldn't create in repository")
		}
		a.cache.Remove(e.EventAggregatorID())
		break
	case app.BeneficiaryUpdated:
		// change party path and fallthrough
		partyPath = beneficiaryParty
		fallthrough
	case app.DebtorUpdated:
		var party app.ThirdParty
		if err = json.Unmarshal(e.EventData(), &party); err != nil {
			return errors.Wrap(err, "couldn't unmarshal data")
		}
		party.PaymentID = e.EventAggregatorID()
		party.Time = e.EventTime()

		if err = a.updateThirdParty(a.Ctx, party, partyPath); err != nil {
			return
		}
		curr := a.cache.Get(e.EventAggregatorID())
		if curr == nil {
			return errors.New("couldn't get from cache")
		}
		curr.Attributes.BeneficiaryParty = payment.ThirdParty{
			PaymentID:     party.PaymentID,
			AccountName:   party.AccountName,
			AccountNumber: party.AccountNumber}

		a.cache.Set(e.EventAggregatorID(), *curr)
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
