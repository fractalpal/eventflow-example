package app

import (
	"context"

	"errors"
)

// ErrPartyNotSupported informs about not supported third party type
var ErrPartyNotSupported = errors.New("third party not supported")

// ErrExists informs about already existing record
var ErrExists = errors.New("already exists")

// ErrExists informs about no results
var ErrNoResults = errors.New("no results")

// Store interface for payments
type Store interface {
	Insert(context.Context, Payment) error
	UpdateThirdParty(context.Context, ThirdParty, string) error
	FindByID(context.Context, string) (Payment, error)
	FindAll(context.Context, int64, int64) ([]Payment, error)
	Delete(context.Context, string) error
}
