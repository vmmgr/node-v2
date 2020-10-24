package gateway

import (
	"context"
	pb "github.com/vmmgr/node/proto/proto-go"
	"github.com/yoneyan/vm_mgr/node/vm"
	"log"
	"time"
)

func (s *server) StopNode(_ context.Context, _ *pb.Null) (*pb.Result, error) {
	log.Println("----StopNode----")
	//Stop Force
	vm.StopProcess()
	timer := time.NewTimer(time.Second * 1)
	<-timer.C
	log.Printf("Stop Node End! ")
	return &pb.Result{Status: true}, nil
}
