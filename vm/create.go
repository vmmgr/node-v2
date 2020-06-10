package vm

import (
	"fmt"
	"github.com/vmmgr/node/db"
	pb "github.com/vmmgr/node/proto/proto-go"
	"strconv"
)

func createVM(data *pb.VMData) error {
	var storage, nic, pci string
	for _, a := range data.GetStorage() {
		storage += strconv.Itoa(int(a.ID)) + ","
	}
	for _, a := range data.GetNIC() {
		nic += strconv.Itoa(int(a.ID)) + ","
	}
	for _, a := range data.GetPCIData() {
		pci += strconv.Itoa(int(a.ID)) + ","
	}
	if r := db.AddDBVM(db.VM{
		GroupID:   int(data.GetGroupID()),
		Name:      data.GetName(),
		CPU:       int(data.GetCPU()),
		Mem:       int(data.GetMem()),
		Storage:   storage,
		NIC:       nic,
		PCI:       pci,
		Status:    0,
		AutoStart: true,
	}); r.Error != nil {
		return fmt.Errorf("DB Error!! ")
	}

	return nil
}
