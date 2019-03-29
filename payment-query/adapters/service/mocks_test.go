// +build test

package service

import (
	"context"
	"errors"

	"github.com/fractalpal/eventflow-example/payment-query/app"
	"github.com/stretchr/testify/mock"
)

// MockedRepository for Repository interface
type MockedRepository struct {
	mock.Mock
}

func (m *MockedRepository) Insert(ctx context.Context, p app.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockedRepository) UpdateThirdParty(ctx context.Context, party app.ThirdParty, partyKey string) error {
	args := m.Called(ctx, party)
	return args.Error(0)
}
func (m *MockedRepository) FindByID(ctx context.Context, id string) (app.Payment, error) {
	args := m.Called(ctx, id)
	p, ok := args.Get(0).(app.Payment)
	if !ok {
		return app.Payment{}, errors.New("couldn't cast to app.Payment")
	}
	return p, args.Error(1)
}
func (m *MockedRepository) FindAll(ctx context.Context, page int64, limit int64) ([]app.Payment, error) {
	args := m.Called(ctx)
	p, ok := args.Get(0).([]app.Payment)
	if !ok {
		return nil, errors.New("couldn't cast to []app.Payment")
	}
	return p, args.Error(1)
}
func (m *MockedRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
