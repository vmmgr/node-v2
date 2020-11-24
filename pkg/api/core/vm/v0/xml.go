package v0

import (
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	nic "github.com/vmmgr/node/pkg/api/core/nic/v0"
	storage "github.com/vmmgr/node/pkg/api/core/storage/v0"
	"github.com/vmmgr/node/pkg/api/core/tool/gen"
	"github.com/vmmgr/node/pkg/api/core/vm"
)

func (h *VMHandler) xmlGenerate() (*libVirtXml.Domain, error) {

	uuid, err := gen.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// storage xmlの生成
	hStorage := storage.NewStorageHandler(storage.StorageHandler{
		Conn:    h.Conn,
		VM:      h.VM,
		Address: &vm.Address{PCICount: 0, DiskCount: 0},
	})
	disks, err := hStorage.XmlGenerate()
	if err != nil {
		return nil, err
	}

	// nic xmlの生成
	hNIC := nic.NewNICHandler(nic.NICHandler{
		Conn:    h.Conn,
		VM:      h.VM,
		Address: hStorage.Address,
	})
	nics, err := hNIC.XmlGenerate()
	if err != nil {
		return nil, err
	}

	domCfg := &libVirtXml.Domain{
		Type: "kvm",
		Memory: &libVirtXml.DomainMemory{
			Value:    h.VM.Memory,
			Unit:     "MB",
			DumpCore: "on",
		},
		VCPU:        &libVirtXml.DomainVCPU{Value: h.VM.VCPU},
		CPU:         &libVirtXml.DomainCPU{Mode: vm.GetCPUMode(h.VM.CPUMode)},
		UUID:        uuid,
		Name:        h.VM.Name,
		Title:       h.VM.Name,
		Description: h.VM.Name,
		OS: &libVirtXml.DomainOS{
			BootDevices: []libVirtXml.DomainBootDevice{{Dev: "hd"}},
			Kernel:      "",
			//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
			//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
			Type: &libVirtXml.DomainOSType{
				Arch:    vm.GetArchConvert(h.VM.OS.Arch),
				Machine: "pc-q35-4.2",
				Type:    "hvm",
			},
		},
		Devices: &libVirtXml.DomainDeviceList{
			Emulator: "/usr/bin/qemu-system-x86_64",
			Inputs: []libVirtXml.DomainInput{
				{Type: "mouse", Bus: "ps2"},
				{Type: "keyboard", Bus: "ps2"},
			},
			Graphics: []libVirtXml.DomainGraphic{
				{
					VNC: &libVirtXml.DomainGraphicVNC{
						Port:      int(h.VM.VNCPort),
						WebSocket: int(h.VM.WebSocketPort),
						Keymap:    h.VM.KeyMap,
						Listen:    "0.0.0.0",
					},
				},
			},
			Disks:      disks,
			Interfaces: nics,
		},
	}

	return domCfg, nil
}
