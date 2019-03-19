package app

import (
	"context"

	"github.com/pkg/errors"
)

var ErrPartyNotSupported = errors.New("third party not supported")

type Service interface {
	Create(context.Context, Payment) (Payment, error)
	UpdateThirdParty(context.Context, ThirdParty, string) error
	Delete(context.Context, string) error
}
