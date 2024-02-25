package room

import (
	"log"

	grpc "github.com/ajinori-256/Plastic-VR-sync-server/api"
	"github.com/google/uuid"
)

type Room struct {
	ClientChan chan Subscriber
	ServerChan chan grpc.ServerMessage
	Data       grpc.Room
}

var rooms map[string]*Room

func CreateRoom(createRoomRequest *grpc.CreateRoomRequest, playerData *grpc.PlayerData) *grpc.CreateRoomResponse {
	if playerData.PlayerId == "" {
		return &grpc.CreateRoomResponse{
			Error: &grpc.Error{Code: grpc.ErrorCode_UNAUTHORIZED, Message: grpc.ErrorCode_UNAUTHORIZED.String()},
		}
	}
	var room Room
	room_id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error %v\n", err)
		return &grpc.CreateRoomResponse{
			Error: &grpc.Error{Code:grpc.ErrorCode_INTERNAL_SERVER_ERROR,Message: grpc.ErrorCode_INTERNAL_SERVER_ERROR.String(),}
		}
	}
	room.Data.RoomId = room_id.String()
	room.Data.RoomOwner = playerData.PlayerId
	room.Data.RoomConfig = createRoomRequest.RoomConfig
	rooms[room.Data.RoomId] = &room
	go actor(&room)
	return &grpc.CreateRoomResponse{
		RoomId: room.Data.RoomId,
		Error: &grpc.Error{
			Code:    grpc.ErrorCode_NONE,
			Message: "Success",
		},
	}
}
