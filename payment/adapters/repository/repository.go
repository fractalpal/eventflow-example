package repository

import (
	"context"
	"database/sql"

	"github.com/fractalpal/eventflow"
	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment/app"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//  Repository
type FlowRepository struct {
	Commander eventflow.Commander
	Store     app.Store
	Publisher eventflow.Publisher
	fields    logrus.Fields
}

func NewFlow(commander eventflow.Commander, store app.Store, publisher eventflow.Publisher) *FlowRepository {
	fields := logrus.Fields{}
	fields["repository"] = []string{"FlowRepository"}
	return &FlowRepository{
		Commander: commander,
		Store:     store,
		Publisher: publisher,
		fields:    fields,
	}
}

func (m *FlowRepository) CommandStorePublish(cmd interface{}, storeFunc func(context.Context, eventflow.Event) error, ctx context.Context) (err error) {
	ev, err := m.Commander.Command(cmd)
	if err != nil {
		return errors.Wrap(err, "couldn't handle event")
	}

	if err = storeFunc(ctx, ev); err != nil {
		if err == sql.ErrNoRows {
			return app.ErrNoRows
		}
		return
	}

	if err = m.Publisher.Publish(ev); err != nil {
		return errors.Wrap(err, "couldn't publish event")
	}
	return
}

func (m *FlowRepository) Delete(ctx context.Context, id string) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	log.AddFields(ctx, m.fields)
	fields := logrus.Fields{}
	fields["id"] = []string{id}
	log.AddFields(ctx, fields)

	l := log.FromContext(ctx)
	l = l.WithFields(m.fields)

	if err = m.CommandStorePublish(app.DeletePaymentCommand{ID: id}, m.Store.Delete, ctx); err != nil {
		return
	}

	l.Debug("deleted payment")
	return
}

func (m *FlowRepository) Save(ctx context.Context, payment app.Payment) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	log.AddFields(ctx, m.fields)
	fields := logrus.Fields{}
	fields["id"] = []string{payment.ID}
	log.AddFields(ctx, fields)

	l := log.FromContext(ctx)
	l = l.WithFields(m.fields)

	if err = m.CommandStorePublish(app.CreatePaymentCommand{Payment: payment}, m.Store.Save, ctx); err != nil {
		return
	}

	l.Debug("payment created")
	return nil
}

func (m *FlowRepository) UpdateBeneficiary(ctx context.Context, party app.ThirdParty) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	l := log.FromContext(ctx)
	l = l.WithFields(m.fields)

	if err = m.CommandStorePublish(app.UpdateBeneficiaryCommand{ThirdParty: party}, m.Store.Update, ctx); err != nil {
		return
	}

	l.Debug("updated beneficiary")
	return nil
}

func (m *FlowRepository) UpdateDebtor(ctx context.Context, party app.ThirdParty) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	l := log.FromContext(ctx)
	l = l.WithFields(m.fields)

	if err = m.CommandStorePublish(app.UpdateDebtorCommand{ThirdParty: party}, m.Store.Update, ctx); err != nil {
		return
	}

	l.Debug("debtor updated")
	return nil
}
