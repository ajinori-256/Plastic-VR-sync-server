package auth

import (
	"github.com/ajinori-256/Plastic-VR-sync-server/api"
)

func Login(loginRequest *grpc.LoginRequest, playerData *grpc.PlayerData) *grpc.LoginResponse {
	playerData.PlayerId = loginRequest.PlayerId
	return &grpc.LoginResponse{
		Error: &grpc.Error{
			Code:    grpc.ErrorCode_NONE,
			Message: "Success",
		},
	}
}
