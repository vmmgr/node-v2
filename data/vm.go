package data

import (
	"context"
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"github.com/yoneyan/vm_mgr/node/manage"
	"github.com/yoneyan/vm_mgr/node/vm"
	pb "github.com/yoneyan/vm_mgr/proto/proto-go"
	"log"
)

func (s *server) CreateVM(ctx context.Context, in *pb.VMData) (*pb.Result, error) {
	fmt.Println("----------CreateVM-----")
	log.Printf("Receive VMID        : %v", in.GetOption().GetId())
	log.Printf("Receive Type        : %v", in.GetType())
	log.Printf("Receive name        : %v", in.GetVmname())
	log.Printf("Receive cpu         : %v", in.GetVcpu())
	log.Printf("Receive mem         : %v", in.GetVmem())
	log.Printf("Receive StoragePath : %v", in.GetOption().StoragePath)
	log.Printf("Receive Storage     : %v", in.GetStorage())
	log.Printf("Receive CDROM       : %v", in.GetCdrom())
	log.Printf("Receive vnc         : %v", in.GetOption().Vnc)
	log.Printf("Receive net         : %v", in.GetVnet())
	log.Printf("Receive autostart   : %v", in.GetOption().Autostart)
	var r vm.CreateVMInformation

	r.ID = int(in.GetOption().Id)
	r.Name = in.GetVmname()
	r.CPU = int(in.GetVcpu())
	r.Mem = int(in.GetVmem())
	r.Storage = in.GetStorage()
	r.CDROM = in.GetOption().CdromPath
	r.Net = in.GetVnet()
	r.VNC = int(in.GetOption().Vnc)
	r.AutoStart = in.GetOption().Autostart

	if in.GetType()/10 == 1 {
		//Disk copy process
		log.Printf("AutoVMProcess")
		go vm.CreateAutoVMProcess(in)
		return &pb.Result{Status: true, Info: "Add process.. wait...."}, nil
	}

	r.StoragePath = manage.StorageProcess(in)

	info, result := vm.CreateVMProcess(&r)
	if result {
		return &pb.Result{Status: true, Info: "ok!!"}, nil
	}
	return &pb.Result{Status: false, Info: info}, nil
}

func (s *server) DeleteVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	fmt.Println("----------DeleteVM-----")
	log.Printf("Receive VMID: %v", in.GetId())

	info, result := vm.DeleteVMProcess(int(in.GetId()))
	if result {
		fmt.Println("Delete success!!")
		return &pb.Result{Status: true, Info: "ok"}, nil
	} else {
		fmt.Println("Delete Failed....")
		return &pb.Result{Status: false, Info: info}, nil
	}
}

func (s *server) StartVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	fmt.Println("----------StartVM-----")
	log.Printf("Receive VMID: %v", in.GetId())
	info, result := vm.StartVMProcess(int(in.GetId()))
	if result {
		return &pb.Result{Status: true, Info: "ok"}, nil
	} else {
		return &pb.Result{Status: false, Info: info}, nil
	}
}

func (s *server) StopVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	fmt.Println("----------StopVM-----")
	log.Printf("Receive VMID: %v", in.GetId())
	info, result := vm.VMStop(int(in.GetId()))
	if result == false {
		fmt.Println("VMStop Error!!")
		return &pb.Result{Status: false, Info: info}, nil
	}

	return &pb.Result{Status: true, Info: "ok!!"}, nil
}

func (s *server) GetVM(ctx context.Context, in *pb.VMID) (*pb.VMData, error) {
	fmt.Println("----------GetVMID-----")
	log.Printf("Receive VMID: %v", in.GetId())
	result, err := db.VMDBGetData(int(in.GetId()))
	if err != nil {
		fmt.Println("Error!!")
		return &pb.VMData{}, fmt.Errorf("Not Found!!")

	}
	return &pb.VMData{
		Option: &pb.Option{
			StoragePath: result.StoragePath,
			Vnc:         int32(result.Vnc),
			Id:          int64(result.ID),
			Autostart:   result.AutoStart,
			Status:      int32(result.Status),
		},
		Vmname: result.Name,
		Vcpu:   int64(result.CPU),
		Vmem:   int64(result.Mem),
		Vnet:   result.Net,
	}, nil
}

func (s *server) GetVMName(ctx context.Context, in *pb.VMName) (*pb.VMData, error) {
	fmt.Println("----------GetVMName-----")
	log.Printf("Receive Name: %v", in.GetVmname())
	id, err := db.VMDBGetVMID(in.GetVmname())
	if err != nil {
		fmt.Println("NotFound VMID !!")
		return &pb.VMData{}, fmt.Errorf("Not Found VMID!!")
	}
	result, err := db.VMDBGetData(id)
	if err != nil {
		fmt.Println("Not Found!!")
		return &pb.VMData{}, fmt.Errorf("Not Found!!")

	}
	return &pb.VMData{
		Option: &pb.Option{
			StoragePath: result.StoragePath,
			Vnc:         int32(result.Vnc),
			Id:          int64(result.ID),
			Autostart:   result.AutoStart,
		},
		Vmname:  result.Name,
		Vcpu:    int64(result.CPU),
		Vmem:    int64(result.Mem),
		Vnet:    result.Net,
		Storage: result.Storage,
	}, nil
}

//func (s *server) GetAllVM(ctx context.Context, in *pb.VMID)  error {
//	log.Println("----GetAllVM----")
//	log.Printf("Receive GetAllVM")
//	fmt.Println(db.VMDBGetAll())
//	return nil
//}

func (s *server) GetAllVM(base *pb.Base, stream pb.Grpc_GetAllVMServer) error {
	log.Println("----GetAllVM----")
	log.Printf("Receive GetAllVM")
	fmt.Println(db.VMDBGetAll())
	result := db.VMDBGetAll()
	for _, a := range result {
		if err := stream.Send(&pb.VMData{Option: &pb.Option{
			Vnc: int32(a.Vnc), Id: int64(a.ID), Autostart: a.AutoStart, Status: int32(a.Status),
			StoragePath: a.StoragePath},
			Vmname: a.Name, Vcpu: int64(a.CPU), Vmem: int64(a.Mem), Vnet: a.Net, Storage: a.Storage}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ShutdownVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	log.Println("----ShutdownVM----")
	log.Printf("Receive ShutdownVM")
	data, err := db.VMDBGetData(int(in.GetId()))
	if err != nil {
		fmt.Println("ID Not Found !!")
		return &pb.Result{Status: false, Info: "ID Not Found!!"}, nil
	}
	if vm.VMShutdown(data.Name) {
		return &pb.Result{Status: true, Info: "ok"}, nil
	} else {
		return &pb.Result{Status: false, Info: "Error: Shutdown Error!!"}, nil
	}
}

func (s *server) ResetVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	log.Println("----RebootVM----")
	log.Printf("Receive RebootVM")
	data, err := db.VMDBGetData(int(in.GetId()))
	if err != nil {
		fmt.Println("ID Not Found !!")
		return &pb.Result{Status: false, Info: "ID Not Found!!"}, nil
	}
	if vm.VMRestart(data.Name) {
		return &pb.Result{Status: true, Info: "ok"}, nil
	} else {
		return &pb.Result{Status: false, Info: "Restart Error!!"}, nil
	}
}

func (s *server) PauseVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	log.Println("----PauseVM----")
	log.Printf("Receive PauseVM")
	data, err := db.VMDBGetData(int(in.GetId()))
	if err != nil {
		fmt.Println("ID Not Found !!")
		return &pb.Result{Status: false, Info: "ID Not Found!!"}, nil
	}
	if vm.VMPause(data.Name) {
		db.VMDBStatusUpdate(data.ID, 2)
		return &pb.Result{Status: true, Info: "ok"}, nil
	} else {
		return &pb.Result{Status: false, Info: "Pause VM Error!!"}, nil
	}
}

func (s *server) ResumeVM(ctx context.Context, in *pb.VMID) (*pb.Result, error) {
	log.Println("----ResumeVM----")
	log.Printf("Receive ResumeVM")
	data, err := db.VMDBGetData(int(in.GetId()))
	if err != nil {
		fmt.Println("ID Not Found !!")
		return &pb.Result{Status: false, Info: "ID Not Found!!"}, nil
	}
	if vm.VMResume(data.Name) {
		return &pb.Result{Status: true, Info: "ok"}, nil
	} else {
		return &pb.Result{Status: false, Info: "Resume VM Error!!"}, nil
	}
}
