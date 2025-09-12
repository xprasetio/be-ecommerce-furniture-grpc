package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/handler"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/repository"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/service"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/auth"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pkg/database"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pkg/grpcmiddleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))

	authRepository := repository.NewAuthRepository(db)

	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)



	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println(" Reflection is registered") // hanya di development
	}

	log.Println("Server started on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
