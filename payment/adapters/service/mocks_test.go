// +build test

package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/fractalpal/eventflow-example/payment/app"
)

// MockedIDProvider for IDProvider interface
type MockedIDProvider struct {
	mock.Mock
}

func (m *MockedIDProvider) ID() string {
	args := m.Called()
	return args.String(0)
}

// MockedValidator for Validator interface
type MockedValidator struct {
	mock.Mock
}

func (m *MockedValidator) Validate(ctx context.Context, p app.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

// MockedRepository for Repository interface
type MockedRepository struct {
	mock.Mock
}

func (m *MockedRepository) Save(ctx context.Context, p app.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockedRepository) UpdateBeneficiary(ctx context.Context, party app.ThirdParty) error {
	args := m.Called(ctx, party)
	return args.Error(0)
}
func (m *MockedRepository) UpdateDebtor(ctx context.Context, party app.ThirdParty) error {
	args := m.Called(ctx, party)
	return args.Error(0)
}
func (m *MockedRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
