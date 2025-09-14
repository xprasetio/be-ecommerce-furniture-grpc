package jwt

import (
	"context"
	"strings"

	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ParseTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}
	bearerToken := md["authorization"]
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}
	if len(bearerToken) == 0 {
		return "", utils.UnauthenticatedResponse()
	}

	tokenSplit := strings.Split(bearerToken[0], " ")

	if len(tokenSplit) != 2 {
		return "", utils.UnauthenticatedResponse()
	}

	if tokenSplit[0] != "Bearer" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}

	return tokenSplit[1], nil
}
