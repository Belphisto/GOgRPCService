package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Belphisto/GOgRPCService/proto"
)

// Отображение истории чата
func DisplayChatHistory(socialClient pb.SocialServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := socialClient.GetFeed(ctx, &pb.FeedRequest{})
	if err != nil {
		log.Fatalf("Ошибка получения истории чата: %v", err)
	}

	fmt.Println("\n📜 История чата:")
	for _, msg := range resp.Messages {
		fmt.Printf("\n📩 Сообщение #%d от %s: %s (❤️ %d лайков)\n", msg.MessageId, msg.Username, msg.Content, msg.LikeCount)
		for _, comment := range msg.Comments {
			fmt.Printf("💬 [%s]: %s\n", comment.Username, comment.Content)
		}
	}
}
