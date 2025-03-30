package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной goroutine
	go func() {
		log.Printf("Server listening at %v\n", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Failed to start the gRPC server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	log.Println("Server is shutting down...")

	// Graceful shutdown с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Остановка сервера с использованием контекста
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("Server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("Graceful shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
