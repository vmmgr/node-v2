package vm

import (
	"fmt"
	"github.com/vmmgr/node/etc"
	"log"
)

func vmReset(id int) error {
	log.Printf("Reset VM: %d\n", id)
	if err := runQEMUMonitorCmd("system_reset", etc.SocketConnectionPath(id)); err != nil {
		log.Println("Error: Reset Error!!")
		return fmt.Errorf("Error: Reset Error ")
	}
	return nil
}

func vmPause(id int) error {
	log.Printf("Pause VM: %d\n", id)
	if err := runQEMUMonitorCmd("stop", etc.SocketConnectionPath(id)); err != nil {
		log.Println("Error: Pause Error!!")
		return fmt.Errorf("Error: Pause Error ")
	}
	return nil
}

func vmResume(id int) error {
	log.Printf("Resume VM: %d\n", id)
	if err := runQEMUMonitorCmd("cont", etc.SocketConnectionPath(id)); err != nil {
		log.Println("Error: Resume Error!!")
		return fmt.Errorf("Error: Resume Error ")
	}
	return nil
}
