package v0

import (
	"fmt"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"log"
	"strconv"
)

func test() {
	var drive uint
	drive = 0
	name := "test"
	log.Println("Start1")

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("failed to connect to qemu")
	}
	defer conn.Close()

	mem, err := conn.GetFreeMemory()
	if err != nil {
		panic(err)
	}
	log.Println("mem: " + strconv.Itoa(int(mem)))
	tmp, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		panic(err)
	}
	log.Println(tmp)

	for _, dom := range tmp {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}

		uuid, err := dom.GetUUIDString()
		if err == nil {
			fmt.Printf("  %s\n", uuid)
		}

		//dom.Shutdown()
		//dom.Create()
		if name == "name" {
			fmt.Println("testtesttest")
			dom.Undefine()
		}

		dom.Free()
	}
	tmp, err = conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		panic(err)
	}

	for _, dom := range tmp {
		dom.Shutdown()
		dom.Free()
	}

	log.Println("Start2")

	domcfg := &libVirtXml.Domain{
		Type:   "kvm",
		Memory: &libVirtXml.DomainMemory{Value: 2048, Unit: "MB", DumpCore: "on"},
		VCPU:   &libVirtXml.DomainVCPU{Value: 1},
		CPU:    &libVirtXml.DomainCPU{Mode: "host-model"},
		//UUID:        uuid.Must(uuid.NewV4()).String(),
		Name:        name,
		Title:       name,
		Description: name,
		Devices: &libVirtXml.DomainDeviceList{
			Disks: []libVirtXml.DomainDisk{
				{
					Source:  &libVirtXml.DomainDiskSource{File: &libVirtXml.DomainDiskSourceFile{File: "/home/yonedayuto/test/test.qcow2"}},
					Target:  &libVirtXml.DomainDiskTarget{Dev: "hda", Bus: "ide"},
					Alias:   &libVirtXml.DomainAlias{Name: "ide0-0-0"},
					Address: &libVirtXml.DomainAddress{Drive: &libVirtXml.DomainAddressDrive{Controller: &drive, Bus: &drive, Target: &drive, Unit: &drive}},
				},
			},
			Interfaces: []libVirtXml.DomainInterface{
				{
					MAC: &libVirtXml.DomainInterfaceMAC{
						Address: "", Type: "", Check: ""},
					Source: &libVirtXml.DomainInterfaceSource{
						Network: &libVirtXml.DomainInterfaceSourceNetwork{
							Network: "default",
						},
					},
					Model:   &libVirtXml.DomainInterfaceModel{Type: "e1000e"},
					Address: nil,
				},
			},
		},
		OS: &libVirtXml.DomainOS{
			BootDevices: []libVirtXml.DomainBootDevice{
				libVirtXml.DomainBootDevice{
					Dev: "hd",
				},
			},
			Kernel: "",
			//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
			//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
			Type: &libVirtXml.DomainOSType{
				Arch: "x86_64",
				Type: "hvm",
			},
		},
	}

	log.Println("End")

	xml, err := domcfg.Marshal()
	if err != nil {
		panic(err)
	}
	domain, err := conn.DomainDefineXML(xml)
	if err != nil {
		panic(err)
	}
	createDomain, err := conn.DomainCreateXML(xml, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(xml)
	fmt.Println(domain)
	fmt.Println(createDomain)
}
