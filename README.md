# GOgRPCService

Local Environment: 
install in go-build and go-env
go get google.golang.org/grpc
go get google.golang.org/protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --go_out=. --go-grpc_out=. proto/reverse.proto