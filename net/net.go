package net

import (
	"fmt"
	"github.com/vmmgr/node/db"
	pb "github.com/vmmgr/node/proto/proto-go"
	"strconv"
	"strings"
)

type result struct {
	Info string
	Err  error
}

func AddNet(data *pb.NetData) result {
	var groupID []string

	for _, id := range data.GetGroupID() {
		groupID = append(groupID, strconv.Itoa(int(id)))
	}
	strings.Join(groupID, ",")

	if r := db.AddDBNet(db.Net{
		GroupID: strings.Join(groupID, ","),
		Name:    data.GetName(),
		VLAN:    int(data.GetVLAN()),
	}); r.Error != nil {
		return result{Err: fmt.Errorf("Error: failed add database ")}
	}
	return result{Info: "OK", Err: nil}
}

func DeleteNet(data *pb.NetData) result {
	if dbData, err := db.SearchDBNet(db.Net{ID: int(data.GetID())}); err != nil {
		return result{Err: fmt.Errorf("Error: failed read db ")}
	} else if dbData.Lock != 0 {
		return result{Err: fmt.Errorf("Error: Locked Net !! ")}
	}

	if r := db.DeleteDBNet(db.Net{ID: int(data.GetID())}); r.Error != nil {
		return result{Err: fmt.Errorf("Error: failed add database ")}
	}
	return result{Info: "OK", Err: nil}
}

func UpdateNet(data *pb.NetData) result {
	dbData, err := db.SearchDBNet(db.Net{ID: int(data.GetID())})
	if err != nil {
		return result{Err: fmt.Errorf("Error: failed read db ")}
	} else if dbData.Lock != 0 {
		return result{Err: fmt.Errorf("Error: Locked Net !! ")}
	}
	if data.VLAN != 0 {
		dbData.VLAN = int(data.GetVLAN())
	}
	if data.Name != "" {
		dbData.Name = data.GetName()
	}
	if data.GetOption() != 0 {
		gid := strings.Split(dbData.GroupID, ",")
		//Add
		if data.GetOption() == 1 {
			for _, a := range data.GetGroupID() {
				same := false
				for _, b := range gid {
					if strconv.Itoa(int(a)) == b {
						same = true
					}
					if same == false {
						gid = append(gid, strconv.Itoa(int(a)))
					}
				}
			}
		}
		//Delete
		if data.GetOption() == 2 {
			var tmp []string
			for _, a := range gid {
				if a != strconv.Itoa(int(data.GetGroupID()[0])) {
					tmp = append(tmp, a)
				}
			}
			gid = tmp
		}
		dbData.GroupID = strings.Join(gid, ",")
	}
	return result{Info: "OK", Err: nil}
}
