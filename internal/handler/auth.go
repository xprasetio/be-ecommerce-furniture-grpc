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

func (s *authHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	validationErros, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErros != nil {
		return &auth.RegisterResponse{
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

func (s *authHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationErros, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErros != nil {
		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationErros),
		}, nil
	}
	//Proses Login
	res, err := s.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *authHandler) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	validationErros, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErros != nil {
		return &auth.LogoutResponse{
			Base: utils.ValidationErrorResponse(validationErros),
		}, nil
	}
	res, err := s.authService.Logout(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *authHandler) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	validationErros, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErros != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse(validationErros),
		}, nil
	}
	res, err := s.authService.ChangePassword(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *authHandler) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	res, err := s.authService.GetProfile(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
