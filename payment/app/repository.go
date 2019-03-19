package app

import (
	"context"
)

type Repository interface {
	Save(context.Context, Payment) error
	UpdateBeneficiary(context.Context, ThirdParty) error
	UpdateDebtor(context.Context, ThirdParty) error
	Delete(context.Context, string) error
}
