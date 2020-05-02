package vm

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/etc"
)

func VMShutdown(name string) bool {
	fmt.Println("Shutdown VM: " + name)
	err := RunQEMUMonitor("system_powerdown", etc.SocketConnectionPath(name))
	if err != nil {
		fmt.Println("Error: Shutdown Error!!")
		return false
	}
	return true
}

func VMRestart(name string) bool {
	fmt.Println("Reset VM: " + name)
	err := RunQEMUMonitor("system_reset", etc.SocketConnectionPath(name))
	if err != nil {
		fmt.Println("Error: Restart Error!!")
		return false
	}
	return true
}

func VMPause(name string) bool {
	fmt.Println("Pause VM: " + name)
	err := RunQEMUMonitor("stop", etc.SocketConnectionPath(name))
	if err != nil {
		fmt.Println("Error: Pause Error!!")
		return false
	}
	return true
}

func VMResume(name string) bool {
	fmt.Println("Resume VM: " + name)
	err := RunQEMUMonitor("cont", etc.SocketConnectionPath(name))
	if err != nil {
		fmt.Println("Error: Resume Error!!")
		return false
	}
	return true
}
