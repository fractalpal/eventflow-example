package app

import (
	"context"

	"github.com/pkg/errors"
)

// ErrPartyNotSupported informs about not supported third party
var ErrPartyNotSupported = errors.New("third party not supported")

// Service
type Service interface {
	Create(context.Context, Payment) (Payment, error)
	UpdateThirdParty(context.Context, ThirdParty, string) error
	Delete(context.Context, string) error
}
