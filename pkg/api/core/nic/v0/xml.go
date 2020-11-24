package v0

import (
	"fmt"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/nic"
)

func (h *NICHandler) XmlGenerate() ([]libVirtXml.DomainInterface, error) {
	var nics []libVirtXml.DomainInterface

	var usedMAC []string

	for _, nicTmp := range h.VM.NIC {
		if nicTmp.MAC == "" {
			mac, err := h.generateMac(usedMAC)
			if err != nil {
				return nics, fmt.Errorf("MAC Address Generate Error: %s", err)
			}
			usedMAC = append(usedMAC, mac)
			nicTmp.MAC = mac
		}

		h.Address.PCICount++

		nics = append(nics, *generateTemplate(nic.GenerateNICXml{
			NIC:           nicTmp,
			AddressNumber: h.Address.PCICount,
		}))
	}

	return nics, nil
}

func generateTemplate(xmlStruct nic.GenerateNICXml) *libVirtXml.DomainInterface {
	//デフォルトはブートディスク(VirtIO)

	domNIC := libVirtXml.DomainInterface{}

	// Bridge
	if xmlStruct.NIC.Type == 0 {
		// defaultでもいけるかもしれない（要確認必要）
		domNIC.Source = &libVirtXml.DomainInterfaceSource{
			Bridge: &libVirtXml.DomainInterfaceSourceBridge{
				Bridge: xmlStruct.NIC.Device,
			},
		}
		// NAT
	} else if xmlStruct.NIC.Type == 1 {
		// defaultでもいけるかもしれない（要確認必要）
		domNIC.Source = &libVirtXml.DomainInterfaceSource{
			Network: &libVirtXml.DomainInterfaceSourceNetwork{
				Network: xmlStruct.NIC.Device,
			},
		}

		// macvtap
	} else if xmlStruct.NIC.Type == 2 {
		domNIC.Source = &libVirtXml.DomainInterfaceSource{
			Direct: &libVirtXml.DomainInterfaceSourceDirect{
				Dev:  xmlStruct.NIC.Device,
				Mode: nic.GetModeName(xmlStruct.NIC.Mode),
			},
		}
	}

	//Driver
	domNIC.Model = &libVirtXml.DomainInterfaceModel{
		Type: nic.GetDriverName(xmlStruct.NIC.Driver),
	}
	//MAC
	domNIC.MAC = &libVirtXml.DomainInterfaceMAC{
		Address: xmlStruct.NIC.MAC,
	}
	//PCI Address
	domNIC.Address = &libVirtXml.DomainAddress{
		PCI: &libVirtXml.DomainAddressPCI{
			Domain:   &[]uint{0}[0],
			Bus:      &[]uint{xmlStruct.AddressNumber}[0],
			Slot:     &[]uint{0}[0],
			Function: &[]uint{0}[0],
		},
	}

	return &domNIC
}
