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
	"github.com/Belphisto/GOgRPCService/socialserver/client/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("🚀 Запуск клиента социальной сети...")

	// Подключаемся к серверу gRPC
	clientConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer clientConn.Close()

	socialClient := pb.NewSocialServiceClient(clientConn)
	reactionsClient := pb.NewReactionsServiceClient(clientConn)

	// Читаем имя пользователя
	username := client.ReadUsername()

	for {
		// Показываем историю чата
		client.DisplayChatHistory(socialClient)

		// Читаем действие пользователя
		action := client.ReadUserAction()

		if action == "5" {
			fmt.Println("👋 Выход из клиента...")
			break
		}

		if action == "1" {
			// Отправка нового сообщения
			fmt.Print("✏ Введите сообщение: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			content := scanner.Text()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err = socialClient.SendMessage(ctx, &pb.MessageRequest{Username: username, Content: content})
			if err != nil {
				log.Fatalf("Ошибка отправки сообщения: %v", err)
			}

		} else {
			// Обработка ID сообщения
			fmt.Print("📩 Введите ID сообщения: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			messageID, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("⚠ Ошибка: введите корректный ID.")
				continue
			}

			if action == "2" {
				client.LikeMessage(reactionsClient, int32(messageID), username)
			} else if action == "3" {
				fmt.Print("💬 Введите комментарий: ")
				scanner.Scan()
				content := scanner.Text()
				client.CommentMessage(reactionsClient, int32(messageID), username, content)
			}
		}
	}
}
