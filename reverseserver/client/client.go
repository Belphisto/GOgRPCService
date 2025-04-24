package reverseserver

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

func main() {
	// Создаём новое gRPC подключение с insecure credentials
	clientConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewReverseServiceClient(clientConn)

	// Ввод данных от пользователя
	fmt.Print("Введите строку для реверса: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput := scanner.Text()

	// Отправляем строку на сервер
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ReverseRequest{Input: userInput}
	res, err := client.ReverseString(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	fmt.Printf("Оригинальная строка: %s\nПеревернутая строка: %s\n", req.Input, res.Output)
}
