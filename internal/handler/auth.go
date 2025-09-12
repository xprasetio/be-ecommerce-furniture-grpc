package handler

import (
	"context"

	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/service"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/utils"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/auth"
)


type authHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.IAuthService
}

func (s *authHandler)  Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {	
	validationErros, err :=utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErros != nil {
		return  &auth.RegisterResponse{ 
			Base: utils.ValidationErrorResponse(validationErros),
		}, nil
	}
	//Proses Register
	res, err := s.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewAuthHandler(authService  service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}