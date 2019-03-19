package app

import (
	"context"

	"database/sql"

	"github.com/fractalpal/eventflow"
)

var ErrNoRows = sql.ErrNoRows

type Store interface {
	Save(context.Context, eventflow.Event) error
	Update(context.Context, eventflow.Event) error
	Delete(context.Context, eventflow.Event) error
}
