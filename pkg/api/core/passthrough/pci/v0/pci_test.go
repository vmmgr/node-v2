package v0

import (
	"fmt"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/node/pkg/api/core/passthrough/pci"
	"os/exec"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	//options := []string{"-vmm", "-D"}
	options := []string{"-nnmm"}

	out, err := exec.Command("lspci", options...).Output()

	//cmd := exec.Command("lspci", options...)
	//
	//out := &bytes.Buffer{}
	//cmd.Stdout = out
	//err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	space := strings.Split(string(out), "\n")
	t.Log(len(space))

	count := 0
	var datas []pci.PCI

	for _, tmp := range space {
		var tmpPCI pci.PCI
		outArray := strings.Split(tmp, "\"")
		for _, oneTmpPCI := range outArray {
			if !(len(oneTmpPCI) == 1 && strings.Contains(oneTmpPCI, string(' '))) && !(len(oneTmpPCI) == 0) {
				if count == 0 {
					tmpPCI.ID = fmt.Sprintf("%s", oneTmpPCI)
				} else if count == 1 {
					tmpPCI.ClassName = fmt.Sprintf("%s", oneTmpPCI)
				} else if count == 2 {
					vendor := strings.Split(oneTmpPCI, "[")
					tmpPCI.VendorName = fmt.Sprintf("%s", vendor[0])
					tmpPCI.VendorID = fmt.Sprintf("%s", vendor[1][:len(vendor[1])-1])
				} else if count == 3 {
					device := strings.Split(oneTmpPCI, "[")
					tmpPCI.DeviceName = fmt.Sprintf("%s", device[0])
					tmpPCI.DeviceID = fmt.Sprintf("%s", device[1][:len(device[1])-1])
				} else if count == 4 {
					tmpPCI.Comment = fmt.Sprintf("%s", oneTmpPCI)
				}
				count++
			}
		}
		datas = append(datas, tmpPCI)
		t.Logf("\n")
		count = 0
	}

	t.Log(datas)

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	nodeDevice, err := conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_PCI_DEV)
	if err != nil {
		t.Fatal(err)
	}

	for _, dev := range nodeDevice {
		a, _ := dev.GetXMLDesc(0)
		t.Log(a)
	}
}
