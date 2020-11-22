package v0

import (
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	nic "github.com/vmmgr/node/pkg/api/core/nic/v0"
	storage "github.com/vmmgr/node/pkg/api/core/storage/v0"
	"github.com/vmmgr/node/pkg/api/core/tool/gen"
	"github.com/vmmgr/node/pkg/api/core/vm"
)

func xmlGenerate(input vm.VirtualMachine) (*libVirtXml.Domain, error) {

	uuid, err := gen.GenerateUUID()
	if err != nil {
		return nil, err
	}

	domCfg := &libVirtXml.Domain{
		Type: "kvm",
		Memory: &libVirtXml.DomainMemory{
			Value:    input.Memory,
			Unit:     "MB",
			DumpCore: "on",
		},
		VCPU:        &libVirtXml.DomainVCPU{Value: input.VCPU},
		CPU:         &libVirtXml.DomainCPU{Mode: vm.GetCPUMode(input.CPUMode)},
		UUID:        uuid,
		Name:        input.Name,
		Title:       input.Name,
		Description: input.Name,
		OS: &libVirtXml.DomainOS{
			BootDevices: []libVirtXml.DomainBootDevice{{Dev: "hd"}},
			Kernel:      "",
			//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
			//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
			Type: &libVirtXml.DomainOSType{
				Arch:    vm.GetArchConvert(input.OS.Arch),
				Machine: "pc-q35-4.2",
				Type:    "hvm",
			},
		},
		Devices: &libVirtXml.DomainDeviceList{
			Graphics: []libVirtXml.DomainGraphic{
				{
					VNC: &libVirtXml.DomainGraphicVNC{
						Port:      int(input.VNCPort),
						WebSocket: int(input.WebSocketPort),
						Keymap:    input.KeyMap,
						Listen:    "address",
					},
				},
			},
		},
	}

	// storage xmlの生成
	disks, address, err := storage.XmlGenerate(input)
	if err != nil {
		return nil, err
	}

	// nic xmlの生成
	nics, address, err := nic.XmlGenerate(input, address)
	if err != nil {
		return nil, err
	}

	// DomainDeviceListをdomCfgに挿入
	domCfg.Devices = &libVirtXml.DomainDeviceList{
		Disks:      disks,
		Interfaces: nics,
	}

	return domCfg, nil
}
