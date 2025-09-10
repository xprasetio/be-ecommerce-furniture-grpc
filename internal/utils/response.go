package utils

import "github.com/xprasetio/be-ecommerce-furniture-grpc.git/pb/common"


func SuccessResponse(message string) *common.BaseResponse{
	return &common.BaseResponse{
		StatusCode: 200,
		Message: message,
	}
}