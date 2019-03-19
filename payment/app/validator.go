package app

import (
	"context"
)

type Validator interface {
	Validate(context.Context, Payment) error
}
