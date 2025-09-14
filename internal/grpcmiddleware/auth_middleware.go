package grpcmiddleware

import (
	"context"
	"log"

	gocache "github.com/patrickmn/go-cache"
	jwtentity "github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/entity/jwt"
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/utils"
	"google.golang.org/grpc"
)



type authMiddleware struct {
	cacheService *gocache.Cache
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println(info.FullMethod)
	tokenStr, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := am.cacheService.Get(tokenStr); ok {
		return nil, utils.UnauthenticatedResponse()
	}
	claims, err := jwtentity.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, err
	}
	ctx = claims.SetToContext(ctx)

	res, err := handler(ctx, req)
	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
