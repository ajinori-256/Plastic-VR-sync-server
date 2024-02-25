package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ajinori-256/Plastic-VR-sync-server/api"
	"github.com/ajinori-256/Plastic-VR-sync-server/auth"
	"github.com/ajinori-256/Plastic-VR-sync-server/room"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Start(stream grpc.App_StartServer) error {
	var playerData grpc.PlayerData
	for {
		clientMessage, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving client message: %v", err)
			return err
		}
		var serverMessage *grpc.ServerMessage
		serverMessage = processClientMessage(clientMessage, &playerData)
		err = stream.Send(serverMessage)
		if err != nil {
			log.Printf("Error sending server message: %v", err)
			return err
		}
	}
}

func processClientMessage(clientMessage *grpc.ClientMessage, playerData *grpc.PlayerData) *grpc.ServerMessage {
	switch clientMessage.Data.(type) {
	case *grpc.ClientMessage_LoginRequest:
		return &grpc.ServerMessage{
			Data: &grpc.ServerMessage_LoginResponse{
				LoginResponse: auth.Login(clientMessage.GetLoginRequest(), playerData),
			},
		}
	case *grpc.ClientMessage_CreateRoomRequest:
		return &grpc.ServerMessage{
			Data: &grpc.ServerMessage_CreateRoomResponse{
				CreateRoomResponse: room.CreateRoom(clientMessage.GetCreateRoomRequest(), playerData),
			},
		}
	case *grpc.ClientMessage_JoinRoomRequest:
		return &grpc.ServerMessage{
			Data: &grpc.ServerMessage_JoinRoomResponse{
				JoinRoomResponse: room.JoinRoom(clientMessage.GetJoinRoomRequest(), playerData),
			},
		}
	case *grpc.ClientMessage_LeaveRoomRequest:
		return &grpc.ServerMessage{
			Data: &grpc.ServerMessage_LeaveRoomResponse{
				LeaveRoomResponse: room.LeaveRoom(clientMessage.GetCreateRoomRequest(), playerData),
			},
		}
	case *grpc.ClientMessage_PlayerDataPush:
		return &grpc.ServerMessage{
			Data: &grpc.ServerMessage_PlayerDataPushResponse{
				PlayerDataPushResponse: room.PlayerDataPush(clientMessage.GetPlayerDataPush(), playerData),
			},
		}
	}
	return nil

}
func main() {
	port := 8000
	listener, err := net.Listen("tcp")
}
