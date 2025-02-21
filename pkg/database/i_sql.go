package database

import "database/sql"

type ISqlDatabase interface {
	Ping() error
	Connect() error
	Close() error
	Query(query string) (*sql.Rows, error)
}
