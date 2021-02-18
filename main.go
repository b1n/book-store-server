package main

import (
	"context"
	"fmt"
	"github.com/b1n/proto-book-store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

func main() {
	s := &Service{}
	startServer(s)
}

type Service struct{}

func (s *Service) GetBook(_ context.Context, request *book_store.GetBookRequest) (*book_store.GetBookResponse, error) {
	response := &book_store.GetBookResponse{
		Id:   request.Id,
		Name: fmt.Sprintf("Test %d", request.Id),
	}
	return response, nil
}

func startServer(service *Service) {
	listener, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Printf("Can't listen TCP port: %s", os.Getenv("GRPC_PORT"))
		log.Println("Error: ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor))

	book_store.RegisterBookStoreServer(server, service)

	log.Printf("Starting gRPC server at: %s", os.Getenv("GRPC_PORT"))

	if err := server.Serve(listener); err != nil {
		log.Printf("Can't start gRPC server at: %s", os.Getenv("GRPC_PORT"))
	}
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	tokens, ok := md["access-token"]
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "PermissionDenied")
	}
	if len(tokens) <= 0 {
		return nil, status.Error(codes.PermissionDenied, "PermissionDenied")
	}
	if tokens[0] != os.Getenv("TOKEN") {
		return nil, status.Error(codes.PermissionDenied, "PermissionDenied")
	}

	return handler(ctx, req)
}
