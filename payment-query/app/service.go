package app

import (
	"context"
)

type Service interface {
	FindByID(context.Context, string) (Payment, error)
	FindAll(context.Context, int64, int64) ([]Payment, error)
}
