package v0

import (
	"fmt"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"log"
)

func XmlGenerate(input vm.VirtualMachine) ([]libVirtXml.DomainDisk, vm.Address, error) {

	var disks []libVirtXml.DomainDisk

	//countの定義＆初期化
	var virtIOCount uint = 0
	var otherCount uint = 0

	var pciAddressCount uint = 0
	var diskAddressCount uint = 0

	log.Println(input.Storage)

	for _, storageTmp := range input.Storage {
		if storageTmp.Path == "" {
			return nil, vm.Address{}, fmt.Errorf("black: storage path")
		}

		var tmpCount uint
		var tmpAddressCount uint

		// VirtIOの場合はVirtIO Countに数字を代入＋加算する
		if storageTmp.Type == 10 {
			tmpCount = virtIOCount
			tmpAddressCount = pciAddressCount
			pciAddressCount++
			virtIOCount++
		} else {
			tmpCount = otherCount
			// VirtIO以外の場合はOther Countに数字を代入＋加算する
			tmpAddressCount = diskAddressCount
			otherCount++
			diskAddressCount++
		}

		fmt.Println(storageTmp)

		disks = append(disks, *generateTemplate(storage.GenerateStorageXml{
			Storage:       storageTmp,
			Number:        tmpCount,
			AddressNumber: tmpAddressCount,
		}))
	}

	log.Println(disks)

	return disks, vm.Address{
		PCICount:  pciAddressCount,
		DiskCount: diskAddressCount,
	}, nil

}

func generateTemplate(xmlStruct storage.GenerateStorageXml) *libVirtXml.DomainDisk {
	//デフォルトはブートディスク(VirtIO)

	domDisk := libVirtXml.DomainDisk{}
	var dev string
	// CDROM
	if xmlStruct.Storage.Type == 0 {
		dev = "sda"

		domDisk.Target = &libVirtXml.DomainDiskTarget{Bus: "sata"}
		domDisk.Address = &libVirtXml.DomainAddress{
			Drive: &libVirtXml.DomainAddressDrive{
				Controller: &[]uint{0}[0],
				Bus:        &[]uint{0}[0],
				Target:     &[]uint{0}[0],
				Unit:       &[]uint{xmlStruct.AddressNumber}[0],
			},
		}
		// Boot Disk(SATA)
	} else if xmlStruct.Storage.Type == 11 {
		dev = "sda"

		domDisk.Target = &libVirtXml.DomainDiskTarget{Bus: "sata"}
		domDisk.Address = &libVirtXml.DomainAddress{
			Drive: &libVirtXml.DomainAddressDrive{
				Controller: &[]uint{0}[0],
				Bus:        &[]uint{0}[0],
				Target:     &[]uint{0}[0],
				Unit:       &[]uint{xmlStruct.AddressNumber}[0],
			},
		} // Boot Disk(IDE)
	} else if xmlStruct.Storage.Type == 12 {
		dev = "sda"

		domDisk.Target = &libVirtXml.DomainDiskTarget{Bus: "ide"}
		domDisk.Address = &libVirtXml.DomainAddress{
			Drive: &libVirtXml.DomainAddressDrive{
				Controller: &[]uint{0}[0],
				Bus:        &[]uint{0}[0],
				Target:     &[]uint{0}[0],
				Unit:       &[]uint{xmlStruct.AddressNumber}[0],
			},
		}
	} else {
		dev = "vda"

		domDisk.Target = &libVirtXml.DomainDiskTarget{Bus: "virtio"}
		domDisk.Address = &libVirtXml.DomainAddress{
			PCI: &libVirtXml.DomainAddressPCI{
				Domain:   &[]uint{0}[0],
				Bus:      &[]uint{xmlStruct.AddressNumber}[0],
				Slot:     &[]uint{0}[0],
				Function: &[]uint{0}[0],
			},
		}
	}

	if xmlStruct.Storage.ReadOnly {
		domDisk.ReadOnly = &libVirtXml.DomainDiskReadOnly{}
	}

	domDisk.Target = &libVirtXml.DomainDiskTarget{Dev: dev[0:2] + string(dev[2]+uint8(xmlStruct.Number))}
	// Driver
	domDisk.Driver = &libVirtXml.DomainDiskDriver{
		Name: "qemu",
		Type: storage.GetDriverName(xmlStruct.Storage.FileType),
	}
	// File Path
	domDisk.Source = &libVirtXml.DomainDiskSource{
		File: &libVirtXml.DomainDiskSourceFile{File: xmlStruct.Storage.Path},
	}

	return &domDisk
}
