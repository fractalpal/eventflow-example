package app

import (
	"context"
)

type Repository interface {
	FindByID(context.Context, string) (Payment, error)
	FindAll(context.Context, int64, int64) ([]Payment, error)
	Insert(context.Context, Payment) error
	UpdateThirdParty(context.Context, ThirdParty, string) error
	Delete(context.Context, string) error
}
