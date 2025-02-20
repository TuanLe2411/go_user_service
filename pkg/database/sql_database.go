package database

type SqlDatabase interface {
	Ping() error
	Connect() error
	Close() error
	Query(string) (interface{}, error)
	Execute(string) (interface{}, error)
}
