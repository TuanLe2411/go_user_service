package mock

import "go-service-demo/pkg/database"

type DatabaseMock struct{}

func NewDatabaseMock() database.Database {
	return &DatabaseMock{}
}

func (d *DatabaseMock) Close() error {
	return nil
}

func (d *DatabaseMock) Connect() error {
	return nil
}

func (d *DatabaseMock) Ping() error {
	return nil
}

func (d *DatabaseMock) Query(string) (interface{}, error) {
	return "Select result", nil
}

func (d *DatabaseMock) Execute(string) (interface{}, error) {
	return "Execute result", nil
}
