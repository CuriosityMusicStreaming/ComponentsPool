package server

import (
	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerConfig struct {
	ServeAddress string
}

func NewGrpcServer(server *grpc.Server, config GrpcServerConfig, logger log.Logger) Server {
	return &grpcServer{
		baseServer: server,
		config:     config,
		logger:     logger,
	}
}

type grpcServer struct {
	baseServer *grpc.Server
	config     GrpcServerConfig
	logger     log.Logger
}

func (g *grpcServer) Serve() error {
	grpcListener, grpcErr := net.Listen("tcp", g.config.ServeAddress)
	if grpcErr != nil {
		return errors.Wrapf(grpcErr, "failed to listen port %s", g.config.ServeAddress)
	}

	g.logger.Info("GRPC Server started")
	grpcErr = g.baseServer.Serve(grpcListener)
	return errors.Wrap(grpcErr, "failed to serve GRPC")
}

func (g *grpcServer) Stop() error {
	g.baseServer.GracefulStop()
	return nil
}
