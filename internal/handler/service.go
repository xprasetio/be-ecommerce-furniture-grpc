package handler

import (
	"context"
	"fmt"

	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/utils"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/service"
)


type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (s *serviceHandler)  HelloWorld(ctx context.Context,request *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", request.Name),
		Base: utils.SuccessResponse("Success"),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}