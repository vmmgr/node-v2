package data

import (
	"context"
	"fmt"
	"github.com/vmmgr/node/db"
	pb "github.com/vmmgr/node/proto/proto-go"
	"github.com/vmmgr/node/vm"
	"log"
	"strconv"
	"strings"
)

func (s *server) AddVM(ctx context.Context, in *pb.VMData) (*pb.Result, error) {
	log.Println("----------AddVM-----")
	log.Printf("Receive Name     : %v", in.GetName())
	log.Printf("Receive GroupID  : %v", in.GetGroupID())
	log.Printf("Receive CPU      : %v", in.GetCPU())
	log.Printf("Receive Mem      : %v", in.GetMem())
	log.Printf("Receive Boot     : %v", in.GetBoot())
	log.Printf("Receive ISO      : %v", in.GetISO())
	log.Printf("Receive Storage  : %v", in.GetStorage())
	log.Printf("Receive Net      : %v", in.GetNet())
	log.Printf("Receive PCI      : %v", in.GetPCIData())
	log.Printf("Receive Auto     : %v", in.GetAutoStart())

	if result := vm.AddVM(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) DeleteVM(ctx context.Context, in *pb.VMData) (*pb.Result, error) {
	fmt.Println("----------DeleteVM-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result := vm.DeleteVM(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) UpdateVM(ctx context.Context, in *pb.VMData) (*pb.Result, error) {
	fmt.Println("----------UpdateVM-----")
	log.Printf("Receive ID: %v", in.GetID())

	if result := vm.UpdateVM(in); result.Err != nil {
		return &pb.Result{Status: false, Info: fmt.Sprint(result.Err)}, nil
	} else {
		return &pb.Result{Status: true, Info: "OK"}, nil
	}
}

func (s *server) GetVM(ctx context.Context, in *pb.VMData) (*pb.VMData, error) {
	fmt.Println("----------GetVM-----")
	log.Printf("Receive ID: %v", in.GetID())
	var storage []*pb.StorageData
	var net []*pb.NetData
	var pci []*pb.PCIData

	if data, err := db.SearchDBVM(db.VM{ID: int(in.GetID())}); err != nil {
		return &pb.VMData{}, err
	} else {
		for _, tmp := range strings.Split(data.Storage, ",") {
			value, _ := strconv.Atoi(tmp)
			storage = append(storage, &pb.StorageData{ID: int64(value)})
		}
		for _, tmp := range strings.Split(data.Net, ",") {
			value, _ := strconv.Atoi(tmp)
			net = append(net, &pb.NetData{ID: int64(value)})
		}
		//for _, tmp := range strings.Split(data.PCI, ",") {
		//	value, _ := strconv.Atoi(tmp)
		//	pci = append(pci, &pb.PCIData{ID: int64(value)})
		//}

		return &pb.VMData{
			ID:      int64(data.ID),
			Name:    data.Name,
			GroupID: int64(data.GroupID),
			CPU:     int32(data.CPU),
			Mem:     int32(data.Mem),
			Storage: storage,
			Net:     net,
			PCIData: pci,
		}, nil
	}
}

func (s *server) GetAllVM(_ *pb.Null, stream pb.Node_GetAllVMServer) error {
	log.Println("----GetAllVM----")
	log.Printf("Receive GetAllVM")
	var storage []*pb.StorageData
	var net []*pb.NetData
	var pci []*pb.PCIData

	if result, err := db.GetAllDBVM(); err != nil {
		return err
	} else {
		log.Println(result)
		for _, data := range result {
			for _, tmp := range strings.Split(data.Storage, ",") {
				value, _ := strconv.Atoi(tmp)
				storage = append(storage, &pb.StorageData{ID: int64(value)})
			}
			for _, tmp := range strings.Split(data.Net, ",") {
				value, _ := strconv.Atoi(tmp)
				net = append(net, &pb.NetData{ID: int64(value)})
			}
			//for _, tmp := range strings.Split(data.PCI, ",") {
			//	value, _ := strconv.Atoi(tmp)
			//	pci = append(pci, &pb.PCIData{ID: int64(value)})
			//}
			if err := stream.Send(&pb.VMData{
				ID:      int64(data.ID),
				Name:    data.Name,
				GroupID: int64(data.GroupID),
				CPU:     int32(data.CPU),
				Mem:     int32(data.Mem),
				Storage: storage,
				Net:     net,
				PCIData: pci,
			}); err != nil {
				return err
			}
		}
		return nil
	}
}
