package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err :=net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	serv := grpc.NewServer()

	log.Println("Server started on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}