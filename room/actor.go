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
	players := make(map[string]Player, room.Data.RoomConfig.MaxPlayers)
	for {
		//add subscriber
		for {
			newSubscriber := <-room.ClientChan
			if newSubscriber.ClientChan != nil {
				players[newSubscriber.Name] = Player{ClientChan: newSubscriber.ClientChan, ServerChan: newSubscriber.ServerChan}
			} else {
				break
			}
		}
		//update Player Data
		for _, id := range maps.Keys(players) {
			for {
				var playerMessage *grpc.ClientMessage = <-players[id].ClientChan
				switch playerMessage.Data.(type) {
				case *grpc.ClientMessage_JoinRoomRequest:
					players[id].ServerChan <- &grpc.ServerMessage{Data: &grpc.ServerMessage_JoinRoomResponse{
						JoinRoomResponse: processJoinRequest(playerMessage.GetJoinRoomRequest(), id, players),
					}}

				case *grpc.ClientMessage_LeaveRoomRequest:
					players[id].ServerChan <- &grpc.ServerMessage{Data: &grpc.ServerMessage_LeaveRoomResponse{}}
				case *grpc.ClientMessage_PlayerDataPush:
					transform := playerMessage.GetPlayerDataPush().Data.GetPlayerTransform()
					players[id].Data.Transform.Positon.X = transform.Positon.X
					players[id].Data.Transform.Positon.Y = transform.Positon.Y
					players[id].Data.Transform.Positon.Z = transform.Positon.Z

					players[id].Data.Transform.Rotation.X = transform.Rotation.X
					players[id].Data.Transform.Rotation.Y = transform.Rotation.Y
					players[id].Data.Transform.Rotation.Z = transform.Rotation.Z
					players[id].Data.Transform.Rotation.W = transform.Rotation.W

					players[id].Data.Transform.Size.X = transform.Size.X
					players[id].Data.Transform.Size.Y = transform.Size.Y
					players[id].Data.Transform.Size.Z = transform.Size.Z

				}
			}
		}

	}
}
func sendPlayersData(players map[string]Player) {
	var players_array []grpc.Player
	for _, keys := range maps.Keys(players) {
		players_array = append(players_array)
	}
	//grpc.ServerDataPush.Data
}
func serverNotification(players map[string]Player, serverMessage *grpc.ServerMessage) {
	for _, playerId := range maps.Keys(players) {
		players[playerId].ServerChan <- serverMessage
	}
}
func processJoinRequest(joinRequest *grpc.JoinRoomRequest, playerId string, players map[string]Player) *grpc.JoinRoomResponse {

	err := checkJoinPermission(joinRequest, playerId)
	if err.Code != grpc.ErrorCode_NONE {
		return &grpc.JoinRoomResponse{Error: err}
	}
	serverNotification(players, &grpc.ServerMessage{Data: &grpc.ServerMessage_JoinRoomNotification{JoinRoomNotification: &grpc.JoinRoomNotification{
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
