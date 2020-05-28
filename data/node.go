package data

import (
	"context"
	//"github.com/vmmgr/node/vm"
	pb "github.com/vmmgr/node/proto/proto-go"
	"log"
	"time"
)

func (s *server) StopNode(_ context.Context, _ *pb.Null) (*pb.Result, error) {
	log.Println("----StopNode----")
	//vm.StopProcess()
	timer := time.NewTimer(time.Second * 1)
	<-timer.C
	log.Printf("Stop Node End! ")
	return &pb.Result{Status: true}, nil
}
