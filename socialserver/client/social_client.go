package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Belphisto/GOgRPCService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func streamMessages(client pb.SocialServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.StreamFeed(ctx, &pb.StreamRequest{})
	if err != nil {
		log.Fatalf("Ошибка подключения к потоку: %v", err)
	}

	fmt.Println("\n🔄 Начинаем потоковое получение сообщений.") // Исправлено

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Ошибка получения сообщения: %v", err)
		}
		fmt.Printf("[%s]: %s\n", msg.Username, msg.Content)
	}
}

func main() {
	clientConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewSocialServiceClient(clientConn)

	fmt.Print("Введите имя пользователя: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	// Запускаем потоковое получение сообщений
	go streamMessages(client)

	for {
		fmt.Print("Введите сообщение (или 'exit' для выхода): ")
		scanner.Scan()
		content := scanner.Text()

		if content == "exit" {
			fmt.Println("Выход из клиента...")
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = client.SendMessage(ctx, &pb.MessageRequest{Username: username, Content: content})
		cancel()

		if err != nil {
			log.Fatalf("Ошибка отправки сообщения: %v", err)
		}
	}
}
