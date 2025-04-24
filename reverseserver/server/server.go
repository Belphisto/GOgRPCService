package reverseserver

import (
	"context"
	"log"
	"net"

	pb "github.com/Belphisto/GOgRPCService/proto"
	"google.golang.org/grpc"
)

type reverseServer struct {
	pb.UnimplementedReverseServiceServer
}

func (s *reverseServer) ReverseString(ctx context.Context, req *pb.ReverseRequest) (*pb.ReverseResponse, error) {
	runes := []rune(req.Input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return &pb.ReverseResponse{Output: string(runes)}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при создании порта: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterReverseServiceServer(server, &reverseServer{})

	log.Println("Запуск сервера на порту 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
