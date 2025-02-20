package sql_db_mock

import "go-service-demo/pkg/database"

type SqlDatabaseMock struct {
}

func NewSqlDatabaseMock() database.SqlDatabase {
	return &SqlDatabaseMock{}
}

func (d *SqlDatabaseMock) Close() error {
	return nil
}

func (d *SqlDatabaseMock) Connect() error {
	return nil
}

func (d *SqlDatabaseMock) Ping() error {
	return nil
}

func (d *SqlDatabaseMock) Query(string) (interface{}, error) {
	return "Select result", nil
}

func (d *SqlDatabaseMock) Execute(string) (interface{}, error) {
	return "Execute result", nil
}
