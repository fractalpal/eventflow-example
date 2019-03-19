package app

import (
	"context"

	"errors"
)

var ErrPartyNotSupported = errors.New("third party not supported")
var ErrExists = errors.New("already exists")
var ErrNoResults = errors.New("no results")

type Store interface {
	Insert(context.Context, Payment) error
	UpdateThirdParty(context.Context, ThirdParty, string) error
	FindByID(context.Context, string) (Payment, error)
	FindAll(context.Context, int64, int64) ([]Payment, error)
	Delete(context.Context, string) error
}
