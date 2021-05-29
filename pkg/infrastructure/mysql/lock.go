package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
)

const timeoutInSeconds = 5

var (
	ErrLockTimeout     = errors.New("timeout is reached for lock")
	ErrLockNotFound    = errors.New("lock not found")
	ErrLockNotAcquired = errors.New("lock not acquired")
)

func NewLock(client Client, lockName string) Lock {
	return Lock{
		client:           client,
		lockName:         lockName,
		timeoutInSeconds: timeoutInSeconds,
	}
}

type Lock struct {
	client           Client
	lockName         string
	timeoutInSeconds int
}

func (l *Lock) Lock() error {
	const sqlQuery = `SELECT GET_LOCK(SUBSTRING(CONCAT(?, '.', DATABASE()), 1, 64), ?)`
	var result int
	err := l.client.Get(&result, sqlQuery, l.lockName, l.timeoutInSeconds)
	if result == 0 && err == nil {
		return ErrLockTimeout
	}
	return errors.WithStack(err)
}

func (l *Lock) Unlock() error {
	const sqlQuery = `SELECT RELEASE_LOCK(SUBSTRING(CONCAT(?, '.', DATABASE()), 1, 64), ?)`
	var result sql.NullInt32
	err := l.client.Get(&result, sqlQuery, l.lockName)
	if err == nil {
		if !result.Valid {
			return ErrLockNotFound
		}
		if result.Int32 == 0 {
			return ErrLockNotAcquired
		}
	}
	return errors.WithStack(err)
}
