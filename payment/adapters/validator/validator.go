package validator

import (
	"context"

	"github.com/fractalpal/eventflow-example/payment/app"
)

type Validator interface {
	Validate(context.Context, app.Payment) error
}

// EmptyValidator. No validation
type EmptyValidator struct {
}

func NewEmpty() *EmptyValidator {
	return &EmptyValidator{}
}

func (v *EmptyValidator) Validate(ctx context.Context, payment app.Payment) error {
	return nil
}
