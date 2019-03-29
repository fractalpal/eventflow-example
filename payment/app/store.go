package app

import (
	"context"

	"database/sql"

	"github.com/fractalpal/eventflow"
)

// ErrNoRows is returned in case of no records found in store
var ErrNoRows = sql.ErrNoRows

// Store interface for payment events
type Store interface {
	Save(context.Context, eventflow.Event) error
	Update(context.Context, eventflow.Event) error
	Delete(context.Context, eventflow.Event) error
}
