package mysql

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	_ "github.com/go-sql-driver/mysql"                   // provides MySQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file" // provides filesystem source
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

const (
	dbDriverName            = "mysql"
	maxReconnectWaitingTime = 15 * time.Second
)

type DSN struct {
	User     string
	Password string
	Host     string
	Database string
}

func (dsn *DSN) String() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true", dsn.User, dsn.Password, dsn.Host, dsn.Database)
}

type Connector interface {
	Open(dsn DSN, maxConnections int) error
	MigrateUp(dsn DSN, migrationsProvider MigrationProvider) error
	Client() *sqlx.DB
	Close() error
}

type connector struct {
	db *sqlx.DB
}

func NewConnector() Connector {
	return &connector{}
}

func (c *connector) MigrateUp(dsn DSN, migrationsProvider MigrationProvider) error {
	db, err := openDb(dsn, 1)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = migrate.Exec(db.DB, dbDriverName, makeMigrationSource(migrationsProvider), migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to migrate")
	}

	return nil
}

func (c *connector) Open(dsn DSN, maxConnections int) error {
	var err error
	c.db, err = openDb(dsn, maxConnections)
	return errors.WithStack(err)
}

func (c *connector) Close() error {
	err := c.db.Close()
	return errors.Wrap(err, "failed to disconnect")
}

func (c *connector) Client() *sqlx.DB {
	return c.db
}

func openDb(dsn DSN, maxConnections int) (*sqlx.DB, error) {
	db, err := sqlx.Open(dbDriverName, dsn.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database")
	}

	db.SetMaxOpenConns(maxConnections)

	err = backoff.Retry(func() error {
		tryError := db.Ping()
		return tryError
	}, newExponentialBackOff())
	if err != nil {
		dbCloseErr := db.Close()
		if dbCloseErr != nil {
			err = errors.Wrap(err, dbCloseErr.Error())
		}
		return nil, errors.Wrapf(err, "failed to ping database")
	}
	return db, errors.WithStack(err)
}

func makeMigrationSource(migrationsProvider MigrationProvider) migrate.MigrationSource {
	return migrate.HttpFileSystemMigrationSource{FileSystem: migrationsProvider.GetDir()}
}

func newExponentialBackOff() *backoff.ExponentialBackOff {
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = maxReconnectWaitingTime
	return exponentialBackOff
}
