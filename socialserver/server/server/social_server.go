package server

import (
	"context"
	"log"
	"sync"

	pb "github.com/Belphisto/GOgRPCService/proto"
)

var SocialServerInstance = &SocialServer{
	messages: []*Message{},
}

// SocialServer реализует сервис общения
type SocialServer struct {
	pb.UnimplementedSocialServiceServer
	messages []*Message
	mu       sync.Mutex
}

// Отправка сообщения
func (s *SocialServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	messageID := int32(len(s.messages) + 1)
	s.messages = append(s.messages, &Message{
		ID:       messageID,
		Username: req.Username,
		Content:  req.Content,
		Likes:    0,
		Comments: []*pb.Comment{},
	})

	log.Printf("📩 Сообщение #%d от %s: %s\n", messageID, req.Username, req.Content)
	return &pb.MessageResponse{Success: true}, nil
}

// Получение истории сообщений
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
