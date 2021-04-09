package mysql

import "net/http"

type MigrationProvider interface {
	GetDir() http.FileSystem
}
