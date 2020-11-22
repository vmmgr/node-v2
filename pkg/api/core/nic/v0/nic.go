package v0

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
