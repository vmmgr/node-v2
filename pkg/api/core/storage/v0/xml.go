package v0

import (
	"fmt"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/storage"
)

func (h *StorageHandler) XmlGenerate() ([]libVirtXml.DomainDisk, error) {

	var disks []libVirtXml.DomainDisk

	//countの定義＆初期化
	var virtIOCount uint = 0
	var otherCount uint = 0

	for _, storageTmp := range h.VM.Storage {
		if storageTmp.Path == "" {
			return nil, fmt.Errorf("black: storage path")
		}

		var Number uint
		var AddressNumber uint

		// VirtIOの場合はVirtIO Countに数字を代入＋加算する
		if storageTmp.Type == 10 {
			Number = virtIOCount
			AddressNumber = h.Address.PCICount
			h.Address.PCICount++
			virtIOCount++
		} else {
			Number = otherCount
			// VirtIO以外の場合はOther Countに数字を代入＋加算する
			AddressNumber = h.Address.DiskCount
			otherCount++
			h.Address.DiskCount++
		}

		disks = append(disks, *generateTemplate(storage.GenerateStorageXml{
			Storage:       storageTmp,
			Number:        Number,
			AddressNumber: AddressNumber,
		}))
	}

	return disks, nil
}

func generateTemplate(xmlStruct storage.GenerateStorageXml) *libVirtXml.DomainDisk {
	//デフォルトはブートディスク(VirtIO)

	domDisk := libVirtXml.DomainDisk{}
	var dev string
	var bus string
	// CDROM
	if xmlStruct.Storage.Type == 1 {
		dev = "sda"
		bus = "sata"
		domDisk.Device = "cdrom"
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
		bus = "sata"
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
		bus = "ide"

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
		bus = "virtio"
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

	domDisk.Target = &libVirtXml.DomainDiskTarget{Bus: bus, Dev: dev[0:2] + string(dev[2]+uint8(xmlStruct.Number))}
	// Driver
	if xmlStruct.Storage.Type == 1 || xmlStruct.Storage.Type == 2 {
		domDisk.Driver = &libVirtXml.DomainDiskDriver{
			Name: "qemu",
			Type: "raw",
		}
	} else {
		domDisk.Driver = &libVirtXml.DomainDiskDriver{
			Name: "qemu",
			Type: storage.GetExtensionName(xmlStruct.Storage.FileType),
		}
	}
	// File Path
	domDisk.Source = &libVirtXml.DomainDiskSource{
		File: &libVirtXml.DomainDiskSourceFile{File: xmlStruct.Storage.Path},
	}

	return &domDisk
}
