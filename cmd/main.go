package main

import (
	"io"
	"net"
	"os"
	"strconv"

	"github.com/go-kit/log"
	pb "github.com/numeris/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"

	"github.com/numeris/internal/config"
	server "github.com/numeris/internal/delivery/grpc"
	"github.com/numeris/internal/service"

	"gorm.io/driver/postgres"
)

func main() {
	cfg := config.LoadConfig()
	var logWriter io.Writer

	if cfg.LogFile == "" {
		logWriter = os.Stderr
	} else {
		logfile, _ := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer logfile.Close()
		logWriter = logfile
	}
	logger := log.NewLogfmtLogger(log.NewSyncWriter(logWriter))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	//Database configs
	dbCfg := cfg.Database
	invservice, err := service.NewInvoiceService(
		service.WithGorm(postgres.Open(dbCfg.ConnectionInfo()), &gorm.Config{}),
		service.WithLogMode(),
		service.WithInvoice(),
	)
	logger.Log("Database connection failure")
	must(err)
	defer invservice.Close()
	invservice.AutoMigrate()

	invsvr := server.NewInvoiceServer(invservice)
	opts := []grpc.ServerOption{}
	if cfg.Tls {
		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			must(err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	server := grpc.NewServer(opts...)

	// Register InvoiceServiceServer with the gRPC server
	pb.RegisterInvoiceServiceServer(server, invsvr)

	listen_port := ":" + strconv.Itoa(int(cfg.Port))
	lis, err := net.Listen("tcp", listen_port)
	if err != nil {
		logger.Log("Net Listening error")
		must(err)
	}
	err = server.Serve(lis)
	if err != nil {
		logger.Log("Error serving gRPC: %v", err)
		must(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
