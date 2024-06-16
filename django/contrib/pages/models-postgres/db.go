// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models_postgres

import (
	"context"
	"database/sql"

	"github.com/Nigel2392/django/contrib/pages/models"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) models.Querier {
	return &Queries{
		db: tx,
	}
}
