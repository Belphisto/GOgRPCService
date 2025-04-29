package client

import (
	"context"
	"fmt"

	pb "github.com/Belphisto/GOgRPCService/proto"
)

// Лайк сообщения
func LikeMessage(reactionsClient pb.ReactionsServiceClient, messageID int32, username string) {
	_, err := reactionsClient.LikeMessage(context.Background(), &pb.LikeRequest{MessageId: messageID, Username: username})
	if err != nil {
		fmt.Printf("⚠ Ошибка лайка: %v\n", err)
		return
	}
	fmt.Printf("✅ Лайк от %s к сообщению #%d\n", username, messageID)
}

// Комментарий к сообщению
func CommentMessage(reactionsClient pb.ReactionsServiceClient, messageID int32, username, content string) {
	_, err := reactionsClient.CommentMessage(context.Background(), &pb.CommentRequest{MessageId: messageID, Username: username, Content: content})
	if err != nil {
		fmt.Printf("⚠ Ошибка комментария: %v\n", err)
		return
	}
	fmt.Printf("✅ Комментарий к сообщению #%d отправлен!\n", messageID)
}
