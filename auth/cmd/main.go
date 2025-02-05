package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/AndreiAvdko/auth/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

// Get ...
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Println("==============================")
	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id:             int64(gofakeit.Number(0, 1000)),
			Name:           gofakeit.Name(),
			Email:          gofakeit.Email(),
			Password:       gofakeit.Password(true, true, true, true, true, 5),
			PassworConfirm: gofakeit.Password(true, true, true, true, true, 5),
			IsAdmin:        gofakeit.Bool(),
			CreatedAt:      timestamppb.New(gofakeit.Date()),
			UpdatedAt:      timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Println("==============================")
	log.Printf("User email: %s", req.GetName())
	log.Printf("User email: %s", req.GetEmail())
	log.Printf("User email: %s", req.GetPassword())
	log.Printf("User email: %s", req.GetPassworConfirm())
	log.Printf("User email: %t", req.GetIsAdmin())

	return &desc.CreateResponse{
		Id: (int64)(gofakeit.Number(0, 1000)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
