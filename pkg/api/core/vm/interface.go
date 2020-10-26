package vm

import (
	"github.com/vmmgr/node/pkg/api/core/nic"
	"github.com/vmmgr/node/pkg/api/core/storage"
)

type VirtualMachine struct {
	Name    string            `json:"name"`
	UUID    string            `json:"uuid"`
	Memory  uint              `json:"memory"`
	CPUMode uint              `json:"cpu_mode"`
	VCPU    uint              `json:"vcpu"`
	OS      OS                `json:"os"`
	NIC     []nic.NIC         `json:"nic"`
	Storage []storage.Storage `json:"storage"`
	Stat    uint              `json:"stat"`
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
