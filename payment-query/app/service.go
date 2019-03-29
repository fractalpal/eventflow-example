package app

import (
	"context"
)

// Service interface for fetching payments
type Service interface {
	FindByID(context.Context, string) (Payment, error)
	FindAll(context.Context, int64, int64) ([]Payment, error)
}
