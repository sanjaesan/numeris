package app

import (
	"github.com/numeris/internal/config"

	server "github.com/numeris/internal/delivery/grpc"
	"github.com/numeris/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	config config.Config
	logger logger.Logger
	svr    *grpc.Server
	invsvr server.InvoiceServer
}

func NewApp(config config.Config, logger logger.Logger, invsvr server.InvoiceServer) (*App, error) {
	opts := []grpc.ServerOption{}
	if config.Tls {
		creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	server := grpc.NewServer(opts...)

	// pb.RegisterInvoiceServiceServer(server, invsvr)

	return &App{
		config: config,
		logger: logger,
		svr:    server,
		invsvr: invsvr,
	}, nil
}
