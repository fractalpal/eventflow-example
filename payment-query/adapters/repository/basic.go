package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment-query/app"
)

//Basic Repository impl., without cache, just store calls
type Basic struct {
	fields logrus.Fields
	store  app.Store
}

func NewBasic(store app.Store) *Basic {
	fields := logrus.Fields{}
	fields["repository"] = []string{"Basic"}
	return &Basic{
		fields: fields,
		store:  store,
	}
}

func (m *Basic) Insert(ctx context.Context, payment app.Payment) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	l := log.FromContext(ctx)
	fields := logrus.Fields{}
	fields["op"] = "create"
	fields["id"] = []string{payment.ID}
	l = l.WithFields(m.fields).WithFields(fields)

	if err = m.store.Insert(ctx, payment); err != nil {
		if err == app.ErrExists {
			l.WithError(err).Info("silently ignore that payment exists")
			err = nil
			return
		}
		return
	}

	l.Debug("payment created")
	return
}

func (m *Basic) UpdateThirdParty(ctx context.Context, thirdParty app.ThirdParty, partyKey string) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	if err = m.store.UpdateThirdParty(ctx, thirdParty, partyKey); err != nil {
		return
	}

	l := log.FromContext(ctx).WithFields(m.fields)
	fields := logrus.Fields{}
	fields["op"] = "updated"
	fields["party"] = partyKey
	l = l.WithFields(fields)
	l.Debug("updated in repository")
	return
}

func (m *Basic) FindByID(ctx context.Context, id string) (payment app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	if payment, err = m.store.FindByID(ctx, id); err != nil {
		return
	}

	l := log.FromContext(ctx)
	fields := logrus.Fields{}
	fields["id"] = []string{id}
	l = l.WithFields(m.fields).WithFields(fields)
	l.Debug("found in repository")
	return
}

func (m *Basic) FindAll(ctx context.Context, page int64, limit int64) (payments []app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)

	if payments, err = m.store.FindAll(ctx, page, limit); err != nil {
		return
	}
	// serve at least empty collection
	if payments == nil {
		payments = []app.Payment{}
	}

	l := log.FromContext(ctx).WithFields(m.fields)
	l.Debug("found list in repository")
	return
}

func (m *Basic) Delete(ctx context.Context, id string) (err error) {
	defer log.AddFieldsForErr(ctx, m.fields, err)
	l := log.FromContext(ctx)
	l.WithFields(m.fields)
	l.WithField("id", id)

	if err = m.store.Delete(ctx, id); err != nil {
		if err == app.ErrExists {
			l.WithError(err).Info("silently ignore that payment not exists")
			err = nil
			return
		}
		return
	}
	l.Debug("deleted from repository")
	return
}
