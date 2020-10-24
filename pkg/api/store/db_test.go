package store

import (
	"fmt"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	InitCreateDB()
}

func TestCreateDatabase(t *testing.T) {
	AddDBVM(VM{
		Name: "TestVM1",
	})
	AddDBVM(VM{
		Name: "TestVM2",
	})
	AddDBVM(VM{
		GroupID: 10,
		Name:    "TestVM3",
		CPU:     2,
		Mem:     2048,
		//Storage:   []int{1, 2},
		//net:       []int{1},
		Status:    0,
		AutoStart: false,
	})
	AddDBStorage(Storage{
		ID:      1,
		GroupID: 10,
		Name:    "TestStorage1",
		Driver:  0,
		Type:    0,
		MaxSize: 1024,
		Lock:    0,
	})
	AddDBStorage(Storage{
		ID:      1,
		GroupID: 10,
		Name:    "TestStorage2",
		Driver:  0,
		Type:    0,
		MaxSize: 2048,
		Lock:    0,
	})
	//AddDBNet(Net{
	//	ID: 0,
	//	GroupID: []int{1},
	//Name:   "TestNetwork1",
	//Driver: 1,
	//Vlan:   190,
	//Status: 0,
	//})

	fmt.Println("=====VM=====")
	fmt.Println(GetAllDBVM())
	fmt.Println("=====Storage=====")
	fmt.Println(GetAllDBStorage())
	fmt.Println("=====Network=====")
	fmt.Println(GetAllDBNet())
	fmt.Println("==========")
}
