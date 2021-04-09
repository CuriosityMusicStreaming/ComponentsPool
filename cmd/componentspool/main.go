package main

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/server"
	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
)

var appID = "UNKNOWN"

func main() {
	logger.SetFormatter(&logger.JSONFormatter{})

	config, err := parseEnv()
	if err != nil {
		logger.Fatal(err)
	}

	err = runService(config)
	if err == server.ErrStopped {
		logger.Info("service is successfully stopped")
	} else if err != nil {
		logger.Fatal(err)
	}
}

func runService(config *config) error {
	dsn := mysql.DSN{
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		Host:     config.DatabaseHost,
		Database: config.DatabaseName,
	}
	connector := mysql.NewConnector()

	err := connector.Open(dsn, config.MaxDatabaseConnections)
	if err != nil {
		return err
	}

	defer connector.Close()

	return nil
}
