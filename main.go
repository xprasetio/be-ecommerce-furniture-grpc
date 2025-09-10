package main

import (
	"context"
	"log"
	"net"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/handler"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/service"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pkg/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func errorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack() // stack trace
			err = status.Errorf(
				codes.Internal,
				"internal server error: %v", r,
			)
		}
	}()
	res, err := handler(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}
		return res , err
}


func main() {
	ctx := context.Background()
	godotenv.Load()
	lis, err :=net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}

	database.ConnectDB(ctx, os.Getenv("DB_URI"))

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			errorMiddleware,
		),
	)

	service.RegisterHelloWorldServiceServer(serv, serviceHandler)

	if os.Getenv("ENVIRONMENT") == "dev" { 
		reflection.Register(serv)
		log.Println(" Reflection is registered") // hanya di development
	}

	log.Println("Server started on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}