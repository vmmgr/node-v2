package v0

import (
	"fmt"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"log"
)

func (h *PCIHandler) GetStr() ([]string, error) {
	var xmls []string

	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_PCI_DEV)
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

func (h *PCIHandler) Get() ([]libVirtXml.NodeDevice, error) {
	var xmls []libVirtXml.NodeDevice

	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_PCI_DEV)
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

		log.Println(xmlStr)
		//var a libVirtXml.NodeDevice
		xml := libVirtXml.NodeDevice{}
		xml.Unmarshal(xmlStr)
		xmls = append(xmls, xml)
	}
	return xmls, nil
}

func (h *PCIHandler) Detach(name string) error {
	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_PCI_DEV)
	if err != nil {
		log.Println(err)
	}

	for _, dev := range nodeDevice {
		xmlStr, err := dev.GetXMLDesc(0)
		if err != nil {
			log.Println(err)
			return err
		}
		xml := libVirtXml.NodeDevice{}
		xml.Unmarshal(xmlStr)
		if xml.Name == name {
			dev.Detach()
			return nil
		}
	}
	return fmt.Errorf("Error: PCI Device is not found ... ")
}

func (h *PCIHandler) Destroy(name string) error {
	nodeDevice, err := h.Conn.ListAllNodeDevices(libvirt.CONNECT_LIST_NODE_DEVICES_CAP_PCI_DEV)
	if err != nil {
		log.Println(err)
	}

	for _, dev := range nodeDevice {
		xmlStr, err := dev.GetXMLDesc(0)
		if err != nil {
			log.Println(err)
			return err
		}
		xml := libVirtXml.NodeDevice{}
		xml.Unmarshal(xmlStr)
		if xml.Name == name {
			dev.Destroy()
			return nil
		}
	}
	return fmt.Errorf("Error: PCI Device is not found ... ")
}
