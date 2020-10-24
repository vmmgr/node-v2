package gateway

import (
	"context"
	"fmt"
	"github.com/vmmgr/node/db"
	pb "github.com/vmmgr/node/proto/proto-go"
	"github.com/vmmgr/node/storage"
	"log"
	"time"
)

func addStorageExec(resultCh chan struct {
	info   string
	result bool
}, data *pb.StorageData, t *storage.Tmp) {
	defer func() { close(resultCh) }()

	r := t.AddStorage(data)
	log.Println(r)

	if r.Err != nil {
		resultCh <- struct {
			info   string
			result bool
		}{info: r.Info + " Error: " + fmt.Sprint(r.Err), result: false}
	} else {
		resultCh <- struct {
			info   string
			result bool
		}{info: "OK", result: true}
	}
	return
}

func (s *server) AddStorage(in *pb.StorageData, stream pb.Node_AddStorageServer) error {
	log.Println("----------AddStorage-----")
	log.Printf("Receive ID       : %v", in.GetID())
	log.Printf("Receive Name     : %v", in.GetName())
	log.Printf("Receive GroupID  : %v", in.GetGroupID())
	log.Printf("Receive Driver   : %v", in.GetDriver())
	log.Printf("Receive Mode     : %v", in.GetMode())
	log.Printf("Receive Path     : %v", in.GetPath())
	log.Printf("Receive MaxSize  : %v", in.GetMaxSize())

	resultCh := make(chan struct {
		info   string
		result bool
	})

	t := &storage.Tmp{}
	go addStorageExec(resultCh, in, t)

	for {
		timer := time.NewTimer(time.Second * 1)
		<-timer.C
		select {
		case d := <-resultCh:
			if err := stream.Send(&pb.Result{
				Info:   d.info,
				Status: d.result,
			}); err != nil {
				return err
			}
			return nil

		default:
			log.Println(t.Info)
			if err := stream.Send(&pb.Result{
				Info: t.Info,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *server) DeleteStorage(_ context.Context, in *pb.StorageData) (*pb.Result, error) {
	log.Println("----------DeleteVM-----")
	log.Printf("Receive VMID: %v", in.GetID())

	if r, _ := db.SearchDBStorage(db.Storage{ID: int(in.GetID() % 100000)}); r.GroupID != int(in.GroupID) {
		return &pb.Result{Status: false, Info: "NG"}, nil
	}
	if result := storage.DeleteStorage(in); result.Err != nil {
		return &pb.Result{Status: false, Info: result.Info + "ErrorLog: " + fmt.Sprint(result.Err)}, nil
	}
	return &pb.Result{Status: true, Info: "ok"}, nil
}

func (s *server) UpdateStorage(_ context.Context, in *pb.StorageData) (*pb.Result, error) {
	log.Println("----------UpdateStorage-----")
	log.Printf("Receive VMID: %v", in.GetID())

	if r, _ := db.SearchDBStorage(db.Storage{ID: int(in.GetID() % 100000)}); r.GroupID != int(in.GroupID) {
		return &pb.Result{Status: false, Info: "NG"}, nil
	}
	if result := storage.UpdateStorage(in); result.Err != nil {
		return &pb.Result{Status: false, Info: fmt.Sprint(result.Err)}, nil
	} else {
		return &pb.Result{Status: true, Info: "OK"}, nil
	}
}

func (s *server) GetStorage(_ context.Context, in *pb.StorageData) (*pb.StorageData, error) {
	log.Println("----------GetStorage-----")
	log.Printf("Receive VMID: %v", in.GetID())

	if result, err := db.SearchDBStorage(db.Storage{ID: int(in.GetID())}); err != nil {
		return &pb.StorageData{}, err
	} else {
		return &pb.StorageData{
			ID:      uint64(result.ID),
			Name:    result.Name,
			GroupID: uint64(result.GroupID),
			Driver:  uint32(result.Driver),
			Mode:    uint32(result.Mode),
			Path:    result.Path,
			MaxSize: uint64(result.MaxSize),
		}, nil
	}
}

func (s *server) GetAllStorage(_ *pb.Null, stream pb.Node_GetAllStorageServer) error {
	log.Println("----GetAllStorage----")
	log.Printf("Receive GetAllVM")

	if result, err := db.GetAllDBStorage(); err != nil {
		return err
	} else {
		log.Println(result)
		for _, data := range result {
			if err := stream.Send(&pb.StorageData{
				ID:      uint64(data.ID),
				Name:    data.Name,
				GroupID: uint64(data.GroupID),
				Driver:  uint32(data.Driver),
				Mode:    uint32(data.Mode),
				Path:    data.Path,
				MaxSize: uint64(data.MaxSize),
			}); err != nil {
				return err
			}
		}
		return nil
	}
}
