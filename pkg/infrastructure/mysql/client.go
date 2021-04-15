package mysql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Client interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)

	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type Transaction interface {
	Client
	Commit() error
	Rollback() error
}

type TransactionalClient interface {
	Client
	BeginTransaction() (Transaction, error)
}

type transactionalClient struct {
	*sqlx.DB
}

func (t *transactionalClient) BeginTransaction() (Transaction, error) {
	return t.Beginx()
}
