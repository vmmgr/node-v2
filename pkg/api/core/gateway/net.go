package gateway

import (
	"context"
	"fmt"
	"github.com/vmmgr/node/db"
	"github.com/vmmgr/node/net"
	pb "github.com/vmmgr/node/proto/proto-go"
	"log"
	"strconv"
	"strings"
)

func (s *server) AddNet(ctx context.Context, in *pb.NetData) (*pb.Result, error) {
	log.Println("----------AddNet-----")
	log.Printf("Receive ID       : %v", in.GetID())
	log.Printf("Receive Name     : %v", in.GetName())
	log.Printf("Receive GroupID  : %v", in.GetGroupID())
	log.Printf("Receive VLAN     : %v", in.GetVLAN())
	log.Printf("Receive Option   : %v", in.GetOption())

	if result := net.AddNet(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) DeleteNet(ctx context.Context, in *pb.NetData) (*pb.Result, error) {
	log.Println("----------DeleteVM-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result := net.DeleteNet(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) UpdateNet(ctx context.Context, in *pb.NetData) (*pb.Result, error) {
	log.Println("----------UpdateNet-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result := net.UpdateNet(in); result.Err != nil {
		return &pb.Result{Status: false, Info: fmt.Sprint(result.Err)}, nil
	} else {
		return &pb.Result{Status: true, Info: "OK"}, nil
	}
}

func (s *server) GetNet(ctx context.Context, in *pb.NetData) (*pb.NetData, error) {
	log.Println("----------GetNet-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result, err := db.SearchDBNet(db.Net{ID: int(in.GetID())}); err != nil {
		return &pb.NetData{}, err
	} else {
		var gid []uint64
		for _, a := range strings.Split(result.GroupID, ",") {
			tmp, _ := strconv.Atoi(a)
			gid = append(gid, uint64(tmp))
		}
		return &pb.NetData{
			ID:      uint64(result.ID),
			Name:    result.Name,
			GroupID: gid,
			VLAN:    uint32(result.VLAN),
			Lock:    result.Lock == 1,
		}, nil
	}
}

func (s *server) GetAllNet(_ *pb.Null, stream pb.Node_GetAllNetServer) error {
	log.Println("----GetAllNet----")
	log.Printf("Receive GetAllVM")

	if result, err := db.GetAllDBNet(); err != nil {
		return err
	} else {
		log.Println(result)
		for _, data := range result {
			var gid []uint64
			for _, a := range strings.Split(data.GroupID, ",") {
				tmp, _ := strconv.Atoi(a)
				gid = append(gid, uint64(tmp))
			}
			if err := stream.Send(&pb.NetData{
				ID:      uint64(data.ID),
				Name:    data.Name,
				GroupID: gid,
				VLAN:    uint32(data.VLAN),
				Lock:    data.Lock == 1,
			}); err != nil {
				return err
			}
		}
		return nil
	}
}
