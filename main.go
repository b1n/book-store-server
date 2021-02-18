package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/b1n/proto-book-store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
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
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("Can't listen TCP port: %s", "8080")
		log.Println("Error: ", err)
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor))

	book_store.RegisterBookStoreServer(server, service)

	log.Printf("Starting gRPC server at: %s", "8080")

	if err := server.Serve(listener); err != nil {
		log.Printf("Can't start gRPC server at: %s", "8080")
	}
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if md["access-token"][0] != "our_super-mega-secret_token" {
		return nil, errors.New("auth error")
	}
	return handler(ctx, req)
}
