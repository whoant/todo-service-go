package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"todo-service/demogrpc/demo"
)

type server struct{}

func (server) GetItemLikes(ctx context.Context, req *demo.GetItemLikesReq) (*demo.ItemLikeResp, error) {
	return &demo.ItemLikeResp{
		Result: map[int32]int32{
			1: 1,
			2: 4,
			3: 6,
		},
	}, nil
}

func main() {
	address := "0.0.0.0:50052"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	fmt.Printf("Server is listening on %v ...", address)
	s := grpc.NewServer()
	demo.RegisterItemLikeServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
