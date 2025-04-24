package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/Belphisto/GOgRPCService/proto"
	"google.golang.org/grpc"
)

type socialServer struct {
	pb.UnimplementedSocialServiceServer
	messages []*pb.MessageRequest
	mu       sync.Mutex
	streams  []pb.SocialService_StreamFeedServer // Список потоков для клиентов
}

func (s *socialServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	s.mu.Lock()
	s.messages = append(s.messages, req)
	for _, stream := range s.streams {
		stream.Send(req) // Отправка нового сообщения всем подписанным клиентам
	}
	s.mu.Unlock()

	log.Printf("Сообщение от %s: %s\n", req.Username, req.Content)
	return &pb.MessageResponse{Success: true}, nil
}

func (s *socialServer) StreamFeed(req *pb.StreamRequest, stream pb.SocialService_StreamFeedServer) error {
	s.mu.Lock()
	messagesCopy := s.messages
	s.streams = append(s.streams, stream) // Регистрируем нового клиента
	s.mu.Unlock()

	// Отправляем клиенту все существующие сообщения
	for _, msg := range messagesCopy {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	// Клиент остаётся подключённым и получает новые сообщения
	for {
		time.Sleep(1 * time.Second) // Ожидание новых сообщений
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Ошибка создания сервера: %v", err)
	}

	server := grpc.NewServer()
	socialSrv := &socialServer{}
	pb.RegisterSocialServiceServer(server, socialSrv)

	log.Println("Сервер социальной сети запущен на порту 50052...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
