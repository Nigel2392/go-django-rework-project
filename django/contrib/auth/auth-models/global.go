package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

var (
	queries Querier
)

func NewQueries(db *sql.DB) Querier {
	var q Querier
	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
		q = &MySQLQueries{
			db: db,
		}
	case *sqlite3.SQLiteDriver:
		q = &SQLiteQueries{
			db: db,
		}
	default:
		panic(fmt.Sprintf("unsupported driver: %T", db.Driver()))
	}

	queries = q
	return q
}

type MySQLQueries struct {
	db DBTX
}

func (q *MySQLQueries) WithTx(tx *sql.Tx) *MySQLQueries {
	return &MySQLQueries{
		db: tx,
	}
}

type SQLiteQueries struct {
	db DBTX
}

func (q *SQLiteQueries) WithTx(tx *sql.Tx) *SQLiteQueries {
	return &SQLiteQueries{
		db: tx,
	}
}
