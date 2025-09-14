package utils

import (
	"github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		Message:    message,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "Validation Error",
		IsError:          true,
		ValidationErrors: validationErrors,
	}
}

func UnauthenticatedResponse() error {
	return status.Error(codes.Unauthenticated, "unauthenticated")
}
