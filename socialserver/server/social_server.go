package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/Belphisto/GOgRPCService/proto"
	"google.golang.org/grpc"
)

type Message struct {
	ID       int32
	Username string
	Content  string
	Likes    int32
	Comments []*pb.Comment
}

type SocialServer struct {
	pb.UnimplementedSocialServiceServer
	messages []*Message
	mu       sync.Mutex
}

func (s *SocialServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	s.mu.Lock()
	messageID := int32(len(s.messages) + 1)
	s.messages = append(s.messages, &Message{
		ID:       messageID,
		Username: req.Username,
		Content:  req.Content,
		Likes:    0,
		Comments: []*pb.Comment{},
	})
	s.mu.Unlock()

	log.Printf("📩 Сообщение #%d от %s: %s\n", messageID, req.Username, req.Content)
	return &pb.MessageResponse{Success: true}, nil
}

func (s *SocialServer) GetFeed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var messages []*pb.MessageRequest
	for _, msg := range s.messages {
		messages = append(messages, &pb.MessageRequest{
			MessageId: msg.ID,
			Username:  msg.Username,
			Content:   msg.Content,
			LikeCount: msg.Likes,
			Comments:  msg.Comments,
		})
	}

	return &pb.FeedResponse{Messages: messages}, nil
}

// ✅ Добавляем `ReactionsServer`
type ReactionsServer struct {
	pb.UnimplementedReactionsServiceServer
	mu sync.Mutex
}

func (s *ReactionsServer) LikeMessage(ctx context.Context, req *pb.LikeRequest) (*pb.LikeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, msg := range socialServer.messages {
		if msg.ID == req.MessageId {
			msg.Likes++
		}
	}

	log.Printf("❤️ Лайк от %s к сообщению #%d\n", req.Username, req.MessageId)
	return &pb.LikeResponse{Success: true, LikeCount: socialServer.messages[req.MessageId-1].Likes}, nil
}

func (s *ReactionsServer) CommentMessage(ctx context.Context, req *pb.CommentRequest) (*pb.CommentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, msg := range socialServer.messages {
		if msg.ID == req.MessageId {
			msg.Comments = append(msg.Comments, &pb.Comment{Username: req.Username, Content: req.Content})
		}
	}

	log.Printf("💬 Комментарий от %s к сообщению #%d: %s\n", req.Username, req.MessageId, req.Content)
	return &pb.CommentResponse{Success: true, Comments: socialServer.messages[req.MessageId-1].Comments}, nil
}

var socialServer = &SocialServer{}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Ошибка создания сервера: %v", err)
	}

	serverInstance := grpc.NewServer()
	pb.RegisterSocialServiceServer(serverInstance, socialServer)          // ✅ Регистрируем SocialService
	pb.RegisterReactionsServiceServer(serverInstance, &ReactionsServer{}) // ✅ Регистрируем ReactionsService

	log.Println("🚀 Сервер запущен на порту 50052...")
	if err := serverInstance.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
