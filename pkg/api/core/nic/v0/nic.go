package v0

import (
	"encoding/xml"
	"fmt"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"log"
	"sort"
	"strconv"
	"strings"
)

var maxMAC = 65535
var startMAC = 10

type NICHandler struct {
	Conn    *libvirt.Connect
	Input   vm.VirtualMachine
	Address *vm.Address
}

func NewNICHandler(handler NICHandler) *NICHandler {
	return &NICHandler{Conn: handler.Conn, Input: handler.Input, Address: handler.Address}
}

func Network() {
	//
	//// Template
	//nic := []libVirtXml.DomainInterface{
	//	{
	//		MAC: &libVirtXml.DomainInterfaceMAC{
	//			Address: "00:00:00:00:00:00:01",
	//		},
	//		Model: &libVirtXml.DomainInterfaceModel{Type: "e1000e"},
	//		Address: &libVirtXml.DomainAddress{
	//			PCI: &libVirtXml.DomainAddressPCI{
	//				Domain:   nil,
	//				Bus:      nil,
	//				Slot:     nil,
	//				Function: nil,
	//			},
	//		},
	//	},
	//}
	//
	//var source libVirtXml.DomainInterfaceSource
	//
	//// NAT (Default)
	//source = libVirtXml.DomainInterfaceSource{
	//	Network: &libVirtXml.DomainInterfaceSourceNetwork{
	//		Network: "default",
	//	},
	//}
	//
	//// Bridge
	//source = libVirtXml.DomainInterfaceSource{
	//	Network: &libVirtXml.DomainInterfaceSourceNetwork{
	//		Bridge: "default",
	//	},
	//}
	//
	//// macvtap (bridge)
	//source = libVirtXml.DomainInterfaceSource{
	//	Direct: &libVirtXml.DomainInterfaceSourceDirect{
	//		Dev:  "vmnet1",
	//		Mode: "bridge",
	//	},
	//}
	//
	//// macvtap (vepa)
	//source = libVirtXml.DomainInterfaceSource{
	//	Direct: &libVirtXml.DomainInterfaceSourceDirect{
	//		Dev:  "vmnet1",
	//		Mode: "vepa",
	//	},
	//}
	//
	//// macvtap (private)
	//source = libVirtXml.DomainInterfaceSource{
	//	Direct: &libVirtXml.DomainInterfaceSourceDirect{
	//		Dev:  "vmnet1",
	//		Mode: "private",
	//	},
	//}
	//
	//// macvtap (passthrough)
	//source = libVirtXml.DomainInterfaceSource{
	//	Direct: &libVirtXml.DomainInterfaceSourceDirect{
	//		Dev:  "vmnet1",
	//		Mode: "passthrough",
	//	},
	//}
}

func (h *NICHandler) generateMac() (string, error) {
	var macs []int

	doms, err := h.Conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return "", err
	}

	// Todo:
	for _, dom := range doms {
		data := libVirtXml.Domain{}
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &data)

		if len(data.Devices.Interfaces) != 0 {
			for _, tmp := range data.Devices.Interfaces {
				mac := strings.Split(tmp.MAC.Address, ":")
				if (mac[0] + mac[1]) == "5254" {
					v, _ := strconv.ParseInt(mac[4]+mac[5], 16, 0)
					macs = append(macs, int(v))
					log.Println(v)
				}
			}
		}
	}

	//昇順に並び替える
	sort.Ints(macs)

	//startMACを定義
	macIndex := startMAC

	for _, m := range macs {
		//Port番号が上限に達する場合、エラーを返す
		if maxMAC <= macIndex {
			return "", fmt.Errorf("Error: max mac address ")
		}
		if macIndex < m {
			break
		}
		macIndex++
	}

	macIndex1 := macIndex / 256
	macIndex2 := macIndex % 256

	// macアドレスを10進数から16進数に変換し、結合
	mac := fmt.Sprintf("%s:%.2x:%.2x", config.Conf.Node.MAC, macIndex1, macIndex2)

	return mac, nil
}
