package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fractalpal/eventflow"
	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment/app"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Postgres(dataSource string) (*sql.DB, error) {
	return sql.Open("postgres", dataSource)
}

func Migration(db *sql.DB, sourceUrl string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, errors.New("couldn't initialize db diver")
	}
	// Would be better get migration from compiled binary
	return migrate.NewWithDatabaseInstance(
		sourceUrl,
		"postgres",
		driver)
}

type PostgresStore struct {
	fields logrus.Fields
	db     *sql.DB
}

func New(db *sql.DB) *PostgresStore {
	fields := logrus.Fields{}
	fields["store"] = []string{"postgres"}
	return &PostgresStore{
		fields: fields,
		db:     db,
	}
}

func (s *PostgresStore) Save(ctx context.Context, event eventflow.Event) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	tx, err := s.db.Begin()
	if err != nil {
		return
	}

	paymentID, ok := event.Columns[app.AggregatePayments].(string)
	if !ok {
		err = fmt.Errorf("expected column '%s' is not a string", app.AggregatePayments)
		return
	}

	if err = s.existsWithID(tx, ctx, paymentID); err == nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return errors.New("event already created")
	}

	if err = s.insert(tx, ctx, event); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	l := log.FromContext(ctx)
	l = l.WithFields(s.fields).WithField("op", "save")
	l.Debug("saved in store")
	return nil
}

func (s *PostgresStore) insert(tx *sql.Tx, ctx context.Context, event eventflow.Event) (err error) {
	query := `INSERT INTO events (id, type, created_at, payload, mapper) VALUES ($1,$2,$3,$4,$5)`
	paymentID, ok := event.Columns[app.AggregatePayments].(string)
	if !ok {
		err = fmt.Errorf("expected column '%s' is not a string", app.AggregatePayments)
		return
	}
	result, err := s.db.Exec(query, paymentID, event.Type, event.Timestamp, event.Data, event.Mapper)
	if err != nil {
		return
	}
	if _, err = result.RowsAffected(); err != nil {
		return
	}
	return
}

func (s *PostgresStore) existsWithID(tx *sql.Tx, ctx context.Context, id string) (err error) {
	query := "SELECT id FROM events WHERE id=$1"
	row := s.db.QueryRowContext(ctx, query, id)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = app.ErrNoRows
			return
		}
		return
	}
	return
}

func (s *PostgresStore) Update(ctx context.Context, event eventflow.Event) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	paymentID, ok := event.Columns[app.AggregatePayments].(string)
	if !ok {
		err = fmt.Errorf("expected column '%s' is not a string", app.AggregatePayments)
		return
	}
	if err = s.existsWithID(tx, ctx, paymentID); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return
	}
	if err = s.insert(tx, ctx, event); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	l := log.FromContext(ctx)
	l = l.WithFields(s.fields).WithField("op", "update")
	l.Debug("updated in store")
	return nil
}

func (s *PostgresStore) Delete(ctx context.Context, event eventflow.Event) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	// at this point we are inserting event even if it can be duplicated
	// any additional deletion for given id should not affect aggregator anyway
	// todo refactor
	if err = s.insert(tx, ctx, event); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	l := log.FromContext(ctx)
	l = l.WithFields(s.fields).WithField("op", "delete")
	l.Debug("deleted in store")
	return nil
}
