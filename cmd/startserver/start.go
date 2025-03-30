package server

import (
	"fmt"
	"log"
	"net"

	"github.com/krekio/TagesTest/config"
	"github.com/krekio/TagesTest/internal/server"
	"github.com/krekio/TagesTest/internal/storage"
	pb "github.com/krekio/TagesTest/protos"
	"google.golang.org/grpc"
)

func Run() {
	cfg := config.NewDefaultConfig()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Server startup error: %v", err)
	}

	fileStorage, err := storage.NewFileStorage(cfg.Server.StoragePath)
	if err != nil {
		log.Fatalf("Failed to initialize the file storage: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterFileServiceServer(grpcServer, server.NewFileServiceServer(fileStorage))

	log.Printf("Server listening at %v\n", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the gRPC server: %v", err)
	}
}
