package vm

import (
	"fmt"
	"github.com/mattn/go-pipeline"
	"github.com/vmmgr/node/db"
	"github.com/vmmgr/node/etc"
	"log"
	"strconv"
)

func AllVMStopForce() error {
	if vm, err := db.GetAllDBVM(); err != nil {
		return err
	} else {
		for _, d := range vm {
			if d.Status == 1 {
				if err := vmStop(d.ID); err != nil {
					log.Println("Error:VM Stop failed")
				}
			}
		}
	}
	return nil
}

func AllVMShutdown() error {
	if vm, err := db.GetAllDBVM(); err != nil {
		return err
	} else {
		for _, d := range vm {
			if d.Status == 1 {
				if err := vmShutdown(d.ID); err != nil {
					log.Println("Error:VM Stop failed")
				}
			}
		}
	}
	return nil
}

func vmStop(id int) error {
	if d, err := db.SearchDBVM(db.VM{ID: id}); err != nil {
		return fmt.Errorf("VM is not found")
	} else {
		if d.Status == 0 {
			log.Println("Power Off State")
		} else if d.Status == 1 {
			log.Println("Power On State")
		} else {
			log.Println("Error: Power State")
			return fmt.Errorf("Error: Power state ")
		}
	}

	//ps axf | grep test|grep qemu  | grep -v grep | awk '{print "kill -9 " $1}' | sudo sh
	log.Println("-----VMStop Command-----")
	if out, err := pipeline.CombinedOutput(
		[]string{"ps", "axf"},
		[]string{"grep", strconv.Itoa(id) + ".sock"},
		[]string{"grep", "qemu"},
		[]string{"grep", "-v", "grep"},
		[]string{"awk", "{print \"kill -9 \" $1}"},
		[]string{"sudo", "sh"},
	); err != nil {
		log.Println("already stop")
	} else {
		log.Printf("%s \n", out)
	}

	if r := db.UpdateDBVM(db.VM{ID: id, Status: 0}); r.Error != nil {
		return fmt.Errorf("Error: DB Error. but vm stopped ")
	} else {
		return nil
	}
}

func vmShutdown(id int) error {
	log.Printf("Shutdown VM: %d\n", id)
	if err := runQEMUMonitorCmd("system_powerdown", etc.SocketConnectionPath(id)); err != nil {
		log.Println("Error: Shutdown Error!!")
		return fmt.Errorf("Error: Shutdown Error ")
	}
	return nil
}
