package vm

import (
	"fmt"
	"github.com/mattn/go-pipeline"
	"github.com/yoneyan/vm_mgr/node/db"
	"github.com/yoneyan/vm_mgr/node/etc"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func RunQEMUMonitor(command, socket string) error {
	//Example:
	//echo "system_powerdown" | socat - unix-connect:/var/run/someapp/vm.sock
	//

	out, err := pipeline.Output(
		[]string{"echo", command},
		[]string{"sudo", "socat", "-", socket},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(out))
	return nil
}
func RunQEMUCmd(command []string) error {
	fmt.Println("----CommandRun")
	//cmd = append(cmd,"-") //Intel VT-d support enable
	cmd := exec.Command("qemu-system-x86_64", command...)

	id, err := db.VMDBGetVMID(command[2])
	if err != nil {
		fmt.Println("DB Read Error")
	}
	db.VMDBStatusUpdate(id, 1)

	//go manage.VMLifeCheck(command[2], cmd)
	go func() {
		cmd.Start()
		fmt.Println("--------------------------------")
		fmt.Println("VMName: "+command[2]+" StartTime: ", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("--------------------------------")
		cmd.Wait()
		db.VMDBStatusUpdate(id, 0)
		fmt.Println("--------------------------------")
		fmt.Println("VMName: "+command[2]+" EndTime: ", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("--------------------------------")
	}()
	return nil
}

func CreateGenerateCmd(c *CreateVMInformation) []string {
	var diskindex int
	var cmd []string
	begin := []string{"-enable-kvm", "-name", c.Name, "-smp", strconv.Itoa(c.CPU), "-m",
		strconv.Itoa(c.Mem), "-monitor", etc.SocketGenerate(c.Name),
		"-vnc", "0.0.0.0:" + strconv.Itoa(c.VNC) + ",websocket=" + strconv.Itoa(c.VNC+7000)}

	cmd = append(cmd, begin...)

	if c.CDROM != "" {
		cmd = append(cmd, "-boot")
		cmd = append(cmd, "order=d")
		index, cdromdata := GenerateCDROMCmd(c.CDROM)
		for _, a := range cdromdata {
			cmd = append(cmd, a)
		}
		diskindex = index
	}

	if c.StoragePath != "" {
		if c.CDROM == "" {
			cmd = append(cmd, "-boot")
			cmd = append(cmd, "order=c")
			diskindex = 0
		}
		diskdata := GenerateDiskCmd(c.StoragePath, diskindex)
		for _, a := range diskdata {
			cmd = append(cmd, a)
		}
	}
	if len(c.Net) != 0 {
		netdata := GenerateNetworkCmd(c.Net)
		//add qemu network command
		for _, a := range netdata {
			cmd = append(cmd, a)
		}
	}

	fmt.Printf("GenerateCommand: ")
	fmt.Println(cmd)

	return cmd
}

func GenerateCDROMCmd(cdrom string) (int, []string) {
	fmt.Println("GenerateCDROMCommand")
	data := strings.Split(cdrom, ",")
	fmt.Println(data)

	var cmd []string
	var index int

	//-drive file=/images/cdrom.iso,index=2,media=cdrom
	for index, a := range data {
		cmd = append(cmd, "-drive")
		cmd = append(cmd, "file="+a+",index="+strconv.Itoa(index)+",media=cdrom")
	}
	return index, cmd
}

//Generate Disk Command
func GenerateDiskCmd(storagepath string, index int) []string {
	fmt.Println("GenerateDiskCommand")
	data := strings.Split(storagepath, ",")
	fmt.Println(data)

	index++
	diskmode := 0
	pathmode := 0
	var cmd []string

	for i, a := range data {
		if i%2 == 0 {
			m, _ := strconv.Atoi(a)
			diskmode = m / 10 //*0
			pathmode = m % 10 //0*
			fmt.Println("disk: " + strconv.Itoa(diskmode) + " path: " + strconv.Itoa(pathmode))
		} else {
			if diskmode == 0 {
				//0_ <= default Disk Mode
				//default disk mount diskmode
				//-drive file=/images/image2.raw,index=1,media=disk
				cmd = append(cmd, "-drive")
				if pathmode != 0 {
					basepath := etc.GetDiskPath(pathmode)
					if basepath == "" {
						fmt.Println("diskpath config error!!")
					}
					cmd = append(cmd, "file="+basepath+"/"+a+",index="+strconv.Itoa(index)+",media=disk")
				} else {
					cmd = append(cmd, "file="+a+",index="+strconv.Itoa(index)+",media=disk")
				}

				index++
			} else if diskmode == 1 {
				//disk mount diskmode is using virtio
				//-drive file=/images/image2.raw,index=1,media=disk,if=virtio
				cmd = append(cmd, "-drive")
				if pathmode != 0 {
					basepath := etc.GetDiskPath(pathmode)
					if basepath == "" {
						fmt.Println("diskpath config error!!")
					}
					cmd = append(cmd, "file="+basepath+"/"+a+",index="+strconv.Itoa(index)+",media=disk,if=virtio")
				} else {
					cmd = append(cmd, "file="+a+",index="+strconv.Itoa(index)+",media=disk,if=virtio")
				}
				index++
			} else {
				fmt.Println("diskmode error!!")
				fmt.Println("diskmode is " + strconv.Itoa(diskmode))
			}
		}
	}
	return cmd
}

//Generate Network Command
func GenerateNetworkCmd(net string) []string {
	fmt.Println("GenerateNetworkCommand")
	data := strings.Split(net, ",")
	fmt.Println(data)
	mode, err := strconv.Atoi(data[0])
	if err != nil {
		log.Println(err)
	}
	var bridge []string
	var mac []string
	for i, a := range data {
		if i > 0 {
			if i%2 == 1 {
				bridge = append(bridge, a)
			} else {
				mac = append(mac, a)
			}
		}
	}

	//verify bridge and mac array length
	if len(bridge) != len(mac) {
		fmt.Println("Warning: NetworkCount Error")
	}

	var cmd []string
	///etc/qemu/bridge.conf <- allow br0
	//1 Network
	//-net nic,macaddr=52:54:01:11:22:33 -net bridge,br=br0
	//2 Network
	//-net nic,macaddr=52:54:01:11:22:33 -net bridge,br=br0 -net nic,macaddr=52:54:02:11:22:33 -net bridge,br=br0
	//3 Network(Default)
	//-nic bridge,br=br1000,mac=.... , -nic bridge,br=br1100,mac=...,model=virtio

	//old setting
	//-netdev tap,helper=/usr/lib/qemu/qemu-bridge-helper,id=br0 -device virtio-net-pci,netdev=br0,id=net1,mac=52:54:85:98:60:93
	//-netdev tap,helper=/usr/lib/qemu/qemu-bridge-helper,id=br100 -device virtio-net-pci,netdev=br100,id=net2,mac=52:54:85:98:60:94
	//-netdev tap,helper=/usr/lib/qemu/qemu-bridge-helper,id=br200 -device virtio-net-pci,netdev=br200,id=net3,mac=52:54:85:98:60:95

	//new setting
	//-nic bridge,br=br1000,mac=.... ,model=e1000 -nic bridge,br=br1100,mac=...,model=virtio
	//not tested for multiple nic enviroment...

	if mode == 0 {
		//e1000
		for i, m := range mac {
			cmd = append(cmd, "-nic")
			cmd = append(cmd, "bridge,br="+bridge[i]+",mac="+m+",model=e1000")
		}
	} else if mode == 1 {
		//virtio
		for i, m := range mac {
			cmd = append(cmd, "-nic")
			cmd = append(cmd, "bridge,br="+bridge[i]+",mac="+m+",model=virtio")
		}
	}

	fmt.Printf("GenerateNetwork: ")
	fmt.Println(cmd)
	return cmd
}
