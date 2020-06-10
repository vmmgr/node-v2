package data

import (
	pb "github.com/vmmgr/node/proto/proto-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50100"

type server struct {
	pb.UnimplementedNodeServer
}

func Server() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
