package vm

import (
	pb "github.com/vmmgr/node/proto/proto-go"
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
	if err := deleteVMProcess(int(data.GetID())); err != nil {
		return result{Err: err}
	}
	return result{Info: "OK", Err: nil}
}

func UpdateVM(data *pb.VMData) result {

	//if cmdResult := deleteVMCmd(storage{path: path}); cmdResult.Err != nil {
	//	return result{Err: fmt.Errorf("Error: delete failed !! ")}
	//}
	//
	//if deleteResult := db.DeleteDBVM(db.VM{ID: int(data.GetID())}); deleteResult.Error != nil {
	//	return result{Err: fmt.Errorf("Error: delete db error !! ")}
	//}
	return result{Info: "not function", Err: nil}
}
