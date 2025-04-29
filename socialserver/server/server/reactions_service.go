package server

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/Belphisto/GOgRPCService/proto"
)

// Реализация сервиса реакций (лайки, комментарии)
type ReactionsServer struct {
	pb.UnimplementedReactionsServiceServer
	mu sync.Mutex
}

// Лайк сообщения
func (s *ReactionsServer) LikeMessage(ctx context.Context, req *pb.LikeRequest) (*pb.LikeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var targetMessage *Message
	for _, msg := range SocialServerInstance.messages {
		if msg.ID == req.MessageId {
			targetMessage = msg
			break
		}
	}

	if targetMessage == nil {
		log.Printf("❌ Ошибка: Сообщение #%d не найдено\n", req.MessageId)
		return nil, fmt.Errorf("сообщение #%d не найдено", req.MessageId)
	}

	targetMessage.Likes++
	log.Printf("❤️ Лайк от %s к сообщению #%d\n", req.Username, req.MessageId)

	return &pb.LikeResponse{Success: true, LikeCount: targetMessage.Likes}, nil
}

// Комментарий к сообщению
func (s *ReactionsServer) CommentMessage(ctx context.Context, req *pb.CommentRequest) (*pb.CommentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем существование сообщения
	if req.MessageId < 1 || int(req.MessageId) > len(SocialServerInstance.messages) {
		log.Printf("❌ Ошибка: Сообщение #%d не найдено\n", req.MessageId)
		return nil, fmt.Errorf("сообщение #%d не найдено", req.MessageId)
	}

	SocialServerInstance.messages[req.MessageId-1].Comments = append(
		SocialServerInstance.messages[req.MessageId-1].Comments,
		&pb.Comment{Username: req.Username, Content: req.Content},
	)

	log.Printf("💬 Комментарий от %s к сообщению #%d: %s\n", req.Username, req.MessageId, req.Content)
	return &pb.CommentResponse{Success: true, Comments: SocialServerInstance.messages[req.MessageId-1].Comments}, nil
}
