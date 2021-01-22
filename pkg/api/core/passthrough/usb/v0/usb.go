package v0

import (
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"log"
)

func (h *USBHandler) GetStr() ([]string, error) {
	var xmls []string

	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_USB_DEV)
	if err != nil {
		log.Println(err)
		return []string{""}, err
	}

	for _, dev := range nodeDevice {
		xmlStr, err := dev.GetXMLDesc(0)
		if err != nil {
			log.Println(err)
			return []string{""}, err
		}
		xmls = append(xmls, xmlStr)
	}
	return xmls, nil
}

func (h *USBHandler) Get() ([]libVirtXml.NodeDevice, error) {
	var xmls []libVirtXml.NodeDevice

	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_USB_DEV)
	if err != nil {
		log.Println(err)
		return xmls, err
	}

	for _, dev := range nodeDevice {
		xmlStr, err := dev.GetXMLDesc(0)
		if err != nil {
			log.Println(err)
			return xmls, err
		}
		xml := libVirtXml.NodeDevice{}
		xml.Unmarshal(xmlStr)
		xmls = append(xmls, xml)
	}
	return xmls, nil
}

func (h *USBHandler) Destory() {
	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_PCI_DEV)
	if err != nil {
		log.Println(err)
	}

	for _, dev := range nodeDevice {
		xmlStr, err := dev.GetXMLDesc(0)
		if err != nil {
			log.Println(err)
		}
		xml := libVirtXml.NodeDevice{}
		xml.Unmarshal(xmlStr)
	}
}
