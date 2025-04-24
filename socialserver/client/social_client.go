package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/Belphisto/GOgRPCService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Функция для отображения истории чата
func displayChatHistory(socialClient pb.SocialServiceClient) {
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

func main() {
	log.Println("🚀 Запуск клиента социальной сети...")

	// Подключаемся к серверу gRPC
	clientConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer clientConn.Close()

	socialClient := pb.NewSocialServiceClient(clientConn)
	reactionsClient := pb.NewReactionsServiceClient(clientConn)

	// Читаем имя пользователя
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	scanner.Scan()
	username := scanner.Text()

	for {
		// Перед каждым действием показываем историю чата
		displayChatHistory(socialClient)

		fmt.Print("\nВыберите действие:\n1 - Отправить сообщение\n2 - Лайкнуть сообщение\n3 - Комментировать сообщение\n4 - Просмотреть чат\n5 - Выход\n")
		scanner.Scan()
		action := scanner.Text()

		if action == "5" {
			fmt.Println("👋 Выход из клиента...")
			break
		}

		if action == "1" {
			// Отправка нового сообщения
			fmt.Print("✏ Введите сообщение: ")
			scanner.Scan()
			content := scanner.Text()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err = socialClient.SendMessage(ctx, &pb.MessageRequest{Username: username, Content: content})
			if err != nil {
				log.Fatalf("Ошибка отправки сообщения: %v", err)
			}

		} else {
			// Лайк или комментарий к сообщению
			fmt.Print("📩 Введите ID сообщения: ")
			scanner.Scan()
			messageID, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("⚠ Ошибка: введите корректный ID.")
				continue
			}

			if action == "2" {
				// Лайк сообщения
				_, err := reactionsClient.LikeMessage(context.Background(), &pb.LikeRequest{MessageId: int32(messageID), Username: username})
				if err != nil {
					log.Fatalf("Ошибка лайка сообщения: %v", err)
				}
				fmt.Printf("✅ Лайк от %s к сообщению #%d\n", username, messageID)

			} else if action == "3" {
				// Комментарий к сообщению
				fmt.Print("💬 Введите комментарий: ")
				scanner.Scan()
				content := scanner.Text()

				_, err := reactionsClient.CommentMessage(context.Background(), &pb.CommentRequest{MessageId: int32(messageID), Username: username, Content: content})
				if err != nil {
					log.Fatalf("Ошибка комментария: %v", err)
				}
				fmt.Printf("✅ Комментарий к сообщению #%d отправлен!\n", messageID)
			}
		}
	}
}
