package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/entity"
	jwtentity "github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/entity/jwt"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/repository"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/utils"
	auth "github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache
}

func (s *authService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if request.Password != request.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password and Confirm Password not match"),
		}, nil
	}
	user, err := s.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exist"),
		}, nil
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  string(bcryptPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: request.FullName,
	}
	err = s.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User is Registered"),
	}, nil
}

func (s *authService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated") // authentication from grpc
		}
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtentity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecommerce-furniture",
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})
	secretKey := os.Getenv("JWT_SECRET_KEY")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login Success"),
		AccessToken: accessToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tokenClaims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	s.cacheService.Set(jwtToken, "", time.Duration(tokenClaims.ExpiresAt.Unix()-time.Now().Unix())*time.Second)
	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout Success"),
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	if request.NewPassword != request.NewPasswordConfirmation {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("New Password and Confirm Password not match"),
		}, nil
	}

	jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tokenClaims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}
	user, err := s.authRepository.GetUserByEmail(ctx, tokenClaims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated") // authentication from grpc
		}
		return nil, err
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(bcryptPassword)
	err = s.authRepository.UpdateUserPassword(ctx, user.Id, user.Password, tokenClaims.FullName)
	if err != nil {
		return nil, err
	}

	return &auth.ChangePasswordResponse{
		Base: utils.SuccessResponse("Change Password Success"),
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.GetProfileResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}

	return &auth.GetProfileResponse{
		Base:        utils.SuccessResponse("Get Profile Success"),
		UserId:      claims.Subject,
		Email:       claims.Email,
		FullName:    claims.FullName,
		RoleCode:    claims.Role,
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}
