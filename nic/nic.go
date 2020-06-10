package nic

import (
	"fmt"
	"github.com/vmmgr/node/db"
	pb "github.com/vmmgr/node/proto/proto-go"
)

type result struct {
	Info string
	Err  error
}

func AddNIC(data *pb.NICData) result {
	if r := db.AddDBNIC(db.NIC{
		GroupID:    int(data.GetGroupID()),
		NetID:      int(data.GetNetID()),
		Name:       data.GetName(),
		Driver:     int(data.GetDriver()),
		MacAddress: generateMacAddress(),
		Lock:       0,
	}); r.Error != nil {
		return result{Err: fmt.Errorf("Error: failed add database ")}
	}
	return result{Info: "OK", Err: nil}
}

func DeleteNIC(data *pb.NICData) result {
	if dbData, err := db.SearchDBNIC(db.NIC{ID: int(data.GetID())}); err != nil {
		return result{Err: fmt.Errorf("Error: failed read db ")}
	} else if dbData.Lock != 0 {
		return result{Err: fmt.Errorf("Error: Locked NIC !! ")}
	}

	if r := db.DeleteDBNIC(db.NIC{ID: int(data.GetID())}); r.Error != nil {
		return result{Err: fmt.Errorf("Error: failed add database ")}
	}
	return result{Info: "OK", Err: nil}
}

func UpdateNIC(data *pb.NICData) result {
	dbData, err := db.SearchDBNIC(db.NIC{ID: int(data.GetID())})
	if err != nil {
		return result{Err: fmt.Errorf("Error: failed read db ")}
	} else if dbData.Lock != 0 {
		return result{Err: fmt.Errorf("Error: Locked NIC !! ")}
	}
	if data.GetMacAddress() != "" {
		dbData.MacAddress = data.GetMacAddress()
	}
	if data.GetName() != "" {
		dbData.Name = data.GetName()
	}
	if data.GetDriver() != 0 {
		dbData.Name = data.GetName()
	}
	if r := db.UpdateDBNIC(dbData); r.Error != nil {
		return result{Info: "OK", Err: nil}
	} else {
		return result{Info: "Error: db update error", Err: r.Error}
	}
}
