package service

import (
	"context"

	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment-query/app"
	"github.com/sirupsen/logrus"
)

type query struct {
	Repository app.Repository
	AppCtx     context.Context
	fields     logrus.Fields
}

func NewQuery(repository app.Repository) *query {
	fields := logrus.Fields{}
	fields["service"] = []string{"query"}
	return &query{
		fields:     fields,
		Repository: repository,
	}
}

func (s *query) FindByID(ctx context.Context, id string) (payment app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	log.AddFields(ctx, s.fields)
	fields := logrus.Fields{}
	fields["id"] = []string{id}
	log.AddFields(ctx, fields)

	payment, err = s.Repository.FindByID(ctx, id)
	if err != nil {
		return
	}

	fields["func"] = "find"
	log.AddFields(ctx, fields)
	l := log.FromContext(ctx)
	l.Debug("found in repository")
	return payment, nil
}

func (s *query) FindAll(ctx context.Context, page, limit int64) (payments []app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	payments, err = s.Repository.FindAll(ctx, page, limit)
	if err != nil {
		return
	}

	log.AddFields(ctx, s.fields)
	fields := logrus.Fields{}
	fields["func"] = "find_all"
	log.AddFields(ctx, fields)
	l := log.FromContext(ctx)
	l.Debug("found in repository")
	return
}
