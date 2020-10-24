package gateway

import (
	"context"
	"fmt"
	"github.com/vmmgr/node/db"
	"github.com/vmmgr/node/nic"
	pb "github.com/vmmgr/node/proto/proto-go"
	"log"
)

func (s *server) AddNIC(ctx context.Context, in *pb.NICData) (*pb.Result, error) {
	log.Println("----------AddNIC-----")
	log.Printf("Receive ID       : %v", in.GetID())
	log.Printf("Receive Name     : %v", in.GetName())
	log.Printf("Receive GroupID  : %v", in.GetGroupID())
	log.Printf("Receive NetID    : %v", in.GetNetID())
	log.Printf("Receive Driver   : %v", in.GetDriver())

	if result := nic.AddNIC(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) DeleteNIC(ctx context.Context, in *pb.NICData) (*pb.Result, error) {
	log.Println("----------DeleteVM-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result := nic.DeleteNIC(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) UpdateNIC(ctx context.Context, in *pb.NICData) (*pb.Result, error) {
	log.Println("----------UpdateNIC-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result := nic.UpdateNIC(in); result.Err != nil {
		return &pb.Result{Status: false, Info: fmt.Sprint(result.Err)}, nil
	} else {
		return &pb.Result{Status: true, Info: "OK"}, nil
	}
}

func (s *server) GetNIC(ctx context.Context, in *pb.NICData) (*pb.NICData, error) {
	log.Println("----------GetNIC-----")
	log.Printf("Receive ID: %v", in.GetID())

	if data, err := db.SearchDBNIC(db.NIC{ID: int(in.GetID())}); err != nil {
		return &pb.NICData{}, err
	} else {
		return &pb.NICData{
			ID:         uint64(data.ID),
			Name:       data.Name,
			GroupID:    uint64(data.GroupID),
			NetID:      uint64(data.GroupID),
			MacAddress: data.MacAddress,
			Driver:     uint32(data.Driver),
			Lock:       data.Lock == 1,
		}, nil
	}
}

func (s *server) GetAllNIC(_ *pb.Null, stream pb.Node_GetAllNICServer) error {
	log.Println("----GetAllNIC----")
	log.Printf("Receive GetAllVM")

	if result, err := db.GetAllDBNIC(); err != nil {
		return err
	} else {
		log.Println(result)
		for _, data := range result {
			if err := stream.Send(&pb.NICData{
				ID:         uint64(data.ID),
				Name:       data.Name,
				GroupID:    uint64(data.GroupID),
				NetID:      uint64(data.GroupID),
				MacAddress: data.MacAddress,
				Driver:     uint32(data.Driver),
				Lock:       data.Lock == 1,
			}); err != nil {
				return err
			}
		}
		return nil
	}
}
