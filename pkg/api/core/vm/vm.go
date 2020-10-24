package vm

import (
	pb "github.com/vmmgr/node/proto/proto-go"
	"strconv"
)

type result struct {
	Info string
	Err  error
}

func AddVM(data *pb.VMData) result {
	if err := createVM(data); err != nil {
		return result{Err: err}
	}

	return result{Info: "OK", Err: nil}
}

func DeleteVM(data *pb.VMData) result {
	if err := deleteVMProcess(int(data.GetID() % 100000)); err != nil {
		return result{Info: "NG", Err: err}
	}
	return result{Info: "OK", Err: nil}
}

func UpdateVM(data *pb.VMData) result {
	//|Model| 0:Status 1:VMUpdateNow 2:VMUpdateAfter
	var info string
	id := data.GetID() % 100000
	boot := data.GetBoot()
	//mode status
	if data.GetMode() == 0 {
		//Status
		//0:Start 1:Stop 2:Shutdown 3:Restart 4:HardReset 5:Pause 6:Resume 10:SnapShot
		if data.GetStatus() == 0 {
			if err := startVM(vmData{id: int(id), boot: int(boot)}); err != nil {
				return result{Info: "failed vm start", Err: err}
			}
			info = "VMID:" + strconv.Itoa(int(id)) + " Start"
		} else if data.GetStatus() == 1 {
			if err := vmStop(int(id)); err != nil {
				return result{Info: "failed vm stop", Err: err}
			}
			info = "VMID:" + strconv.Itoa(int(id)) + " Stop"
		} else if data.GetStatus() == 2 {
			if err := vmShutdown(int(id)); err != nil {
				return result{Info: "failed vm shutdown", Err: err}
			}
			info = "VMID:" + strconv.Itoa(int(id)) + " Shutdown"
		} else if data.GetStatus() == 3 {
			//if err := vmReset(int(id)); err != nil {
			//	return result{Info: "failed vm restart", Err: err}
			//}
			//info = "VMID:" + strconv.Itoa(int(id)) + " Restart"
		} else if data.GetStatus() == 4 {
			if err := vmReset(int(id)); err != nil {
				return result{Info: "failed vm restart", Err: err}
			}
			info = "VMID:" + strconv.Itoa(int(id)) + " Restart"
		} else if data.GetStatus() == 5 {
			if err := vmPause(int(id)); err != nil {
				return result{Info: "failed vm pause", Err: err}
			}
			info = "VMID:" + strconv.Itoa(int(id)) + " Pause"
		} else if data.GetStatus() == 6 {
			if err := vmResume(int(id)); err != nil {
				return result{Info: "failed vm resume", Err: err}
			}
			info = "VMID:" + strconv.Itoa(int(id)) + " Resume"
		}
	}

	//if cmdResult := deleteVMCmd(storage{path: path}); cmdResult.Err != nil {
	//	return result{Err: fmt.Errorf("Error: delete failed !! ")}
	//}
	//
	//if deleteResult := store.DeleteDBVM(store.VM{ID: int(vmData.GetID())}); deleteResult.Error != nil {
	//	return result{Err: fmt.Errorf("Error: delete store error !! ")}
	//}
	return result{Info: info, Err: nil}
}
