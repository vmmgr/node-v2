package vm

import (
	"github.com/mattn/go-pipeline"
	"github.com/vmmgr/node/db"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func runQEMUMonitorCmd(command, socket string) error {
	//echo "system_powerdown" | socat - unix-connect:/var/run/someapp/vm.sock
	if out, err := pipeline.Output(
		[]string{"echo", command},
		[]string{"sudo", "socat", "-", socket},
	); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println(string(out))
		return nil
	}
}
func runQEMUCmd(command []string) error {
	log.Println("----CommandRun")
	//cmd = append(cmd,"-") //Intel VT-d support enable
	cmd := exec.Command("qemu-system-x86_64", command...)
	id, err := strconv.Atoi(command[2])
	if err != nil {
		return err
	}

	if result := db.UpdateDBVM(db.VM{ID: id, Status: 1}); result.Error != nil {
		return err
	}

	go func() {
		cmd.Start()
		log.Println("--------------------------------")
		log.Println("VMName: "+command[2]+" VMStart  Time: ", time.Now().Format("2006-01-02 15:04:05"))
		log.Println("--------------------------------")
		cmd.Wait()
		db.UpdateDBVM(db.VM{ID: id, Status: 1})
		log.Println("--------------------------------")
		log.Println("VMName: "+command[2]+" VMStop   Time: ", time.Now().Format("2006-01-02 15:04:05"))
		log.Println("--------------------------------")
	}()
	return nil
}
