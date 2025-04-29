package server

import pb "github.com/Belphisto/GOgRPCService/proto"

// Message — структура для хранения сообщений
type Message struct {
	ID       int32
	Username string
	Content  string
	Likes    int32
	Comments []*pb.Comment
}
