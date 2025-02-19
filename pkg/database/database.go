package database

type Database interface {
	Ping() error
	Connect() error
	Close() error
	Query(string) (interface{}, error)
	Execute(string) (interface{}, error)
}
