package mysql

import (
	"database/sql"
	"go-service-demo/pkg/database"

	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	conn string
	db   *sql.DB
}

func NewMySql(conn string) database.ISqlDatabase {
	return &MySql{
		conn: conn,
	}
}

func (m *MySql) Connect() error {
	db, err := sql.Open("mysql", m.conn)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *MySql) Close() error {
	return m.db.Close()
}

func (m *MySql) Ping() error {
	return m.db.Ping()
}

func (m *MySql) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
