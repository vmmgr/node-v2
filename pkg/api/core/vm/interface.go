package vm

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/nic"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/tool/cloudinit"
)

type VirtualMachine struct {
	Name           string              `json:"name"`
	UUID           string              `json:"uuid"`
	Memory         uint                `json:"memory"`
	CPUMode        uint                `json:"cpu_mode"` //0:custom 1:host-model 2:pass-through
	VCPU           uint                `json:"vcpu"`
	OS             OS                  `json:"os"`
	VNCPort        uint                `json:"vnc_port"`
	WebSocketPort  uint                `json:"websocket_port"`
	KeyMap         string              `json:"keymap"`
	NIC            []nic.NIC           `json:"nic"`
	Storage        []storage.VMStorage `json:"storage"`
	CloudInit      cloudinit.CloudInit `json:"cloudinit"`
	CloudInitApply bool                `json:"cloudinit_apply"`
	Stat           uint                `json:"stat"`
}

type Detail struct {
	VM   libvirtxml.Domain `json:"vm"`
	Stat uint              `json:"stat"`
}

type OS struct {
	Boot   string `json:"boot"`
	Kernel string `json:"kernel"`
	Arch   uint   `json:"arch"`
	Type   string `json:"type"`
}

type Address struct {
	PCICount  uint
	DiskCount uint
}

type VirtualMachineStop struct {
	Force bool `json:"force"`
}

func GetCPUMode(mode uint) string {
	//デフォルトではcustomModeになる
	//https://access.redhat.com/documentation/ja-jp/red_hat_enterprise_linux/6/html/virtualization_administration_guide/sect-libvirt-dom-xml-cpu-model-top

	if mode == 1 {
		return "host-model"
	} else if mode == 2 {
		return "host-passthrough"
	}
	return "custom"
}

func GetArchConvert(arch uint) string {
	if arch == 32 {
		return "i686"
	}
	return "x86_64"
}

type ResultStatus struct {
	Status    int    `json:"status"`
	StatusStr string `json:"status_str"`
}
