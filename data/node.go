package data

import (
	"context"
	"github.com/yoneyan/vm_mgr/node/vm"
	pb "github.com/yoneyan/vm_mgr/proto/proto-go"
	"log"
	"time"
)

func (s *server) StopNode(ctx context.Context, in *pb.NodeID) (*pb.Result, error) {
	log.Println("----StopNode----")
	vm.StopProcess()
	timer := time.NewTimer(time.Second * 1)
	<-timer.C
	log.Printf("Node End! ")
	return &pb.Result{Status: true}, nil
}
