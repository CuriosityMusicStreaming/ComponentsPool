package main

import (
	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/amqp"
	jsonlog "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/logger"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/server"
	_ "github.com/go-sql-driver/mysql"
	stdlog "log"
)

var appID = "UNKNOWN"

func main() {
	logger, err := initLogger()
	if err != nil {
		stdlog.Fatal("failed to initialize logger")
	}

	config, err := parseEnv()
	if err != nil {
		logger.FatalError(err)
	}

	err = runService(config, logger)
	if err == server.ErrStopped {
		logger.Info("service is successfully stopped")
	} else if err != nil {
		logger.FatalError(err)
	}
}

func runService(config *config, logger log.MainLogger) error {
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

	_ = amqp.NewAMQPConnection(&amqp.Config{
		User:     config.AMQPUser,
		Password: config.AMQPPassword,
		Host:     config.AMQPHost,
	}, logger)

	return nil
}

func initLogger() (log.MainLogger, error) {
	return jsonlog.NewLogger(&jsonlog.Config{AppName: appID}), nil
}
