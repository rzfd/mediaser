package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcServer "github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/pkg/pb"
)

func main() {
	log.Println("🔔 Starting Notification Microservice...")

	// Initialize notification service
	notificationService := grpcServer.NewMockNotificationService()

	// Create gRPC server
	lis, err := net.Listen("tcp", ":"+getEnv("GRPC_PORT", "9093"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	
	// Register notification service
	notificationGRPCServer := grpcServer.NewNotificationGRPCServer(notificationService)
	pb.RegisterNotificationServiceServer(grpcSrv, notificationGRPCServer)

	// Enable reflection for development
	reflection.Register(grpcSrv)

	// Graceful shutdown
	go func() {
		log.Printf("Notification Service listening on port %s", getEnv("GRPC_PORT", "9093"))
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Notification Service...")
	grpcSrv.GracefulStop()
	log.Println("Notification Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 