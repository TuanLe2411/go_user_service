package database

import (
	"context"
	"database/sql"
)

type Database interface {
	Ping() error
	Connect() error
	Close() error
	QueryRows(query string, args ...any) (*sql.Rows, error, context.CancelFunc)
	QueryRow(query string, args ...any) (*sql.Row, error, context.CancelFunc)
	Exec(query string, args ...any) (sql.Result, error, context.CancelFunc)
}
