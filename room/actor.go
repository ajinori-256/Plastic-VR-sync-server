package room

import (
	"log"

	grpc "github.com/ajinori-256/Plastic-VR-sync-server/api"
	"golang.org/x/exp/maps"
)

type Subscriber struct {
	ClientChan chan *grpc.ClientMessage
	ServerChan chan *grpc.ServerMessage
	Name       string
}
type Player struct {
	ClientChan chan *grpc.ClientMessage
	ServerChan chan *grpc.ServerMessage
	Data       grpc.PlayerData
}

func actor(room *Room) {
	var players grpc.ServerDataPush
	var clientChans map[string]chan *grpc.ClientMessage
	var serverChans map[string]chan *grpc.ServerMessage
	for {
		//add subscriber
		for {
			newSubscriber := <-room.ClientChan
			if newSubscriber.ClientChan != nil {
				clientChans[newSubscriber.Name] = newSubscriber.ClientChan
				serverChans[newSubscriber.Name] = newSubscriber.ServerChan
				players.Data[newSubscriber.Name] = new(grpc.PlayerData)
			} else {
				break
			}
		}
		//Process ClientMessage
		for _, id := range maps.Keys(players.Data) {
			for {
				var playerMessage *grpc.ClientMessage = <-clientChans[id]
				if playerMessage == nil {
					break
				}
				switch playerMessage.Data.(type) {
				case *grpc.ClientMessage_JoinRoomRequest:
					serverChans[id] <- &grpc.ServerMessage{Data: &grpc.ServerMessage_JoinRoomResponse{
						JoinRoomResponse: processJoinRequest(playerMessage.GetJoinRoomRequest(), id, &players, clientChans, serverChans),
					}}
				case *grpc.ClientMessage_LeaveRoomRequest:
					serverChans[id] <- &grpc.ServerMessage{Data: &grpc.ServerMessage_LeaveRoomResponse{}}
				case *grpc.ClientMessage_PlayerDataPush:
					transform := playerMessage.GetPlayerDataPush().Data.GetTransform()
					players.Data[id].Transform = transform
				}
			}
		}
		serverNotification(serverChans, &grpc.ServerMessage{Data: &grpc.ServerMessage_ServerDataPush{ServerDataPush: &players}})
	}
}
func serverNotification(serverChans map[string]chan *grpc.ServerMessage, serverMessage *grpc.ServerMessage) {
	for _, playerId := range maps.Keys(serverChans) {
		serverChans[playerId] <- serverMessage
	}
}
func processJoinRequest(joinRequest *grpc.JoinRoomRequest, playerId string, players *grpc.ServerDataPush, clientChans map[string]chan *grpc.ClientMessage, serverChans map[string]chan *grpc.ServerMessage) *grpc.JoinRoomResponse {
	err := checkJoinPermission(joinRequest, playerId)
	if err.Code != grpc.ErrorCode_NONE {
		return &grpc.JoinRoomResponse{Error: err}
	}
	serverNotification(serverChans, &grpc.ServerMessage{Data: &grpc.ServerMessage_JoinRoomNotification{JoinRoomNotification: &grpc.JoinRoomNotification{
		RoomId:   joinRequest.RoomId,
		PlayerId: playerId,
	},
	}})
	return &grpc.JoinRoomResponse{
		RoomId: joinRequest.RoomId,
		Error: &grpc.Error{
			Code:    grpc.ErrorCode_NONE,
			Message: grpc.ErrorCode_NONE.String(),
		},
	}
}
func checkJoinPermission(playerMessage *grpc.JoinRoomRequest, playerId string) *grpc.Error {
	//place holder
	return &grpc.Error{Code: grpc.ErrorCode_NONE, Message: grpc.ErrorCode_NONE.String()}
}
