package app

import (
	"context"
)

// Validator for validates payments
type Validator interface {
	Validate(context.Context, Payment) error
}
