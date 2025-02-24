package database

import "database/sql"

type IDatabase interface {
	Ping() error
	Connect() error
	Close() error
	Query(query string) (*sql.Rows, error)
}
