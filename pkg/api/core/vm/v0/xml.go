package v0

import (
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	nic "github.com/vmmgr/node/pkg/api/core/nic/v0"
	storageInterface "github.com/vmmgr/node/pkg/api/core/storage"
	storage "github.com/vmmgr/node/pkg/api/core/storage/v0"
	"github.com/vmmgr/node/pkg/api/core/tool/cloudinit"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"github.com/vmmgr/node/pkg/api/core/tool/file"
	"github.com/vmmgr/node/pkg/api/core/tool/gen"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"log"
	"os"
	"strconv"
)

func (h *VMHandler) xmlGenerate() (*libVirtXml.Domain, error) {

	uuid, err := gen.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// nic xmlの生成
	hNIC := nic.NewNICHandler(nic.NICHandler{
		Conn:    h.Conn,
		VM:      h.VM,
		Address: &vm.Address{PCICount: 0, DiskCount: 0},
	})
	nics, err := hNIC.XmlGenerate()
	if err != nil {
		return nil, err
	}

	// CloudInit周りの処理
	if h.VM.CloudInitApply {
		directory := config.Conf.Storage[0].Path + "/" + uuid

		if !file.ExistsCheck(directory) {
			if err := os.Mkdir(directory, 0755); err != nil {
				log.Println(err)
				return nil, err
			}
		}

		for i, a := range nics {
			h.VM.CloudInit.NetworkConfig.Config[i].MacAddress = a.MAC.Address
			h.VM.CloudInit.NetworkConfig.Config[i].Name = "eth" + strconv.Itoa(i)
		}

		log.Println(h.VM.CloudInit.NetworkConfig)
		log.Println(h.VM.CloudInit.UserData)

		hCloudInit := cloudinit.NewCloudInitHandler(cloudinit.CloudInit{
			DirPath:       directory,
			MetaData:      cloudinit.MetaData{InstanceID: h.VM.UUID, LocalHostName: h.VM.Name},
			UserData:      h.VM.CloudInit.UserData,
			NetworkConfig: h.VM.CloudInit.NetworkConfig,
		})

		err := hCloudInit.Generate()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		h.VM.Storage = append(h.VM.Storage, storageInterface.VMStorage{
			Type:     1,
			Path:     directory + "/cloudinit.img",
			ReadOnly: true,
		})
	}

	log.Println("test1")
	log.Println(hNIC.Address)

	// storage xmlの生成
	hStorage := storage.NewStorageHandler(storage.StorageHandler{
		Conn:    h.Conn,
		VM:      h.VM,
		Address: hNIC.Address,
	})
	disks, err := hStorage.XmlGenerate()
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
		Features: &libVirtXml.DomainFeatureList{
			ACPI: &libVirtXml.DomainFeature{},
			APIC: &libVirtXml.DomainFeatureAPIC{},
		},
		OS: &libVirtXml.DomainOS{
			BootDevices: []libVirtXml.DomainBootDevice{{Dev: "hd"}},
			Kernel:      "",
			//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
			//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
			Type: &libVirtXml.DomainOSType{
				Arch:    vm.GetArchConvert(h.VM.OS.Arch),
				Machine: config.Conf.Node.Machine,
				Type:    "hvm",
			},
		},
		Devices: &libVirtXml.DomainDeviceList{
			Emulator: config.Conf.Node.Emulator,
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
			Videos: []libVirtXml.DomainVideo{
				{
					Model: libVirtXml.DomainVideoModel{
						Type:    "qxl",
						Heads:   1,
						Ram:     65536,
						VRam:    65536,
						VGAMem:  16384,
						Primary: "yes",
					},
					Alias: &libVirtXml.DomainAlias{Name: "video0"},
				},
			},
			Disks:      disks,
			Interfaces: nics,
		},
	}

	return domCfg, nil
}
