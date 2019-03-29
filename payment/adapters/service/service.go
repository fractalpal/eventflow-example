package service

import (
	"context"

	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment/app"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type IDProvider interface {
	ID() string
}

type UUIDProvider struct {
}

func (p *UUIDProvider) ID() string {
	return uuid.New().String()
}

type Svc struct {
	idProvider IDProvider
	validator  app.Validator
	repository app.Repository
	fields     logrus.Fields
}

func New(idProvider IDProvider, repo app.Repository, validator app.Validator) app.Service {
	fields := logrus.Fields{}
	fields["service"] = []string{"payment"}
	svc := &Svc{
		idProvider: idProvider,
		validator:  validator,
		repository: repo,
		fields:     fields,
	}
	return svc
}

func (s *Svc) Create(ctx context.Context, payment app.Payment) (pay app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)
	// this implementation generates ID and overrides any incoming
	id := s.idProvider.ID()
	payment.ID = id
	fields := logrus.Fields{}
	fields["id"] = []string{payment.ID}
	log.AddFields(ctx, fields)
	l := log.FromContext(ctx)

	// for business logic validation
	if err = s.validator.Validate(ctx, payment); err != nil {
		err = errors.Wrap(err, "cannot validate payment")
		return app.Payment{}, err
	}

	if err = s.repository.Save(ctx, payment); err != nil {
		err = errors.Wrap(err, "cannot create in repository")
		return app.Payment{}, err
	}

	l.Debug("payment created")

	return payment, nil
}

func (s *Svc) Delete(ctx context.Context, id string) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	fields := logrus.Fields{}
	fields["id"] = []string{id}
	log.AddFields(ctx, fields)
	l := log.FromContextWithFields(ctx, s.fields)

	if err = s.repository.Delete(ctx, id); err != nil {
		return
	}
	l.Debug("payment deleted")
	return nil
}

func (s *Svc) UpdateThirdParty(ctx context.Context, thirdParty app.ThirdParty, partyKey string) (err error) {

	switch partyKey {
	case "beneficiary":
		if err = s.updateBeneficiaryParty(ctx, thirdParty); err != nil {
			return
		}
		break
	case "debtor":
		if err = s.updateDebtorParty(ctx, thirdParty); err != nil {
			return
		}
		break
	default:
		return app.ErrPartyNotSupported
	}

	return
}

func (s *Svc) updateBeneficiaryParty(ctx context.Context, party app.ThirdParty) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	fields := logrus.Fields{}
	fields["third_party"] = []string{"beneficiary"}
	log.AddFields(ctx, fields)
	l := log.FromContextWithFields(ctx, s.fields)

	if err = s.repository.UpdateBeneficiary(ctx, party); err != nil {
		if errors.Cause(err) == app.ErrNoRows {
			return
		}
		err = errors.Wrap(err, "couldn't update in repository")
		return
	}

	l.Debug("beneficiary updated")
	return nil
}

func (s *Svc) updateDebtorParty(ctx context.Context, party app.ThirdParty) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	fields := logrus.Fields{}
	fields["third_party"] = []string{"debtor"}
	log.AddFields(ctx, fields)
	l := log.FromContextWithFields(ctx, s.fields)

	if err = s.repository.UpdateDebtor(ctx, party); err != nil {
		if errors.Cause(err) == app.ErrNoRows {
			return
		}
		err = errors.Wrap(err, "couldn't update in repository")
		return
	}

	l.Debug("debtor updated")
	return nil
}
