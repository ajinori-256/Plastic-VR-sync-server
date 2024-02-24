package main
import (
  "context"
  "net"
  "fmt"
  "log"
  "google.golang.org/grpc"
  "grpc/grpc"
)
type server struct{}

func (s *server) Start(stream grpc.App_StartServer) error {
  for {
    clientMessage,err := stream.Recv()
  }
}
func main(){
  port := 8000
  listener, err := net.Listen("tcp")
}
