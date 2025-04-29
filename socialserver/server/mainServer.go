package main

import (
	"log"
	"net"

	pb "github.com/Belphisto/GOgRPCService/proto"
	"github.com/Belphisto/GOgRPCService/socialserver/server/server"
	"google.golang.org/grpc"
)

var socialServer = &server.SocialServer{}

func main() {
	listener, err := net.Listen("tcp", server.ServerPort)
	if err != nil {
		log.Fatalf("Ошибка создания сервера: %v", err)
	}

	serverInstance := grpc.NewServer()
	pb.RegisterSocialServiceServer(serverInstance, server.SocialServerInstance)  // ✅ Регистрируем SocialService
	pb.RegisterReactionsServiceServer(serverInstance, &server.ReactionsServer{}) // ✅ Регистрируем ReactionsService

	log.Println("🚀 Сервер запущен на порту", server.ServerPort)
	if err := serverInstance.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
