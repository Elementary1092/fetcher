package main

import (
	"context"
	"fmt"
	"github.com/elem1092/fetcher/internal/adapters/api"
	"github.com/elem1092/fetcher/internal/adapters/db"
	"github.com/elem1092/fetcher/internal/config"
	fetch "github.com/elem1092/fetcher/pkg/client/grpc"
	"github.com/elem1092/fetcher/pkg/client/postgre"
	"github.com/elem1092/fetcher/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Staring Fetcher service")

	cfg := config.GetConfiguration()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelFunc()

	dbClient, dbErr := postgre.NewClient(ctx, cfg.DBConfig)
	if dbErr != nil {
		logger.Fatalf("Failed to connect to the database due to: %v", dbErr)
	}
	logger.Info("Connected to the database")

	storage := db.NewPostgreSQLStorage(dbClient, logger)

	service := api.NewService(logger, storage, cfg.ServerConfig.BaseURL)

	var mutex sync.Mutex
	pages, err := strconv.ParseInt(cfg.ServerConfig.Pages, 10, 32)
	if err != nil {
		logger.Warnf("Unable to parse pages amount from the configuration. Defaulting to 5")
		pages = 5
	}

	logger.Info("Starting server and listener")
	server := api.NewServer(service, logger, mutex, int32(pages))
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.ServerConfig.Address, cfg.ServerConfig.Port))
	if err != nil {
		logger.Fatalf("Unable to start server due to: %v", err)
	}

	logger.Infof("Initializing gRPC server on %s:%s", cfg.ServerConfig.Address, cfg.ServerConfig.Port)
	grpcServer := grpc.NewServer()
	fetch.RegisterFetchServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Fatalf("Unable to serve the server due to: %v", grpcServer.Serve(listener))
}
